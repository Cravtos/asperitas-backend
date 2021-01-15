package handlers

import (
	"context"
	gql2 "github.com/cravtos/asperitas-backend/business/data/gql"
	"github.com/cravtos/asperitas-backend/foundation/web"
	"github.com/cravtos/asperitas-backend/foundation/web/gql"
	"github.com/graphql-go/graphql"
	"net/http"
)

//todo think about names for everything

type GraphQLGroup struct {
	A      gql2.Access
	schema graphql.Schema
}

func (gqlg *GraphQLGroup) handle(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// get query
	opts := gql.NewRequestOptions(r)

	ctx = context.WithValue(ctx, gql2.Key, gqlg.A)
	// execute graphql query
	params := graphql.Params{
		Schema:         gqlg.schema,
		RequestString:  opts.Query,
		VariableValues: opts.Variables,
		OperationName:  opts.OperationName,
		Context:        ctx,
	}

	result := graphql.Do(params)
	web.Respond(ctx, w, result, http.StatusOK)
	return nil
}
