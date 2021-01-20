package mid

import (
	"context"
	"log"
	"net/http"

	"github.com/cravtos/asperitas-backend/foundation/web"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// Run the next handler and catch any propagated error.
			if err := handler(ctx, w, r); err != nil {
				if web.IsErrorResponseGQL(err) {
					for _, er := range err.(*web.ResponseGQL).PrivateErrors {
						log.Printf("ERROR: executing utilgql %v\n", er)
					}
					if er := err.(*web.ResponseGQL).SendingError; er != nil {
						// Log the error.
						log.Printf("ERROR: %v", er)

						// Respond to the error.
						if err := web.RespondError(ctx, w, err); err != nil {
							return err
						}
					}
				} else {
					// Log the error.
					log.Printf("ERROR: %v", err)

					// Respond to the error.
					if err := web.RespondError(ctx, w, err); err != nil {
						return err
					}

					// If we receive the shutdown err we need to return it
					// back to the base handler to shutdown the service.
					if ok := web.IsShutdown(err); ok {
						return err
					}
				}
			}

			// The error has been handled so we can stop propagating it.
			return nil
		}

		return h
	}

	return m
}
