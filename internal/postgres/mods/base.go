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

func WithIn(column string, values []int64) QueryModifier {
	if len(values) == 0 {
		return func(q *bun.SelectQuery) *bun.SelectQuery {
			return q
		}
	}

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

func WithOrderByIDsIn(idColumnName string, ids []int64) QueryModifier {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		if len(ids) == 0 {
			return q
		}

		return q.OrderExpr(fmt.Sprintf("%s IN (?) DESC", idColumnName), bun.In(ids))
	}
}

func WithOrderByIDsAndCount(idColumnName string, ids []int64) QueryModifier {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		if len(ids) == 0 {
			return WithOrderByCount()(q)
		}

		return q.OrderExpr(fmt.Sprintf("%s IN (?) DESC, rc.cnt DESC", idColumnName), bun.In(ids))
	}
}

func WithOrderByCount() QueryModifier {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Order("rc.cnt DESC")
	}
}

func WithWhereGroup(concatenator string, modifiers ...QueryModifier) QueryModifier {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.WhereGroup(concatenator, func(q *bun.SelectQuery) *bun.SelectQuery {
			for _, modifier := range modifiers {
				modifier(q)
			}
			return q
		})
	}
}
