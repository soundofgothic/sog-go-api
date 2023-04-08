package postgres

import (
	"github.com/uptrace/bun"
	"soundofgothic.pl/backend/domain"
)

type GuildRepository struct {
	commonRepository[domain.Guild]
}

func NewGuildRepository(db *bun.DB) *GuildRepository {
	return &GuildRepository{
		commonRepository: commonRepository[domain.Guild]{
			db: db,
		},
	}
}

func (g *postgresRepositoryStorage) Guild() domain.GuildService {
	return NewGuildRepository(g.db)
}
