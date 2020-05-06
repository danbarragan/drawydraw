package models

import (
	"time"

	"github.com/patrickmn/go-cache"
)

const gameStatusKey = "GAME_STATUS"

var (
	// 20 min TTL, purges every 5 min
	memCache = cache.New(20*time.Minute, 5*time.Minute)
)

// GameStage defines what are the individual stages that make up the game
type GameStage uint64 // Probably want to change this to str someday to make the stored representation friendlier

const (
	// WaitingForPlayers - A group is created and the host is waiting for players
	WaitingForPlayers GameStage = iota
	// InitialPromptCreation - Players are entering their initial prompts
	InitialPromptCreation
)

// Player contains all the information relevant to a game's participant
type Player struct {
	Name   string
	Host   bool
	Points uint64
}

// Prompt is to be implemented...
type Prompt struct {
}

// GameState contains all data that represents the state of the game at any point
type GameState struct {
	GroupName    string
	Players      []*Player
	CurrentStage GameStage
}

// GetGameState returns the current game state for a given room name
func GetGameState(roomName string) *GameState {
	state, found := memCache.Get(gameStatusKey)
	if found {
		return state.(*GameState)
	}
	return nil
}

// Save persists the game state
func (state *GameState) Save() error {
	memCache.Set(gameStatusKey, state, cache.DefaultExpiration)
	return nil
}
