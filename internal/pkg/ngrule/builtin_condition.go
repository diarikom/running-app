package ngrule

import "fmt"

type PrimitiveCondition struct {
	ParameterName      string `json:"param"`
	ComparisonOperator string `json:"operator"`
	NextLogicOperator  string `json:"next_op"`
}

func (c *PrimitiveCondition) NextLogic() string {
	return RenderOperator(c.NextLogicOperator)
}

func (c *PrimitiveCondition) RenderExp(varName string, getterType string) string {
	return fmt.Sprintf("%s.Get%s(\"%s\") %s", varName, getterType, c.ParameterName, RenderOperator(c.ComparisonOperator))
}

type IntCondition struct {
	*PrimitiveCondition
	RefValue int64 `json:"ref_value"`
}

func (c *IntCondition) Render(varName string) string {
	return c.RenderExp(varName, "Int") + fmt.Sprintf(" %d", c.RefValue)
}

func (c *IntCondition) GetParams() []FactParam {
	p := NewIntParam()
	p.Name = c.ParameterName
	return []FactParam{p}
}

type FloatCondition struct {
	*PrimitiveCondition
	RefValue float64 `json:"ref_value"`
}

func (c *FloatCondition) Render(varName string) string {
	// TODO: Determine float number formatting
	floatFmt := "%f"
	return c.RenderExp(varName, "Float") + fmt.Sprintf(" "+floatFmt, c.RefValue)
}

func (c *FloatCondition) GetParams() []FactParam {
	p := NewFloatParam()
	p.Name = c.ParameterName
	return []FactParam{p}
}

type StringCondition struct {
	*PrimitiveCondition
	RefValue string `json:"ref_value"`
}

func (c *StringCondition) Render(varName string) string {
	return c.RenderExp(varName, "String") + " \"" + c.RefValue + "\""
}

func (c *StringCondition) GetParams() []FactParam {
	p := NewStringParam()
	p.Name = c.ParameterName
	return []FactParam{p}
}

type BooleanCondition struct {
	*PrimitiveCondition
	RefValue bool `json:"ref_value"`
}

func (c *BooleanCondition) Render(varName string) string {
	return c.RenderExp(varName, "Bool") + fmt.Sprintf(" %t", c.RefValue)
}

func (c *BooleanCondition) GetParams() []FactParam {
	p := NewBooleanParam()
	p.Name = c.ParameterName
	return []FactParam{p}
}

type BetweenIntCondition struct {
	ParameterPrefix   string `json:"param_prefix"`
	NextLogicOperator string `json:"next_op"`
}

func (c *BetweenIntCondition) Render(varName string) string {
	return fmt.Sprintf("(%s.GetInt(\"%s_ref\") >= %s.GetInt(\"%s_start\") && %s.GetInt(\"%s_ref\") <= %s.GetInt(\"%s_end\"))", varName, c.ParameterPrefix,
		varName, c.ParameterPrefix, varName, c.ParameterPrefix, varName, c.ParameterPrefix)
}

func (c *BetweenIntCondition) NextLogic() string {
	return RenderOperator(c.NextLogicOperator)
}

func (c *BetweenIntCondition) GetParams() []FactParam {
	pStart := NewIntParam()
	pStart.Name = c.ParameterPrefix + "_start"
	pEnd := NewIntParam()
	pEnd.Name = c.ParameterPrefix + "_end"
	pRef := NewIntParam()
	pRef.Name = c.ParameterPrefix + "_ref"
	return []FactParam{pStart, pEnd, pRef}
}

type GroupCondition struct {
	InnerConditions   ConditionArray `json:"conditions"`
	NextLogicOperator string         `json:"next_op"`
}

func (c *GroupCondition) Render(varName string) string {
	// Render inner condition
	inner, err := c.InnerConditions.Render(varName)
	if err != nil {
		return ""
	}

	// Format string
	return fmt.Sprintf("(%s)", inner)
}

func (c *GroupCondition) NextLogic() string {
	return RenderOperator(c.NextLogicOperator)
}

func (c *GroupCondition) GetParams() []FactParam {
	return c.InnerConditions.GetParams()
}
