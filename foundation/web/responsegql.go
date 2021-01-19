package web

import (
	"context"
	"github.com/graphql-go/graphql"
	"net/http"
)

type ResponseGQL struct {
	SendingError error
	Errors       []error
	Data         *graphql.Result
}

func (err ResponseGQL) Error() string {
	str := ""
	for _, s := range err.Errors {
		str += s.Error()
		str += "\n"
	}
	return str
}

func IsErrorResponseGQL(err error) bool {
	_, ok := err.(*ResponseGQL)
	return ok
}

// RespondGQL sends GQL reponse back to the client.
func RespondGQL(ctx context.Context, w http.ResponseWriter, gqlData *ResponseGQL) error {
	respError := Respond(ctx, w, gqlData.Data, http.StatusOK)
	gqlData.SendingError = respError
	return gqlData
}
