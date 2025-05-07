package rjson

import (
	"encoding/json"
	"net/http"

	"github.com/enhanced-tools/errors"
	"github.com/enhanced-tools/errors/opts"
)

func HandleError(w http.ResponseWriter, err error) {
	enErr := errors.Enhance(err)
	errors.Enhance(err).Log("rest")
	var statusCode opts.StatusCode = 500
	enErr.GetOpt(&statusCode)
	jsonData := errors.AsJSON(errors.Enhance(err))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(statusCode))
	if err := json.NewEncoder(w).Encode(jsonData); err != nil {
		errors.Enhance(err).Log("rest")
		return
	}
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
		HandleError(w, errors.Enhance(err))
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
