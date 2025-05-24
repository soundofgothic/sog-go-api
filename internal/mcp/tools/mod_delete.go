package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/samber/lo"
	"soundofgothic.pl/backend/internal/domain"
)

func init() {
	registerTool(NewModDelete)
}

type ModDelete struct {
	domain.Repositories
}

func NewModDelete(repositories domain.Repositories) *ModDelete {
	return &ModDelete{
		Repositories: repositories,
	}
}

func (c *ModDelete) Tool() mcp.Tool {
	return mcp.NewTool(
		"delete_mod",
		mcp.WithDescription("Deletes an existing dialogue modification project by its unique ID."),
		mcp.WithNumber("id", mcp.Description("The unique identifier of the mod to be deleted. This field is required."), mcp.Required()),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:          "delete mod",
			IdempotentHint: lo.ToPtr(true),
			OpenWorldHint:  lo.ToPtr(true),
		}),
	)
}

type deleteModParams struct {
	ID int64 `json:"id" validate:"required"`
}

func (c *ModDelete) Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params, err := MapTo[deleteModParams](request.Params.Arguments)
	if err != nil {
		return nil, err
	}

	err = c.Repositories.Mod().Delete(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(fmt.Sprintf("deleted mod id: %d", params.ID)), nil
}
