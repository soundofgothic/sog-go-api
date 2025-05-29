package postgres

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
	"soundofgothic.pl/backend/internal/domain"
	"soundofgothic.pl/backend/internal/postgres/mods"
)

type AlternativeRepository struct {
	commonRepository[domain.Alternative]
}

func NewAlternativeRepository(db *bun.DB) *AlternativeRepository {
	return &AlternativeRepository{
		commonRepository: commonRepository[domain.Alternative]{
			db: db,
		},
	}
}

func (g *postgresRepositoryStorage) Alternative() domain.AlternativeService {
	return NewAlternativeRepository(g.db)
}

func (a *AlternativeRepository) Persist(ctx context.Context, alternative *domain.Alternative) error {
	_, err := a.db.NewInsert().
		Model(alternative).
		On("CONFLICT (mod_id, recording_id) DO UPDATE").
		Set("transcript = EXCLUDED.transcript").
		Set("state = EXCLUDED.state").
		Set("wave = EXCLUDED.wave").
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to persist alternative: %w", err)
	}
	return nil
}

func (a *AlternativeRepository) List(ctx context.Context, query domain.AlternativeOptions) ([]domain.Alternative, int64, error) {
	queryMods := []mods.QueryModifier{
		mods.WithRelations("Recording.Game", "Recording.NPC", "Recording.Guild", "Recording.Voice", "Recording.SourceFile"),
		mods.WithSearchOptions(query),
		mods.WithOrderByIDsIn("r.id", query.IDs),
	}

	queryMods = append(queryMods, recordingOptsToMods(query.RecordingSearchOptions)...)

	return a.commonRepository.List(ctx, queryMods...)
}

func (a *AlternativeRepository) Get(ctx context.Context, modID, recordingID int64) (*domain.Alternative, error) {
	return a.commonRepository.Get(ctx,
		mods.WithRelations("R.Game", "R.NPC", "R.Guild", "R.Voice", "R.SourceFile"),
		mods.WithExactMatch("mod_id", modID),
		mods.WithExactMatch("recording_id", recordingID),
	)
}

func (a *AlternativeRepository) Delete(ctx context.Context, modID, recordingID int64) error {
	_, err := a.db.NewDelete().Model((*domain.Alternative)(nil)).Where("mod_id = ? AND recording_id = ?", modID, recordingID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete alternative: %w", err)
	}
	return nil
}
