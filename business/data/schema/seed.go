package schema

import (
	"github.com/jmoiron/sqlx"
)

// Seed runs the set of seed-data queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func Seed(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(seeds); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

// seeds is a string constant containing all of the queries needed to get the
// db seeded to a useful state for development.
//
// Note that database servers besides PostgreSQL may not support running
// multiple queries as part of the same execution so this single large constant
// may need to be broken up.
const seeds = `
-- Create admin and regular User with password "gophers"
INSERT INTO users (user_id, name, password_hash, date_created) VALUES
	('5cf37266-3473-4006-984f-9325122678b7', 'Admin Gopher', '$2a$10$1ggfMVZV6Js0ybvJufLRUOWHS5f6KneuP0XwwHpJ8L8ipdry9f2/a', '2019-03-24 00:00:00'),
	('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'User Gopher', '$2a$10$9/XASPKBbJKVfCAZKDH.UuhsuALDr5vVm6VrYA9VFR8rccK86C1hW', '2019-03-24 00:00:00')
	ON CONFLICT DO NOTHING;

INSERT INTO posts (post_id, score, views, type, title, category, payload, date_created, user_id) VALUES
	('a2b0639f-2cc6-44b8-b97b-15d69dbb511e', 1, 50, 'url', 'testpost',  'music', 'https://exmaple.com/', '2019-01-01 00:00:01.000001+00', '5cf37266-3473-4006-984f-9325122678b7'),
	('72f8b983-3eb4-48db-9ed0-e45cc6bd716b', 1, 75, 'text', 'secondpost', 'funny', 'hahatext', '2019-01-01 00:00:02.000001+00', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f')
	ON CONFLICT DO NOTHING;

INSERT INTO votes (post_id, user_id, vote) VALUES
	('a2b0639f-2cc6-44b8-b97b-15d69dbb511e', '5cf37266-3473-4006-984f-9325122678b7', 1),
	('72f8b983-3eb4-48db-9ed0-e45cc6bd716b', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 1)
	ON CONFLICT DO NOTHING;

INSERT INTO comments (comment_id, user_id, post_id, body, date_created) VALUES
	('a2b0639f-2cc6-44b8-b97b-15d69dbb5123', '5cf37266-3473-4006-984f-9325122678b7', '72f8b983-3eb4-48db-9ed0-e45cc6bd716b', 'That is great!', '2019-01-01 00:00:05.000001+00'),
	('72f8b983-3eb4-48db-9ed0-e45cc6bd7321', '45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'a2b0639f-2cc6-44b8-b97b-15d69dbb511e', 'awful', '2019-01-01 00:00:06.000001+00')
	ON CONFLICT DO NOTHING;
`

// DeleteAll runs the set of Drop-table queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func DeleteAll(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(deleteAll); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

// deleteAll is used to clean the database between tests.
const deleteAll = `
DELETE FROM posts;
DELETE FROM users;`
