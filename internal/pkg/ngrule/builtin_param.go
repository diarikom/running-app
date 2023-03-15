package ngrule

import "strings"

type PrimitiveParam struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Valid bool   `json:"-"`
}

func (p PrimitiveParam) GetName() string {
	return p.Name
}

func (p *PrimitiveParam) SetName(name string) {
	p.Name = name
}

func (p PrimitiveParam) SetType(t string) {
	p.Type = t
}

func (p PrimitiveParam) GetType() string {
	return p.Type
}

func (p PrimitiveParam) IsValid() bool {
	return p.Valid
}

type StringParam struct {
	PrimitiveParam
	Value string `json:"value"`
}

func (s *StringParam) SetValue(v interface{}) {
	// Check value is an int
	switch c := v.(type) {
	case string:
		s.Value = c
		s.Valid = true
	default:
		s.Valid = false
	}
}

func (s *StringParam) GetValue() interface{} {
	return s.Value
}

type BooleanParam struct {
	PrimitiveParam
	Value bool `json:"value"`
}

func (s *BooleanParam) SetValue(v interface{}) {
	// Check value is an int
	switch c := v.(type) {
	case bool:
		s.Value = c
	case string:
		s.Value = strings.ToLower(c) == "true"
	default:
		s.Valid = false
		return
	}

	s.Valid = true
}

func (s *BooleanParam) GetValue() interface{} {
	return s.Value
}

func NewStringParam() *StringParam {
	return &StringParam{
		PrimitiveParam: PrimitiveParam{Type: StringType, Valid: true},
	}
}

func NewBooleanParam() *BooleanParam {
	return &BooleanParam{
		PrimitiveParam: PrimitiveParam{Type: BooleanType, Valid: true},
	}
}
