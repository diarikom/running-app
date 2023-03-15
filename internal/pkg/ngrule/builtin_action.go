package ngrule

import (
	"fmt"
	"strings"
)

const ParamValueRefPrefix = "$param:"

type FactMapAction struct {
	Target    string      `json:"target"`
	Operator  string      `json:"-"`
	Value     interface{} `json:"value"`
	ValueType string      `json:"value_type"`
}

func (a *FactMapAction) GetTargetName() string {
	return a.Target
}

func (a *FactMapAction) GetValue() interface{} {
	return a.Value
}

func (a *FactMapAction) GetTargets() []FactParam {
	pc, ok := paramTC[a.ValueType]
	if !ok {
		return nil
	}

	p := pc()
	p.SetName(a.Target)
	return []FactParam{p}
}

func (a *FactMapAction) Render(varName string) string {
	// Get value
	val := a.RenderValue(varName)

	// Render expression
	return fmt.Sprintf("%s.%s(\"%s\", %s)", varName, RenderMathFunc(a.Operator), a.Target, val)
}

func (a *FactMapAction) RenderValue(varName string) string {
	// Get value
	val := a.Value

	// Assert to string
	tmpStr, ok := val.(string)

	// If not a string, render
	if !ok {
		return fmt.Sprintf("%#v", val)
	}

	// If do not have params prefix, return string
	if hasParamPrefix := strings.HasPrefix(tmpStr, ParamValueRefPrefix); !hasParamPrefix {
		return fmt.Sprintf("%#v", val)
	}

	// Get parameter reference
	paramRef := strings.TrimPrefix(tmpStr, ParamValueRefPrefix)

	// Render value expression
	valExp := fmt.Sprintf("%s.%s(\"%s\")", varName, RenderGetterType(a.ValueType), paramRef)
	return valExp
}

func RenderGetterType(typeName string) string {
	switch typeName {
	case IntType:
		return "GetInt"
	case FloatType:
		return "GetFloat"
	case StringType:
		return "GetString"
	case BooleanType:
		return "GetBool"
	}
	return ""
}
