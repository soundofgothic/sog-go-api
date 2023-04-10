package rest

import (
	"github.com/enhanced-tools/errors"
	"github.com/ggicci/httpin"
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
	r.Use(middleware.Recoverer)
	r.Use(middlewares.CORS)
	env := &restEnvironment{
		repositories: repositories,
	}
	r.With(httpin.NewInput(RecordingListInput{})).Get("/recordings", env.recordingsList)
	r.With(httpin.NewInput(NPCListInput{})).Get("/npcs", env.npcList)
	r.With(httpin.NewInput(GuildListInput{})).Get("/guilds", env.guildsList)
	r.With(httpin.NewInput(SourceFilesListInput{})).Get("/source_files", env.sourcefilesList)
	r.Get("/games", env.gamesList)
	r.Get("/voices", env.voicesList)
}
