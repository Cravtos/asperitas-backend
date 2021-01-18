package post

import (
	"time"
)

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

// Info generalizes text and link posts
type Info interface {
	Info()
}

// InfoText represents an individual text post which is sent to user.
type InfoText struct {
	Type             string    `json:"type"`
	ID               string    `json:"id"`
	Score            int       `json:"score"`
	Views            int       `json:"views"`
	Title            string    `json:"title"`
	Category         string    `json:"category"`
	Payload          string    `json:"text"`
	DateCreated      time.Time `json:"created"`
	Author           Author    `json:"author"`
	Votes            []Vote    `json:"votes"`
	Comments         []Comment `json:"comments"`
	UpvotePercentage int       `json:"upvotePercentage"`
}

// InfoLink represents an individual link post which is sent to user.
type InfoLink struct {
	Type             string    `json:"type"`
	ID               string    `json:"id"`
	Score            int       `json:"score"`
	Views            int       `json:"views"`
	Title            string    `json:"title"`
	Payload          string    `json:"url"`
	Category         string    `json:"category"`
	DateCreated      time.Time `json:"created"`
	Author           Author    `json:"author"`
	Votes            []Vote    `json:"votes"`
	Comments         []Comment `json:"comments"`
	UpvotePercentage int       `json:"upvotePercentage"`
}

func (it InfoText) Info() {}
func (il InfoLink) Info() {}

// Comment represents info about comments for the post prepared to be sent to user.
type Comment struct {
	DateCreated time.Time `json:"created"`
	Author      Author    `json:"author"`
	Body        string    `json:"body"`
	ID          string    `json:"id"`
}

// NewPost is what we require from users when adding a PostSet.
type NewPost struct {
	Type     string `json:"type" default:"link"`
	Title    string `json:"title" validate:"required"`
	Category string `json:"category" validate:"required"`
	Text     string `json:"text"`
	URL      string `json:"url"`
}

// NewComment is what we require from users when adding a Comment.
type NewComment struct {
	Text string `json:"comment" validate:"required"`
}
