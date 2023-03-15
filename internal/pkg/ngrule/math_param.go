package ngrule

const (
	AddOp = iota
	SubtractOp
	MultiplyOp
	DivideOp
)

type MathFactParam interface {
	Add(v interface{})
	Subtract(v interface{})
	Multiply(v interface{})
	Divide(v interface{})
}

type IntParam struct {
	PrimitiveParam
	Value int64 `json:"value"`
}

func (i *IntParam) SetValue(v interface{}) {
	result, ok := ConvertInt(v)
	i.Value = result
	i.Valid = ok
}

func (i *IntParam) GetValue() interface{} {
	return i.Value
}

func (i *IntParam) Add(v interface{}) {
	addedVal, ok := ConvertInt(v)
	if !ok {
		i.Valid = false
		return
	}

	i.Value += addedVal
	i.Valid = true
}

func (i *IntParam) Subtract(v interface{}) {
	addedVal, ok := ConvertInt(v)
	if !ok {
		i.Valid = false
		return
	}

	i.Value -= addedVal
	i.Valid = true
}

func (i *IntParam) Multiply(v interface{}) {
	addedVal, ok := ConvertInt(v)
	if !ok {
		i.Valid = false
		return
	}

	i.Value *= addedVal
	i.Valid = true
}

func (i *IntParam) Divide(v interface{}) {
	addedVal, ok := ConvertInt(v)
	if !ok {
		i.Valid = false
		return
	}

	i.Value /= addedVal
	i.Valid = true
}

type FloatParam struct {
	PrimitiveParam
	Value float64 `json:"value"`
}

func (f *FloatParam) SetValue(v interface{}) {
	result, ok := ConvertFloat(v)
	f.Value = result
	f.Valid = ok
}

func (f *FloatParam) GetValue() interface{} {
	return f.Value
}

func (f *FloatParam) Add(v interface{}) {
	addedVal, ok := ConvertFloat(v)
	if !ok {
		f.Valid = false
		return
	}

	f.Value += addedVal
	f.Valid = true
}

func (f *FloatParam) Subtract(v interface{}) {
	addedVal, ok := ConvertFloat(v)
	if !ok {
		f.Valid = false
		return
	}

	f.Value -= addedVal
	f.Valid = true
}

func (f *FloatParam) Multiply(v interface{}) {
	addedVal, ok := ConvertFloat(v)
	if !ok {
		f.Valid = false
		return
	}

	f.Value *= addedVal
	f.Valid = true
}

func (f *FloatParam) Divide(v interface{}) {
	addedVal, ok := ConvertFloat(v)
	if !ok {
		f.Valid = false
		return
	}

	f.Value /= addedVal
	f.Valid = true
}

func NewIntParam() *IntParam {
	return &IntParam{
		PrimitiveParam: PrimitiveParam{Type: IntType, Valid: true},
	}
}

func NewFloatParam() *FloatParam {
	return &FloatParam{
		PrimitiveParam: PrimitiveParam{Type: FloatType, Valid: true},
	}
}
