package postgres

import (
	"context"

	"github.com/enhanced-tools/errors"
	"github.com/uptrace/bun"
	"soundofgothic.pl/backend/internal/domain"
)

type TTSVoiceRepository struct {
	commonRepository[domain.TTSVoice]
}

func NewTTSVoiceRepository(db *bun.DB) *TTSVoiceRepository {
	return &TTSVoiceRepository{
		commonRepository: commonRepository[domain.TTSVoice]{
			db: db,
		},
	}
}

func (g *postgresRepositoryStorage) TTSVoice() domain.TTSVoiceService {
	return NewTTSVoiceRepository(g.db)
}

func (v *TTSVoiceRepository) Persist(ctx context.Context, voice *domain.TTSVoice) error {
	_, err := v.db.NewInsert().
		Model(voice).
		On("CONFLICT (id) DO UPDATE").
		Set("name = EXCLUDED.name").
		Exec(ctx)

	return errors.Enhance(err)
}

func (v *TTSVoiceRepository) List(ctx context.Context) ([]domain.TTSVoice, error) {
	voices, _, err := v.commonRepository.List(ctx)
	return voices, err
}

func (v *TTSVoiceRepository) Get(ctx context.Context, id int64) (*domain.TTSVoice, error) {
	return v.Get(ctx, id)
}
