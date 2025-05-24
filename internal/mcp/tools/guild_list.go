package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/samber/lo"
	"soundofgothic.pl/backend/internal/domain"
)

func init() {
	registerTool(NewGuildList)
}

type GuildList struct {
	Repositories domain.Repositories
}

func NewGuildList(repositories domain.Repositories) *GuildList {
	return &GuildList{
		Repositories: repositories,
	}
}

func (c *GuildList) Tool() mcp.Tool {
	return mcp.NewTool(
		"list_guilds",
		mcp.WithDescription("Retrieves a list of game guilds, with optional filters for specific games, associated entities, or search queries."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "list guilds",
			ReadOnlyHint: lo.ToPtr(true),
		}),
		mcp.WithString("query", mcp.Description("Search term to filter guilds by their name or code name."), mcp.DefaultString("")),
		mcp.WithNumber("page", mcp.Description("Page number for pagination of results."), mcp.DefaultNumber(1)),
		mcp.WithNumber("page_size", mcp.Description("Number of guilds to return per page."), mcp.DefaultNumber(50)),
		mcp.WithArray("ids", mcp.Items(map[string]any{
			"type": "number",
		}), mcp.Description("A list of specific Guild IDs to retrieve.")),
		mcp.WithNumber("voice_id", mcp.Description("Filter guilds associated with a specific Voice ID (original voice actor). Set to zero or omit to skip."), mcp.DefaultNumber(0)),
		mcp.WithNumber("script_id", mcp.Description("Filter guilds associated with a specific Script ID. Set to zero or omit to skip."), mcp.DefaultNumber(0)),
		mcp.WithNumber("game_id", mcp.Description("Filter guilds belonging to a specific Game ID. Set to zero or omit to skip."), mcp.DefaultNumber(0)),
	)
}

type GuildListParams struct {
	Query  string  `json:"query"`
	Page   int64   `json:"page"`
	Size   int64   `json:"page_size"`
	Ids    []int64 `json:"ids"`
	Voice  int64   `json:"voice_id"`
	Script int64   `json:"script_id"`
	Game   int64   `json:"game_id"`
}

func (c *GuildList) Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params, err := MapTo[GuildListParams](request.Params.Arguments)
	if err != nil {
		return nil, err
	}

	guilds, total, err := c.Repositories.Guild().List(ctx, domain.GuildSearchOptions{
		Query:     params.Query,
		Page:      params.Page,
		PageSize:  params.Size,
		IDs:       params.Ids,
		VoiceIDs:  ToSlice(params.Voice),
		ScriptIDs: ToSlice(params.Script),
		GameIDs:   ToSlice(params.Game),
	})
	if err != nil {
		return nil, err
	}

	guildsString, err := json.Marshal(map[string]any{
		"total":  total,
		"guilds": guilds,
	})
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(guildsString)), nil
}
