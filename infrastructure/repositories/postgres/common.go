package postgres

import (
	"context"

	"github.com/enhanced-tools/errors"
	"github.com/uptrace/bun"
	"soundofgothic.pl/backend/infrastructure/repositories/postgres/mods"
)

type Model any

type commonRepository[T Model] struct {
	db *bun.DB
}

func (c *commonRepository[T]) FindByID(ctx context.Context, id int64) (T, error) {
	panic("not implemented")
}

func (c *commonRepository[T]) List(ctx context.Context, modifiers ...mods.QueryModifier) ([]T, int64, error) {
	var model T
	results := make([]T, 0)
	baseQuery := c.db.NewSelect().Model(&model)
	for _, modifier := range modifiers {
		baseQuery = baseQuery.Apply(modifier)
	}
	count, err := baseQuery.ScanAndCount(ctx, &results)
	return results, int64(count), errors.Enhance(err)
}
