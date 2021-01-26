package graph

import (
	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/jmoiron/sqlx"
	"log"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	log  *log.Logger
	db   *sqlx.DB
	Auth *auth.Auth
}
