package domain

import (
	"context"

	"github.com/uptrace/bun"
)

type Recording struct {
	bun.BaseModel `bun:"recordings,alias:r"`

	ID           int64        `json:"id" bun:"id,pk"`
	Wave         string       `json:"wave" bun:"wave,notnull"`
	Transcript   string       `json:"transcript" bun:"transcript,notnull"`
	GameID       int64        `json:"gameID" bun:"game_id,notnull"`
	Game         *Game        `json:"game,omitempty" bun:"rel:belongs-to,join:game_id=id"`
	SourceFileID int64        `json:"sourceFileID" bun:"source_file_id,notnull"`
	SourceFile   *SourceFile  `json:"sourceFile,omitempty" bun:"rel:belongs-to"`
	NPCID        *int64       `json:"npcID" bun:"npc_id"`
	NPC          *NPC         `json:"npc,omitempty" bun:"rel:belongs-to"`
	GuildID      *int64       `json:"guildID" bun:"guild_id"`
	Guild        *Guild       `json:"guild,omitempty" bun:"rel:belongs-to"`
	VoiceID      *int64       `json:"voiceID" bun:"voice_id"`
	Voice        *Voice       `json:"voice,omitempty" bun:"rel:belongs-to"`
	Title        *string      `json:"title,omitempty" bun:"title"`
	Alternative  *Alternative `bun:"rel:has-one,join:id=recording_id"`
}

type RecordingSearchOptions struct {
	IDs          []int64 `search:"type:in;columns:r.id;"`
	Query        string
	Page         int64   `search:"type:page;"`
	PageSize     int64   `search:"type:pageSize;"`
	SourceFileID []int64 `search:"type:in;columns:r.source_file_id;"`
	GameIDs      []int64 `search:"type:in;columns:r.game_id;"`
	NPCIDs       []int64 `search:"type:in;columns:r.npc_id;"`
	GuildIDs     []int64 `search:"type:in;columns:r.guild_id;"`
	VoiceIDs     []int64 `search:"type:in;columns:r.voice_id;"`
}

type RecordingGetOptions struct {
	GameID int64  `search:"type:exact;columns:r.game_id;"`
	Wave   string `search:"type:exact;columns:r.wave;"`
}

type RecordingService interface {
	List(ctx context.Context, query RecordingSearchOptions) ([]Recording, int64, error)
	Get(ctx context.Context, query RecordingGetOptions) (*Recording, error)
}
