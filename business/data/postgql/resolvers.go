package postgql

import (
	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/cravtos/asperitas-backend/business/data/db"
	"github.com/cravtos/asperitas-backend/business/data/user"
	"github.com/cravtos/asperitas-backend/foundation/web"
	"github.com/google/uuid"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

//todo: move out authentication
func Hello(p graphql.ResolveParams) (interface{}, error) {
	return "World", nil
}

func postTitle(p graphql.ResolveParams) (interface{}, error) {
	if src, ok := p.Source.(Info); ok {
		return src.Title, nil
	}
	return nil, web.NewShutdownError("info missing from context")
}

func postType(p graphql.ResolveParams) (interface{}, error) {
	if src, ok := p.Source.(Info); ok {
		if src.Type == "url" {
			return "link", nil
		}
		return src.Type, nil
	}
	return nil, web.NewShutdownError("info missing from context")
}

func anyPost(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, web.NewShutdownError("postGQL missing from context")
	}
	dbs := db.NewDBset(a.log, a.db)

	postsDB, err := dbs.SelectAllPosts(p.Context)
	if err != nil {
		return nil, newPrivateError(err.Error())
	}

	post := convertPost(postsDB[0])
	post, err = a.fillInfo(p.Context, post)
	if err != nil {
		return nil, newPrivateError(err.Error())
	}
	return post, nil
}

func posts(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, web.NewShutdownError("postGQL missing from context")
	}
	dbs := db.NewDBset(a.log, a.db)

	category, userID := parseCatAndUser(p.Args["category"], p.Args["user_id"])
	postsDB, err := dbs.ObtainPosts(p.Context, category, userID)
	if err != nil {
		return nil, newPrivateError(err.Error())
	}

	posts := convertPosts(postsDB)
	posts, err = a.fillInfos(p.Context, posts)
	if err != nil {
		return nil, newPrivateError(err.Error())
	}
	return posts, nil
}

func postURL(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Info)
	if !ok {
		return nil, web.NewShutdownError("info missing from context")
	}
	if src.Type != "url" {
		return nil, newPrivateError("provided post is not link post")
	}
	return src.Payload, nil
}

func postText(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Info)
	if !ok {
		return nil, web.NewShutdownError("info missing from context")
	}
	if src.Type != "text" {
		return nil, newPrivateError("provided post is not text post")
	}
	return src.Payload, nil
}

func postID(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Info)
	if !ok {
		return nil, web.NewShutdownError("info missing from context")
	}
	return src.ID, nil
}

func postScore(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Info)
	if !ok {
		return nil, web.NewShutdownError("info missing from context")
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
	return nil, web.NewShutdownError("info missing from context")
}

func postCategory(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Info)
	if !ok {
		return nil, web.NewShutdownError("info missing from context")
	}
	return src.Category, nil
}

func postDateCreated(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Info)
	if !ok {
		return nil, web.NewShutdownError("info missing from context")
	}
	return src.DateCreated, nil
}

func postAuthor(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Info)
	if !ok {
		return nil, web.NewShutdownError("info missing from context")
	}
	return src.Author, nil
}

func postUpvotePercentage(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Info)
	if !ok {
		return nil, web.NewShutdownError("info missing from context")
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
		return nil, web.NewShutdownError("user missing from context")
	}
	return src.Username, nil
}

func authorID(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(*Author)
	if !ok {
		return nil, web.NewShutdownError("user missing from context")
	}
	return src.ID, nil
}

func voteVote(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Vote)
	if !ok {
		return nil, web.NewShutdownError("vote missing from context")
	}
	return src.Vote, nil
}

func voteUserID(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Vote)
	if !ok {
		return nil, web.NewShutdownError("vote missing from context")
	}
	return src.UserID, nil
}

func commentID(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Comment)
	if !ok {
		return nil, web.NewShutdownError("comment missing from context")
	}
	return src.ID, nil
}

func commentBody(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Comment)
	if !ok {
		return nil, web.NewShutdownError("comment missing from context")
	}
	return src.Body, nil
}

