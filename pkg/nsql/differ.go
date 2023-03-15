package nsql

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	// Tag keys
	DbTagKey   = "db"
	DiffTagKey = "diff"

	// DiffItem Tag Values
	IdentifierTag            = "id"
	IgnoredTag               = "-"
	RequiredTag              = "required"
	ConcurrencyControllerTag = "cc"

	// Separator
	CommaSeparator = ", "
	AndSeparator   = " AND "
)

type DiffComparator interface {
	DiffValue() interface{}
}

type Differ struct {
	tableName      string
	logTableName   string
	logPkCol       string
	changelogCol   string
	sourceType     reflect.Type
	idCols         []Field
	trackedCols    map[string]Field
	requiredCols   []Field
	anonCols       map[string][]Field
	concurrencyCol *Field
}

type DiffItem struct {
	Value interface{}
	Field
}

type Diff struct {
	Values      []DiffItem
	Ids         []DiffItem
	Count       int
	TrackedCols []string
}

type DifferOpt struct {
	Sample       interface{}
	TableName    string
	LogTableName string
	LogPKCol     string
	ChangelogCol string
}

type DifferQuery struct {
	Query string
	Args  []interface{}
}

// PrepareDiffer reflect sample types and create a reusable difference finder instance
func PrepareDiffer(opt DifferOpt) *Differ {
	// Reflect type
	t := reflect.TypeOf(opt.Sample)
	// Init new diff
	d := Differ{
		tableName:  opt.TableName,
		sourceType: t,
		anonCols:   make(map[string][]Field),
	}

	// Set log table name
	if opt.LogTableName == "" {
		d.logTableName = opt.TableName + "_log"
	} else {
		d.logTableName = opt.LogTableName
	}

	// Set log pk name
	if opt.LogPKCol == "" {
		d.logPkCol = "log_id"
	} else {
		d.logPkCol = opt.LogPKCol
	}

	// Set changelog name
	if opt.ChangelogCol == "" {
		d.changelogCol = "changelog"
	} else {
		d.changelogCol = opt.ChangelogCol
	}

	// Read tag
	err := d.analyzeDiffStruct(t)
	if err != nil {
		panic(fmt.Errorf("nsql: unable to analyze struct (%s)", err))
	}

	// Return differ
	return &d
}

func (d *Differ) GetFields(trackedColumns []string) []Field {
	// Init fields
	fields := make([]Field, 0)

	// Search by column
	haystack := d.trackedCols

	/// Iterate columns
	for _, needle := range trackedColumns {
		// Search
		if f, ok := haystack[needle]; ok {
			fields = append(fields, f)
		}
	}

	return fields
}

func (d *Differ) GetDiff(updatedInstance interface{}, overriddenKeys []string) (diff *Diff, err error) {
	// Reflect instance
	val := reflect.ValueOf(updatedInstance)

	// Get type
	t := val.Type()

	// Check with differ source
	if t != d.sourceType {
		err = errors.New("nsql: instance struct is different type with Differ")
		return nil, err
	}

	// Get fields from overridden keys
	fields := d.GetFields(overriddenKeys)

	// Init changes
	var changes []DiffItem
	var trackedCols []string

	// Get changed values
	for _, f := range fields {
		// Get key
		k := f.Name

		// Retrieve Field value
		fieldValue := val.FieldByName(k)

		// Retrieve  value instance
		fieldValInstance := fieldValue.Interface()

		// If anonymous, extract struct
		if f.Anonymous {
			// extract diff
			anonChanges := getAnonValues(fieldValue)
			// merge diff
			changes = append(changes, anonChanges...)
		} else {
			// If field is anonymous, retrieve
			changes = append(changes, DiffItem{
				Value: fieldValInstance,
				Field: f,
			})

			// If field is tracked, push columns to tracked cols
			if f.Tracked {
				trackedCols = append(trackedCols, f.Column)
			}
		}
	}

	// Get identifier values
	ids := getValues(d.idCols, val)

	// Get Required values
	requiredValues := getValues(d.requiredCols, val)

	// Count changes
	count := len(changes)

	// Merge required values and changes
	changes = append(changes, requiredValues...)

	// Create diff
	diff = &Diff{
		Values:      changes,
		Ids:         ids,
		Count:       count,
		TrackedCols: trackedCols,
	}

	return diff, nil
}

