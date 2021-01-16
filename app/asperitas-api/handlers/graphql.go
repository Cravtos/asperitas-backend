package handlers

import (
	"context"
	"github.com/cravtos/asperitas-backend/business/data/postgql"
	"github.com/cravtos/asperitas-backend/foundation/web"
	"github.com/cravtos/asperitas-backend/foundation/web/gql"
	"github.com/graphql-go/graphql"
	"net/http"
)

type PostGroupGQL struct {
	P      postgql.PostGQL
	schema graphql.Schema
}

func (gqlg *PostGroupGQL) handle(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// get query
	opts := gql.NewRequestOptions(r)

	ctx = context.WithValue(ctx, postgql.Key, gqlg.P)
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
