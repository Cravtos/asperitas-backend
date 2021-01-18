package postgql

import "time"

// Info represents an individual post before sending it to user
type Info struct {
	ID          string
	Views       int
	Type        string
	Title       string
	Category    string
	Payload     string
	DateCreated time.Time
	UserID      string
	Author      *Author
	Votes       []Vote
	Comments    []Comment
}

// Author represents info about authorType
type Author struct {
	Username string
	ID       string
}

// Vote represents info about user vote.
type Vote struct {
	UserID string
	Vote   int
}

// Comment represents info about comments for the post prepared to be sent to user.
type Comment struct {
	PostID      string
	DateCreated time.Time
	Author      Author
	Body        string
	ID          string
}
