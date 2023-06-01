package postgres

import (
	"context"

	"github.com/uptrace/bun"
	"soundofgothic.pl/backend/domain"
	"soundofgothic.pl/backend/infrastructure/repositories/postgres/mods"
)

type RecordingRepository struct {
	commonRepository[domain.Recording]
}

func NewRecordingRepository(db *bun.DB) *RecordingRepository {
	return &RecordingRepository{
		commonRepository: commonRepository[domain.Recording]{
			db: db,
		},
	}
}

func (g *postgresRepositoryStorage) Recording() domain.RecordingService {
	return NewRecordingRepository(g.db)
}

func (rc *RecordingRepository) List(ctx context.Context, opts domain.RecordingSearchOptions) ([]domain.Recording, int64, error) {
	return rc.commonRepository.List(ctx,
		mods.WithRelations("Game", "NPC", "Guild", "Voice", "SourceFile"),
		mods.WithSearchOptions(opts),
	)
}

func (rc *RecordingRepository) Get(ctx context.Context, opts domain.RecordingGetOptions) (*domain.Recording, error) {
	return rc.commonRepository.Get(ctx,
		mods.WithRelations("Game", "NPC", "Guild", "Voice", "SourceFile"),
		mods.WithSearchOptions(opts),
	)
}
