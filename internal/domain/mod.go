package domain

import (
	"context"

	"github.com/uptrace/bun"
)

type Mod struct {
	bun.BaseModel `bun:"mods,alias:m"`

	ID     int64 `json:"id" bun:"id,pk,autoincrement"`
	GameID int64 `json:"gameID" bun:"game_id,notnull"`
	Game   *Game `json:"game,omitempty" bun:"rel:belongs-to,join:game_id=id"`
	// Name is the name of the mod.
	Name string `json:"name" bun:"name,notnull"`
	// Description is the description of the mod.
	Description string `json:"description" bun:"description,notnull"`
	// Version is the version of the mod.
	Version string `json:"version" bun:"version,notnull"`
}

type ModService interface {
	Persist(ctx context.Context, mod *Mod) error
	List(ctx context.Context) ([]Mod, error)
	Get(ctx context.Context, id int64) (*Mod, error)
	Delete(ctx context.Context, id int64) error
}
