package postgres

import (
	"context"

	"github.com/uptrace/bun"
	"soundofgothic.pl/backend/internal/domain"
	"soundofgothic.pl/backend/internal/postgres/mods"
)

type VoiceRepository struct {
	commonRepository[domain.Voice]
}

func NewVoiceRepository(db *bun.DB) *VoiceRepository {
	return &VoiceRepository{
		commonRepository: commonRepository[domain.Voice]{
			db: db,
		},
	}
}

func (g *postgresRepositoryStorage) Voice() domain.VoiceService {
	return NewVoiceRepository(g.db)
}

func (vc *VoiceRepository) List(ctx context.Context, opts domain.VoiceOptions) ([]domain.Voice, error) {
	result, _, err := vc.commonRepository.List(ctx,
		mods.WithWhereGroup(" AND ",
			mods.WithRecordingsCount("v", "voice_id",
				mods.NewMergedRestrictions(
					mods.NewFieldRestriction(opts.ScriptIDs, "source_file_id"),
					mods.NewFieldRestriction(opts.NPCIDs, "npc_id"),
					mods.NewFieldRestriction(opts.GameIDs, "game_id"),
					mods.NewFieldRestriction(opts.GuildIDs, "guild_id"),
				),
			),
			mods.WithSearchOptions(opts),
		),
		mods.WithWhereGroup(" OR ", mods.WithIn("v.id", opts.IDs)),
		mods.WithOrderByIDsAndCount("v.id", opts.IDs),
	)
	return result, err
}
