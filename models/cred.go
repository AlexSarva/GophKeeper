package models

import (
	"AlexSarva/GophKeeper/crypto"
	"time"

	"github.com/google/uuid"
)

// Cred represents credentials (login / password) information that stored in database
type Cred struct {
	ID      uuid.UUID `json:"id" db:"id"`
	Title   string    `json:"title" db:"title"`
	Login   string    `json:"login" db:"login"`
	Passwd  string    `json:"passwd" db:"passwd"`
	Notes   string    `json:"notes,omitempty" db:"notes"`
	Created time.Time `json:"created" db:"created"`
	Changed *nullTime `json:"changed,omitempty" db:"changed"`
}

// NewCred represents credentials (login / password) that posted by user in service
type NewCred struct {
	ID     uuid.UUID
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	Title  string    `json:"title" db:"title"`
	Login  string    `json:"login" db:"login"`
	Passwd string    `json:"passwd" db:"passwd"`
	Notes  string    `json:"notes,omitempty" db:"notes"`
}

// Encrypt cipher values (login / password)
func (nc *NewCred) Encrypt(cryptorizer *crypto.Cryptorizer) error {
	cryptLogin, cryptLoginErr := cryptorizer.Cryptorizer.Encrypt(nc.Login)
	if cryptLoginErr != nil {
		return cryptLoginErr
	}
	cryptPasswd, cryptPasswdErr := cryptorizer.Cryptorizer.Encrypt(nc.Passwd)
	if cryptPasswdErr != nil {
		return cryptPasswdErr
	}
	nc.Login = cryptLogin
	nc.Passwd = cryptPasswd
	return nil
}

// Decrypt decipher values (login / password)
func (c *Cred) Decrypt(cryptorizer *crypto.Cryptorizer) error {
	decryptLogin, decryptLoginErr := cryptorizer.Cryptorizer.Decrypt(c.Login)
	if decryptLoginErr != nil {
		return decryptLoginErr
	}
	decryptPasswd, decryptPasswdErr := cryptorizer.Cryptorizer.Decrypt(c.Passwd)
	if decryptPasswdErr != nil {
		return decryptPasswdErr
	}
	c.Login = decryptLogin
	c.Passwd = decryptPasswd
	return nil
}
