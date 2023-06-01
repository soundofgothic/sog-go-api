package mods

import (
	"reflect"
	"strings"

	"github.com/uptrace/bun"
)

type searchTag struct {
	Columns []string
	Type    string
}

func searchTagFromString(s string) searchTag {
	var columns []string
	var searchType string
	for _, component := range strings.Split(s, ";") {
		if component == "" {
			continue
		}
		componentParts := strings.Split(component, ":")
		if len(componentParts) != 2 {
			panic("invalid search tag")
		}
		label := componentParts[0]
		value := componentParts[1]
		switch label {
		case "columns":
			columns = strings.Split(value, ",")
		case "type":
			searchType = value
		}
	}
	return searchTag{
		Columns: columns,
		Type:    searchType,
	}
}

func WithSearchOptions(options any) QueryModifier {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		reflectType := reflect.TypeOf(options)
		reflectValue := reflect.ValueOf(options)
		if reflectType.Kind() == reflect.Ptr {
			reflectType = reflectType.Elem()
			reflectValue = reflectValue.Elem()
		}
		if reflectType.Kind() != reflect.Struct {
			panic("you need to pass a struct")
		}
		page := 0
		pageSize := 0
		for i := 0; i < reflectType.NumField(); i++ {
			field := reflectType.Field(i)
			value := reflectValue.Field(i)
			searchTagString := field.Tag.Get("search")
			if searchTagString == "" {
				continue
			}
			if value.IsZero() {
				continue
			}
			searchTag := searchTagFromString(searchTagString)
			switch searchTag.Type {
			case "in":
				if len(searchTag.Columns) == 0 {
					panic("in search needs columns")
				}
				if value.Kind() != reflect.Slice {
					panic("in search needs a slice")
				}
				if len(searchTag.Columns) == 1 {
					q.Apply(WithIn(searchTag.Columns[0], value.Interface()))
					continue
				}
				panic("in search with multiple columns not implemented yet")
			case "exact":
				if len(searchTag.Columns) == 0 {
					panic("exact search needs columns")
				}
				if value.Kind() != reflect.String && value.Kind() != reflect.Int64 {
					panic("exact search needs a string or int64")
				}
				if len(searchTag.Columns) == 1 {
					q.Apply(WithExactMatch(searchTag.Columns[0], value.Interface()))
					continue
				}
				panic("exact search with multiple columns not implemented yet")
			case "like":
				if len(searchTag.Columns) == 0 {
					panic("like search needs columns")
				}
				if len(searchTag.Columns) == 1 {
					q.Apply(WithTextSearch(searchTag.Columns[0], value.String()))
					continue
				}
				matches := make(map[string]any)
				for _, column := range searchTag.Columns {
					matches[column] = value.Interface()
				}
				q.Apply(WithOrLikes(matches))
			case "page":
				page = int(value.Int())
			case "pageSize":
				pageSize = int(value.Int())
			}
		}
		if page != 0 && pageSize != 0 {
			q.Apply(WithPaging(int64(page), int64(pageSize)))
		}
		return q
	}
}
