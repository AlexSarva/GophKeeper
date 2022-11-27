package models

import (
	"AlexSarva/GophKeeper/crypto"
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
