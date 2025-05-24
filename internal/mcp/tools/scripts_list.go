package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/samber/lo"
	"soundofgothic.pl/backend/internal/domain"
)

func init() {
	registerTool(NewScriptList)
}

type ScriptList struct {
	Repositories domain.Repositories
}

func NewScriptList(repositories domain.Repositories) *ScriptList {
	return &ScriptList{
		Repositories: repositories,
	}
}

func (c *ScriptList) Tool() mcp.Tool {
	return mcp.NewTool(
		"list_scripts",
		mcp.WithDescription("Retrieves a list of game scripts (dialogue files/sets), with options to filter by game, associated NPCs, guilds, voice actors, script type, or a search query."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "list scripts",
			ReadOnlyHint: lo.ToPtr(true),
		}),
		mcp.WithString("query", mcp.Description("Search term to filter scripts by their filename, path, or content snippets within the dialogue set."), mcp.DefaultString("")),
		mcp.WithNumber("page", mcp.Description("Page number for pagination of results."), mcp.DefaultNumber(1)),
		mcp.WithNumber("page_size", mcp.Description("Number of scripts to return per page."), mcp.DefaultNumber(50)),
		mcp.WithArray("ids", mcp.Items(map[string]any{
			"type": "number",
		}), mcp.Description("Filter by a list of specific Script IDs (Source File IDs).")),
		mcp.WithString("type", mcp.Description("Filter scripts by their type (e.g., 'svm' for NPC sound snippets, 'guild' for guild-specific dialogues, 'mission' for quest-related dialogues). Refer to game documentation for specific type values."), mcp.DefaultString("")),
		mcp.WithArray("game_ids", mcp.Items(map[string]any{
			"type": "number",
		}), mcp.Description("Filter scripts belonging to a list of specific Game IDs.")),
		mcp.WithArray("guild_ids", mcp.Items(map[string]any{
			"type": "number",
		}), mcp.Description("Filter scripts primarily involving NPCs from a list of specific Guild IDs.")),
		mcp.WithArray("npc_ids", mcp.Items(map[string]any{
			"type": "number",
		}), mcp.Description("Filter scripts containing dialogues spoken by NPCs in a list of specific NPC IDs.")),
		mcp.WithArray("voice_ids", mcp.Items(map[string]any{
			"type": "number",
		}), mcp.Description("Filter scripts containing dialogues spoken by voice actors in a list of specific Voice IDs.")),
	)
}

type ScriptListParams struct {
	Query    string  `json:"query"`
	Page     int64   `json:"page"`
	Size     int64   `json:"page_size"`
	Ids      []int64 `json:"ids"`
	Type     string  `json:"type"`
	GameIDs  []int64 `json:"game_ids"`
	GuildIDs []int64 `json:"guild_ids"`
	NPCIDs   []int64 `json:"npc_ids"`
	VoiceIDs []int64 `json:"voice_ids"`
}

func (c *ScriptList) Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params, err := MapTo[ScriptListParams](request.Params.Arguments)
	if err != nil {
		return nil, err
	}

	scripts, total, err := c.Repositories.SourceFile().List(ctx, domain.SourceFileSearchOptions{
		IDs:      params.Ids,
		Query:    params.Query,
		Page:     params.Page,
		PageSize: params.Size,
		Type:     params.Type,
		GameIDs:  params.GameIDs,
		GuildIDs: params.GuildIDs,
		NPCIDs:   params.NPCIDs,
		VoiceIDs: params.VoiceIDs,
	})
	if err != nil {
		return nil, err
	}

	resultStr, err := json.Marshal(map[string]any{
		"total":   total,
		"scripts": scripts,
	})

	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(resultStr)), nil
}
