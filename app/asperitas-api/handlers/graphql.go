package handlers

import (
	"context"
	"fmt"
	"github.com/cravtos/asperitas-backend/business/auth"
	"github.com/cravtos/asperitas-backend/business/data/postgql"
	"github.com/cravtos/asperitas-backend/foundation/utilgql"
	"github.com/cravtos/asperitas-backend/foundation/web"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"net/http"
)

type PostGroupGQL struct {
	P      postgql.PostGQL
	schema graphql.Schema
	auth   *auth.Auth
}

func (gqlg *PostGroupGQL) handle(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// get query
	opts := utilgql.NewRequestOptions(r)

	ctx = context.WithValue(ctx, postgql.KeyPostGQL, gqlg.P)
	ctx = context.WithValue(ctx, postgql.KeyAuth, gqlg.auth)
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

	privateErrs := make([]error, 0)
	for i := range result.Errors {
		err, ok := result.Errors[i].OriginalError().(*gqlerrors.Error)
		if ok {
			err2 := err.OriginalError
			switch err2.(type) {
			case *web.Shutdown:
				return err
			case *postgql.PrivateError:
				privateErrs = append(privateErrs, err2)
				result.Errors[i].Message = "Internal server error"
				result.Errors[i].Locations = nil
				result.Errors[i].Path = nil
			default:
				fmt.Println("hey  ", err)
			}
		}
	}
	return web.RespondGQL(ctx, w, &web.ResponseGQL{
		PrivateErrors: privateErrs,
		Data:          result,
	})
}
