package domain

import (
	"context"

	"github.com/uptrace/bun"
)

type Voice struct {
	bun.BaseModel `bun:"voices,alias:v"`

	ID   int64  `json:"id" bun:"id,pk"`
	Name string `json:"name" bun:"name,notnull"`
}

type VoiceService interface {
	List(ctx context.Context) ([]Voice, error)
}
