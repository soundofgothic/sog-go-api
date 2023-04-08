package postgres

import (
	"github.com/uptrace/bun"
	"soundofgothic.pl/backend/domain"
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
	return &GameRepository{}
}
