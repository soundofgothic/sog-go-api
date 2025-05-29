package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"

	"soundofgothic.pl/backend/internal/config"
	"soundofgothic.pl/backend/internal/postgres"
	"soundofgothic.pl/backend/internal/rest"
)

func run() int {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to config file")
	flag.StringVar(&configPath, "c", "", "Path to config file (shorthand)")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cfg, err := config.LoadConfigFromFile(configPath)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		return 1
	}

	r := chi.NewRouter()
	repositories, err := postgres.NewPostgresRepositories(postgres.WithAuth(
		postgres.DBAuth{
			Host:     cfg.Postgres.Host,
			Port:     cfg.Postgres.Port,
			Username: cfg.Postgres.User,
			Password: cfg.Postgres.Password,
			Name:     cfg.Postgres.Database,
		},
	))
	if err != nil {
		slog.Error("failed to connect to Postgres", "error", err)
		return 1
	}
	rest.RegisterBackendEndpoints(r, repositories)
	log.Printf("Listening on %s", cfg.Address)
	if err := http.ListenAndServe(cfg.Address, r); err != nil {
		slog.Error("failed to start server", "address", cfg.Address, "error", err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run())
}
