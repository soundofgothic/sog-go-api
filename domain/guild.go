package domain

import "github.com/uptrace/bun"

type Guild struct {
	bun.BaseModel `bun:"guilds,alias:guilds"`

	ID          int64  `json:"id" bun:"id,pk"`
	Name        string `json:"name" bun:"name,notnull"`
	InGameID    int64  `json:"inGameID" bun:"in_game_id,notnull"`
	InGameAlias string `json:"inGameAlias" bun:"in_game_alias,notnull"`
	GameID      int64  `json:"game_id" bun:"game_id,notnull"`
	Game        *Game  `json:"game,omitempty" bun:"rel:belongs-to"`
}

type GuildService interface{}
