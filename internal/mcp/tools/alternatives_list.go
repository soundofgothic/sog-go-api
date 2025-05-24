package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/samber/lo"
	"soundofgothic.pl/backend/internal/domain"
)

func init() {
	registerTool(NewAlternativeList)
}

type AlternativeList struct {
	Repositories domain.Repositories
}

func NewAlternativeList(repositories domain.Repositories) *AlternativeList {
	return &AlternativeList{
		Repositories: repositories,
	}
}

func (c *AlternativeList) Tool() mcp.Tool {
	return mcp.NewTool(
		"list_alternatives",
		mcp.WithDescription("Retrieves a list of alternatives for dialogue recordings, with options to filter by alternative properties and properties of the associated original recordings."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "list alternatives",
			ReadOnlyHint: lo.ToPtr(true),
		}),
		mcp.WithNumber("mod_id", mcp.Description("Mod ID to filter alternatives. Set to zero or omit to skip filtering by mod."), mcp.Required()),
		mcp.WithString("alternative_query", mcp.Description("Query string to filter alternatives by their transcribed text."), mcp.DefaultString("")),
		mcp.WithString("state", mcp.Description("Filter alternatives by their current state (e.g., 'pending', 'approved', 'rejected')."), mcp.DefaultString("")),
		mcp.WithNumber("tts_voice_id", mcp.Description("Filter alternatives by the ID of the Text-To-Speech voice assigned to them. Set to zero or omit to skip this filter."), mcp.DefaultNumber(0)),
		// RecordingSearchOptions fields
		mcp.WithArray("ids", mcp.Items(map[string]any{"type": "number"}), mcp.Description("Filter alternatives based on a list of specific associated Recording IDs.")),
		mcp.WithString("recording_query", mcp.Description("Filter alternatives based on a query string matching the transcript of their associated original recordings."), mcp.DefaultString("")),
		mcp.WithNumber("page", mcp.Description("Page number for pagination of results."), mcp.DefaultNumber(1)),
		mcp.WithNumber("page_size", mcp.Description("Number of alternatives to return per page."), mcp.DefaultNumber(50)),
		mcp.WithArray("source_file_id", mcp.Items(map[string]any{"type": "number"}), mcp.Description("Filter alternatives based on a list of Source File IDs of their associated original recordings.")),
		mcp.WithArray("game_ids", mcp.Items(map[string]any{"type": "number"}), mcp.Description("Filter alternatives based on a list of Game IDs to which their associated original recordings belong.")),
		mcp.WithArray("npc_ids", mcp.Items(map[string]any{"type": "number"}), mcp.Description("Filter alternatives based on a list of NPC IDs associated with their original recordings.")),
		mcp.WithArray("guild_ids", mcp.Items(map[string]any{"type": "number"}), mcp.Description("Filter alternatives based on a list of Guild IDs associated with their original recordings.")),
		mcp.WithArray("voice_ids", mcp.Items(map[string]any{"type": "number"}), mcp.Description("Filter alternatives based on a list of Voice IDs (original voice actor) associated with their original recordings.")),
	)
}

type AlternativeListParams struct {
	ModID            int64   `json:"mod_id" validate:"required"`
	AlternativeQuery string  `json:"alternative_query"`
	State            string  `json:"state"`
	TTSVoiceID       int64   `json:"tts_voice_id"`
	IDs              []int64 `json:"ids"`
	RecordingQuery   string  `json:"recording_query"`
	Page             int64   `json:"page"`
	PageSize         int64   `json:"page_size"`
	SourceFileID     []int64 `json:"source_file_id"`
	GameIDs          []int64 `json:"game_ids"`
	NPCIDs           []int64 `json:"npc_ids"`
	GuildIDs         []int64 `json:"guild_ids"`
	VoiceIDs         []int64 `json:"voice_ids"`
}

func (c *AlternativeList) Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params, err := MapTo[AlternativeListParams](request.Params.Arguments)
	if err != nil {
		return nil, err
	}

	options := domain.AlternativeOptions{
		RecordingSearchOptions: domain.RecordingSearchOptions{
			IDs:          params.IDs,
			Query:        params.RecordingQuery,
			Page:         params.Page,
			PageSize:     params.PageSize,
			SourceFileID: params.SourceFileID,
			GameIDs:      params.GameIDs,
			NPCIDs:       params.NPCIDs,
			GuildIDs:     params.GuildIDs,
			VoiceIDs:     params.VoiceIDs,
		},
		ModID:      params.ModID,
		Query:      params.AlternativeQuery,
		State:      params.State,
		TTSVoiceID: params.TTSVoiceID,
	}

	alternatives, total, err := c.Repositories.Alternative().List(ctx, options)
	if err != nil {
		return nil, err
	}

	responseString, err := json.Marshal(map[string]any{
		"total":        total,
		"alternatives": alternatives,
	})
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(responseString)), nil
}
