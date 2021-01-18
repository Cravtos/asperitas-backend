package handlers

import (
	"context"
	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/cravtos/asperitas-backend/business/data/postgql"
	"github.com/cravtos/asperitas-backend/foundation/web"
	"github.com/cravtos/asperitas-backend/foundation/web/gql"
	"github.com/graphql-go/graphql"
	"net/http"
)

type PostGroupGQL struct {
	P      postgql.PostGQL
	schema graphql.Schema
	a      *auth.Auth
}

func (gqlg *PostGroupGQL) handle(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// get query
	opts := gql.NewRequestOptions(r)

	ctx = context.WithValue(ctx, postgql.KeyPostGQL, gqlg.P)
	ctx = context.WithValue(ctx, postgql.KeyAuth, gqlg.a)
	ctx = context.WithValue(ctx, postgql.KeyAuthHeader, r.Header.Get("authorization"))
	// execute graphql query
	params := graphql.Params{
		Schema:         gqlg.schema,
		RequestString:  opts.Query,
		VariableValues: opts.Variables,
		OperationName:  opts.OperationName,
		Context:        ctx,
	}

	result := graphql.Do(params)
	return web.Respond(ctx, w, result, http.StatusOK)
}
