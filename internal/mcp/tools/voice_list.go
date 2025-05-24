package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/samber/lo"
	"soundofgothic.pl/backend/internal/domain"
)

func init() {
	registerTool(NewVoiceList)
}

type VoiceList struct {
	Repositories domain.Repositories
}

func NewVoiceList(repositories domain.Repositories) *VoiceList {
	return &VoiceList{
		Repositories: repositories,
	}
}

func (c *VoiceList) Tool() mcp.Tool {
	return mcp.NewTool(
		"voice_list",
		mcp.WithDescription("Retrieves a list of unique original voice actors/talents found in game recordings, with options to filter by game, associated NPCs, guilds, scripts, or a search query."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "list voices",
			ReadOnlyHint: lo.ToPtr(true),
		}),
		mcp.WithString("query", mcp.Description("Search term to filter voice actors by their known name or character names they voiced."), mcp.DefaultString("")),
		mcp.WithArray("ids", mcp.Items(map[string]any{
			"type": "number",
		}), mcp.Description("Filter by a list of specific Voice IDs (original voice actor IDs).")),
		mcp.WithArray("guild_ids", mcp.Items(map[string]any{
			"type": "number",
		}), mcp.Description("Filter voice actors who voiced NPCs belonging to a list of specific Guild IDs.")),
		mcp.WithArray("npc_ids", mcp.Items(map[string]any{
			"type": "number",
		}), mcp.Description("Filter voice actors who voiced NPCs in a list of specific NPC IDs.")),
		mcp.WithArray("script_ids", mcp.Items(map[string]any{
			"type": "number",
		}), mcp.Description("Filter voice actors whose voices appear in a list of specific Script IDs (Source File IDs).")),
		mcp.WithArray("game_ids", mcp.Items(map[string]any{
			"type": "number",
		}), mcp.Description("Filter voice actors who have voiced characters in a list of specific Game IDs.")),
	)
}

type VoiceListParams struct {
	Query     string  `json:"query"`
	IDs       []int64 `json:"ids"`
	GuildIDs  []int64 `json:"guild_ids"`
	NPCIDs    []int64 `json:"npc_ids"`
	ScriptIDs []int64 `json:"script_ids"`
	GameIDs   []int64 `json:"game_ids"`
}

func (c *VoiceList) Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params, err := MapTo[VoiceListParams](request.Params.Arguments)
	if err != nil {
		return nil, err
	}

	voices, err := c.Repositories.Voice().List(ctx, domain.VoiceOptions{
		Query:     params.Query,
		IDs:       params.IDs,
		GuildIDs:  params.GuildIDs,
		NPCIDs:    params.NPCIDs,
		ScriptIDs: params.ScriptIDs,
		GameIDs:   params.GameIDs,
	})
	if err != nil {
		return nil, err
	}

	resultStr, err := json.Marshal(map[string]any{
		"voices": voices,
	})

	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(resultStr)), nil
}
