package postgres

import (
	"context"

	"github.com/uptrace/bun"
	"soundofgothic.pl/backend/domain"
)

type VoiceRepository struct {
	commonRepository[domain.Voice]
}

func NewVoiceRepository(db *bun.DB) *VoiceRepository {
	return &VoiceRepository{
		commonRepository: commonRepository[domain.Voice]{
			db: db,
		},
	}
}

func (g *postgresRepositoryStorage) Voice() domain.VoiceService {
	return NewVoiceRepository(g.db)
}

func (vc *VoiceRepository) List(ctx context.Context) ([]domain.Voice, error) {
	result, _, err := vc.commonRepository.List(ctx)
	return result, err
}
