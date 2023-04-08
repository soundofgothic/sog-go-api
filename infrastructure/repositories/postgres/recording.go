package postgres

import (
	"context"

	"github.com/uptrace/bun"
	"soundofgothic.pl/backend/domain"
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

func (rc *RecordingRepository) SearchByText(ctx context.Context, query string, page int64, pageSize int64) ([]domain.Recording, int64, error) {
	return rc.commonRepository.List(ctx,
		WithTextSearch("transcript", query),
		WithRelations("Game", "NPC", "Guild", "Voice"),
		WithPaging(page, pageSize),
	)
}
