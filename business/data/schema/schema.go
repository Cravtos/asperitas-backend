// Package schema contains the database schema, migrations and seeding data.
package schema

import (
	"github.com/dimiro1/darwin"
	"github.com/jmoiron/sqlx"
)

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sqlx.DB) error {
	driver := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	d := darwin.New(driver, migrations, nil)
	return d.Migrate()
}

// migrations contains the queries needed to construct the database schema.
// Entries should never be removed once they have been run in production.
//
// Using constants in a .go file is an easy way to ensure the schema is part
// of the compiled executable and avoids pathing issues with the working
// directory. It has the downside that it lacks syntax highlighting and may be
// harder to read for some cases compared to using .sql files. You may also
// consider a combined approach using a tool like packr or go-bindata.
var migrations = []darwin.Migration{
	{
		Version:     1.0,
		Description: "Create table users",
		Script: `
CREATE TABLE users (
	user_id       UUID,
	name          TEXT,
	password_hash TEXT,
	date_created  TIMESTAMP,

	PRIMARY KEY (user_id)
);`,
	},
	{
		Version:     1.1,
		Description: "Create table posts",
		Script: `
CREATE TABLE posts (
	post_id          UUID,
	score            INT,
	views            INT,
	type             TEXT,
	title            TEXT,
	payload          TEXT,
	category         TEXT,
	date_created     TIMESTAMP,
	
	user_id			 UUID references users(user_id),

	PRIMARY KEY (post_id)
);`,
	},
	{
		Version:     1.2,
		Description: "Create table votes",
		Script: `
CREATE TABLE votes (
	post_id          UUID references posts(post_id),
	user_id          UUID references users(user_id),
	vote             INT
);`,
	},
	{
		Version: 1.3,
		Script: `
CREATE TABLE comments (
	comment_id       UUID,
	post_id          UUID references posts(post_id),
	user_id          UUID references users(user_id),
	body             TEXT,
	date_created     TIMESTAMP,

	PRIMARY KEY (comment_id)
);`,
		Description: "Create table comments",
	},
}
