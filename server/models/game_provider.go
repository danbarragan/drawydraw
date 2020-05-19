package models

// GameProvider defines the interface different providers of game storage implement
type GameProvider interface {
	LoadGame(groupName string) *Game
	SaveGame(game *Game) error
}

var gameProvider GameProvider = nil

// GetGameProvider gets the provider to be used for loading and saving games
func GetGameProvider() GameProvider {
	if gameProvider != nil {
		return gameProvider
	}
	return getMemcacheGameProviderInstance()
}

// SetGameProvider changes the provider to be used for loading and saving games
func SetGameProvider(provider GameProvider) {
	gameProvider = provider
}
