package post

import "time"


// Author represents info about author
type Author struct {
	Username string `json:"username"`
	ID       string `json:"id"`
}

// Vote represents info about user vote.
type Vote struct {
	User string `json:"user"`
	Vote int    `json:"vote"`
}

// Comment represents info about comments for the post.
type Comment struct {
	DateCreated time.Time `json:"created"`
	Author      Author    `json:"author"`
	Body        string    `json:"body"`
	ID          string    `json:"id"`
}

// todo: find more appropriate name
// Comment represents info about comments in database.
type CommentDB struct {
	DateCreated time.Time `json:"created"`
	UserID      string    `json:"user_id"`
	Body        string    `json:"body"`
	ID          string    `json:"id"`
}

// todo: find more appropriate name not starting with "Post"
// PostDB represents an individual post in database.
type PostDB struct {
	ID               string    `db:"post_id" json:"id"`
	Score            int       `db:"score" json:"score"`
	Views            int       `db:"views" json:"views"`
	Type             string    `db:"type" json:"type" default:"text"`
	Title            string    `db:"title" json:"title"`
	URL              string    `db:"url" json:"url"`
	Category         string    `db:"category" json:"category"`
	Text             string    `db:"text" json:"text"`
	DateCreated      time.Time `db:"date_created" json:"created"`
	UserID			 string	   `db:"user_id" json:"user_id"`
}

// todo: split into two (InfoText, InfoLink or smth like that)
// Info represents an individual post which is sent to user.
type Info struct {
	ID               string    `json:"id"`
	Score            int       `json:"score"`
	Views            int       `json:"views"`
	Type             string    `son:"type" default:"text"`
	Title            string    `json:"title"`
	URL              string    `json:"url"`
	Category         string    `json:"category"`
	Text             string    `json:"text"`
	DateCreated      time.Time `json:"created"`
	Author           Author    `json:"author"`
	Votes            []Vote    `json:"votes"`
	Comments         []Comment `json:"comments"`
	UpvotePercentage int       `json:"upvotePercentage"`
}

// todo: validation on text or url
// NewPost is what we require from users when adding a Post.
type NewPost struct {
	Type             string    `db:"type" json:"type" default:"link"`
	Title            string    `db:"title" json:"title" validate:"required"`
	Category         string    `db:"category" json:"category" validate:"required"`
	Text             string    `db:"text" json:"text"`
	URL              string    `db:"url" json:"url"`
}

type NewComment struct {
	Text string `json:"comment" validate:"required"`
}