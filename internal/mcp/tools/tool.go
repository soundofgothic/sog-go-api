package tools

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
	"soundofgothic.pl/backend/internal/domain"
	"soundofgothic.pl/backend/internal/mcp"
)

var gothicTools []func(r domain.Repositories) mcp.Tool

func registerTool[T mcp.Tool](tool func(r domain.Repositories) T) {
	gothicTools = append(gothicTools, func(r domain.Repositories) mcp.Tool {
		return tool(r)
	})
}

func NewGothicToolsPack(r domain.Repositories) []mcp.Tool {
	tools := make([]mcp.Tool, 0, len(gothicTools))
	for _, tool := range gothicTools {
		tools = append(tools, tool(r))
	}
	return tools
}

var (
	validate     *validator.Validate
	validateOnce sync.Once
)

func MapTo[T any](params map[string]any) (T, error) {
	var v T
	val, err := json.Marshal(params)
	if err != nil {
		return v, fmt.Errorf("failed to marshal params: %w", err)
	}
	if err := json.Unmarshal(val, &v); err != nil {
		return v, fmt.Errorf("failed to unmarshal params: %w", err)
	}

	validateOnce.Do(func() {
		validate = validator.New()
	})
	if err := validate.Struct(v); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errMsgs := make([]string, 0, len(validationErrors))
			for _, ve := range validationErrors {
				detail := ve.ActualTag()
				if ve.Param() != "" {
					detail = fmt.Sprintf("%s=%s", ve.ActualTag(), ve.Param())
				}
				errMsgs = append(errMsgs, fmt.Sprintf("%s: %s", ve.Field(), detail))
			}
			return v, fmt.Errorf("%s", strings.Join(errMsgs, ", "))
		}
		return v, fmt.Errorf("validation failed: %w", err)
	}
	return v, nil
}

func ToSlice[T any](v T) []T {
	if reflect.ValueOf(v).IsZero() {
		return nil
	}
	return []T{v}
}
