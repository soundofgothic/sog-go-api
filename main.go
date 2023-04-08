package main

import (
	"log"
	"net/http"
	"os"

	"github.com/enhanced-tools/errors"
	"github.com/go-chi/chi/v5"
	"soundofgothic.pl/backend/infrastructure/repositories/postgres"
	"soundofgothic.pl/backend/ports/rest"
)

func run() int {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cfg, err := LoadConfig()
	if err != nil {
		errors.Enhance(err).Log()
		return 1
	}
	errors.Manager().SetDefaultLogger(errors.CustomLogger(
		errors.WithErrorFormatter(errors.MultilineFormatter),
		errors.WithSaveStack(true),
		errors.WithStackTraceFormatter(errors.MultilineStackTraceFormatter),
	))
	r := chi.NewRouter()
	repositories, err := postgres.NewPostgresRepositories(postgres.WithAuth(
		postgres.DBAuth{
			Address:  cfg.Postgres.Host,
			Port:     cfg.Postgres.Port,
			Username: cfg.Postgres.User,
			Password: cfg.Postgres.Password,
			Name:     cfg.Postgres.Database,
		},
	))
	if err != nil {
		errors.Enhance(err).Log()
		return 1
	}
	rest.RegisterEndpoints(r, repositories)
	log.Printf("Listening on %s", cfg.Address)
	if err := http.ListenAndServe(cfg.Address, r); err != nil {
		errors.Enhance(err).Log()
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}
