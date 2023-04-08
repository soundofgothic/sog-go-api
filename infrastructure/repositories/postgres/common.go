package postgres

import (
	"context"
	"fmt"

	"github.com/enhanced-tools/errors"
	"github.com/uptrace/bun"
)

type Model any

type commonRepository[T Model] struct {
	db *bun.DB
}

func (c *commonRepository[T]) FindByID(ctx context.Context, id int64) (T, error) {
	panic("not implemented")
}

type QueryModifier func(*bun.SelectQuery) *bun.SelectQuery

func WithPaging(page int64, pageSize int64) QueryModifier {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Offset(int((page - 1) * pageSize)).Limit(int(pageSize))
	}
}

func WithTextSearch(column string, query string) QueryModifier {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where(fmt.Sprintf("%s LIKE ?", column), fmt.Sprintf("%%%s%%", query))
	}
}

func WithRelations(relations ...string) QueryModifier {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		for _, relation := range relations {
			q.Relation(relation)
		}
		return q
	}
}

func (c *commonRepository[T]) List(ctx context.Context, modifiers ...QueryModifier) ([]T, int64, error) {
	var model T
	results := make([]T, 0)
	baseQuery := c.db.NewSelect().Model(&model)
	for _, modifier := range modifiers {
		baseQuery = baseQuery.Apply(modifier)
	}
	count, err := baseQuery.ScanAndCount(ctx, &results)
	return results, int64(count), errors.Enhance(err)
}
