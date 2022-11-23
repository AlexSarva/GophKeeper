package models

import (
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Title      string    `json:"title" db:"title"`
	CardNumber string    `json:"card_number" db:"card_number"`
	CardOwner  string    `json:"card_owner" db:"card_owner"`
	CardExp    string    `json:"card_exp" db:"card_exp"`
	Notes      string    `json:"notes,omitempty" db:"notes"`
	Created    time.Time `json:"created" db:"created"`
	Changed    *nullTime `json:"changed,omitempty" db:"changed"`
}

type NewCard struct {
	ID         uuid.UUID
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	Title      string    `json:"title" db:"title"`
	CardNumber string    `json:"card_number" db:"card_number"`
	CardOwner  string    `json:"card_owner" db:"card_owner"`
	CardExp    string    `json:"card_exp" db:"card_exp"`
	Notes      string    `json:"notes,omitempty" db:"notes"`
}
