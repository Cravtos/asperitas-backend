package post

import "time"


// Author represents info about author
type Author struct {
	Username string `json:"username"`
	ID       string `json:"id"`
}

// Votes represents info about user votes.
type Votes struct {
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

// Info represents an individual post.
type Info struct {
	Score            int       `db:"score" json:"score"`
	Views            int       `db:"views" json:"views"`
	Type             string    `db:"type" json:"type" default:"link"`
	Title            string    `db:"title" json:"title"`
	URL              string    `db:"url" json:"url"`
	//Author           Author    `db:"author" json:"author"`
	Category         string    `db:"category" json:"category"`
	Votes            []Votes   `db:"votes" json:"votes"`
	Comments         []Comment `db:"comments" json:"comments"`
	DateCreated      time.Time `db:"date_created" json:"created"`
	UpvotePercentage int       `json:"upvotePercentage"`
	Text             string    `db:"text" json:"text"`
	ID               string    `db:"post_id" json:"id"`
}

// NewPost is what we require from users when adding a Post.
type NewPost struct {
	Type             string    `db:"type" json:"type" default:"link"`
	Title            string    `db:"title" json:"title" validate:"required"`
	URL              string    `db:"url" json:"url" validate:"required"`
	Category         string    `db:"category" json:"category" validate:"required"`
	Text             string    `db:"text" json:"text"`
}