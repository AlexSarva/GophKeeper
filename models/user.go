package models

import (
	"time"

	"github.com/google/uuid"
)

type Status struct {
	Result string `json:"result"`
}

type PasswordCheck struct {
	ID       uuid.UUID `db:"id"`
	Password string    `db:"passwd"`
}

type User struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Username string    `json:"username" db:"username"`
	Email    string    `json:"email" db:"email"`
	Password string    `json:"password" db:"passwd"`
	Token    string    `json:"token" db:"token"`
	TokenExp time.Time `json:"token_expires" db:"token_expires"`
}

type UserLogin struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"passwd"`
}

type Token struct {
	ID       uuid.UUID `json:"user_id" db:"id"`
	Username string    `json:"username" db:"username"`
	Email    string    `json:"email" db:"email"`
	Token    string    `json:"token" db:"token"`
	Admin    bool      `json:"admin" db:"is_admin"`
	Created  time.Time `json:"created" db:"created"`
}

type UserInfo struct {
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
}
