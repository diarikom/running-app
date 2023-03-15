package nstr

import (
	"strconv"
)

func ParseFloat32(input string, defaultValue float32) float32 {
	o, err := strconv.ParseFloat(input, 32)
	if err != nil {
		return defaultValue
	}
	return float32(o)
}

func ParseFloat64(input string, defaultValue float64) float64 {
	o, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return defaultValue
	}
	return o
}

func ParseInt(input string, defaultValue int) int {
	o, err := strconv.Atoi(input)
	if err != nil {
		return defaultValue
	}
	return o
}

func ParseInt8(input string, defaultValue int8) int8 {
	o, err := strconv.ParseInt(input, 10, 8)
	if err != nil {
		return defaultValue
	}
	return int8(o)
}

func ParseInt16(input string, defaultValue int16) int16 {
	o, err := strconv.ParseInt(input, 10, 16)
	if err != nil {
		return defaultValue
	}
	return int16(o)
}

func ParseInt64(input string, defaultValue int64) int64 {
	o, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return defaultValue
	}
	return o
}
