package postgql

import (
	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/cravtos/asperitas-backend/business/data/posts"
	"github.com/cravtos/asperitas-backend/business/data/users"
	"github.com/cravtos/asperitas-backend/foundation/web"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

func postTitle(p graphql.ResolveParams) (interface{}, error) {
	if src, ok := p.Source.(posts.Info); ok {
		return src.Title, nil
	}
	return nil, web.NewShutdownError("info missing from context")
}

func postType(p graphql.ResolveParams) (interface{}, error) {
	if src, ok := p.Source.(posts.Info); ok {
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
	ps := posts.New(a.log, a.db)

	infos, err := ps.Query(p.Context)
	if err != nil {
		return nil, newPrivateError(err)
	}
	if len(infos) == 0 {
		return nil, newPublicError(errors.New("there is no infos at all"))
	}
	return infos[0], nil
}

func postsRes(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, web.NewShutdownError("postGQL missing from context")
	}
	ps := posts.New(a.log, a.db)

	category, userID := parseCatAndUser(p.Args["category"], p.Args["user_id"])
	infos, err := ps.ObtainPosts(p.Context, category, userID)
	if err != nil {
		return nil, newPrivateError(err)
	}
	return infos, nil
}

func postURL(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(posts.Info)
	if !ok {
		return nil, web.NewShutdownError("info missing from context")
	}
	if src.Type != "url" && src.Type != "link" {
		return nil, newPrivateError(errors.New("provided post is not link post"))
	}
	return src.Payload, nil
}

func postText(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(posts.Info)
	if !ok {
		return nil, web.NewShutdownError("info missing from context")
	}
	if src.Type != "text" {
		return nil, newPrivateError(errors.New("provided post is not text post"))
	}
	return src.Payload, nil
}

func postID(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(posts.Info)
	if !ok {
		return nil, web.NewShutdownError("info missing from context")
	}
	return src.ID, nil
}

func postScore(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(posts.Info)
	if !ok {
		return nil, web.NewShutdownError("info missing from context")
	}

	return src.Score, nil
}

func postViews(p graphql.ResolveParams) (interface{}, error) {
	if src, ok := p.Source.(posts.Info); ok {
		return src.Views, nil
	}
	return nil, web.NewShutdownError("info missing from context")
}

func postCategory(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(posts.Info)
	if !ok {
		return nil, web.NewShutdownError("info missing from context")
	}
	return src.Category, nil
}

func postDateCreated(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(posts.Info)
	if !ok {
		return nil, web.NewShutdownError("info missing from context")
	}
	return src.DateCreated, nil
}

func postAuthor(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(posts.Info)
	if !ok {
		return nil, web.NewShutdownError("info missing from context")
	}
	return src.Author, nil
}

func postUpvotePercentage(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(posts.Info)
	if !ok {
		return nil, web.NewShutdownError("info missing from context")
	}

	return src.UpvotePercentage, nil
}

func authorUsername(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(*posts.Author)
	if !ok {
		return nil, web.NewShutdownError("users missing from context")
	}
	return src.Username, nil
}

func authorID(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(*posts.Author)
	if !ok {
		return nil, web.NewShutdownError("users missing from context")
	}
	return src.ID, nil
}

func voteVote(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(posts.Vote)
	if !ok {
		return nil, web.NewShutdownError("vote missing from context")
	}
	return src.Vote, nil
}

func voteUserID(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(posts.Vote)
	if !ok {
		return nil, web.NewShutdownError("vote missing from context")
	}
	return src.UserID, nil
}

func commentID(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(posts.Comment)
	if !ok {
		return nil, web.NewShutdownError("comment missing from context")
	}
	return src.ID, nil
}

func commentBody(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(posts.Comment)
	if !ok {
		return nil, web.NewShutdownError("comment missing from context")
	}
	return src.Body, nil
}

func commentAuthor(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(posts.Comment)
	if !ok {
		return nil, web.NewShutdownError("comment missing from context")
	}
	return &src.Author, nil
}

func commentDateCreated(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(posts.Comment)
	if !ok {
		return nil, web.NewShutdownError("comment missing from context")
	}
	return src.DateCreated, nil
}

func postVotes(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(posts.Info)
	if !ok {
		return nil, web.NewShutdownError("post missing from context")
	}

	return src.Votes, nil
}

func postComments(p graphql.ResolveParams) (interface{}, error) {
	src, ok := p.Source.(posts.Info)
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
	src, ok := p.Source.(*posts.Author)
	if !ok {
		return nil, web.NewShutdownError("author missing from context")
	}
	ps := posts.New(a.log, a.db)

	category, userID := parseCatAndUser(p.Args["category"], src.ID)
	infos, err := ps.ObtainPosts(p.Context, category, userID)
	if err != nil {
		return nil, newPrivateError(err)
	}

	return infos, nil
}

func voteUser(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, web.NewShutdownError("postGQL missing from context")
	}
	src, ok := p.Source.(posts.Vote)
	if !ok {
		return nil, web.NewShutdownError("vote missing from context")
	}
	ps := posts.New(a.log, a.db)

	author, err := ps.AuthorByID(p.Context, src.UserID)
	if err != nil {
		return nil, newPrivateError(err)
	}
	//a.log.Println(author)
	return &author, nil
}

