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
	GameID       int64       `json:"gameID" bun:"game_id,notnull"`
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

type RecordingSearchOptions struct {
	Query        string `search:"type:like;columns:transcript;"`
	Page         int64  `search:"type:page;"`
	PageSize     int64  `search:"type:pageSize;"`
	GameID       int64  `search:"type:exact;columns:game_id;"`
	SourceFileID int64  `search:"type:exact;columns:source_file_id;"`
	NPCID        int64  `search:"type:exact;columns:npc_id;"`
	GuildID      int64  `search:"type:exact;columns:guild_id;"`
	VoiceID      int64  `search:"type:exact;columns:r.voice_id;"`
}

type RecordingService interface {
	List(ctx context.Context, query RecordingSearchOptions) ([]Recording, int64, error)
}