func commentAuthor(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Comment)
	if !ok {
		return nil, web.NewShutdownError("comment missing from context")
	}
	return src.Author, nil
}

func commentDateCreated(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Comment)
	if !ok {
		return nil, web.NewShutdownError("comment missing from context")
	}
	return src.DateCreated, nil
}

func postVotes(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Info)
	if !ok {
		return nil, web.NewShutdownError("post missing from context")
	}

	return src.Votes, nil
}

func postComments(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(Info)
	if !ok {
		return nil, web.NewShutdownError("post missing from context")
	}
	return src.Comments, nil
}

func authorPosts(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, web.NewShutdownError("postGQL missing from context")
	}
	src, ok := p.Source.(*Author)
	if !ok {
		return nil, web.NewShutdownError("author missing from context")
	}
	dbs := db.NewDBset(a.log, a.db)

	category, userID := parseCatAndUser(p.Args["category"], src.ID)
	postsDB, err := dbs.ObtainPosts(p.Context, category, userID)
	if err != nil {
		return nil, newPrivateError(err.Error())
	}

	posts := convertPosts(postsDB)
	posts, err = a.fillInfos(p.Context, posts)
	if err != nil {
		return nil, newPrivateError(err.Error())
	}
	return posts, nil
}

func voteUser(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, web.NewShutdownError("postGQL missing from context")
	}
	src, ok := p.Source.(Vote)
	if !ok {
		return nil, web.NewShutdownError("vote missing from context")
	}
	dbs := db.NewDBset(a.log, a.db)

	authorDB, err := dbs.GetUserByID(p.Context, src.UserID)
	if err != nil {
		return nil, newPrivateError(err.Error())
	}
	author := convertUser(authorDB)
	return author, nil
}

func commentPost(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, web.NewShutdownError("postGQL missing from context")
	}
	src, ok := p.Source.(Comment)
	if !ok {
		return nil, web.NewShutdownError("comment missing from context")
	}
	dbs := db.NewDBset(a.log, a.db)

	postDB, err := dbs.GetPostByID(p.Context, src.PostID)
	if err != nil {
		return nil, newPrivateError(err.Error())
	}
	post := convertPost(postDB)
	return post, nil
}

func post(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, web.NewShutdownError("postGQL missing from context")
	}
	dbs := db.NewDBset(a.log, a.db)

	postID, _ := p.Args["post_id"].(string)
	if _, err := uuid.Parse(postID); err != nil {
		return nil, newPublicError(ErrInvalidPostID.Error())
	}

	postDB, err := dbs.GetPostByID(p.Context, postID)
	if err != nil {
		if err == db.ErrPostNotFound {
			return nil, newPublicError(err.Error())
		}
		return nil, newPrivateError(err.Error())
	}
	post := convertPost(postDB)
	post, err = a.fillInfo(p.Context, post)
	if err != nil {
		return nil, newPrivateError(err.Error())
	}
	return post, nil
}

func postCreate(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, web.NewShutdownError("postGQL missing from context")
	}

	au, ok := p.Context.Value(KeyAuth).(*auth.Auth)
	if !ok {
		return nil, web.NewShutdownError("auth missing from context")
	}

	claims, err := au.ValidateString(p.Context.Value(KeyAuthHeader).(string))
	if err != nil {
		if err == auth.ErrExpectedBearer {
			return nil, newPublicError(err.Error())
		}
		return nil, newPrivateError(err.Error())
	}

	v, ok := p.Context.Value(web.KeyValues).(*web.Values)
	if !ok {
		return nil, web.NewShutdownError("web value missing from context")
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
		return nil, newPrivateError(err.Error())
	}

	if err := dbs.InsertVote(p.Context, newPost.ID, newPost.UserID, 1); err != nil {
		return nil, newPrivateError(err.Error())
	}
	p.Args["post_id"] = newPost.ID
	return post(p)
}

