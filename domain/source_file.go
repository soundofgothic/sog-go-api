package domain

import (
	"context"

	"github.com/uptrace/bun"
)

type SourceFile struct {
	bun.BaseModel `bun:"source_files,alias:sfs"`

	ID     int64  `json:"id" bun:"id,pk"`
	Name   string `json:"name" bun:"name,notnull"`
	Type   string `json:"type" bun:"type,notnull"`
	GameID int64  `json:"gameID" bun:"game_id,notnull"`
	Game   *Game  `json:"game,omitempty" bun:"rel:belongs-to"`

	Count int64 `json:"count" bun:",scanonly"`
}

type SourceFileSearchOptions struct {
	IDs      []int64
	Query    string `search:"type:like;columns:sfs.name;"`
	Page     int64  `search:"type:page;"`
	PageSize int64  `search:"type:pageSize;"`
	Type     string `search:"type:exact;columns:type;"`

	GameIDs  []int64
	GuildIDs []int64
	NPCIDs   []int64
	VoiceIDs []int64
}

type SourceFileService interface {
	List(ctx context.Context, options SourceFileSearchOptions) ([]SourceFile, int64, error)
}
