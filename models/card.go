package models

import (
	"AlexSarva/GophKeeper/crypto"
	"AlexSarva/GophKeeper/utils"
	"errors"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrNotValidCardNumber = errors.New("not valid card number")
var ErrNotValidCardOwner = errors.New("wrong value for card owner field")
var ErrNotValidCardExp = errors.New("wrong value for card expiration date")

// Card represents credit card information that stored in database
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

// NewCard represents credit card information that posted by user in service
type NewCard struct {
	ID         uuid.UUID
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	Title      string    `json:"title" db:"title"`
	CardNumber string    `json:"card_number" db:"card_number"`
	CardOwner  string    `json:"card_owner" db:"card_owner"`
	CardExp    string    `json:"card_exp" db:"card_exp"`
	Notes      string    `json:"notes,omitempty" db:"notes"`
}

// CheckValid format logic check values of fields
func (nc *NewCard) CheckValid() error {
	if !utils.CheckValidCardNumber(nc.CardNumber) {
		return ErrNotValidCardNumber
	}
	reExp := regexp.MustCompile(`^(0[1-9]|1[0-2])\/?([0-9]{4}|[0-9]{2})$`)
	if !reExp.MatchString(nc.CardExp) {
		log.Fatalln(nc.CardExp)
		return ErrNotValidCardExp
	}
	reOwner := regexp.MustCompile(`(^[a-zA-z]+[\.\s]{0,2}?[a-zA-z]+[\.\s]?$)`)
	if !reOwner.MatchString(nc.CardOwner) {
		return ErrNotValidCardOwner
	}
	nc.CardOwner = strings.ToUpper(nc.CardOwner)
	return nil
}

// Encrypt cipher values (card number, card owner and expiration date)
func (nc *NewCard) Encrypt(cryptorizer *crypto.Cryptorizer) error {
	cryptCardNum, cryptCardNumErr := cryptorizer.Cryptorizer.Encrypt(nc.CardNumber)
	if cryptCardNumErr != nil {
		return cryptCardNumErr
	}
	cryptCardOwner, cryptCardOwnerErr := cryptorizer.Cryptorizer.Encrypt(nc.CardOwner)
	if cryptCardOwnerErr != nil {
		return cryptCardOwnerErr
	}
	cryptCardExp, cryptCardExpErr := cryptorizer.Cryptorizer.Encrypt(nc.CardExp)
	if cryptCardExpErr != nil {
		return cryptCardExpErr
	}
	nc.CardNumber = cryptCardNum
	nc.CardOwner = cryptCardOwner
	nc.CardExp = cryptCardExp
	return nil
}

// Decrypt decipher values (card number, card owner and expiration date)
func (c *Card) Decrypt(cryptorizer *crypto.Cryptorizer) error {
	decryptCardNum, decryptCardNumErr := cryptorizer.Cryptorizer.Decrypt(c.CardNumber)
	if decryptCardNumErr != nil {
		return decryptCardNumErr
	}
	decryptCardOwner, decryptCardOwnerErr := cryptorizer.Cryptorizer.Decrypt(c.CardOwner)
	if decryptCardOwnerErr != nil {
		return decryptCardOwnerErr
	}
	decryptCardExp, decryptCardExpErr := cryptorizer.Cryptorizer.Decrypt(c.CardExp)
	if decryptCardExpErr != nil {
		return decryptCardExpErr
	}
	c.CardNumber = decryptCardNum
	c.CardOwner = decryptCardOwner
	c.CardExp = decryptCardExp
	return nil
}
