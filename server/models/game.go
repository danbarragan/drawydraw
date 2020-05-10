package models

import (
	"errors"
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
	// Todo: Probably worth having a sort of device id in case two players register with the same name
	Name   string `json:"name"`
	Host   bool   `json:"host"`
	Points uint64 `json:"points"`
}

// Game contains all data that represents the game at any point
type Game struct {
	GroupName    string    `json:"groupName"`
	Players      []*Player `json:"players"`
	CurrentState GameState `json:"currentState"`
	HostPlayer   string    `json:"hostPlayer"`
}

// Todo: Put SaveGame/LoadGame methods behind an interface to faciliate unit tests

// LoadGame returns the current game for a given group name
func LoadGame(groupName string) *Game {
	state, found := memCache.Get(groupName)
	if found {
		return state.(*Game)
	}
	return nil
}

// AddPlayer adds a player to the game (if that player isn't there already)
func (game *Game) AddPlayer(player *Player) error {
	// First check if the player is already in the game and no-op if that's the case
	for _, currentPlayer := range game.Players {
		if currentPlayer.Name == player.Name {
			return nil
		}
	}
	game.Players = append(game.Players, player)
	return nil
}

// GetHostName Gets the name of the game's host
func (game *Game) GetHostName() (string, error) {
	for _, currentPlayer := range game.Players {
		if currentPlayer.Host {
			return currentPlayer.Name, nil
		}
	}
	return "", errors.New("game has no host")
}

// SaveGame persists the game
func SaveGame(state *Game) error {
	memCache.Set(state.GroupName, state, cache.DefaultExpiration)
	return nil
}
