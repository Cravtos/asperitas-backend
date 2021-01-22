package users

import (
	"time"
)

// Info represents an individual users.
type Info struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	PasswordHash []byte    `json:"-"`
	DateCreated  time.Time `json:"date_created"`
}

// NewUser contains information needed to create a new User.
type NewUser struct {
	Name     string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
