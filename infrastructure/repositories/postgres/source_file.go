package postgres

import (
	"github.com/uptrace/bun"
	"soundofgothic.pl/backend/domain"
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
