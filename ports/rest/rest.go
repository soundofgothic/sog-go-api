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
	r.Use(middleware.Recoverer)
	r.Use(middlewares.CORS)
	env := &restEnvironment{
		repositories: repositories,
	}
	r.With(middlewares.ValidatedInput(RecordingListInput{})).Get("/recordings", env.recordingsList)
	r.With(middlewares.ValidatedInput(RecordingGetInput{})).Get("/recordings/{gameID}/{wave}", env.recordingsGet)
	r.With(middlewares.ValidatedInput(NPCListInput{})).Get("/npcs", env.npcList)
	r.With(middlewares.ValidatedInput(GuildListInput{})).Get("/guilds", env.guildsList)
	r.With(middlewares.ValidatedInput(SourceFilesListInput{})).Get("/source_files", env.sourcefilesList)
	r.Get("/games", env.gamesList)
	r.Get("/voices", env.voicesList)
}
