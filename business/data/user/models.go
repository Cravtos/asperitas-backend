package user

import (
	"time"

	"github.com/lib/pq"
)

// Info represents an individual user.
type Info struct {
	ID           string         `db:"user_id" json:"id"`
	Name         string         `db:"name" json:"name"`
	Roles        pq.StringArray `db:"roles" json:"roles"`
	PasswordHash []byte         `db:"password_hash" json:"-"`
	DateCreated  time.Time      `db:"date_created" json:"date_created"`
}

// NewUser contains information needed to create a new User.
type NewUser struct {
	Name            string   `json:"name" validate:"required"`
	Roles           []string `json:"roles" validate:"required"`
	Password        string   `json:"password" validate:"required"`
}