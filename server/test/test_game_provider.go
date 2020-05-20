package test

import (
	"drawydraw/models"
	"testing"
)

// TestGameProvider facilitates testing by having a simple implementation that doesn't
// involve caches or external calls.
type TestGameProvider struct {
	games map[string]*models.Game
}

func NewTestGameProvider() *TestGameProvider {
	return &TestGameProvider{games: map[string]*models.Game{}}
}

func (provider *TestGameProvider) LoadGame(groupName string) *models.Game {
	return provider.games[groupName]
}

func (provider *TestGameProvider) SaveGame(game *models.Game) error {
	provider.games[game.GroupName] = game
	return nil
}

// SetupTestGameProvider sets up a clean test game provider and tears it down after the test finishes
func SetupTestGameProvider(t *testing.T) {
	previousProvider := models.GetGameProvider()
	models.SetGameProvider(NewTestGameProvider())
	t.Cleanup(func() {
		models.SetGameProvider(previousProvider)
	})
}
