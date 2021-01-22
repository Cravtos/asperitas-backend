package posts

import (
	"time"
)

// Info represents an individual posts before sending it to users
type Info struct {
	ID               string
	Views            int
	Type             string
	Title            string
	Category         string
	Payload          string
	Score            int
	UpvotePercentage int
	DateCreated      time.Time
	UserID           string
	Author           *Author
	Votes            []Vote
	Comments         []Comment
}

// Author represents info about authorType
type Author struct {
	Username string
	ID       string
}

// Vote represents info about users vote.
type Vote struct {
	UserID string
	Vote   int
}

// Comment represents info about comments for the posts prepared to be sent to users.
type Comment struct {
	PostID      string
	DateCreated time.Time
	Author      Author
	Body        string
	ID          string
}

// NewPost is what we require from users when adding a PostSetup.
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
