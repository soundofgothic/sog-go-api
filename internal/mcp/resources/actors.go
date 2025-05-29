package resources

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"soundofgothic.pl/backend/internal/domain"
)

type ActorResource struct {
	repositories domain.Repositories
}

func NewActorResource(repositories domain.Repositories) *ActorResource {
	return &ActorResource{
		repositories: repositories,
	}
}

func (a *ActorResource) Resource() mcp.Resource {
	return mcp.NewResource(
		"actors:///",
		"Voice Actors",
		mcp.WithResourceDescription("Polish Voice actors in the Gothic universe. This data also include counts of recordings with them in the game"),
		mcp.WithMIMEType("application/json"),
	)
}

func (a *ActorResource) Handler(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	voices, err := a.repositories.Voice().List(ctx, domain.VoiceOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list voices: %w", err)
	}

	contents := make([]mcp.ResourceContents, len(voices))
	for i, voice := range voices {
		text, err := json.Marshal(voice)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal voice %d: %w", voice.ID, err)
		}

		contents[i] = mcp.TextResourceContents{
			URI:      fmt.Sprintf("actors:///%d", voice.ID),
			MIMEType: "application/json",
			Text:     string(text),
		}
	}

	return contents, nil
}
