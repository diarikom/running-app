package ngrule

import (
	"encoding/json"
)

var conditionTC = map[string]ConditionConstructor{
	IntType: func() Condition {
		return &IntCondition{}
	},
	FloatType: func() Condition {
		return &FloatCondition{}
	},
	StringType: func() Condition {
		return &StringCondition{}
	},
	BooleanType: func() Condition {
		return &BooleanCondition{}
	},
	GroupRef: func() Condition {
		return &GroupCondition{
			InnerConditions: NewConditionArray(),
		}
	},
	BetweenIntRef: func() Condition {
		return &BetweenIntCondition{}
	},
}

type Condition interface {
	Render(varName string) string
	NextLogic() string
	GetParams() []FactParam
}

type ConditionConstructor func() Condition

type ConditionArray struct {
	coll []Condition
}

func (c *ConditionArray) Render(varName string) (string, error) {
	// Count condition
	end := len(c.coll) - 1

	// Iterate conditions
	var result string
	for k, v := range c.coll {
		// Render condition rule
		result += v.Render(varName)

		// If not end of collection, render next logic operator
		if k < end {
			nextOp := v.NextLogic()
			if nextOp == "" {
				return "", NewError(ErrNoNextLogicOp, "next logic operator is not provided")
			}

			result += " " + nextOp + " "
		}
	}

	return result, nil
}

func (c *ConditionArray) UnmarshalJSON(data []byte) error {
	// Parse as c temporary Condition object placeholder
	var p []struct {
		Type    string          `json:"type"`
		Options json.RawMessage `json:"options"`
	}
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	// Prepare collection
	c.coll = make([]Condition, len(p))

	// Iterate ConditionItem array
	for k, v := range p {
		// Parse by type
		conditionConstructor, ok := conditionTC[v.Type]
		if !ok {
			return NewError(ErrConditionUnsupported, "unsupported Condition type")
		}

		// Construct item
		item := conditionConstructor()

		// Parse json
		err = json.Unmarshal(v.Options, item)
		if err != nil {
			return err
		}

		// Push to c
		c.coll[k] = item
	}

	return nil
}

func (c *ConditionArray) AddType(name string, constructor ConditionConstructor) {
	conditionTC[name] = constructor
}

func (c *ConditionArray) RemoveType(name string) {
	// If type is available, delete
	_, ok := conditionTC[name]
	if ok {
		delete(conditionTC, name)
	}
}

func (c *ConditionArray) GetParams() []FactParam {
	// Get parameters from each condition
	params := make([]FactParam, 0)
	for _, v := range c.coll {
		// Get item parameters
		itemParams := v.GetParams()
		// Merge parameters
		params = MergeParams(params, itemParams)
	}
	return params
}

func NewConditionArray() ConditionArray {
	return ConditionArray{
		coll: make([]Condition, 0),
	}
}
