package db

import (
	"context"
	"database/sql"
	"github.com/cravtos/asperitas-backend/foundation/database"
	"github.com/pkg/errors"
)

// GetUserByID obtains Author using ID from database
func (setup DBset) GetUserByID(ctx context.Context, ID string) (UserDB, error) {
	const qAuthor = `SELECT user_id, name FROM users WHERE user_id = $1`

	//setup.log.Printf("%s: %s", "db.getUserByID", database.Log(qAuthor, ID))

	var user UserDB
	if err := setup.db.GetContext(ctx, &user, qAuthor, ID); err != nil {
		if err == sql.ErrNoRows {
			return UserDB{}, ErrUserNotFound
		}
		return UserDB{}, err
	}
	return user, nil
}

// GetUserByName returns first found UserID with a given name
func (setup DBset) GetUserByName(ctx context.Context, name string) (UserDB, error) {
	const qAuthor = `SELECT user_id, name FROM users WHERE name = $1`

	//setup.log.Printf("%s: %s", "db.getAuthorByID", database.Log(qAuthor, name))

	var users []UserDB
	if err := setup.db.SelectContext(ctx, &users, qAuthor, name); err != nil {
		if err == sql.ErrNoRows {
			return UserDB{}, ErrUserNotFound
		}
		return UserDB{}, errors.Wrap(err, "selecting users by name")
	}
	return users[0], nil
}

// SelectUsersByName returns all Users with a given name
func (setup DBset) SelectUsersByName(ctx context.Context, name string) ([]UserDB, error) {
	const qAuthor = `SELECT user_id, name FROM users WHERE name = $1`

	//setup.log.Printf("%s: %s", "posts.helpers.getAuthorByID", database.Log(qAuthor, name))

	var users []UserDB
	if err := setup.db.SelectContext(ctx, &users, qAuthor, name); err != nil {
		return nil, errors.Wrap(err, "selecting users")
	}
	return users, nil
}

func (setup DBset) CreateUser(ctx context.Context, usr FullUserDB) error {
	const q = `
	INSERT INTO users
		(user_id, name, password_hash, date_created)
	VALUES
		($1, $2, $3, $4)`

	setup.log.Printf("%s: %s", "users.Create",
		database.Log(q, usr.ID, usr.Name, usr.PasswordHash, usr.DateCreated),
	)

	if _, err := setup.db.ExecContext(ctx, q, usr.ID, usr.Name, usr.PasswordHash, usr.DateCreated); err != nil {
		return errors.Wrap(err, "inserting users")
	}
	return nil
}

// GetFullUserByName returns first found UserID with a given name
func (setup DBset) GetFullUserByName(ctx context.Context, name string) (FullUserDB, error) {
	const qAuthor = `SELECT * FROM users WHERE name = $1`

	//setup.log.Printf("%s: %s", "db.getAuthorByID", database.Log(qAuthor, name))

	var users []FullUserDB
	if err := setup.db.SelectContext(ctx, &users, qAuthor, name); err != nil {
		if err == sql.ErrNoRows {
			return FullUserDB{}, ErrUserNotFound
		}
		return FullUserDB{}, errors.Wrap(err, "selecting users")
	}
	return users[0], nil
}
