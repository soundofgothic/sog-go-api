package domain

import (
	"context"

	"github.com/uptrace/bun"
)

type Game struct {
	bun.BaseModel `bun:"games,alias:g"`

	ID   int64  `json:"id" bun:"id,pk"`
	Name string `json:"name" bun:"name,notnull"`
}

type GameService interface {
	List(ctx context.Context) ([]Game, error)
}
