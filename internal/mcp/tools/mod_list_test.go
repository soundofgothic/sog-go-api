package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"soundofgothic.pl/backend/internal/domain"
	"soundofgothic.pl/backend/internal/test"
)

func TestModList_Handler_ReturnsMods(t *testing.T) {
	repos := test.NewRepositories()
	modList := NewModList(repos)

	request := NewMCPRequest(nil)
	result, err := modList.Handler(context.Background(), request)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	var mods []domain.Mod
	err = json.Unmarshal([]byte(AsTextContent(result)), &mods)
	assert.NoError(t, err)
}
