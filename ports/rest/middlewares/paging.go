package middlewares

import (
	"context"
	"net/http"
	"strconv"

	"github.com/enhanced-tools/errors"
	"github.com/enhanced-tools/errors/opts"
	"soundofgothic.pl/backend/ports/rest/rjson"
)

var (
	ErrPageParamError = errors.Template().With(
		opts.Title("Page param error"),
		opts.StatusCode(400),
	)
	ErrPageOptionsMissing = errors.Template().With(
		opts.Title("Page options missing"),
		opts.StatusCode(500),
	)
)

type PageOptions struct {
	PageSize int64
	Page     int64
}

type PageCtxKey string

const PageOptionsKey PageCtxKey = "pageOptions"

func Paging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			err      error
			pageSize int64 = 50
			page     int64 = 1
		)
		if pageSizeStr := r.URL.Query().Get("pageSize"); pageSizeStr != "" {
			pageSize, err = strconv.ParseInt(pageSizeStr, 10, 64)
			if err != nil {
				rjson.HandleError(w, ErrPageParamError.FromEmpty())
				return
			}
		}
		if pageStr := r.URL.Query().Get("page"); pageStr != "" {
			page, err = strconv.ParseInt(pageStr, 10, 64)
			if err != nil {
				rjson.HandleError(w, ErrPageParamError.FromEmpty())
				return
			}
		}
		pageOptions := rjson.PageOptions{
			PageSize: pageSize,
			Page:     page,
		}
		r = r.WithContext(context.WithValue(r.Context(), PageOptionsKey, pageOptions))
		handler.ServeHTTP(w, r)
	})
}

func GetPageOptions(ctx context.Context) (*rjson.PageOptions, error) {
	pageOptions, ok := ctx.Value(PageOptionsKey).(rjson.PageOptions)
	if !ok {
		return nil, ErrPageOptionsMissing.FromEmpty()
	}
	return &pageOptions, nil
}
