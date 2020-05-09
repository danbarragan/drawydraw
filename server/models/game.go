package models

import (
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	// 20 min TTL, purges every 5 min
	memCache = cache.New(20*time.Minute, 5*time.Minute)
)

// GameState defines what are the individual states that make up the game
type GameState string

const (
	// WaitingForPlayers - A group is created and the host is waiting for players
	WaitingForPlayers GameState = "WaitingForPlayers"
	// InitialPromptCreation - Players are entering their initial prompts
	InitialPromptCreation GameState = "InitialPromptCreation"
)

// Player contains all the information relevant to a game's participant
type Player struct {
	Name   string `json:"name"`
	Host   bool   `json:"host"`
	Points uint64 `json:"points"`
}

// Game contains all data that represents the game at any point
type Game struct {
	GroupName    string    `json:"groupName"`
	Players      []*Player `json:"players"`
	CurrentState GameState `json:"currentState"`
}

// Todo: Put these methods behind an interface to faciliate unit tests

// LoadGame returns the current game for a given room name
func LoadGame(roomName string) *Game {
	state, found := memCache.Get(roomName)
	if found {
		return state.(*Game)
	}
	return nil
}

// SaveGame persists the game
func SaveGame(state *Game) error {
	memCache.Set(state.GroupName, state, cache.DefaultExpiration)
	return nil
}