func postDelete(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, web.NewShutdownError("postGQL missing from context")
	}

	au, ok := p.Context.Value(KeyAuth).(*auth.Auth)
	if !ok {
		return nil, web.NewShutdownError("auth missing from context")
	}

	dbs := db.NewDBset(a.log, a.db)

	claims, err := au.ValidateString(p.Context.Value(KeyAuthHeader).(string))
	if err != nil {
		if err == auth.ErrExpectedBearer {
			return nil, newPublicError(err.Error())
		}
		return nil, newPrivateError(err.Error())
	}

	postID := p.Args["post_id"].(string)
	if _, err := uuid.Parse(postID); err != nil {
		return nil, ErrInvalidPostID
	}

	postDB, err := dbs.GetPostByID(p.Context, postID)
	if err != nil {
		if err == db.ErrPostNotFound {
			return nil, newPublicError(err.Error())
		}
		return nil, newPrivateError(err.Error())
	}

	if claims.User.ID != postDB.UserID {
		return nil, ErrForbidden
	}

	if err := dbs.DeletePost(p.Context, postID); err != nil {
		return nil, newPrivateError(err.Error())
	}
	return post(p)
}

func commentCreate(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, web.NewShutdownError("postGQL missing from context")
	}

	au, ok := p.Context.Value(KeyAuth).(*auth.Auth)
	if !ok {
		return nil, web.NewShutdownError("auth missing from context")
	}

	claims, err := au.ValidateString(p.Context.Value(KeyAuthHeader).(string))
	if err != nil {
		if err == auth.ErrExpectedBearer {
			return nil, newPublicError(err.Error())
		}
		return nil, newPrivateError(err.Error())
	}

	v, ok := p.Context.Value(web.KeyValues).(*web.Values)
	if !ok {
		return nil, web.NewShutdownError("web value missing from context")
	}

	dbs := db.NewDBset(a.log, a.db)

	postID := p.Args["post_id"].(string)
	if _, err := uuid.Parse(postID); err != nil {
		return nil, ErrInvalidPostID
	}

	if err := dbs.CheckPost(p.Context, postID); err != nil {
		if err == db.ErrPostNotFound {
			return nil, newPublicError(err.Error())
		}
		return nil, newPrivateError(err.Error())
	}

	ncDB := db.CommentDB{
		DateCreated: v.Now,
		PostID:      postID,
		AuthorID:    claims.User.ID,
		Body:        p.Args["text"].(string),
		ID:          uuid.New().String(),
	}

	if err := dbs.CreateComment(p.Context, ncDB); err != nil {
		return nil, newPrivateError(err.Error())
	}
	return post(p)
}

func commentDelete(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, web.NewShutdownError("postGQL missing from context")
	}

	au, ok := p.Context.Value(KeyAuth).(*auth.Auth)
	if !ok {
		return nil, web.NewShutdownError("auth missing from context")
	}

	claims, err := au.ValidateString(p.Context.Value(KeyAuthHeader).(string))
	if err != nil {
		if err == auth.ErrExpectedBearer {
			return nil, newPublicError(err.Error())
		}
		return nil, newPrivateError(err.Error())
	}

	dbs := db.NewDBset(a.log, a.db)

	postID := p.Args["post_id"].(string)
	if _, err := uuid.Parse(postID); err != nil {
		return nil, ErrInvalidPostID
	}
	commentID := p.Args["comment_id"].(string)
	if _, err := uuid.Parse(commentID); err != nil {
		return nil, ErrInvalidCommentID
	}

	if err := dbs.CheckPost(p.Context, postID); err != nil {
		if err == db.ErrPostNotFound {
			return nil, newPublicError(err.Error())
		}
		return nil, newPrivateError(err.Error())
	}

	commentDB, err := dbs.GetCommentByID(p.Context, commentID)
	if err != nil {
		if err == db.ErrPostNotFound {
			return nil, newPublicError(err.Error())
		}
		return nil, newPrivateError(err.Error())
	}

	if claims.User.ID != commentDB.AuthorID {
		return nil, ErrForbidden
	}

	if err := dbs.DeleteComment(p.Context, commentID); err != nil {
		return nil, newPrivateError(err.Error())
	}

	return post(p)
}

