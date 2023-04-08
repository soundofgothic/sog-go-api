package domain

import (
	"context"

	"github.com/uptrace/bun"
)

type Recording struct {
	bun.BaseModel `bun:"recordings,alias:r"`

	ID           int64       `json:"id" bun:"id,pk"`
	Wave         string      `json:"wave" bun:"wave,notnull"`
	Transcript   string      `json:"transcript" bun:"transcript,notnull"`
	GameID       int64       `json:"game_id" bun:"game_id,notnull"`
	Game         *Game       `json:"game,omitempty" bun:"rel:belongs-to,join:game_id=id"`
	SourceFileID int64       `json:"sourceFileID" bun:"source_file_id,notnull"`
	SourceFile   *SourceFile `json:"sourceFile,omitempty" bun:"rel:belongs-to"`
	NPCID        *int64      `json:"npcID" bun:"npc_id"`
	NPC          *NPC        `json:"npc,omitempty" bun:"rel:belongs-to"`
	GuildID      *int64      `json:"guildID" bun:"guild_id"`
	Guild        *Guild      `json:"guild,omitempty" bun:"rel:belongs-to"`
	VoiceID      *int64      `json:"voiceID" bun:"voice_id"`
	Voice        *Voice      `json:"voice,omitempty" bun:"rel:belongs-to"`
	Title        *string     `json:"title,omitempty" bun:"title"`
}

type RecordingService interface {
	SearchByText(ctx context.Context, query string, page int64, pageSize int64) ([]Recording, int64, error)
}
