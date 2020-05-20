package models

import (
	"sync"
)

var once sync.Once

// GameProvider defines the interface different providers of game storage implement
type GameProvider interface {
	LoadGame(groupName string) *Game
	SaveGame(game *Game) error
}

var (
	gameProvider GameProvider = nil
)

// GetGameProvider gets the provider to be used for loading and saving games
func GetGameProvider() GameProvider {
	once.Do(func() {
		gameProvider = createMemcacheGameProvider()
	})
	return gameProvider
}

// SetGameProvider changes the provider to be used for loading and saving games.
// This should not be called outside test code
func SetGameProvider(provider GameProvider) {
	gameProvider = provider
}
