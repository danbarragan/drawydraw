package models

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	// 20 min TTL, purges every 5 min
	memCache = cache.New(20*time.Minute, 5*time.Minute)
)

// GameStage defines what are the individual stages that make up the game
type GameStage string

const (
	// WaitingForPlayers - A group is created and the host is waiting for players
	WaitingForPlayers GameStage = "WaitingForPlayers"
	// InitialPromptCreation - Players are entering their initial prompts
	InitialPromptCreation GameStage = "InitialPromptCreation"
)

// Player contains all the information relevant to a game's participant
type Player struct {
	Name   string `json:"name"`
	Host   bool   `json:"host"`
	Points uint64 `json:"points"`
}

// GameState contains all data that represents the state of the game at any point
type GameState struct {
	GroupName    string    `json:"groupName"`
	Players      []*Player `json:"players"`
	CurrentStage GameStage `json:"currentStage"`
}

// Todo: Put these methods behind an interface to faciliate unit tests

// LoadGameState returns the current game state for a given room name
func LoadGameState(roomName string) *GameState {
	state, found := memCache.Get(roomName)
	if found {
		return state.(*GameState)
	}
	return nil
}

// SaveGameState persists the game state
func SaveGameState(state *GameState) error {
	memCache.Set(state.GroupName, state, cache.DefaultExpiration)
	return nil
}
