package tools

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"soundofgothic.pl/backend/internal/test"
)

func NewMCPRequest(params map[string]any) mcp.CallToolRequest {
	return mcp.CallToolRequest{
		Params: struct {
			Name      string         `json:"name"`
			Arguments map[string]any `json:"arguments,omitempty"`
			Meta      *mcp.Meta      `json:"_meta,omitempty"`
		}{
			Arguments: params,
		},
	}
}

func AsTextContent(result *mcp.CallToolResult) string {
	if len(result.Content) == 0 {
		return ""
	}
	if len(result.Content) > 1 {
		return ""
	}

	val, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		return ""
	}
	return val.Text
}

func TestModCreate_Handler_Success(t *testing.T) {
	repos := test.NewRepositories()
	modCreate := NewModCreate(repos)

	params := map[string]any{
		"name":        "Test Mod",
		"description": "A test mod",
		"version":     "1.0.0",
	}
	request := NewMCPRequest(params)
	result, err := modCreate.Handler(context.Background(), request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, result.Content[0], "mod id:")
}

func TestModCreate_Handler_MissingRequired(t *testing.T) {
	repos := test.NewRepositories()
	modCreate := NewModCreate(repos)

	params := map[string]any{
		"name": "Test Mod",
		// missing description and version
	}
	request := NewMCPRequest(params)
	result, err := modCreate.Handler(context.Background(), request)
	assert.Error(t, err)
	assert.Nil(t, result)
}
