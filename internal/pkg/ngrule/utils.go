package ngrule

import (
	"sort"
	"strconv"
	"strings"
)

const (
	StringType         = "string"
	StringArrayType    = "string_array"
	IntType            = "int"
	FloatType          = "float"
	BooleanType        = "boolean"
	GroupRef           = "group"
	BetweenIntRef      = "between_int"
	AssignActionType   = "assign"
	AddActionType      = "add"
	SubtractActionType = "subtract"
	MultiplyActionType = "multiply"
	DivideActionType   = "divide"
)

func RenderOperator(op string) string {
	switch strings.ToLower(op) {
	case "eq", "==":
		return "=="
	case "neq", "!=":
		return "=="
	case "gt", ">":
		return ">"
	case "gte", ">=":
		return ">="
	case "lt", "<":
		return "<"
	case "lte", "<=":
		return "<="
	case "and", "&&":
		return "&&"
	case "or", "||":
		return "||"
	}
	return ""
}

func ConvertInt(v interface{}) (int64, bool) {
	var result int64
	switch c := v.(type) {
	case int:
		result = int64(c)
	case int8:
		result = int64(c)
	case int16:
		result = int64(c)
	case int32:
		result = int64(c)
	case int64:
		result = c
	default:
		return 0, false
	}
	return result, true
}

func ConvertFloat(v interface{}) (float64, bool) {
	var result float64
	switch c := v.(type) {
	case float32:
		result = float64(c)
	case float64:
		result = c
	default:
		return 0, false
	}
	return result, true
}

type FactParamSorter []FactParam

func (f FactParamSorter) Len() int {
	return len(f)
}

func (f FactParamSorter) Less(i, j int) bool {
	return f[i].GetName() < f[j].GetName()
}

func (f FactParamSorter) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f FactParamSorter) GetIndexByName(name string) int {
	// Get index
	idx := sort.Search(len(f), func(i int) bool {
		return f[i].GetName() >= name
	})

	// If not found, return -1
	if f[idx].GetName() != name {
		return -1
	}

	// Return index
	return idx
}

func MergeParams(a, b []FactParam) []FactParam {
	// Cast to sorter type
	as := FactParamSorter(a)

	// If as is empty, return b
	if len(as) == 0 {
		return b
	}

	// Sort a if not sorted
	if !sort.IsSorted(as) {
		sort.Sort(as)
	}

	// Iterate b
	for _, v := range b {
		// Get index
		idx := as.GetIndexByName(v.GetName())

		// If not found, then push to a
		if idx == -1 {
			as = append(as, v)
		}
	}

	return as
}

// ParseInt64 converts interface into int64
func ParseInt64(input interface{}, defaultValue int64) int64 {
	switch v := input.(type) {
	case float64:
		return int64(v)
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int32:
		return int64(v)
	case string:
		tmp, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return defaultValue
		}
		return tmp
	}

	return defaultValue
}

// ParseInt converts interface into int
func ParseInt(input interface{}, defaultValue int) int {
	switch v := input.(type) {
	case float64:
		return int(v)
	case int8:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case string:
		tmp, err := strconv.ParseInt(v, 10, 32)
		if err != nil {
			return defaultValue
		}
		return int(tmp)
	}
	return defaultValue
}
