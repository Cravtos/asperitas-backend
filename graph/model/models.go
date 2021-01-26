package model

import (
	"time"
)

type Info interface {
	IsInfo()
}

func (PostLink) IsInfo() {}

func (PostText) IsInfo() {}

type AuthData struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type Author struct {
	Username string `json:"username"`
	AuthorID string `json:"author_id"`
}

type Comment struct {
	Body        string    `json:"body"`
	CommentID   string    `json:"comment_id"`
	Author      *Author   `json:"author"`
	DateCreated time.Time `json:"date_created"`
}

type PostLink struct {
	PostID           string     `json:"post_id"`
	Title            string     `json:"title"`
	Type             PostType   `json:"type"`
	Score            int        `json:"score"`
	Views            int        `json:"views"`
	Category         Category   `json:"category"`
	DateCreated      time.Time  `json:"date_created"`
	UpvotePercentage int        `json:"upvote_percentage"`
	Author           *Author    `json:"author"`
	Votes            []*Vote    `json:"votes"`
	Comments         []*Comment `json:"Comments"`
	URL              string     `json:"url"`
}

type PostText struct {
	PostID           string     `json:"post_id"`
	Title            string     `json:"title"`
	Type             PostType   `json:"type"`
	Score            int        `json:"score"`
	Views            int        `json:"views"`
	Category         Category   `json:"category"`
	DateCreated      time.Time  `json:"date_created"`
	UpvotePercentage int        `json:"upvote_percentage"`
	Author           *Author    `json:"author"`
	Votes            []*Vote    `json:"votes"`
	Comments         []*Comment `json:"Comments"`
	Text             string     `json:"text"`
}

type User struct {
	Username string `json:"username"`
	UserID   string `json:"user_id"`
}

type Vote struct {
	Vote     int    `json:"vote"`
	AuthorID string `json:"author_id"`
}

//todo think how to implement these resolvers
func (v *Vote) Author() *Author {
	return &Author{}
}

func (c *Comment) Post() Info {
	return PostLink{}
}

func (a *Author) Posts() []Info {
	return make([]Info, 0)
}
