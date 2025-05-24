package tools

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"soundofgothic.pl/backend/internal/domain"
	"soundofgothic.pl/backend/internal/test"
)

func TestModDelete_Integration(t *testing.T) {
	repos := test.NewRepositories()
	modService := repos.Mod()
	ctx := context.Background()

	// Create a mod to delete
	mod := &domain.Mod{
		Name:        "Delete Me",
		Description: "To be deleted",
		Version:     "1.0.0",
		GameID:      2,
	}
	err := modService.Persist(ctx, mod)
	assert.NoError(t, err)
	assert.NotZero(t, mod.ID)

	// Now delete it using the tool
	tool := NewModDelete(repos)
	params := map[string]any{"id": mod.ID}
	request := NewMCPRequest(params)
	result, err := tool.Handler(ctx, request)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, AsTextContent(result), "deleted mod id:")

	// Ensure it is actually deleted
	modAfter, err := modService.Get(ctx, mod.ID)
	assert.Error(t, err)
	assert.Nil(t, modAfter)
}
