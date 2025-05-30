package postgres

import (
	"context"

	"github.com/uptrace/bun"
	"soundofgothic.pl/backend/internal/domain"
)

type GameRepository struct {
	commonRepository[domain.Game]
}

func NewGameRepository(db *bun.DB) *GameRepository {
	return &GameRepository{
		commonRepository: commonRepository[domain.Game]{
			db: db,
		},
	}
}

func (g *postgresRepositoryStorage) Game() domain.GameService {
	return NewGameRepository(g.db)
}

func (gc *GameRepository) List(ctx context.Context) ([]domain.Game, error) {
	games, _, err := gc.commonRepository.List(ctx)
	return games, err
}
