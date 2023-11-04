package mods

import (
	"fmt"
	"strings"

	"github.com/uptrace/bun"
)

type CountRestriction struct {
	Column string
	Values []int64
}

type CountRestrictor interface {
	Restrictions() (string, []any)
}

type FieldRestriction struct {
	field string
	ids   []int64
}

func (f FieldRestriction) Restrictions() (string, []any) {
	if len(f.ids) == 0 {
		return "", nil
	}

	return fmt.Sprintf("WHERE %s IN (?)", f.field), []any{bun.In(f.ids)}
}

func NewFieldRestriction(ids []int64, field string) FieldRestriction {
	return FieldRestriction{
		ids:   ids,
		field: field,
	}
}

type MergedRestrictions struct {
	restrictions []CountRestrictor
}

func NewMergedRestrictions(restrictions ...CountRestrictor) MergedRestrictions {
	return MergedRestrictions{
		restrictions: restrictions,
	}
}

func (m MergedRestrictions) Restrictions() (string, []any) {
	if len(m.restrictions) == 0 {
		return "", nil
	}

	restrictionBuilder := strings.Builder{}
	args := make([]any, 0, len(m.restrictions))

	for _, restriction := range m.restrictions {
		rawRestrictionString, restrictionArgs := restriction.Restrictions()
		if rawRestrictionString == "" {
			continue
		}

		if restrictionBuilder.Len() == 0 {
			restrictionBuilder.WriteString("WHERE")
		} else {
			restrictionBuilder.WriteString(" AND")
		}

		restriction, _ := strings.CutPrefix(rawRestrictionString, "WHERE")

		restrictionBuilder.WriteString(fmt.Sprintf(" %s", restriction))
		args = append(args, restrictionArgs...)
	}

	return restrictionBuilder.String(), args
}

func WithRecordingsCount(table, column string, countRestriction ...CountRestrictor) QueryModifier {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		restructionQuery := ""
		args := []any{}

		if len(countRestriction) > 0 {
			restructionQuery, args = countRestriction[0].Restrictions()
		}

		return q.Join(fmt.Sprintf("LEFT JOIN (SELECT r.%s, count(1) cnt FROM recordings r %s GROUP BY r.%s) AS rc ON %s.id = rc.%s", column, restructionQuery, column, table, column), args...).
			Where("rc.cnt > 0").
			Order("rc.cnt DESC").
			ColumnExpr(fmt.Sprintf("%s.*, rc.cnt as count", table))
	}
}