func (d *Differ) Compare(old, new interface{}, changedCols []string) (diff *Diff, err error) {
	// Reflect value
	oldVal := reflect.ValueOf(old)
	newVal := reflect.ValueOf(new)

	// check old and new type with source
	if oldVal.Type() != d.sourceType || newVal.Type() != d.sourceType {
		err = errors.New("nsql: old and new struct is different type with sources")
		return
	}

	// Check id
	if e := d.validateId(oldVal, newVal); e != nil {
		err = e
		return
	}

	// Track changes values by key changes
	changes, trackedCols := d.compareFields(oldVal, newVal, changedCols)

	// Get Required values
	requiredValues := getValues(d.requiredCols, newVal)

	// Get identifier values
	ids := getValues(d.idCols, newVal)

	// Count changes
	count := len(changes)

	// Merge required values and changes
	changes = append(changes, requiredValues...)

	// Create diff
	diff = &Diff{
		Values:      changes,
		Ids:         ids,
		Count:       count,
		TrackedCols: trackedCols,
	}

	return
}

func (d *Differ) UpdateQuery(diff *Diff) (q string, args []interface{}, err error) {
	// Join update query
	updateQuery, updateArgs := joinUpdateQuery(diff.Values, CommaSeparator)

	// Join where query
	whereQuery, whereArgs := joinUpdateQuery(diff.Ids, AndSeparator)

	// Assemble query
	q = fmt.Sprintf("UPDATE %s SET %s WHERE %s", d.tableName, updateQuery, whereQuery)

	// Merge arguments
	args = append(updateArgs, whereArgs...)
	return
}

func (d *Differ) UpdateQuerySafe(diff *Diff, oldVersion interface{}) (string, []interface{}, error) {
	q, args, err := d.UpdateQuery(diff)
	if err != nil {
		return "", nil, err
	}

	// If concurrency control is set, then add concurrency control in where clause
	if d.concurrencyCol != nil {
		q = fmt.Sprintf("%s AND %s = ?", q, d.concurrencyCol.Name)
		args = append(args, oldVersion)
	}

	// If concurrency control is available, set
	return q, args, nil
}

func (d *Differ) InsertLogQuery(diff *Diff, logId interface{}, changelog Changelog) (q string, args []interface{}, err error) {
	// Init args
	args = []interface{}{logId, changelog}

	// Initiate query
	q = "INSERT INTO " + d.logTableName + "(" + d.logPkCol + ", " + d.changelogCol + ", %s) VALUES (?, ?, %s)"

	// Join values
	insertCols, bindVars, insertArgs := joinInsertQuery(diff.Values, diff.Ids)

	// Assemble query
	q = fmt.Sprintf(q, insertCols, bindVars)

	// Merge arguments
	args = append(args, insertArgs...)

	return
}

func (d *Differ) analyzeDiffStruct(t reflect.Type) (err error) {
	// Init tmp vars
	var idCols, reqCols []Field
	var concurrencyCol *Field
	trackedCols := make(map[string]Field)

	// Iterate fields
	for i := 0; i < t.NumField(); i++ {
		// Get Field
		f := t.Field(i)

		// Get tags
		dbTag := f.Tag.Get(DbTagKey)
		diffTag := f.Tag.Get(DiffTagKey)

		// If db or diff tag is ignored, continue next Field
		if dbTag == IgnoredTag || diffTag == IgnoredTag {
			continue
		}

		// Set database Column mapping
		col := Field{Name: f.Name, Anonymous: f.Anonymous}
		col.setColumn(dbTag)

		// Split between tags
		tags := strings.Split(diffTag, ",")

		// Iterate tags
		tagCount := 0
		for _, vTag := range tags {
			// Group diff tag
			switch vTag {
			case IdentifierTag:
				idCols = append(idCols, col)
				tagCount++
			case RequiredTag:
				reqCols = append(reqCols, col)
				tagCount++
			case ConcurrencyControllerTag:
				concurrencyCol = &col
				tagCount++
			}
		}

		// If tag count is 0, then set as tracked columns
		if tagCount == 0 {
			// Get override mapping name if set
			var trackingTag string
			if diffTag != "" {
				trackingTag = diffTag
			} else {
				trackingTag = dbTag
			}
			// Set track flag
			col.Tracked = true

			// Push
			trackedCols[trackingTag] = col
		}

		// If Column is not anonymous, extract anonymous columns and map it
		if col.Anonymous {
			// Get type
			ct := f.Type

			// Extract fields
			switch ct.Kind() {
			case reflect.Struct, reflect.Ptr:
				// Get db fields
				dbFields := getDbFields(ct)
				// If db fields is available, set to anon Field map
				if len(dbFields) > 0 {
					d.anonCols[col.Name] = dbFields
				}
			default:
				return errors.New("nsql: anonymous Field must be a struct or ptr")
			}
		}
	}

	// If identifier Field is not set, return error
	if len(idCols) == 0 {
		return errors.New("nsql: no identifier Field is set")
	}

	// Set fields
	d.idCols = idCols
	d.requiredCols = reqCols
	d.trackedCols = trackedCols
	d.concurrencyCol = concurrencyCol
	return nil
}

