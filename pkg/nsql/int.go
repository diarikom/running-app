package nsql

import (
	"database/sql"
	"encoding/json"
)

// String add functionality to handle null types of JSON strings
type Int64 struct {
	sql.NullInt64
}

func NewInt64(input int64) Int64 {
	return Int64{
		NullInt64: sql.NullInt64{
			Int64: input,
			Valid: true,
		},
	}
}

func (s Int64) MarshalJSON() (raw []byte, err error) {
	if s.Valid {
		return json.Marshal(s.Int64)
	}
	return json.Marshal(Null)
}

func (s *Int64) UnmarshalJSON(b []byte) error {
	// If value is a null, then set to invalid
	if Null == string(b) {
		s.Valid = false
		s.Int64 = 0
		return nil
	}

	var tmp int64
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}
	// Else, set receiver values
	s.Valid = true
	s.Int64 = tmp
	return nil
}

// newInt returns pointer of a variable that contains integer
func newInt(i int) *int {
	return &i
}
