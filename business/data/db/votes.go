package db

import (
	"context"
	"github.com/pkg/errors"
)

// SelectVotesByPostID returns slice of Vote for a single post
func (setup DBset) SelectVotesByPostID(ctx context.Context, ID string) ([]VoteDB, error) {
	const qVotes = `SELECT user_id, vote FROM votes WHERE post_id = $1`

	//setup.log.Printf("%s: %s", "db.selectVotesByPostID", database.Log(qVotes, ID))

	var votes []VoteDB
	if err := setup.db.SelectContext(ctx, &votes, qVotes, ID); err != nil {
		return nil, errors.Wrap(err, "selecting votes")
	}
	return votes, nil
}

// InsertVote adds one row to votes database
func (setup DBset) InsertVote(ctx context.Context, postID string, userID string, vote int) error {
	const qVote = `INSERT INTO votes (post_id, user_id, vote) VALUES ($1, $2, $3)`

	//setup.log.Printf("%s: %s", "db.insertVote", database.Log(qVote, postID, userID, vote))

	if _, err := setup.db.ExecContext(ctx, qVote, postID, userID, vote); err != nil {
		return errors.Wrap(err, "inserting Vote")
	}
	return nil
}

// CheckVote shows whether vote to given post by user exists in database.
// It returns an error if vote does not exist.
func (setup DBset) CheckVote(ctx context.Context, postID string, userID string) error {
	const qCheckExist = `SELECT COUNT(*) FROM votes WHERE post_id = $1 AND user_id = $2`

	//setup.log.Printf("%s: %s", "db.checkVote", database.Log(qCheckExist, postID, userID))

	var exist int
	if err := setup.db.GetContext(ctx, &exist, qCheckExist, postID, userID); err != nil {
		return errors.Wrap(err, "checking if vote exists")
	}

	if exist == 0 {
		return ErrVoteNotFound
	}
	return nil
}

// UpdateVote changes specific vote value
func (setup DBset) UpdateVote(ctx context.Context, postID string, userID string, vote int) error {
	const qUpdateVote = `UPDATE votes SET vote = $3 WHERE post_id = $1 AND user_id = $2`

	//setup.log.Printf("%s: %s", "db.updateVote", database.Log(qUpdateVote, postID, userID, vote))

	if _, err := setup.db.ExecContext(ctx, qUpdateVote, postID, userID, vote); err != nil {
		return errors.Wrap(err, "updating vote")
	}
	return nil
}

// DeleteVote deletes vote
func (setup DBset) DeleteVote(ctx context.Context, postID string, userID string) error {
	const qDeleteVote = `DELETE FROM votes WHERE post_id = $1 AND user_id = $2`

	//setup.log.Printf("%s: %s", "db.deleteVote", database.Log(qDeleteVote, postID, userID))

	if _, err := setup.db.ExecContext(ctx, qDeleteVote, postID, userID); err != nil {
		return errors.Wrapf(err, "deleting vote on %s from %s", postID, userID)
	}
	return nil
}
