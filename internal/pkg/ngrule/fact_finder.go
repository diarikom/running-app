package ngrule

import (
	"github.com/stoewer/go-strcase"
	"hash/crc32"
	"reflect"
	"runtime"
	"strconv"
)

type FactFinderFn func(facts *FactMap, userId string) error

func NewFactFinderMap() *FactFinderMap {
	return &FactFinderMap{ParamMapFn: make(map[string]FactFinderFn)}
}

type FactFinderMap struct {
	ParamMapFn map[string]FactFinderFn
}

// RegisterParamFn add Fact Finder function for parameter
func (f *FactFinderMap) RegisterParamFn(paramName string, fn FactFinderFn) {
	// Get key name
	paramKey := GetParamFnKey(paramName)

	// Set function
	f.ParamMapFn[paramKey] = fn
}

func (f *FactFinderMap) FindUniqueFunctions(params []FactParam) ([]FactFinderFn, error) {
	// Init function map
	fnMap := make(map[string]FactFinderFn)

	for _, v := range params {
		// Get key name
		paramKey := GetParamFnKey(v.GetName())

		// Get fact finder function for param
		fn, ok := f.ParamMapFn[paramKey]
		if !ok {
			return nil, NewError(ErrParamFactFinderNotRegistered, "fact finder func for param is not registered")
		}

		// Get function name
		fnName := f.GetFnName(fn)

		// If function is not yet added, add function to map
		_, ok = fnMap[fnName]
		if !ok {
			fnMap[fnName] = fn
		}
	}

	// Convert to array
	fnArr := make([]FactFinderFn, len(fnMap))
	idx := 0

	for _, v := range fnMap {
		fnArr[idx] = v
		idx += 1
	}

	return fnArr, nil
}

func (f *FactFinderMap) FindFacts(m *FactMap, params []FactParam, userId string) error {
	// Get unique function
	fns, err := f.FindUniqueFunctions(params)
	if err != nil {
		return err
	}

	// Execute functions
	for _, fn := range fns {
		// Execute fact finder function
		err := fn(m, userId)
		if err != nil {
			return NewError(ErrParamFactFinderExecFail, "error occurred while calling fact finder function "+f.GetFnName(fn))
		}
	}
	return nil
}

func (f *FactFinderMap) GetFnName(fn interface{}) string {
	// Get name
	fnName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()

	// Hash function name
	checksum := crc32.ChecksumIEEE([]byte(fnName))

	// Return crc32 string
	return strconv.FormatUint(uint64(checksum), 16)
}

func GetParamFnKey(paramName string) string {
	return strcase.UpperCamelCase(paramName)
}
