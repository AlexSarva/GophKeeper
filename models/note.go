package models

import (
	"AlexSarva/GophKeeper/crypto"
	"time"

	"github.com/google/uuid"
)

// Note represents notes information that stored in database
type Note struct {
	ID      uuid.UUID `json:"id" db:"id"`
	Title   string    `json:"title" db:"title"`
	Note    string    `json:"note" db:"note"`
	Created time.Time `json:"created" db:"created"`
	Changed *nullTime `json:"changed,omitempty" db:"changed"`
}

// NewNote represents notes information that posted by user in service
type NewNote struct {
	ID     uuid.UUID
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	Title  string    `json:"title" db:"title"`
	Note   string    `json:"note" db:"note"`
}

// Encrypt cipher values (note text)
func (nn *NewNote) Encrypt(cryptorizer *crypto.Cryptorizer) error {
	cryptNote, cryptNoteNumErr := cryptorizer.Cryptorizer.Encrypt(nn.Note)
	if cryptNoteNumErr != nil {
		return cryptNoteNumErr
	}
	nn.Note = cryptNote
	return nil
}

// Decrypt decipher values (note text)
func (n *Note) Decrypt(cryptorizer *crypto.Cryptorizer) error {
	decryptNote, decryptNoteNumErr := cryptorizer.Cryptorizer.Decrypt(n.Note)
	if decryptNoteNumErr != nil {
		return decryptNoteNumErr
	}
	n.Note = decryptNote
	return nil
}
