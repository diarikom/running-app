package nhttp

import (
	"github.com/diarikom/running-app/running-app-api/pkg/nstr"
	"net/http"
	"net/url"
	"time"
)

type Values struct {
	url.Values
}

func (v *Values) Bool(key string) bool {
	input := v.Get(key)
	if input == "true" || input == "1" {
		return true
	}
	return false
}

func (v *Values) Int(key string, defaultValue int) int {
	input := v.Get(key)
	return nstr.ParseInt(input, defaultValue)
}

func (v *Values) Int8(key string, defaultValue int8) int8 {
	input := v.Get(key)
	return nstr.ParseInt8(input, defaultValue)
}

func (v *Values) Int16(key string, defaultValue int16) int16 {
	input := v.Get(key)
	return nstr.ParseInt16(input, defaultValue)
}

func (v *Values) Float32(key string, defaultValue float32) float32 {
	input := v.Get(key)
	return nstr.ParseFloat32(input, defaultValue)
}

func (v *Values) Float64(key string, defaultValue float64) float64 {
	input := v.Get(key)
	return nstr.ParseFloat64(input, defaultValue)
}

func (v *Values) Time(key string, layout string) time.Time {
	input := v.Get(key)
	t, _ := time.Parse(layout, input)
	return t
}

func ParseForm(r *http.Request) (*Values, error) {
	if err := r.ParseForm(); err != nil {
		return nil, ErrBadRequest
	}
	// Return values
	return &Values{Values: r.Form}, nil
}
