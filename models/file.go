package models

import (
	"AlexSarva/GophKeeper/crypto/symmetric"
	"time"

	"github.com/google/uuid"
)

type File struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Title    string    `json:"title" db:"title"`
	File     []byte    `json:"file" db:"file"`
	FileName string    `json:"file_name" db:"file_name"`
	Notes    string    `json:"notes,omitempty" db:"notes"`
	Created  time.Time `json:"created" db:"created"`
	Changed  *nullTime `json:"changed,omitempty" db:"changed"`
}

type NewFile struct {
	ID       uuid.UUID
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	Title    string    `json:"title" db:"title"`
	FileName string    `json:"file_name" db:"file_name"`
	File     []byte    `json:"file" db:"file"`
	Notes    string    `json:"notes,omitempty" db:"notes"`
}

func (nf *NewFile) SymEncrypt(symmCrypt *symmetric.SymmetricCrypto) {
	cryptFile := symmCrypt.Encrypt(nf.File)
	nf.File = cryptFile
}

func (f *File) SymDecrypt(symCrypt *symmetric.SymmetricCrypto) error {
	cryptFile, cryptFileErr := symCrypt.Decrypt(f.File)
	if cryptFileErr != nil {
		return cryptFileErr
	}
	f.File = cryptFile
	return nil
}