func commentPost(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, web.NewShutdownError("postGQL missing from context")
	}
	src, ok := p.Source.(posts.Comment)
	if !ok {
		return nil, web.NewShutdownError("comment missing from context")
	}
	ps := posts.New(a.log, a.db)

	info, err := ps.QueryByID(p.Context, src.PostID)
	if err != nil {
		return nil, newPrivateError(err)
	}

	return info, nil
}

func postRes(p graphql.ResolveParams) (interface{}, error) {
	a, ok := p.Context.Value(KeyPostGQL).(PostGQL)
	if !ok {
		return nil, web.NewShutdownError("postGQL missing from context")
	}
	ps := posts.New(a.log, a.db)

	postID, _ := p.Args["post_id"].(string)
	info, err := ps.QueryByID(p.Context, postID)
	if err != nil {
		switch err {
		case posts.ErrInvalidPostID:
			return nil, newPublicError(ErrInvalidPostID)
		case posts.ErrPostNotFound:
			return nil, newPublicError(err)
		default:
			return nil, newPrivateError(err)
		}
	}
	return info, nil
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
			return nil, newPublicError(err)
		}
		return nil, newPrivateError(err)
	}

	ps := posts.New(a.log, a.db)
	np := posts.NewPost{
		Title:    p.Args["title"].(string),
		Type:     p.Args["type"].(string),
		Category: p.Args["category"].(string),
		Text:     p.Args["payload"].(string),
		URL:      p.Args["payload"].(string),
	}

	v, ok := p.Context.Value(web.KeyValues).(*web.Values)
	if !ok {
		return nil, web.NewShutdownError("web value missing from context")
	}
	info, err := ps.Create(p.Context, claims, np, v.Now)
	if err != nil {
		return nil, newPrivateError(err)
	}

	return info, nil
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

	ps := posts.New(a.log, a.db)

	claims, err := au.ValidateString(p.Context.Value(KeyAuthHeader).(string))
	if err != nil {
		if err == auth.ErrExpectedBearer {
			return nil, newPublicError(err)
		}
		return nil, newPrivateError(err)
	}

	postID := p.Args["post_id"].(string)
	info, err := ps.Delete(p.Context, claims, postID)
	if err != nil {
		switch err {
		case posts.ErrInvalidPostID:
			return nil, ErrInvalidPostID
		case posts.ErrPostNotFound:
			return nil, newPublicError(err)
		case posts.ErrForbidden:
			return nil, ErrForbidden
		default:
			return nil, newPrivateError(err)
		}
	}
	return info, nil
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
			return nil, newPublicError(err)
		}
		return nil, newPrivateError(err)
	}

	v, ok := p.Context.Value(web.KeyValues).(*web.Values)
	if !ok {
		return nil, web.NewShutdownError("web value missing from context")
	}

	ps := posts.New(a.log, a.db)

	postID := p.Args["post_id"].(string)

	nc := posts.NewComment{Text: p.Args["text"].(string)}
	info, err := ps.CreateComment(p.Context, claims, nc, postID, v.Now)
	if err != nil {
		switch err {
		case posts.ErrInvalidPostID:
			return nil, ErrInvalidPostID
		case posts.ErrPostNotFound:
			return nil, newPublicError(err)
		default:
			return nil, newPrivateError(err)
		}
	}
	return info, nil
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
			return nil, newPublicError(err)
		}
		return nil, newPrivateError(err)
	}

	ps := posts.New(a.log, a.db)

	postID := p.Args["post_id"].(string)
	commentID := p.Args["comment_id"].(string)

	info, err := ps.DeleteComment(p.Context, claims, postID, commentID)
	if err != nil {
		switch err {
		case posts.ErrInvalidPostID:
			return nil, ErrInvalidPostID
		case posts.ErrInvalidCommentID:
			return nil, ErrInvalidCommentID
		case posts.ErrPostNotFound:
			return nil, newPublicError(err)
		case posts.ErrCommentNotFound:
			return nil, newPublicError(err)
		case posts.ErrForbidden:
			return nil, ErrForbidden
		default:
			return nil, newPrivateError(err)
		}
	}

	return info, nil
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
			return nil, newPublicError(err)
		}
		return nil, newPrivateError(err)
	}

	ps := posts.New(a.log, a.db)

	postID := p.Args["post_id"].(string)
	info, err := ps.Vote(p.Context, claims, postID, 1)
	if err != nil {
		switch err {
		case posts.ErrInvalidPostID:
			return nil, ErrInvalidPostID
		default:
			return nil, newPrivateError(err)
		}
	}

	return info, nil
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
			return nil, newPublicError(err)
		}
		return nil, newPrivateError(err)
	}

	ps := posts.New(a.log, a.db)

	postID := p.Args["post_id"].(string)
	info, err := ps.Vote(p.Context, claims, postID, -1)
	if err != nil {
		switch err {
		case posts.ErrInvalidPostID:
			return nil, ErrInvalidPostID
		default:
			return nil, newPrivateError(err)
		}
	}
	return info, nil
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
			return nil, newPublicError(err)
		}
		return nil, newPrivateError(err)
	}

	ps := posts.New(a.log, a.db)

	postID := p.Args["post_id"].(string)
	info, err := ps.Unvote(p.Context, claims, postID)
	if err != nil {
		switch err {
		case posts.ErrInvalidPostID:
			return nil, ErrInvalidPostID
		default:
			return nil, newPrivateError(err)
		}
	}
	return info, nil
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
	us := users.New(a.log, a.db)

	nu := users.NewUser{
		Name:     p.Args["name"].(string),
		Password: p.Args["password"].(string),
	}

	_, err := us.Create(p.Context, nu, v.Now)
	if err != nil {
		return nil, newPrivateError(errors.Wrapf(err, "unable to create users with name %s", nu.Name))
	}

	claims, err := us.Authenticate(p.Context, nu.Name, nu.Password, v.Now)
	if err != nil {
		switch err {
		case users.ErrAuthenticationFailure:
			return nil, newPublicError(err)
		default:
			return nil, newPrivateError(errors.Wrapf(err, "unable to authenticate users with name %s", nu.Name))
		}
	}

	kid := au.GetKID()
	Token, err := au.GenerateToken(kid, claims)
	if err != nil {
		return nil, newPrivateError(errors.Wrapf(err, "generating token"))
	}

	return auth.Data{
		Token: Token,
		User:  claims.User,
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
	Password := p.Args["password"].(string)
	claims, err := users.New(a.log, a.db).Authenticate(p.Context, Name, Password, v.Now)
	if err != nil {
		switch err {
		case users.ErrAuthenticationFailure:
			return nil, newPublicError(err)
		default:
			return nil, newPrivateError(errors.Wrapf(err, "unable to authenticate users with name %s", Name))
		}
	}

	kid := au.GetKID()
	Token, err := au.GenerateToken(kid, claims)
	if err != nil {
		return nil, newPrivateError(errors.Wrapf(err, "generating token"))
	}

	return auth.Data{
		Token: Token,
		User:  claims.User,
	}, nil
}

func userID(p graphql.ResolveParams) (interface{}, error) {
	if src, ok := p.Source.(auth.User); ok {
		return src.ID, nil
	}
	return nil, web.NewShutdownError("auth.UserID missing from context")
}

func username(p graphql.ResolveParams) (interface{}, error) {
	if src, ok := p.Source.(auth.User); ok {
		return src.Username, nil
	}
	return nil, web.NewShutdownError("auth.UserID missing from context")
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
