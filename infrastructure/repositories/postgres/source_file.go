package postgres

import (
	"context"

	"github.com/uptrace/bun"
	"soundofgothic.pl/backend/domain"
	"soundofgothic.pl/backend/infrastructure/repositories/postgres/mods"
)

type SourceFileRepository struct {
	commonRepository[domain.SourceFile]
}

func NewSourceFileRepository(db *bun.DB) *SourceFileRepository {
	return &SourceFileRepository{
		commonRepository: commonRepository[domain.SourceFile]{
			db: db,
		},
	}
}

func (g *postgresRepositoryStorage) SourceFile() domain.SourceFileService {
	return NewSourceFileRepository(g.db)
}

func (sc *SourceFileRepository) List(ctx context.Context, opts domain.SourceFileSearchOptions) ([]domain.SourceFile, int64, error) {
	return sc.commonRepository.List(ctx,
		mods.WithRecordingsCount("sfs", "source_file_id"),
		mods.WithSearchOptions(opts),
	)
}