func upvote(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, web.NewShutdownError("postGQL missing from context")
	}

	au, ok := p.Context.Value(KeyAuth).(*auth.Auth)
	if !ok {
		return nil, web.NewShutdownError("auth missing from context")
	}

	claims, err := au.ValidateString(p.Context.Value(KeyAuthHeader).(string))
	if err != nil {
		if err == auth.ErrExpectedBearer {
			return nil, newPublicError(err.Error())
		}
		return nil, newPrivateError(err.Error())
	}

	dbs := db.NewDBset(a.log, a.db)

	postID := p.Args["post_id"].(string)
	if _, err := uuid.Parse(postID); err != nil {
		return nil, ErrInvalidPostID
	}
	if err := dbs.CheckPost(p.Context, postID); err != nil {
		if err == db.ErrPostNotFound {
			return nil, newPublicError(err.Error())
		}
		return nil, newPrivateError(err.Error())
	}

	if err := dbs.CheckVote(p.Context, postID, claims.User.ID); err != nil {
		if err != db.ErrVoteNotFound {
			return nil, newPrivateError(err.Error())
		}
		if err := dbs.InsertVote(p.Context, postID, claims.User.ID, 1); err != nil {
			return nil, newPrivateError(err.Error())
		}
	} else {
		if err := dbs.UpdateVote(p.Context, postID, claims.User.ID, 1); err != nil {
			return nil, newPrivateError(err.Error())
		}
	}

	return post(p)
}

func downvote(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, web.NewShutdownError("postGQL missing from context")
	}

	au, ok := p.Context.Value(KeyAuth).(*auth.Auth)
	if !ok {
		return nil, web.NewShutdownError("auth missing from context")
	}

	claims, err := au.ValidateString(p.Context.Value(KeyAuthHeader).(string))
	if err != nil {
		if err == auth.ErrExpectedBearer {
			return nil, newPublicError(err.Error())
		}
		return nil, newPrivateError(err.Error())
	}

	dbs := db.NewDBset(a.log, a.db)

	postID := p.Args["post_id"].(string)
	if _, err := uuid.Parse(postID); err != nil {
		return nil, ErrInvalidPostID
	}
	if err := dbs.CheckPost(p.Context, postID); err != nil {
		if err == db.ErrPostNotFound {
			return nil, newPublicError(err.Error())
		}
		return nil, newPrivateError(err.Error())
	}

	if err := dbs.CheckVote(p.Context, postID, claims.User.ID); err != nil {
		if err != db.ErrVoteNotFound {
			return nil, newPrivateError(err.Error())
		}
		if err := dbs.InsertVote(p.Context, postID, claims.User.ID, 0); err != nil {
			return nil, newPrivateError(err.Error())
		}
	} else {
		if err := dbs.UpdateVote(p.Context, postID, claims.User.ID, 0); err != nil {
			return nil, newPrivateError(err.Error())
		}
	}

	return post(p)
}

func unvote(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, web.NewShutdownError("postGQL missing from context")
	}

	au, ok := p.Context.Value(KeyAuth).(*auth.Auth)
	if !ok {
		return nil, web.NewShutdownError("auth missing from context")
	}

	claims, err := au.ValidateString(p.Context.Value(KeyAuthHeader).(string))
	if err != nil {
		if err == auth.ErrExpectedBearer {
			return nil, newPublicError(err.Error())
		}
		return nil, newPrivateError(err.Error())
	}

	dbs := db.NewDBset(a.log, a.db)

	postID := p.Args["post_id"].(string)
	if _, err := uuid.Parse(postID); err != nil {
		return nil, ErrInvalidPostID
	}
	if err := dbs.CheckPost(p.Context, postID); err != nil {
		if err == db.ErrPostNotFound {
			return nil, newPublicError(err.Error())
		}
		return nil, newPrivateError(err.Error())
	}

	if err := dbs.CheckVote(p.Context, postID, claims.User.ID); err != nil {
		if err != db.ErrVoteNotFound {
			return nil, newPrivateError(err.Error())
		} else {
			return post(p)
		}
	}
	if err := dbs.DeleteVote(p.Context, postID, claims.User.ID); err != nil {
		return nil, newPrivateError(err.Error())
	}

	return post(p)
}

