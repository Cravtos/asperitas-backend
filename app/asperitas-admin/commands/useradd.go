package commands

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cravtos/asperitas-backend/business/data/user"
	"github.com/cravtos/asperitas-backend/foundation/database"
	"github.com/pkg/errors"
)

// UserAdd adds new users into the database.
func UserAdd(traceID string, log *log.Logger, cfg database.Config, name, password string) error {
	if name == "" || password == "" {
		fmt.Println("help: useradd <name> <password>")
		return ErrHelp
	}

	db, err := database.Open(cfg)
	if err != nil {
		return errors.Wrap(err, "connect database")
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u := user.New(log, db)

	nu := user.NewUser{
		Name:            name,
		Password:        password,
	}

	usr, err := u.Create(ctx, traceID, nu, time.Now())
	if err != nil {
		return errors.Wrap(err, "create user")
	}

	fmt.Println("user id:", usr.ID)
	return nil
}
