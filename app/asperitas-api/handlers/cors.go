package handlers

import (
	"context"
	"github.com/cravtos/asperitas-backend/foundation/web"
	"log"
	"net/http"
)

type corsGroup struct {
	log *log.Logger
}

// allow responds with "Access-Control-Allow-Methods" set to methods and
// "Access-Control-Allow-Headers" set to "authorization, content-type"
func (cg corsGroup) allow(methods ...string) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return web.AllowMethods(ctx, w, methods)
	}
}
