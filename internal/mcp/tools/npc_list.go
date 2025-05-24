package tools

import (
	"context"
	"encoding/json"

	"github.com/enhanced-tools/errors"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/samber/lo"
	"soundofgothic.pl/backend/internal/domain"
)

func init() {
	registerTool(NewNpcList)
}

type NpcList struct {
	domain.Repositories
}

func NewNpcList(repositories domain.Repositories) *NpcList {
	return &NpcList{
		Repositories: repositories,
	}
}

func (c *NpcList) Tool() mcp.Tool {
	return mcp.NewTool(
		"list_npcs",
		mcp.WithDescription("Retrieves a list of Non-Player Characters (NPCs), with options to filter by game, guild, associated voice/script, or search query."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "list npcs",
			ReadOnlyHint: lo.ToPtr(true),
		}),
		mcp.WithString("query", mcp.Description("Search term to filter NPCs by their name or internal code name."), mcp.DefaultString("")),
		mcp.WithNumber("page", mcp.Description("Page number for pagination of results."), mcp.DefaultNumber(1)),
		mcp.WithNumber("page_size", mcp.Description("Number of NPCs to return per page."), mcp.DefaultNumber(50)),
		mcp.WithNumber("guild_id", mcp.Description("Filter NPCs belonging to a specific Guild ID. Set to zero or omit to skip."), mcp.DefaultNumber(0)),
		mcp.WithNumber("voice_id", mcp.Description("Filter NPCs voiced by a specific Voice ID (original voice actor). Set to zero or omit to skip."), mcp.DefaultNumber(0)),
		mcp.WithNumber("script_id", mcp.Description("Filter NPCs associated with a specific Script ID. Set to zero or omit to skip."), mcp.DefaultNumber(0)),
		mcp.WithArray("ids", mcp.Items(map[string]any{
			"type": "number",
		}), mcp.Description("A list of specific NPC IDs to retrieve.")),
	)
}

type NpcListParams struct {
	Query  string  `json:"query"`
	Page   int64   `json:"page"`
	Size   int64   `json:"page_size"`
	Guild  int64   `json:"guild_id"`
	Voice  int64   `json:"voice_id"`
	Script int64   `json:"script_id"`
	Ids    []int64 `json:"ids"`
}

func (c *NpcList) Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params, err := MapTo[NpcListParams](request.Params.Arguments)
	if err != nil {
		return nil, errors.Enhance(err)
	}

	npcs, total, err := c.Repositories.NPC().List(ctx, domain.NPCSearchOptions{
		Query:     params.Query,
		Page:      params.Page,
		PageSize:  params.Size,
		GuildID:   ToSlice(params.Guild),
		VoiceID:   ToSlice(params.Voice),
		ScriptIDs: ToSlice(params.Script),
		IDs:       params.Ids,
	})
	if err != nil {
		return nil, errors.Enhance(err)
	}

	npcsString, err := json.Marshal(map[string]any{
		"total": total,
		"npcs":  npcs,
	})
	if err != nil {
		return nil, errors.Enhance(err)
	}

	return mcp.NewToolResultText(string(npcsString)), nil
}
