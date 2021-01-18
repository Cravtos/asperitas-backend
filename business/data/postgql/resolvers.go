package postgql

import (
	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/cravtos/asperitas-backend/business/data/db"
	"github.com/cravtos/asperitas-backend/foundation/web"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

//todo: move out authentication
//TODO think about error handling
func Hello(p graphql.ResolveParams) (interface{}, error) {
	return "World", nil
}

func postTitle(p graphql.ResolveParams) (interface{}, error) {
	if src, ok := p.Source.(Info); ok {
		return src.Title, nil
	}
	return nil, errors.New("post missing from context")
}

func postType(p graphql.ResolveParams) (interface{}, error) {
	if src, ok := p.Source.(Info); ok {
		if src.Type == "url" {
			return "link", nil
		}
		return src.Type, nil
	}
	return nil, errors.New("post missing from context")
}

func anyPost(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, errors.New("postGQL missing from context")
	}
	dbs := db.NewDBset(a.log, a.db)

	postsDB, err := dbs.SelectAllPosts(p.Context)
	if err != nil {
		return nil, err
	}

	post := convertPost(postsDB[0])
	post, err = a.fillInfo(p.Context, post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func posts(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, errors.New("postGQL missing from context")
	}
	dbs := db.NewDBset(a.log, a.db)

	category, userID := parseCatAndUser(p.Args["category"], p.Args["user_id"])
	postsDB, err := dbs.ObtainPosts(p.Context, category, userID)
	if err != nil {
		return nil, err
	}

	posts := convertPosts(postsDB)
	posts, err = a.fillInfos(p.Context, posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func postURL(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Info)
	if !ok {
		return nil, errors.New("post missing from context")
	}
	if src.Type != "url" {
		return nil, errors.New("provided post is not link post")
	}
	return src.Payload, nil
}

func postText(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Info)
	if !ok {
		return nil, errors.New("post missing from context")
	}
	if src.Type != "text" {
		return nil, errors.New("provided post is not text post")
	}
	return src.Payload, nil
}

func postID(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Info)
	if !ok {
		return nil, errors.New("post missing from context")
	}
	return src.ID, nil
}

func postScore(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Info)
	if !ok {
		return nil, errors.New("post missing from context")
	}

	score := 0
	for _, vote := range src.Votes {
		score += vote.Vote
	}
	return score, nil
}

func postViews(p graphql.ResolveParams) (interface{}, error) {
	if src, ok := p.Source.(Info); ok {
		return src.Views, nil
	}
	return nil, errors.New("post missing from context")
}

func postCategory(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Info)
	if !ok {
		return nil, errors.New("post missing from context")
	}
	return src.Category, nil
}

func postDateCreated(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Info)
	if !ok {
		return nil, errors.New("post missing from context")
	}
	return src.DateCreated, nil
}

func postAuthor(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Info)
	if !ok {
		return nil, errors.New("post missing from context")
	}
	return src.Author, nil
}

func postUpvotePercentage(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Info)
	if !ok {
		return nil, errors.New("post missing from context")
	}

	if len(src.Votes) == 0 {
		return 0, nil
	}
	var positive float32

	for _, vote := range src.Votes {
		if vote.Vote == 1 {
			positive++
		}
	}

	return int(positive / float32(len(src.Votes)) * 100), nil
}

func authorUsername(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(*Author)
	if !ok {
		return nil, errors.New("author missing from context")
	}
	return src.Username, nil
}

func authorID(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(*Author)
	if !ok {
		return nil, errors.New("author missing from context")
	}
	return src.ID, nil
}

func voteVote(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Vote)
	if !ok {
		return nil, errors.New("vote missing from context")
	}
	return src.Vote, nil
}

func voteUserID(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Vote)
	if !ok {
		return nil, errors.New("vote missing from context")
	}
	return src.UserID, nil
}

func commentID(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Comment)
	if !ok {
		return nil, errors.New("comment missing from context")
	}
	return src.ID, nil
}

func commentBody(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Comment)
	if !ok {
		return nil, errors.New("comment missing from context")
	}
	return src.Body, nil
}

