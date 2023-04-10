package mods

import (
	"fmt"

	"github.com/uptrace/bun"
)

func WithRecordingsCount(table, column string) QueryModifier {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Join(fmt.Sprintf("LEFT JOIN (SELECT r.%s, count(1) cnt FROM recordings r GROUP BY r.%s) AS rc ON %s.id = rc.%s", column, column, table, column)).
			Where("rc.cnt > 0"). // TODO: fix database and remove this line
			Order("rc.cnt DESC").
			ColumnExpr(fmt.Sprintf("%s.*, rc.cnt as count", table))
	}
}
