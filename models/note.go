package models

import (
	"AlexSarva/GophKeeper/crypto"
	"time"

	"github.com/google/uuid"
)

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

func (nn *NewNote) Encrypt(cryptorizer *crypto.Cryptorizer) error {
	cryptNote, cryptNoteNumErr := cryptorizer.Cryptorizer.Encrypt(nn.Note)
	if cryptNoteNumErr != nil {
		return cryptNoteNumErr
	}
	nn.Note = cryptNote
	return nil
}

func (n *Note) Decrypt(cryptorizer *crypto.Cryptorizer) error {
	decryptNote, decryptNoteNumErr := cryptorizer.Cryptorizer.Decrypt(n.Note)
	if decryptNoteNumErr != nil {
		return decryptNoteNumErr
	}
	n.Note = decryptNote
	return nil
}
