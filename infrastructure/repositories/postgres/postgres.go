package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/enhanced-tools/errors"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/uptrace/bun/dialect/pgdialect"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/extra/bundebug"
	"soundofgothic.pl/backend/domain"
)

type postgresRepositoryStorage struct {
	db *bun.DB
}

type options struct {
	auth DBAuth
}

type DBAuth struct {
	Address  string
	Port     int
	Name     string
	Username string
	Password string
}

type Option func(*options)

func WithAuth(auth DBAuth) Option {
	return func(o *options) {
		o.auth = auth
	}
}

func NewPostgresRepositories(opts ...Option) (domain.Repositories, error) {
	options := options{
		auth: DBAuth{
			Address:  "localhost",
			Port:     5432,
			Name:     "postgres",
			Username: "postgres",
			Password: "postgres",
		},
	}
	for _, opt := range opts {
		opt(&options)
	}
	conn := pgdriver.NewConnector(
		pgdriver.WithAddr(fmt.Sprintf("%s:%d", options.auth.Address, options.auth.Port)),
		pgdriver.WithUser(options.auth.Username),
		pgdriver.WithPassword(options.auth.Password),
		pgdriver.WithDatabase(options.auth.Name),
		pgdriver.WithInsecure(true),
	)
	db := bun.NewDB(sql.OpenDB(conn), pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Enhance(err)
	}
	return &postgresRepositoryStorage{
		db: db,
	}, nil
}
