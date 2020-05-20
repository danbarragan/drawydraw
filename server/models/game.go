package models

import (
	"errors"
)

// GameState defines what are the individual states that make up the game
type GameState string

const (
	// WaitingForPlayers - A group is created and the host is waiting for players
	WaitingForPlayers GameState = "WaitingForPlayers"
	// InitialPromptCreation - Players are entering their initial prompt
	InitialPromptCreation GameState = "InitialPromptCreation"
	// DrawingsInProgress - Players currently drawing a prompt
	DrawingsInProgress GameState = "DrawingsInProgress"
	// Voting - Players currently drawing a prompt
	Voting GameState = "Voting"
)

// Player contains all the information relevant to a game's participant
type Player struct {
	// Todo: Probably worth having a sort of device id in case two players register with the same name
	Name           string
	Host           bool
	Points         uint64
	AssignedPrompt *Prompt
}

type Prompt struct {
	Author     string
	Noun       string
	Adjectives []string
}

type Drawing struct {
	ImageData string
	Author    string
}

// Game contains all data that represents the game at any point
type Game struct {
	GroupName    string
	Players      []*Player
	CurrentState GameState
	HostPlayer   string
	Prompts      []*Prompt
	Drawings     []*Drawing
}

// AddPlayer adds a player to the game (if that player isn't there already)
func (game *Game) AddPlayer(player *Player) error {
	// First check if the player is already in the game and no-op if that's the case
	for _, currentPlayer := range game.Players {
		if currentPlayer.Name == player.Name {
			return nil
		}
	}
	// Set the host name if the player joining the game is the host
	if player.Host {
		game.HostPlayer = player.Name
	}
	game.Players = append(game.Players, player)
	return nil
}

// AddPrompt adds a player's prompt to the game
func (game *Game) AddPrompt(prompt *Prompt) error {
	game.Prompts = append(game.Prompts, prompt)
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
