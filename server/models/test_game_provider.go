package models

// TestGameProvider facilitates testing by having a simple implementation that doesn't
// involve caches or external calls. It's meant for testing, don't use it outside of that!
type TestGameProvider struct {
	games map[string]*Game
}

func NewTestGameProvider() *TestGameProvider {
	return &TestGameProvider{games: map[string]*Game{}}
}

func (provider *TestGameProvider) LoadGame(groupName string) *Game {
	return provider.games[groupName]
}

func (provider *TestGameProvider) SaveGame(game *Game) error {
	provider.games[game.GroupName] = game
	return nil
}
