package domain

import (
	"context"

	"github.com/uptrace/bun"
)

type Voice struct {
	bun.BaseModel `bun:"voices,alias:v"`

	ID    int64  `json:"id" bun:"id,pk"`
	Name  string `json:"name" bun:"name,notnull"`
	Count int64  `json:"count" bun:",scanonly"`
}

type VoiceOptions struct {
	GuildIDs  []int64
	NPCIDs    []int64
	ScriptIDs []int64
	GameIDs   []int64
}

type VoiceService interface {
	List(ctx context.Context, opts VoiceOptions) ([]Voice, error)
}
