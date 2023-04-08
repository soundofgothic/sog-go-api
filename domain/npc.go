package domain

import "github.com/uptrace/bun"

type NPC struct {
	bun.BaseModel `bun:"npcs,alias:npcs"`

	ID          int64  `json:"id" bun:"id,pk"`
	Name        string `json:"name" bun:"name,notnull"`
	InGameID    int64  `json:"inGameID" bun:"in_game_id,notnull"`
	InGameAlias string `json:"inGameAlias" bun:"in_game_alias,notnull"`
	GameID      int64  `json:"game_id" bun:"game_id,notnull"`
	Game        *Game  `json:"game,omitempty" bun:"rel:belongs-to"`
	VoiceID     int64  `json:"voice_id" bun:"voice_id,notnull"`
	Voice       *Voice `json:"voice,omitempty" bun:"rel:belongs-to"`
	GuildID     int64  `json:"guild_id" bun:"guild_id,notnull"`
	Guild       *Guild `json:"guild,omitempty" bun:"rel:belongs-to"`
}

type NPCService interface{}
