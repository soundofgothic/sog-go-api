package postgres

import (
	"context"

	"github.com/uptrace/bun"
	"soundofgothic.pl/backend/internal/domain"
	"soundofgothic.pl/backend/internal/postgres/mods"
)

type SourceFileRepository struct {
	commonRepository[domain.SourceFile]
}

func NewSourceFileRepository(db *bun.DB) *SourceFileRepository {
	return &SourceFileRepository{
		commonRepository: commonRepository[domain.SourceFile]{
			db: db,
		},
	}
}

func (g *postgresRepositoryStorage) SourceFile() domain.SourceFileService {
	return NewSourceFileRepository(g.db)
}

func (sc *SourceFileRepository) List(ctx context.Context, opts domain.SourceFileSearchOptions) ([]domain.SourceFile, int64, error) {
	return sc.commonRepository.List(ctx,
		mods.WithWhereGroup(" AND ",
			mods.WithRecordingsCount("sfs", "source_file_id",
				mods.NewMergedRestrictions(
					mods.NewFieldRestriction(opts.GuildIDs, "guild_id"),
					mods.NewFieldRestriction(opts.NPCIDs, "npc_id"),
					mods.NewFieldRestriction(opts.VoiceIDs, "voice_id"),
					mods.NewFieldRestriction(opts.GameIDs, "game_id"),
				)),
			mods.WithSearchOptions(opts),
		),
		mods.WithWhereGroup(" OR ", mods.WithIn("sfs.id", opts.IDs)),
		mods.WithOrderByIDsAndCount("sfs.id", opts.IDs),
	)
}
