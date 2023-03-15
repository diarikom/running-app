package pqx

import (
	"bytes"
	"encoding/json"
	"github.com/lib/pq"
	"time"
)

const Null = `null`

func NewTime(t time.Time) Time {
	return Time{
		NullTime: pq.NullTime{
			Time:  t,
			Valid: isTimeValid(t),
		},
	}
}

// Time extends functionality of pq.NullTime to serialize Time into epoch in JSON
type Time struct {
	pq.NullTime
}

/// MarshalJSON implements json Marshal interface
///
/// Time will be written in epoch value
/// If Time is not Valid, then it will be marshalled to 0
func (t Time) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return bytes.NewBufferString(Null).Bytes(), nil
	}
	return json.Marshal(t.Time.Unix())
}

// UnmarshalJSON implements json Unmarshal interface
///
/// Time will be parsed from epoch value
/// If Time is null, then it will unmarshalled to Valid = false
func (t *Time) UnmarshalJSON(b []byte) error {
	// Marshal epoch
	var ep int64
	err := json.Unmarshal(b, &ep)
	if err != nil {
		return err
	}
	// If result epoch is 0 and bytes is equal null, then set Valid to false
	if ep == 0 && Null == string(b) {
		t.Valid = false
		return nil
	}
	// Set time
	t.Time = time.Unix(ep, 0)
	t.Valid = true
	return nil
}

func isTimeValid(t time.Time) bool {
	if t.Unix() > 0 {
		return true
	}
	return false
}

func NullTimeUnix(t pq.NullTime) (epoch int64) {
	if t.Valid {
		epoch = t.Time.Unix()
	}
	return
}
