package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/samber/lo"
	"soundofgothic.pl/backend/internal/domain"
)

func init() {
	registerTool(NewModCreate)
}

type ModCreate struct {
	domain.Repositories
}

func NewModCreate(repositories domain.Repositories) *ModCreate {
	return &ModCreate{
		Repositories: repositories,
	}
}

func (c *ModCreate) Tool() mcp.Tool {
	return mcp.NewTool(
		"create_or_update_mod",
		mcp.WithDescription("Creates a new voice mod or updates an existing one for a Gothic game."),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("The name of the mod. This field is required."),
		),
		mcp.WithString("description", mcp.Description("A detailed description of the mod. This field is required."), mcp.Required()),
		mcp.WithString("version", mcp.Description("The version number of the mod (e.g., \"1.0.0\"). This field is required."), mcp.Required()),
		mcp.WithNumber("id", mcp.Description("The ID of the mod to update. If omitted or set to zero, a new mod will be created.")),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:          "create or update mod",
			IdempotentHint: lo.ToPtr(true),
			OpenWorldHint:  lo.ToPtr(true),
		}),
	)
}

type createModParams struct {
	ID          int64  `json:"id"`
	Version     string `json:"version" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (c *ModCreate) Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params, err := MapTo[createModParams](request.Params.Arguments)
	if err != nil {
		return nil, err
	}

	mod := domain.Mod{
		ID:          params.ID,
		Name:        params.Name,
		Description: params.Description,
		Version:     params.Version,
		GameID:      2,
	}

	err = c.Repositories.Mod().Persist(ctx, &mod)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("mod id: %d", mod.ID)), nil
}
