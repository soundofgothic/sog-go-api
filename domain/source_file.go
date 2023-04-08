package domain

import "github.com/uptrace/bun"

type SourceFile struct {
	bun.BaseModel `bun:"source_files,alias:sfs"`

	ID     int64  `json:"id" bun:"id,pk"`
	Name   string `json:"name" bun:"name,notnull"`
	Type   string `json:"type" bun:"type,notnull"`
	GameID int64  `json:"game_id" bun:"game_id,notnull"`
	Game   *Game  `json:"game,omitempty" bun:"rel:belongs-to"`
}

type SourceFileService interface{}
