package pqx

import (
	"bytes"
	"encoding/json"
	"github.com/lib/pq"
	"time"
)

const defaultDateLayout = "2006-01-02"

type DateOpt struct {
	Input  string
	Layout string
}

func ParseDate(opt DateOpt) Date {
	// Init Date
	dob := Date{}

	// Set date
	if opt.Layout == "" {
		dob.Layout = defaultDateLayout
	}

	// Parse date of birth
	if opt.Input != "" {
		t, err := time.Parse(dob.Layout, opt.Input)
		if err != nil {
			return dob
		}

		dob.NullTime = pq.NullTime{
			Time:  t,
			Valid: true,
		}
	}

	return dob
}

// Time extends functionality of pq.NullTime to serialize Time into epoch in JSON
type Date struct {
	pq.NullTime
	Layout string
}

/// MarshalJSON implements json Marshal interface
///
/// Time will be written in epoch value
/// If Time is not Valid, then it will be marshalled to 0
func (t Date) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return bytes.NewBufferString(Null).Bytes(), nil
	}

	return json.Marshal(t.Time.Format(t.getLayout()))
}

// UnmarshalJSON implements json Unmarshal interface
///
/// Time will be parsed from epoch value
/// If Time is null, then it will unmarshalled to Valid = false
func (t *Date) UnmarshalJSON(b []byte) error {
	// Marshal epoch
	var dateStr string
	err := json.Unmarshal(b, &dateStr)
	if err != nil {
		t.Time = time.Time{}
		t.Valid = false
		return nil
	}

	// Parse time
	dt, err := time.Parse(t.getLayout(), dateStr)
	if err != nil {
		t.Time = time.Time{}
		t.Valid = false
		return nil
	}

	// Set time
	t.Time = dt
	t.Valid = true
	return nil
}

func (t Date) getLayout() string {
	if t.Layout == "" {
		return defaultDateLayout
	}
	return t.Layout
}

func (t Date) DiffValue() interface{} {
	if !t.Valid {
		return ""
	}

	return t.Time.Format(t.getLayout())
}
