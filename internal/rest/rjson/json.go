package rjson

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/enhanced-tools/errors"
)

func InternalError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"internal_error": "An internal error occurred. Please try again later.",
	})

	slog.Error("Internal server error", "error", err)
}

func RespondWithJSON(w http.ResponseWriter, data any, statusCode ...int) {
	w.Header().Set("Content-Type", "application/json")
	if len(statusCode) > 0 {
		w.WriteHeader(statusCode[0])
		return
	}
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		InternalError(w, errors.Enhance(err))
		return
	}
}

type PageOptions struct {
	PageSize int64
	Page     int64
}

type PagedResponse[T any] struct {
	Total    int64 `json:"total"`
	Page     int64 `json:"page"`
	PageSize int64 `json:"pageSize"`
	Results  []T   `json:"results"`
}

func NewPagedResponse[T any](page PageOptions, total int64, results []T) PagedResponse[T] {
	return PagedResponse[T]{
		Total:    total,
		Page:     page.Page,
		PageSize: page.PageSize,
		Results:  results,
	}
}
