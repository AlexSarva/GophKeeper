package models

import (
	"database/sql/driver"
	"time"
)

type nullTime struct {
	Time  time.Time `json:"time"`
	Valid bool      `json:"valid"` // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *nullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt *nullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}
