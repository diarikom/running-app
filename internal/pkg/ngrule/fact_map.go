package ngrule

import "fmt"

var paramTC = map[string]ParamConstructor{
	IntType: func() FactParam {
		return NewIntParam()
	},
	FloatType: func() FactParam {
		return NewFloatParam()
	},
	StringType: func() FactParam {
		return NewStringParam()
	},
	BooleanType: func() FactParam {
		return NewBooleanParam()
	},
	StringArrayType: func() FactParam {
		return NewStringArrayParam()
	},
}

type FactMap struct {
	Params map[string]FactParam
}

func NewFactMap(arr []FactParam) FactMap {
	pm := make(map[string]FactParam)
	for _, v := range arr {
		pm[v.GetName()] = v
	}

	return FactMap{Params: pm}
}

func (m *FactMap) AddType(name string, constructorFn ParamConstructor) {
	paramTC[name] = constructorFn
}

func (m *FactMap) GetInt(key string) int64 {
	p := m.Params[key]
	tmp := p.GetValue()
	v, ok := tmp.(int64)
	if ok {
		return v
	}
	panic(NewError(ErrInvalidIntValue, fmt.Sprintf("cannot %s from FactParam (Name=%s, Type=%s)", "GetInt", p.GetType(), p.GetName())))
}

func (m *FactMap) GetFloat(key string) float64 {
	p := m.Params[key]
	tmp := p.GetValue()
	v, ok := tmp.(float64)
	if ok {
		return v
	}
	panic(NewError(ErrInvalidFloatValue, fmt.Sprintf("cannot %s from FactParam (Name=%s, Type=%s)", "GetFloat", p.GetType(), p.GetName())))
}

func (m *FactMap) GetString(key string) string {
	return fmt.Sprintf("%v", m.Params[key].GetValue())
}

func (m *FactMap) GetBool(key string) bool {
	p := m.Params[key]
	tmp := p.GetValue()
	v, ok := tmp.(bool)
	if ok {
		return v
	}
	panic(NewError(ErrInvalidBooleanValue, fmt.Sprintf("cannot %s from FactParam (Name=%s, Type=%s)", "GetBool", p.GetType(), p.GetName())))
}

func (m *FactMap) Set(key string, val interface{}) {
	m.Params[key].SetValue(val)
}

func (m *FactMap) Assign(key string, val FactParam) {
	m.Params[key] = val
}

func (m *FactMap) Add(key string, value interface{}) {
	m.Math(key, AddOp, value)
}

func (m *FactMap) Subtract(key string, value interface{}) {
	m.Math(key, SubtractOp, value)
}

func (m *FactMap) Multiply(key string, value interface{}) {
	m.Math(key, MultiplyOp, value)
}

func (m *FactMap) Divide(key string, value interface{}) {
	m.Math(key, DivideOp, value)
}

func (m *FactMap) Math(key string, op int, value interface{}) {
	switch v := value.(type) {
	case int, int8, int16, int32, int64, float32, float64:
		p, ok := m.Params[key].(MathFactParam)
		if !ok {
			return
		}

		switch op {
		case AddOp:
			p.Add(v)
		case SubtractOp:
			p.Subtract(v)
		case MultiplyOp:
			p.Multiply(v)
		case DivideOp:
			p.Divide(v)
		}
	}
}

func (m *FactMap) ArrayPush(key string, value interface{}) {
	// Get Parameters that implements ArrayFactParam
	p, ok := m.Params[key].(ArrayFactParam)
	if !ok {
		return
	}

	p.Push(value)
}

func RenderMathFunc(op string) string {
	switch op {
	case "assign", "=":
		return "Set"
	case "add", "+":
		return "Add"
	case "subtract", "-":
		return "Subtract"
	case "multiply", "*":
		return "Multiply"
	case "divide", "/":
		return "Divide"
	default:
		return ""
	}
}
