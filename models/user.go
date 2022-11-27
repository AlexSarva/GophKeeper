package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents user information in service
type User struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Username string    `json:"username" db:"username"`
	Email    string    `json:"email" db:"email"`
	Password string    `json:"password,omitempty" db:"passwd"`
	Token    string    `json:"token" db:"token"`
	TokenExp time.Time `json:"token_expires" db:"token_expires"`
}

// UserRegister represents information than used for register user in service
type UserRegister struct {
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"passwd"`
}

// UserLogin represents information than used for login user in service
type UserLogin struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"passwd"`
}
