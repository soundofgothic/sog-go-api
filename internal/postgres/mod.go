package postgres

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
	"soundofgothic.pl/backend/internal/domain"
)

type ModRepository struct {
	commonRepository[domain.Mod]
}

func NewModRepository(db *bun.DB) *ModRepository {
	return &ModRepository{
		commonRepository: commonRepository[domain.Mod]{
			db: db,
		},
	}
}

func (g *postgresRepositoryStorage) Mod() domain.ModService {
	return NewModRepository(g.db)
}

func (m *ModRepository) Persist(ctx context.Context, mod *domain.Mod) error {
	_, err := m.db.NewInsert().
		Model(mod).
		On("CONFLICT (id) DO UPDATE").
		Set("name = EXCLUDED.name").
		Set("description = EXCLUDED.description").
		Set("version = EXCLUDED.version").
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to persist mod: %w", err)
	}
	return nil
}

func (m *ModRepository) List(ctx context.Context) ([]domain.Mod, error) {
	mods, _, err := m.commonRepository.List(ctx)
	return mods, err
}

func (m *ModRepository) Get(ctx context.Context, id int64) (*domain.Mod, error) {
	return m.FindByID(ctx, id)
}