func (d *Differ) compareFields(oldVal, newVal reflect.Value, changedCols []string) (changes []DiffItem, trackedCols []string) {
	for _, c := range changedCols {
		// Get key
		f, isTracked := d.trackedCols[c]
		if !isTracked {
			continue
		}

		k := f.Name

		// Retrieve Field value
		newFieldValue := newVal.FieldByName(k)
		oldInstance := oldVal.FieldByName(k).Interface()
		newInstance := newFieldValue.Interface()

		// Determine comparator
		var isDiff bool

		// Get comparator interface
		oc, ocOk := oldInstance.(DiffComparator)
		nc, ncOk := newInstance.(DiffComparator)

		if ncOk && ocOk {
			isDiff = oc.DiffValue() != nc.DiffValue()
		} else {
			// Else
			isDiff = oldInstance != newInstance
		}

		// If old and new value has differences, push to delta
		if isDiff {
			// If anonymous, extract struct
			if f.Anonymous {
				// extract diff
				anonChanges := getAnonValues(newFieldValue)
				// merge diff
				changes = append(changes, anonChanges...)
			} else {
				// If field is anonymous, retrieve
				changes = append(changes, DiffItem{
					Value: newInstance,
					Field: f,
				})

				// If field is tracked, push columns to tracked cols
				if f.Tracked {
					trackedCols = append(trackedCols, f.Column)
				}
			}
		}
	}
	return
}

// validateId validate identifier has match value between old and new instance
func (d *Differ) validateId(old, new reflect.Value) error {
	// Check if identifier is not match
	for _, v := range d.idCols {
		oldId := old.FieldByName(v.Name).Interface()
		newId := new.FieldByName(v.Name).Interface()
		if oldId != newId {
			return errors.New("nsql: old and new identifier value is not match")
		}
	}
	return nil
}

// equalQuery generates Column query with bindvar
func equalQuery(col string) string {
	return col + " = ?"
}

func getValues(fields []Field, newVal reflect.Value) (values []DiffItem) {
	for _, f := range fields {
		// Get key
		k := f.Name

		// Retrieve Field value
		fv := newVal.FieldByName(k)

		// If anonymous, extract struct
		if f.Anonymous {
			// extract diff
			anonValues := getAnonValues(fv)
			// merge diff
			values = append(values, anonValues...)
		} else {
			// If field is anonymous, retrieve
			values = append(values, DiffItem{
				Value: fv.Interface(),
				Field: f,
			})
		}
	}
	return
}

func getAnonValues(v reflect.Value) (changes []DiffItem) {
	// Get db fields
	fields := getDbFields(v.Type())

	// If value kind is pointer, get element
	if v.Type().Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Get values
	for _, f := range fields {
		// Get field by name
		fv := v.FieldByName(f.Name)

		// Get value
		diff := DiffItem{
			Value: fv.Interface(),
			Field: f,
		}

		changes = append(changes, diff)
	}

	return
}

func getDbFields(t reflect.Type) (r []Field) {
	// If kind is ptr, get types
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Iterate fields in type struct
	for i := 0; i < t.NumField(); i++ {
		// Get Field
		f := t.Field(i)

		// Get tag
		dbTag := f.Tag.Get(DbTagKey)

		// If db is ignored, continue next Field
		if dbTag == IgnoredTag {
			continue
		}

		// Init Column
		col := Field{Name: f.Name}

		// Set Column mapping
		col.setColumn(dbTag)

		// Push to result
		r = append(r, col)
	}
	return
}

// joinUpdateQuery assemble update query from fields and get arguments
func joinUpdateQuery(fields []DiffItem, sep string) (q string, args []interface{}) {
	// Init index
	i := 0
	// Get first fields
	ff := fields[i]
	// Set query
	q += equalQuery(ff.Column)
	// Set first args
	args = append(args, ff.Value)
	// increment index
	i++
	// Iterate updates
	for i < len(fields) {
		// Get fields
		f := fields[i]
		// Append separator where query
		q += sep + equalQuery(f.Column)
		// Append arguments
		args = append(args, f.Value)
		// Next index
		i++
	}
	return
}

func joinInsertQuery(values, ids []DiffItem) (q, bindvar string, args []interface{}) {
	fields := append(ids, values...)
	// Init index
	i := 0
	// Get first fields
	ff := fields[i]
	// append query and bindvars
	q += ff.Column
	bindvar += "?"
	// Set first args
	args = append(args, ff.Value)
	// increment index
	i++
	// Iterate updates
	for i < len(fields) {
		// Get fields
		f := fields[i]
		// Append separator where query
		q += CommaSeparator + f.Column
		bindvar += CommaSeparator + "?"
		// Append arguments
		args = append(args, f.Value)
		// Next index
		i++
	}
	return
}

type Changelog []string

func (e *Changelog) Scan(src interface{}) error {
	return ScanJSON(src, e)
}

func (e Changelog) Value() (driver.Value, error) {
	return json.Marshal(e)
}
