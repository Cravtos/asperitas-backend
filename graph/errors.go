package graph

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/cravtos/asperitas-backend/foundation/web"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"log"
	"os"
	"syscall"
)

type setup struct {
	log      *log.Logger
	shutdown chan os.Signal
}

func NewSetup(log *log.Logger, shutdown chan os.Signal) *setup {
	return &setup{log: log, shutdown: shutdown}
}

type PublicError struct {
	message string
}

type PrivateError struct {
	message string
}

//newPublicError returns new error that will be shown to client
func newPublicError(msg error) *PublicError {
	return &PublicError{message: msg.Error()}
}

//newPublicError returns new error that won`t be shown to client
func newPrivateError(msg error) *PrivateError {
	return &PrivateError{message: msg.Error()}
}

//implementing built-in interface error
func (err *PublicError) Error() string {
	return err.message
}

//implementing built-in interface error
func (err *PrivateError) Error() string {
	return err.message
}

// SignalShutdown is used to gracefully Shutdown the app when an integrity
// issue is identified.
func (s *setup) SignalShutdown() {
	s.shutdown <- syscall.SIGTERM
}

func (s *setup) ErrorPresenter(ctx context.Context, e error) *gqlerror.Error {
	var err *gqlerror.Error
	switch e.(type) {
	case *PublicError:
		err = graphql.DefaultErrorPresenter(ctx, e)
	case *PrivateError:
		s.log.Println("private err: ", e)
		err = graphql.DefaultErrorPresenter(ctx, e)
		err.Message = "Internal Server Error"
		err.Path = nil
	case *web.Shutdown:
		s.log.Println("shutdown err: ", e)
		err = graphql.DefaultErrorPresenter(ctx, e)
		err.Message = "Internal Server Error"
		err.Path = nil
		s.SignalShutdown()
	default:
		err = graphql.DefaultErrorPresenter(ctx, e)
	}

	return err
}
