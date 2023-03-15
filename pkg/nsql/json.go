package nsql

import (
	"encoding/json"
	"errors"
)

// ScanJSON is a generic scanner function to parse json from row data
func ScanJSON(src interface{}, target interface{}) error {
	// If source is nil, set target to nil
	if src == nil {
		target = nil
		return nil
	}
	// Assert source to byte
	source, ok := src.([]byte)
	if !ok {
		return errors.New("nsql: type assertion to byte failed")
	}
	// Unmarshal to target
	err := json.Unmarshal(source, target)
	if err != nil {
		return err
	}
	return nil
}
