package migrations

import (
	"context"
	"embed"

	"github.com/enhanced-tools/errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

//go:embed *.sql
var sqlMigrations embed.FS

type DBProvider interface {
	DB() *bun.DB
}

func Migrate(ctx context.Context, p DBProvider) error {
	db := p.DB()

	migrations := migrate.NewMigrations()
	if err := migrations.Discover(sqlMigrations); err != nil {
		return errors.Enhance(err)
	}

	migrator := migrate.NewMigrator(db, migrations)
	if err := migrator.Init(ctx); err != nil {
		return errors.Enhance(err)
	}

	if _, err := migrator.Migrate(ctx); err != nil {
		return errors.Enhance(err)
	}

	return nil

}
