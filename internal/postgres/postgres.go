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
)

type postgresRepositoryStorage struct {
	db *bun.DB
}

type options struct {
	auth DBAuth
}

type DBAuth struct {
	Host     string
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

func NewPostgresRepositories(opts ...Option) (*postgresRepositoryStorage, error) {
	options := options{
		auth: DBAuth{
			Host:     "localhost",
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
		pgdriver.WithAddr(fmt.Sprintf("%s:%d", options.auth.Host, options.auth.Port)),
		pgdriver.WithUser(options.auth.Username),
		pgdriver.WithPassword(options.auth.Password),
		pgdriver.WithDatabase(options.auth.Name),
		pgdriver.WithInsecure(true),
	)
	db := bun.NewDB(sql.OpenDB(conn), pgdialect.New())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
	PINGLOOP:
		for err != nil {
			select {
			case <-ctx.Done():
				return nil, errors.Enhance(err)
			default:
				time.Sleep(1 * time.Second)
				if err = db.PingContext(ctx); err == nil {
					break PINGLOOP
				}
			}
		}
	}

	return &postgresRepositoryStorage{
		db: db,
	}, nil
}

func (r *postgresRepositoryStorage) DB() *bun.DB {
	return r.db
}
