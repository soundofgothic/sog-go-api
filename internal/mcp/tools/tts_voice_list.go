package tools

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/samber/lo"
	"soundofgothic.pl/backend/internal/domain"
)

func init() {
	registerTool(NewTTSVoiceList)
}

type TTSVoiceList struct {
	Repositories domain.Repositories
}

func NewTTSVoiceList(repositories domain.Repositories) *TTSVoiceList {
	return &TTSVoiceList{
		Repositories: repositories,
	}
}

func (c *TTSVoiceList) Tool() mcp.Tool {
	return mcp.NewTool(
		"list_tts_voices",
		mcp.WithDescription("Retrieves a list of available Text-To-Speech (TTS) voices that can be used for generating dialogue alternatives in new mods."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "list TTS voices",
			ReadOnlyHint: lo.ToPtr(true),
		}),
		// TODO: Add parameters for filtering TTS voices (e.g., by language_code, gender, provider, query)
		// and update the Handler and add a TTSVoiceListParams struct accordingly.
	)
}

func (c *TTSVoiceList) Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	voices, err := c.Repositories.TTSVoice().List(ctx)
	if err != nil {
		return nil, err
	}

	voicesString, err := json.Marshal(map[string]any{
		"tts_voices": voices,
	})
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(voicesString)), nil
}
