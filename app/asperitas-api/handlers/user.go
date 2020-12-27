package handlers

import (
	"context"
	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/cravtos/asperitas-backend/business/data/user"
	"github.com/cravtos/asperitas-backend/foundation/web"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"time"
)

type userGroup struct {
	user user.User
	auth *auth.Auth
}

func (ug userGroup) register(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "handlers.userGroup.register")
	defer span.End()

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	var nu user.NewUser
	if err := web.Decode(r, &nu); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	_, err := ug.user.Create(ctx, v.TraceID, nu, time.Now())
	if err != nil {
		return err
	}

	claims, err := ug.user.Authenticate(ctx, v.TraceID, time.Now(), nu.Name, nu.Password)
	if err != nil {
		return err
	}

	var tkn struct {
		Token string `json:"token"`
	}
	kid := ug.auth.GetKID()
	tkn.Token, err = ug.auth.GenerateToken(kid, claims)
	if err != nil {
		return errors.Wrapf(err, "generating token")
	}

	return web.Respond(ctx, w, tkn, http.StatusOK)
}

func (ug userGroup) login(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.SpanFromContext(ctx).Tracer().Start(ctx, "handlers.userGroup.login")
	defer span.End()

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	u := struct {
		Name string
		Password string
	}{}

	if err := web.Decode(r, &u); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	claims, err := ug.user.Authenticate(ctx, v.TraceID, time.Now(), u.Name, u.Password)
	if err != nil {
		return err
	}

	var tkn struct {
		Token string `json:"token"`
	}
	kid := ug.auth.GetKID()
	tkn.Token, err = ug.auth.GenerateToken(kid, claims)
	if err != nil {
		return errors.Wrapf(err, "generating token")
	}

	return web.Respond(ctx, w, tkn, http.StatusOK)
}