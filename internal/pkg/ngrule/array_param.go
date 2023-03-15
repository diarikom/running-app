package ngrule

import "sort"

type ArrayFactParam interface {
	GetAt(index int) interface{}
	Merge(v interface{})
	Union(v interface{})
	Push(v interface{})
	SetAt(index int, v interface{})
	RemoveAt(index int)
}

func NewStringArrayParam() *StringArrayParam {
	return &StringArrayParam{
		PrimitiveParam: PrimitiveParam{Type: StringArrayType, Valid: true},
		Value:          make([]string, 0),
	}
}

type StringArrayParam struct {
	PrimitiveParam
	Value []string `json:"value"`
}

func (s *StringArrayParam) GetAt(index int) interface{} {
	// Check for out of bounds
	if index >= len(s.Value) {
		return nil
	}

	return s.Value[index]
}

func (s *StringArrayParam) Union(v interface{}) {
	// Assert type to string array
	strArr, ok := v.([]string)
	if !ok {
		s.Valid = false
		return
	}

	var arr sort.StringSlice = s.Value

	// Iterate strArr
	for _, v := range strArr {
		// Check if string already exist
		idx := arr.Search(v)

		// If found, then continue
		if arr[idx] == v {
			continue
		}

		// Push
		arr = append(arr, v)
	}

	// Set value
	s.Value = arr
	s.Valid = true
}

func (s *StringArrayParam) Merge(v interface{}) {
	// Assert type to string array
	strArr, ok := v.([]string)
	if !ok {
		s.Valid = false
		return
	}

	// Merge
	s.Value = append(s.Value, strArr...)
}

func (s *StringArrayParam) Push(v interface{}) {
	// Assert type to string
	str, ok := v.(string)
	if !ok {
		s.Valid = false
		return
	}

	// Push
	s.Value = append(s.Value, str)
	s.Valid = true
}

func (s *StringArrayParam) SetAt(index int, v interface{}) {
	// Assert type to string
	str, ok := v.(string)
	if !ok {
		s.Valid = false
		return
	}

	// Set to index
	s.Value[index] = str
	s.Valid = true
}

func (s *StringArrayParam) RemoveAt(index int) {
	// Check index is not out of bounds
	if index >= len(s.Value) {
		s.Valid = false
		return
	}

	// Remove index
	s.Value[index] = s.Value[len(s.Value)-1]
	s.Value = s.Value[:len(s.Value)-1]
	s.Valid = true
}

func (s *StringArrayParam) SetValue(v interface{}) {
	// Check value is an string
	strArr, ok := v.([]string)
	if !ok {
		s.Valid = false
		return
	}

	s.Value = strArr
}

func (s *StringArrayParam) GetValue() interface{} {
	return s.Value
}
