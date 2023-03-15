package ngrule

import "encoding/json"

type FactParam interface {
	GetName() string
	SetName(name string)
	GetType() string
	SetType(name string)
	IsValid() bool
	SetValue(v interface{})
	GetValue() interface{}
}

type ParamConstructor func() FactParam

type ParamArray struct {
	coll []FactParam
}

func (c *ParamArray) UnmarshalJSON(data []byte) error {
	// Parse as c temporary Condition object placeholder
	var p []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	err := json.Unmarshal(data, &p)
	if err != nil {
		return err
	}

	// Prepare collection
	c.coll = make([]FactParam, len(p))

	// Iterate Param item array
	for k, v := range p {
		// Parse by type
		con, ok := paramTC[v.Type]
		if !ok {
			return NewError(ErrParamUnsupported, "unsupported Param type")
		}

		// Construct item
		item := con()
		item.SetName(v.Name)
		item.SetType(v.Type)

		// Push to c
		c.coll[k] = item
	}

	return nil
}
