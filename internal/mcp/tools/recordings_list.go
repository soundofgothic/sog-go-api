package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/samber/lo"
	"soundofgothic.pl/backend/internal/domain"
)

func init() {
	registerTool(NewRecordingsList)
}

type RecordingsList struct {
	Repositories domain.Repositories
}

func NewRecordingsList(repositories domain.Repositories) *RecordingsList {
	return &RecordingsList{
		Repositories: repositories,
	}
}

func (c *RecordingsList) Tool() mcp.Tool {
	return mcp.NewTool(
		"list_recording", // Consider renaming to list_recordings for consistency
		mcp.WithDescription("Retrieves a list of original dialogue recordings, with extensive filtering options based on recording metadata and associated game entities."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "list recordings",
			ReadOnlyHint: lo.ToPtr(true),
		}),
		mcp.WithArray("ids", mcp.Items(map[string]any{
			"type": "number",
		}), mcp.Description("Filter recordings by a list of specific Recording IDs.")),
		mcp.WithString("query", mcp.Description("Search term to filter recordings by their transcribed text."), mcp.DefaultString("")),
		mcp.WithNumber("page", mcp.Description("Page number for pagination of results."), mcp.DefaultNumber(1), mcp.Min(1)),
		mcp.WithNumber("page_size", mcp.Description("Number of recordings to return per page."), mcp.DefaultNumber(20), mcp.Max(20)),
		mcp.WithArray("source_file_id", mcp.Items(map[string]any{
			"type": "number",
		}), mcp.Description("Filter recordings by a list of Source File IDs they originate from.")),
		mcp.WithArray("game_ids", mcp.Items(map[string]any{
			"type": "number",
		}), mcp.Description("Filter recordings belonging to a list of specific Game IDs.")),
		mcp.WithArray("npc_ids", mcp.Items(map[string]any{
			"type": "number",
		}), mcp.Description("Filter recordings spoken by NPCs in a list of specific NPC IDs.")),
		mcp.WithArray("guild_ids", mcp.Items(map[string]any{
			"type": "number",
		}), mcp.Description("Filter recordings associated with NPCs belonging to a list of specific Guild IDs.")),
		mcp.WithArray("voice_ids", mcp.Items(map[string]any{
			"type": "number",
		}), mcp.Description("Filter recordings by a list of specific Voice IDs (original voice actor).")),
	)
}

type RecordingListParams struct {
	IDs          []int64 `json:"ids"`
	Query        string  `json:"query"`
	Page         int64   `json:"page"`
	PageSize     int64   `json:"page_size" validate:"max=20"`
	SourceFileID []int64 `json:"source_file_id"`
	GameIDs      []int64 `json:"game_ids"`
	NPCIDs       []int64 `json:"npc_ids"`
	GuildIDs     []int64 `json:"guild_ids"`
	VoiceIDs     []int64 `json:"voice_ids"`
}

func (c *RecordingsList) Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params, err := MapTo[RecordingListParams](request.Params.Arguments)
	if err != nil {
		return nil, err
	}

	recordings, total, err := c.Repositories.Recording().List(ctx, domain.RecordingSearchOptions{
		IDs:          params.IDs,
		Query:        params.Query,
		Page:         params.Page,
		PageSize:     params.PageSize,
		SourceFileID: params.SourceFileID,
		GameIDs:      params.GameIDs,
		NPCIDs:       params.NPCIDs,
		GuildIDs:     params.GuildIDs,
		VoiceIDs:     params.VoiceIDs,
	})
	if err != nil {
		return nil, err
	}

	resultStr, err := json.Marshal(map[string]any{
		"total":      total,
		"recordings": recordings,
	})

	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(resultStr)), nil
}
