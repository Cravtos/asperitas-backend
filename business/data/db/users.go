package db

import (
	"context"
	"database/sql"
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

// GetUserByName returns first found User with a given name
func (setup DBset) GetUserByName(ctx context.Context, name string) (UserDB, error) {
	const qAuthor = `SELECT user_id, name FROM users WHERE name = $1`

	//setup.log.Printf("%s: %s", "db.getAuthorByID", database.Log(qAuthor, name))

	var users []UserDB
	if err := setup.db.SelectContext(ctx, &users, qAuthor, name); err != nil {
		return UserDB{}, errors.Wrap(err, "selecting users")
	}
	return users[0], nil
}

// SelectUsersByName returns all Users with a given name
func (setup DBset) SelectUsersByName(ctx context.Context, name string) ([]UserDB, error) {
	const qAuthor = `SELECT user_id, name FROM users WHERE name = $1`

	//setup.log.Printf("%s: %s", "post.helpers.getAuthorByID", database.Log(qAuthor, name))

	var users []UserDB
	if err := setup.db.SelectContext(ctx, &users, qAuthor, name); err != nil {
		return nil, errors.Wrap(err, "selecting users")
	}
	return users, nil
}
