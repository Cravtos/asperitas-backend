package handlers

import (
	"github.com/cravtos/asperitas-backend/foundation/web"
	"log"
	"net/http"
	"context"
)

type corsGroup struct {
	log *log.Logger
}


func (cg corsGroup) allow(methods ...string) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return web.AllowMethods(ctx, w, methods)
	}
}