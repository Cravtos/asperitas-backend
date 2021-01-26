package rest

import (
	"context"
	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/cravtos/asperitas-backend/business/data/users"
	"github.com/cravtos/asperitas-backend/foundation/web"
	"github.com/pkg/errors"
	"net/http"
)

type UserGroup struct {
	User users.User
	Auth *auth.Auth
}

func (ug UserGroup) Register(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return errors.New("web values missing from context")
	}

	var nu users.NewUser
	if err := web.Decode(r, &nu); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	_, err := ug.User.Create(ctx, nu, v.Now)
	if err != nil {
		return errors.Wrapf(err, "unable to create users with name %s", nu.Name)
	}

	claims, err := ug.User.Authenticate(ctx, nu.Name, nu.Password, v.Now)
	if err != nil {
		switch err {
		case users.ErrAuthenticationFailure:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "unable to authenticate users with name %s", nu.Name)
		}
	}

	var tkn struct {
		Token string `json:"token"`
	}
	// todo: consider HS256
	kid := ug.Auth.GetKID()
	tkn.Token, err = ug.Auth.GenerateToken(kid, claims)
	if err != nil {
		return errors.Wrapf(err, "generating token")
	}

	return web.Respond(ctx, w, tkn, http.StatusOK)
}

func (ug UserGroup) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return errors.New("web values missing from context")
	}

	u := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}

	if err := web.Decode(r, &u); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	claims, err := ug.User.Authenticate(ctx, u.Username, u.Password, v.Now)
	if err != nil {
		switch err {
		case users.ErrAuthenticationFailure:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "unable to authenticate users with name %s", u.Username)
		}
	}

	var tkn struct {
		Token string `json:"token"`
	}
	kid := ug.Auth.GetKID()
	tkn.Token, err = ug.Auth.GenerateToken(kid, claims)
	if err != nil {
		return errors.Wrapf(err, "generating token")
	}

	return web.Respond(ctx, w, tkn, http.StatusOK)
}
