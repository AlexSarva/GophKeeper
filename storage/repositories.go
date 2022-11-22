package storage

import (
	"AlexSarva/GophKeeper/models"
	"errors"

	"github.com/google/uuid"
)

// ErrDuplicatePK error that occurs when adding exists user or order number
var ErrDuplicatePK = errors.New("duplicate PK")

// ErrNoValues error that occurs when no values selected from database
var ErrNoValues = errors.New("no values from select")

// Database primary interface for all types of databases
type Database interface {
	Ping() bool
	NewNote(note *models.NewNote) (models.Note, error)
	AllNotes(userID uuid.UUID) ([]models.Note, error)
	GetNote(noteID uuid.UUID, userID uuid.UUID) (models.Note, error)
	EditNote(note models.NewNote) (models.Note, error)
	DeleteNote(noteID uuid.UUID, userID uuid.UUID) error
}
