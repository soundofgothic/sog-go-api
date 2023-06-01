package mods

import (
	"fmt"
	"strings"

	"github.com/uptrace/bun"
)

type QueryModifier func(*bun.SelectQuery) *bun.SelectQuery

func WithPaging(page int64, pageSize int64) QueryModifier {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Offset(int((page - 1) * pageSize)).Limit(int(pageSize))
	}
}

func WithTextSearch(column string, query string) QueryModifier {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where(fmt.Sprintf("LOWER(%s) LIKE ?", column), fmt.Sprintf("%%%s%%", strings.ToLower(query)))
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

func WithExactMatch(column string, value any) QueryModifier {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where(fmt.Sprintf("%s = ?", column), value)
	}
}

func WithIn(column string, values any) QueryModifier {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where(fmt.Sprintf("%s IN (?)", column), bun.In(values))
	}
}

func WithExactMatches(matches map[string]any) QueryModifier {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		for column, value := range matches {
			q.Where(fmt.Sprintf("%s = ?", column), value)
		}
		return q
	}
}

func WithOrLikes(matches map[string]any) QueryModifier {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.WhereGroup(" AND ", func(q *bun.SelectQuery) *bun.SelectQuery {
			first := true
			for column, value := range matches {
				if first {
					q.Where(fmt.Sprintf("%s LIKE ?", column), fmt.Sprintf("%%%s%%", value))
					first = false
				} else {
					q.WhereOr(fmt.Sprintf("%s LIKE ?", column), fmt.Sprintf("%%%s%%", value))
				}
			}
			return q
		})
	}
}
