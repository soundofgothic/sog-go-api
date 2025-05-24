package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/samber/lo" // Added import for lo
	"soundofgothic.pl/backend/internal/domain"
)

func init() {
	registerTool(NewModList)
}

type ModList struct {
	domain.Repositories
}

func NewModList(repositories domain.Repositories) *ModList {
	return &ModList{
		Repositories: repositories,
	}
}

func (c *ModList) Tool() mcp.Tool {
	return mcp.NewTool(
		"list_mods",
		mcp.WithDescription("Retrieves a list of existing dialogue modification projects."),
		// TODO: Add parameters for filtering (e.g., by game_id, query, status) and pagination (page, page_size)
		// and update the Handler and add a ModListParams struct accordingly.
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "list mods",
			ReadOnlyHint: lo.ToPtr(true),
		}),
	)
}

func (c *ModList) Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	mods, err := c.Repositories.Mod().List(ctx)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(mods)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(string(data)), nil
}
