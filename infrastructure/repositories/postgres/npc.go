package postgres

import (
	"context"

	"github.com/uptrace/bun"
	"soundofgothic.pl/backend/domain"
	"soundofgothic.pl/backend/infrastructure/repositories/postgres/mods"
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

func (nc *NPCRepository) List(ctx context.Context, options domain.NPCSearchOptions) ([]domain.NPC, int64, error) {
	return nc.commonRepository.List(ctx,
		mods.WithWhereGroup(" AND ",
			mods.WithRecordingsCount("npcs", "npc_id",
				mods.NewFieldRestriction(options.ScriptIDs, "source_file_id"),
			),
			mods.WithSearchOptions(options),
		),
		mods.WithWhereGroup(" OR ", mods.WithIn("npcs.id", options.IDs)),
		mods.WithOrderByIDsAndCount("npcs.id", options.IDs),
	)
}
