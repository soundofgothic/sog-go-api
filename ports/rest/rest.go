package rest

import (
	"github.com/enhanced-tools/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"soundofgothic.pl/backend/domain"
	"soundofgothic.pl/backend/ports/rest/middlewares"
)

type restEnvironment struct {
	repositories domain.Repositories
}

func RegisterEndpoints(r *chi.Mux, repositories domain.Repositories) {
	errors.Manager().RegisterLogger("rest", errors.CustomLogger(
		errors.WithErrorFormatter(errors.MultilineFormatter),
		errors.WithSaveStack(true),
		errors.WithStackTraceFormatter(errors.MultilineStackTraceFormatter),
	))

	r.Use(middleware.Logger)
	env := &restEnvironment{
		repositories: repositories,
	}
	r.With(middlewares.Paging).Get("/recordings", env.recordingsList)
}
