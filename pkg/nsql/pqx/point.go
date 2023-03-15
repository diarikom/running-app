package pqx

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"github.com/jackc/pgx/pgtype"
)

// Point add functionality to handle null types of JSON strings
type Point struct {
	Lat   float64
	Lng   float64
	Valid bool
}

type pointJSON struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Scan implements the database/sql Scanner interface.
func (p *Point) Scan(src interface{}) error {
	// Get address of point
	var point pgtype.Point

	// Scan
	err := (&point).Scan(src)
	if err != nil {
		return err
	}

	if point.Status == pgtype.Present {
		p.Valid = true
		p.Lat = point.P.Y
		p.Lng = point.P.X
	} else {
		p.Valid = false
		p.Lat = 0
		p.Lng = 0
	}

	return nil
}

// Value implements the database/sql/driver Valuer interface.
func (p Point) Value() (driver.Value, error) {
	// Set status
	var status pgtype.Status
	if p.Valid {
		status = pgtype.Present
	} else {
		status = pgtype.Null
	}

	// Convert back to pgtype.Point
	src := pgtype.Point{
		P: pgtype.Vec2{
			X: p.Lng,
			Y: p.Lat,
		},
		Status: status,
	}

	return (&src).Value()
}

func (p Point) MarshalJSON() (raw []byte, err error) {
	if !p.Valid {
		return bytes.NewBufferString(Null).Bytes(), nil
	}

	return json.Marshal(pointJSON{
		Lat: p.Lat,
		Lng: p.Lng,
	})
}

func (p *Point) UnmarshalJSON(b []byte) error {
	// If marshalled value is a null
	if Null == string(b) {
		p.Lat = 0
		p.Lng = 0
		p.Valid = false
		return nil
	}

	// Parse value
	var tmp pointJSON
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	// Set values
	p.Lat = tmp.Lat
	p.Lng = tmp.Lng
	p.Valid = true

	return nil
}
