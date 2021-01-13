// Package user contains user related CRUD functionality.
package user

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/cravtos/asperitas-backend/foundation/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrNotFound is used when a specific User is requested but does not exist.
	ErrNotFound = errors.New("not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrAuthenticationFailure occurs when a user attempts to authenticate but
	// anything goes wrong.
	ErrAuthenticationFailure = errors.New("authentication failed")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("attempted action is not allowed")
)

// User manages the set of API's for user access.
type User struct {
	log *log.Logger
	db  *sqlx.DB
}

// New constructs a User for api access.
func New(log *log.Logger, db *sqlx.DB) User {
	return User{
		log: log,
		db:  db,
	}
}

// Create inserts a new user into the database.
func (u User) Create(ctx context.Context, nu NewUser, now time.Time) (Info, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return Info{}, errors.Wrap(err, "generating password hash")
	}

	usr := Info{
		ID:           uuid.New().String(),
		Name:         nu.Name,
		PasswordHash: hash,
		DateCreated:  now,
	}

	const q = `
	INSERT INTO users
		(user_id, name, password_hash, date_created)
	VALUES
		($1, $2, $3, $4)`

	u.log.Printf("%s: %s", "user.Create",
		database.Log(q, usr.ID, usr.Name, usr.PasswordHash, usr.DateCreated),
	)

	if _, err = u.db.ExecContext(ctx, q, usr.ID, usr.Name, usr.PasswordHash, usr.DateCreated); err != nil {
		return Info{}, errors.Wrap(err, "inserting user")
	}

	return usr, nil
}

// Delete removes a user from the database.
func (u User) Delete(ctx context.Context, claims auth.Claims, userID string) error {

	if _, err := uuid.Parse(userID); err != nil {
		return ErrInvalidID
	}

	// If you are looking to delete someone other than yourself.
	if claims.User.ID != userID {
		return ErrForbidden
	}

	const q = `
	DELETE FROM
		users
	WHERE
		user_id = $1`

	u.log.Printf("%s: %s", "user.Delete",
		database.Log(q, userID),
	)

	if _, err := u.db.ExecContext(ctx, q, userID); err != nil {
		return errors.Wrapf(err, "deleting user %s", userID)
	}

	return nil
}

// Query retrieves a list of existing users from the database.
func (u User) Query(ctx context.Context, traceID string, pageNumber int, rowsPerPage int) ([]Info, error) {

	const q = `
	SELECT
		*
	FROM
		users
	ORDER BY
		user_id
	OFFSET $1 ROWS FETCH NEXT $2 ROWS ONLY`

	offset := (pageNumber - 1) * rowsPerPage

	u.log.Printf("%s: %s: %s", traceID, "user.Query",
		database.Log(q, offset, rowsPerPage),
	)

	var users []Info
	if err := u.db.SelectContext(ctx, &users, q, offset, rowsPerPage); err != nil {
		return nil, errors.Wrap(err, "selecting users")
	}

	return users, nil
}

// QueryByID gets the specified user from the database.
func (u User) QueryByID(ctx context.Context, claims auth.Claims, userID string) (Info, error) {

	if _, err := uuid.Parse(userID); err != nil {
		return Info{}, ErrInvalidID
	}

	// If you are looking to retrieve someone other than yourself.
	if claims.User.ID != userID {
		return Info{}, ErrForbidden
	}

	const q = `
	SELECT
		*
	FROM
		users
	WHERE 
		user_id = $1`

	u.log.Printf("%s: %s", "user.QueryByID",
		database.Log(q, userID),
	)

	var usr Info
	if err := u.db.GetContext(ctx, &usr, q, userID); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, ErrNotFound
		}
		return Info{}, errors.Wrapf(err, "selecting user %q", userID)
	}

	return usr, nil
}

// QueryByName gets the specified user from the database by username.
func (u User) QueryByName(ctx context.Context, claims auth.Claims, name string) (Info, error) {

	const q = `
	SELECT
		*
	FROM
		users
	WHERE
		name = $1`

	u.log.Printf("%s: %s", "user.QueryByName",
		database.Log(q, name),
	)

	var usr Info
	if err := u.db.GetContext(ctx, &usr, q, name); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, ErrNotFound
		}
		return Info{}, errors.Wrapf(err, "selecting user %q", name)
	}

	// If you are looking to retrieve someone other than yourself.
	if claims.User.ID != usr.ID {
		return Info{}, ErrForbidden
	}

	return usr, nil
}

// Authenticate finds a user by their name and verifies their password. On
// success it returns a Claims Info representing this user. The claims can be
// used to generate a token for future authentication.
func (u User) Authenticate(ctx context.Context, name, password string, now time.Time) (auth.Claims, error) {

	const q = `
	SELECT
		*
	FROM
		users
	WHERE
		name = $1`

	u.log.Printf("%s: %s", "user.Authenticate",
		database.Log(q, name),
	)

	var usr Info
	if err := u.db.GetContext(ctx, &usr, q, name); err != nil {

		// Normally we would return ErrNotFound in this scenario but we do not want
		// to leak to an unauthenticated user which emails are in the system.
		if err == sql.ErrNoRows {
			return auth.Claims{}, ErrAuthenticationFailure
		}

		return auth.Claims{}, errors.Wrap(err, "selecting single user")
	}

	// Compare the provided password with the saved hash. Use the bcrypt
	// comparison function so it is cryptographically secure.
	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(password)); err != nil {
		return auth.Claims{}, ErrAuthenticationFailure
	}

	// If we are this far the request is valid. Create some claims for the user
	// and generate their token.
	claims := auth.Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(time.Hour).Unix(),
			IssuedAt:  now.Unix(),
		},
		User: auth.User{
			Username: usr.Name,
			ID:       usr.ID,
		},
	}

	return claims, nil
}
