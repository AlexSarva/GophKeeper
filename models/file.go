package models

import (
	"AlexSarva/GophKeeper/crypto/cryptoblock"
	"time"

	"github.com/google/uuid"
)

// File represents file information that stored in database
type File struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Title    string    `json:"title" db:"title"`
	File     []byte    `json:"file" db:"file"`
	FileName string    `json:"file_name" db:"file_name"`
	Notes    string    `json:"notes,omitempty" db:"notes"`
	Created  time.Time `json:"created" db:"created"`
	Changed  *nullTime `json:"changed,omitempty" db:"changed"`
}

// NewFile represents file information that posted by user in service
type NewFile struct {
	ID       uuid.UUID
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	Title    string    `json:"title" db:"title"`
	FileName string    `json:"file_name" db:"file_name"`
	File     []byte    `json:"file" db:"file"`
	Notes    string    `json:"notes,omitempty" db:"notes"`
}

// Encrypt cipher values (file content)
func (nf *NewFile) Encrypt(symmCrypt *cryptoblock.AEADCrypto) {
	cryptFile := symmCrypt.Encrypt(nf.File)
	nf.File = cryptFile
}

// Decrypt decipher values (file content)
func (f *File) Decrypt(symCrypt *cryptoblock.AEADCrypto) error {
	cryptFile, cryptFileErr := symCrypt.Decrypt(f.File)
	if cryptFileErr != nil {
		return cryptFileErr
	}
	f.File = cryptFile
	return nil
}
