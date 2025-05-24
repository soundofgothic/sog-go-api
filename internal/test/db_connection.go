package test

import (
	"soundofgothic.pl/backend/internal/config"
	"soundofgothic.pl/backend/internal/domain"
	"soundofgothic.pl/backend/internal/postgres"
)

func NewRepositories() domain.Repositories {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
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
		panic(err)
	}

	return repositories
}
