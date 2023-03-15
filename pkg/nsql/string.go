package nsql

import (
	"database/sql"
	"encoding/json"
)

const Null = `null`

// String add functionality to handle null types of JSON strings
type String struct {
	sql.NullString
}

func NewString(str string) String {
	return String{
		NullString: sql.NullString{
			String: str,
			Valid:  str != "",
		},
	}
}

func (s String) MarshalJSON() (raw []byte, err error) {
	var str string
	if s.Valid {
		str = s.String
	}
	return json.Marshal(str)
}

func (s *String) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	// If marshalled value is empty string, but equals to null set to null
	if str == "" && Null == string(b) {
		s.Valid = false
		return nil
	}
	// Else, set receiver values
	s.Valid = true
	s.String = str
	return nil
}
