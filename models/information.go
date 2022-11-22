package models

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
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

type Note struct {
	ID      uuid.UUID `json:"id" db:"id"`
	Title   string    `json:"title" db:"title"`
	Note    string    `json:"note" db:"note"`
	Created time.Time `json:"created" db:"created"`
	Changed *nullTime `json:"changed,omitempty" db:"changed"`
}

type NewNote struct {
	ID     uuid.UUID
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	Title  string    `json:"title" db:"title"`
	Note   string    `json:"note" db:"note"`
}
