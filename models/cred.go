package models

import (
	"time"

	"github.com/google/uuid"
)

type Cred struct {
	ID      uuid.UUID `json:"id" db:"id"`
	Title   string    `json:"title" db:"title"`
	Login   string    `json:"login" db:"login"`
	Passwd  string    `json:"passwd" db:"passwd"`
	Notes   string    `json:"notes,omitempty" db:"notes"`
	Created time.Time `json:"created" db:"created"`
	Changed *nullTime `json:"changed,omitempty" db:"changed"`
}

type NewCred struct {
	ID     uuid.UUID
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	Title  string    `json:"title" db:"title"`
	Login  string    `json:"login" db:"login"`
	Passwd string    `json:"passwd" db:"passwd"`
	Notes  string    `json:"notes,omitempty" db:"notes"`
}
