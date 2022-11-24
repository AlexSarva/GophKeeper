package models

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID      uuid.UUID `json:"id" db:"id"`
	Title   string    `json:"title" db:"title"`
	File    []byte    `json:"file" db:"file"`
	Notes   string    `json:"notes,omitempty" db:"notes"`
	Created time.Time `json:"created" db:"created"`
	Changed *nullTime `json:"changed,omitempty" db:"changed"`
}

type NewFile struct {
	ID     uuid.UUID
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	Title  string    `json:"title" db:"title"`
	File   []byte    `json:"file" db:"file"`
	Notes  string    `json:"notes,omitempty" db:"notes"`
}

type NewClientFile struct {
	Title string `json:"title" db:"title"`
	File  string `json:"file" db:"file"`
	Notes string `json:"notes,omitempty" db:"notes"`
}
