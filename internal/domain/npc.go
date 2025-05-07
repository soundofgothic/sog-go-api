package domain

import (
	"context"

	"github.com/uptrace/bun"
)

type NPC struct {
	bun.BaseModel `bun:"npcs,alias:npcs"`

	ID          int64  `json:"id" bun:"id,pk"`
	Name        string `json:"name" bun:"name,notnull"`
	InGameID    int64  `json:"inGameID" bun:"in_game_id,notnull"`
	InGameAlias string `json:"inGameAlias" bun:"in_game_alias,notnull"`
	GameID      int64  `json:"gameID" bun:"game_id,notnull"`
	Game        *Game  `json:"game,omitempty" bun:"rel:belongs-to"`
	VoiceID     int64  `json:"voiceID" bun:"voice_id,notnull"`
	Voice       *Voice `json:"voice,omitempty" bun:"rel:belongs-to"`
	GuildID     int64  `json:"guildID" bun:"guild_id,notnull"`
	Guild       *Guild `json:"guild,omitempty" bun:"rel:belongs-to"`
	Count       int64  `json:"count" bun:",scanonly"`
}

type NPCSearchOptions struct {
	IDs         []int64
	Query       string  `search:"type:like;columns:npcs.name,npcs.in_game_alias;"`
	Page        int64   `search:"type:page;"`
	PageSize    int64   `search:"type:pageSize;"`
	GameID      []int64 `search:"type:in;columns:game_id;"`
	VoiceID     []int64 `search:"type:in;columns:voice_id;"`
	GuildID     []int64 `search:"type:in;columns:guild_id;"`
	InGameAlias string  `search:"type:exact;columns:in_game_alias;"`

	ScriptIDs []int64
}

type NPCService interface {
	List(ctx context.Context, searchOptions NPCSearchOptions) ([]NPC, int64, error)
}
