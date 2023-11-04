package domain

import (
	"context"

	"github.com/uptrace/bun"
)

type Guild struct {
	bun.BaseModel `bun:"guilds,alias:guilds"`

	ID          int64  `json:"id" bun:"id,pk"`
	Name        string `json:"name" bun:"name,notnull"`
	InGameID    int64  `json:"inGameID" bun:"in_game_id,notnull"`
	InGameAlias string `json:"inGameAlias" bun:"in_game_alias,notnull"`
	GameID      int64  `json:"gameID" bun:"game_id,notnull"`
	Game        *Game  `json:"game,omitempty" bun:"rel:belongs-to"`

	Count int64 `json:"count" bun:",scanonly"`
}

type GuildSearchOptions struct {
	Query       string  `search:"type:like;columns:guilds.name,guilds.in_game_alias;"`
	Page        int64   `search:"type:page;"`
	PageSize    int64   `search:"type:pageSize;"`
	GameIDs     []int64 `search:"type:in;columns:game_id;"`
	InGameAlias string  `search:"type:exact;columns:in_game_alias;"`

	VoiceIDs  []int64
	ScriptIDs []int64
	NPCIDs    []int64
}

type GuildService interface {
	List(ctx context.Context, searchOptions GuildSearchOptions) ([]Guild, int64, error)
}
