package ngrule

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

const RuleFormat = `rule %s "%s" salience %d {
  when
    %s
  then
    %s
}`

type Rule struct {
	Code         string         `json:"code"`
	Description  string         `json:"description"`
	Priority     int            `json:"priority"`
	Conditions   ConditionArray `json:"conditions"`
	Actions      ActionArray    `json:"actions"`
	Params       ParamArray     `json:"params"`
	VariableName string         `json:"-"`
	SkipRetract  bool           `json:"-"`
}

func (r *Rule) Render() (string, error) {
	// Render conditions
	conditionSyntax, err := r.Conditions.Render(r.VariableName)
	if err != nil {
		return "", err
	}

	// Render action
	actionSyntax, err := r.Actions.Render(r.VariableName)
	if err != nil {
		return "", err
	}

	// If skip retract is false, add Retract function to action
	if !r.SkipRetract {
		actionSyntax += fmt.Sprintf("\n    Retract(\"%s\");", r.Code)
	}

	// Render action
	return fmt.Sprintf(RuleFormat, r.Code, r.Description, r.Priority, conditionSyntax, actionSyntax), nil
}

func (r *Rule) GetParams() []FactParam {
	conditionParams := r.Conditions.GetParams()
	params := MergeParams(conditionParams, r.Params.coll)
	return params
}

func (r *Rule) GetTargets() []FactParam {
	return r.Actions.GetTargets()
}

// Scan implements the database/sql Scanner interface.
func (r *Rule) Scan(src interface{}) error {
	// Assert source to byte
	source, ok := src.([]byte)
	if !ok {
		return NewError(ErrSqlScan, "unable to scan non []byte type")
	}
	// Unmarshal to target
	err := json.Unmarshal(source, r)
	if err != nil {
		return err
	}
	return nil
}

// Value implements the database/sql/driver Valuer interface.
func (r *Rule) Value() (driver.Value, error) {
	return json.Marshal(r)
}
