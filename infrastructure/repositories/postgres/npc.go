package postgres

import (
	"github.com/uptrace/bun"
	"soundofgothic.pl/backend/domain"
)

type NPCRepository struct {
	commonRepository[domain.NPC]
}

func NewNPCRepository(db *bun.DB) *NPCRepository {
	return &NPCRepository{
		commonRepository: commonRepository[domain.NPC]{
			db: db,
		},
	}
}

func (g *postgresRepositoryStorage) NPC() domain.NPCService {
	return NewNPCRepository(g.db)
}
