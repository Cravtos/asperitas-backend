package handlers

import (
	"context"
	"github.com/cravtos/asperitas-backend/foundation/web"
	"github.com/cravtos/asperitas-backend/foundation/web/gql"
	"github.com/graphql-go/graphql"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

type GraphQLGroup struct {
	log    *log.Logger
	db     *sqlx.DB
	schema graphql.Schema
}

func (gqlg *GraphQLGroup) handle(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// get query
	opts := gql.NewRequestOptions(r)

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
