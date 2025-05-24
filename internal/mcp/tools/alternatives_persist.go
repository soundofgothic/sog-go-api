package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"soundofgothic.pl/backend/internal/domain"
)

func init() {
	registerTool(NewAlternativePersist)
}

type AlternativePersist struct {
	Repositories domain.Repositories
}

func NewAlternativePersist(repositories domain.Repositories) *AlternativePersist {
	return &AlternativePersist{
		Repositories: repositories,
	}
}

func (c *AlternativePersist) Tool() mcp.Tool {
	return mcp.NewTool(
		"persist_alternatives",
		mcp.WithDescription("Creates or updates alternatives for dialogue recordings. This tool accepts a list of alternatives to be persisted in the database."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title: "persist alternatives",
		}),
		mcp.WithArray("alternatives", mcp.Items(map[string]any{
			"type": "object",
			"properties": map[string]any{
				"mod_id": map[string]any{
					"type":        "number",
					"description": "The ID of the mod this alternative belongs to. Must be a positive integer.",
				},
				"recording_id": map[string]any{
					"type":        "number",
					"description": "The ID of the original recording for which this alternative is being created/updated. Must be a positive integer.",
				},
				"tts_voice_id": map[string]any{
					"type":        "number",
					"description": "The ID of the Text-To-Speech voice to be used for this alternative. Must be a positive integer.",
				},
				"transcript": map[string]any{
					"type":        "string",
					"description": "The transcribed text of the alternative dialogue. Cannot be empty.",
				},
			},
		}), mcp.Description("A list of alternative objects to be created or updated.")),
	)
}

type PeristedAlternative struct {
	ModID       int64  `json:"mod_id" validate:"required"`
	RecordingID int64  `json:"recording_id" validate:"required"`
	TTSVoiceID  int64  `json:"tts_voice_id" validate:"required"`
	Transcript  string `json:"transcript" validate:"required,min=1"`
}

type AlternativePersistParams struct {
	Alternatives []PeristedAlternative `json:"alternatives" validate:"required,dive"`
}

func (c *AlternativePersist) Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params, err := MapTo[AlternativePersistParams](request.Params.Arguments)
	if err != nil {
		return nil, err
	}

	for _, alt := range params.Alternatives {
		if err := c.Repositories.Alternative().Persist(ctx, &domain.Alternative{
			ModID:       alt.ModID,
			RecordingId: alt.RecordingID,
			TTSVoiceId:  alt.TTSVoiceID,
			Transcript:  alt.Transcript,
			State:       domain.AlternativeStateNone,
		}); err != nil {
			return nil, err
		}
	}

	return mcp.NewToolResultText("success"), nil
}
