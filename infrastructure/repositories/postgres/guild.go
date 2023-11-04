package postgres

import (
	"context"

	"github.com/uptrace/bun"
	"soundofgothic.pl/backend/domain"
	"soundofgothic.pl/backend/infrastructure/repositories/postgres/mods"
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

func (gc *GuildRepository) List(ctx context.Context, options domain.GuildSearchOptions) ([]domain.Guild, int64, error) {
	return gc.commonRepository.List(ctx,
		mods.WithRecordingsCount("guilds", "guild_id",
			mods.NewMergedRestrictions(
				mods.NewFieldRestriction(options.VoiceIDs, "voice_id"),
				mods.NewFieldRestriction(options.ScriptIDs, "source_file_id"),
				mods.NewFieldRestriction(options.NPCIDs, "npc_id"),
			),
		),
		mods.WithSearchOptions(options),
	)
}