func commentAuthor(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Comment)
	if !ok {
		return nil, errors.New("comment missing from context")
	}
	return src.Author, nil
}

func commentDateCreated(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Comment)
	if !ok {
		return nil, errors.New("comment missing from context")
	}
	return src.DateCreated, nil
}

func postVotes(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Info)
	if !ok {
		return nil, errors.New("post missing from context")
	}

	return src.Votes, nil
}

func postComments(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Info)
	if !ok {
		return nil, errors.New("post missing from context")
	}
	return src.Comments, nil
}

func authorPosts(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, errors.New("postGQL missing from context")
	}
	src, ok := p.Source.(*Author)
	if !ok {
		return nil, errors.New("author missing from context")
	}
	dbs := db.NewDBset(a.log, a.db)

	category, userID := parseCatAndUser(p.Args["category"], src.ID)
	postsDB, err := dbs.ObtainPosts(p.Context, category, userID)
	if err != nil {
		return nil, err
	}

	posts := convertPosts(postsDB)
	posts, err = a.fillInfos(p.Context, posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func voteUser(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, errors.New("postGQL missing from context")
	}
	src, ok := p.Source.(Vote)
	if !ok {
		return nil, errors.New("vote missing from context")
	}
	dbs := db.NewDBset(a.log, a.db)

	authorDB, err := dbs.GetUserByID(p.Context, src.UserID)
	if err != nil {
		return nil, err
	}
	author := convertUser(authorDB)
	return author, nil
}

func commentPost(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, errors.New("postGQL missing from context")
	}
	src, ok := p.Source.(Comment)
	if !ok {
		return nil, errors.New("comment missing from context")
	}
	dbs := db.NewDBset(a.log, a.db)

	postDB, err := dbs.GetPostByID(p.Context, src.PostID)
	if err != nil {
		return nil, err
	}
	post := convertPost(postDB)
	return post, nil
}

func post(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, errors.New("postGQL missing from context")
	}
	dbs := db.NewDBset(a.log, a.db)

	postDB, err := dbs.GetPostByID(p.Context, p.Args["post_id"].(string))
	if err != nil {
		return nil, err
	}
	post := convertPost(postDB)
	post, err = a.fillInfo(p.Context, post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

//todo replace RequestErrors with new GQL errors
func postCreate(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, errors.New("postGQL missing from context")
	}
	authStr, ok := p.Context.Value(KeyAuthHeader).(string)
	if !ok {
		err := errors.New("expected authorization header format: bearer <token>")
		return nil, web.NewRequestError(err, http.StatusUnauthorized)
	}
	v, ok := p.Context.Value(web.KeyValues).(*web.Values)
	if !ok {
		return nil, web.NewShutdownError("web value missing from context")
	}
	// Parse the authorization header.
	parts := strings.Split(authStr, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		err := errors.New("expected authorization header format: bearer <token>")
		return nil, web.NewRequestError(err, http.StatusUnauthorized)
	}

	au, ok := p.Context.Value(KeyAuth).(*auth.Auth)
	if !ok {
		return nil, errors.New("au.Auth missing from context")
	}

	// Validate the token is signed by us.
	claims, err := au.ValidateToken(parts[1])
	if err != nil {
		return nil, web.NewRequestError(err, http.StatusUnauthorized)
	}

	dbs := db.NewDBset(a.log, a.db)
	newPost := db.PostDB{
		ID:          uuid.New().String(),
		Views:       0,
		Title:       p.Args["title"].(string),
		Type:        p.Args["type"].(string),
		Category:    p.Args["category"].(string),
		Payload:     p.Args["payload"].(string),
		DateCreated: v.Now,
		UserID:      claims.User.ID,
	}

	if err := dbs.InsertPost(p.Context, newPost); err != nil {
		return nil, err
	}

	if err := dbs.InsertVote(p.Context, newPost.ID, newPost.UserID, 1); err != nil {
		return nil, err
	}
	p.Args["post_id"] = newPost.ID
	return post(p)
}
