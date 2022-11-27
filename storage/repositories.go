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

	NewCard(card *models.NewCard) (models.Card, error)
	AllCards(userID uuid.UUID) ([]models.Card, error)
	GetCard(cardID uuid.UUID, userID uuid.UUID) (models.Card, error)
	EditCard(card models.NewCard) (models.Card, error)
	DeleteCard(cardID uuid.UUID, userID uuid.UUID) error

	NewCred(cred *models.NewCred) (models.Cred, error)
	AllCreds(userID uuid.UUID) ([]models.Cred, error)
	GetCred(credID uuid.UUID, userID uuid.UUID) (models.Cred, error)
	EditCred(cred models.NewCred) (models.Cred, error)
	DeleteCred(credID uuid.UUID, userID uuid.UUID) error

	NewFile(file *models.NewFile) (models.File, error)
	AllFiles(userID uuid.UUID) ([]models.File, error)
	GetFile(cardID uuid.UUID, userID uuid.UUID) (models.File, error)
	EditFile(file *models.NewFile) (models.File, error)
	DeleteFile(fileID uuid.UUID, userID uuid.UUID) error
}
