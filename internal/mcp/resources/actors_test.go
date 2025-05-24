package resources

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/require"
	"soundofgothic.pl/backend/internal/test"
)

func TestActorsResource(t *testing.T) {
	repositories := test.NewRepositories()

	resource := NewActorResource(repositories)
	content, err := resource.Handler(t.Context(), mcp.ReadResourceRequest{})
	require.NoError(t, err)
	require.NotEmpty(t, content)
}