func register(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, web.NewShutdownError("postGQL missing from context")
	}

	v, ok := p.Context.Value(web.KeyValues).(*web.Values)
	if !ok {
		return nil, web.NewShutdownError("web values missing from context")
	}

	au, ok := p.Context.Value(KeyAuth).(*auth.Auth)
	if !ok {
		return nil, web.NewShutdownError("auth missing from context")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(p.Args["password"].(string)), bcrypt.DefaultCost)
	if err != nil {
		return nil, web.NewShutdownError("generating password hash: " + err.Error())
	}

	usr := db.FullUserDB{
		ID:           uuid.New().String(),
		Name:         p.Args["name"].(string),
		PasswordHash: hash,
		DateCreated:  v.Now,
	}

	dbs := db.NewDBset(a.log, a.db)

	dbs.CreateUser(p.Context, usr)

	claims, err := user.New(a.log, a.db).Authenticate(p.Context, usr.Name, p.Args["password"].(string), v.Now)
	if err != nil {
		switch err {
		case user.ErrAuthenticationFailure:
			return nil, newPublicError(err.Error())
		default:
			return nil, newPrivateError(errors.Wrapf(err, "unable to authenticate user with name %s", usr.Name).Error())
		}
	}

	// todo: consider HS256
	kid := au.GetKID()
	Token, err := au.GenerateToken(kid, claims)
	if err != nil {
		return nil, newPrivateError(errors.Wrapf(err, "generating token").Error())
	}

	return auth.Data{
		Token: Token,
		User: auth.User{
			Username: usr.Name,
			ID:       usr.ID,
		},
	}, nil
}

func signIn(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, web.NewShutdownError("postGQL missing from context")
	}
	v, ok := p.Context.Value(web.KeyValues).(*web.Values)
	if !ok {
		return nil, web.NewShutdownError("web values missing from context")
	}

	au, ok := p.Context.Value(KeyAuth).(*auth.Auth)
	if !ok {
		return nil, web.NewShutdownError("auth missing from context")
	}
	Name := p.Args["name"].(string)
	claims, err := user.New(a.log, a.db).Authenticate(p.Context, Name, p.Args["password"].(string), v.Now)
	if err != nil {
		switch err {
		case user.ErrAuthenticationFailure:
			return nil, newPublicError(err.Error())
		default:
			return nil, newPrivateError(errors.Wrapf(err, "unable to authenticate user with name %s", Name).Error())
		}
	}

	// todo: consider HS256
	kid := au.GetKID()
	Token, err := au.GenerateToken(kid, claims)
	if err != nil {
		return nil, newPrivateError(errors.Wrapf(err, "generating token").Error())
	}

	return auth.Data{
		Token: Token,
		User: auth.User{
			Username: claims.User.Username,
			ID:       claims.User.ID,
		},
	}, nil
}

func userID(p graphql.ResolveParams) (interface{}, error) {
	if src, ok := p.Source.(auth.User); ok {
		return src.ID, nil
	}
	return nil, web.NewShutdownError("auth.User missing from context")
}

func username(p graphql.ResolveParams) (interface{}, error) {
	if src, ok := p.Source.(auth.User); ok {
		return src.Username, nil
	}
	return nil, web.NewShutdownError("auth.User missing from context")
}

func authUser(p graphql.ResolveParams) (interface{}, error) {
	if src, ok := p.Source.(auth.Data); ok {
		return src.User, nil
	}
	return nil, web.NewShutdownError("auth.Data missing from context")
}

func token(p graphql.ResolveParams) (interface{}, error) {
	if src, ok := p.Source.(auth.Data); ok {
		return src.Token, nil
	}
	return nil, web.NewShutdownError("auth.Data missing from context")
}
