package postgres

import (
	"context"

	"github.com/enhanced-tools/errors"
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

	return errors.Enhance(err)
}

func (m *ModRepository) List(ctx context.Context) ([]domain.Mod, error) {
	mods, _, err := m.commonRepository.List(ctx)
	return mods, err
}

func (m *ModRepository) Get(ctx context.Context, id int64) (*domain.Mod, error) {
	return m.commonRepository.FindByID(ctx, id)
}

func (m *ModRepository) Delete(ctx context.Context, id int64) error {
	return m.commonRepository.Delete(ctx, id)
}
