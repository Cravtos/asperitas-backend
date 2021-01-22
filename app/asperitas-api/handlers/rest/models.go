package rest

import "time"

// Author represents info about author
type Author struct {
	Username string `json:"username"`
	ID       string `json:"id"`
}

// Vote represents info about users vote.
type Vote struct {
	User string `json:"users"`
	Vote int    `json:"vote"`
}

// Comment represents info about comments for the posts prepared to be sent to users.
type Comment struct {
	DateCreated time.Time `json:"created"`
	Author      Author    `json:"author"`
	Body        string    `json:"body"`
	ID          string    `json:"id"`
}

// Info generalizes text and link posts
type Info interface {
	Info()
}

// InfoText represents an individual text posts which is sent to users.
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

// InfoLink represents an individual link posts which is sent to users.
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
