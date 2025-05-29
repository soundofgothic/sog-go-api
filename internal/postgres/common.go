package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/uptrace/bun"
	"soundofgothic.pl/backend/internal/postgres/mods"
)

type Model any

type commonRepository[T Model] struct {
	db *bun.DB
}

func (c *commonRepository[T]) FindByID(ctx context.Context, id int64) (*T, error) {
	return c.Get(ctx, mods.WithExactMatch("id", id))
}

func (c *commonRepository[T]) List(ctx context.Context, modifiers ...mods.QueryModifier) ([]T, int64, error) {
	var model T
	results := make([]T, 0)
	baseQuery := c.db.NewSelect().Model(&model)
	for _, modifier := range modifiers {
		baseQuery = baseQuery.Apply(modifier)
	}
	count, err := baseQuery.ScanAndCount(ctx, &results)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list records: %w", err)
	}
	return results, int64(count), nil
}

var (
	ErrMultipleResults = errors.New("multiple results found")
	ErrResultNotFound  = errors.New("result not found")
)

func (c *commonRepository[T]) Get(ctx context.Context, modifiers ...mods.QueryModifier) (*T, error) {
	results, _, err := c.List(ctx, modifiers...)
	if err != nil {
		return nil, err
	}
	if len(results) > 1 {
		return nil, ErrMultipleResults
	}
	if len(results) == 0 {
		return nil, ErrResultNotFound
	}
	return &results[0], nil
}

func (c *commonRepository[T]) Delete(ctx context.Context, id int64) error {
	_, err := c.db.NewDelete().Model((*T)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}
	return nil
}
