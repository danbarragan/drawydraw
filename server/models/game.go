package models

import (
	"fmt"
	"hash/fnv"
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
	// DecoyPromptCreation - Players are creating decoy prompts for a drawing they didn't make
	DecoyPromptCreation GameState = "DecoyPromptCreation"
	// Voting - Players currently drawing a prompt
	Voting GameState = "Voting"
	// Scoring - Players are shown the current scores
	Scoring GameState = "Scoring"
)

// Player contains all the information relevant to a game's participant
type Player struct {
	// Todo: Probably worth having a sort of device id in case two players register with the same name
	Name           string
	Host           bool
	Points         uint64
	AssignedPrompt *Prompt
}

// Prompt is a set of a noun and adjectives that describes a drawing someone will make or has made
type Prompt struct {
	Identifier string
	Author     string
	Noun       string
	Adjectives []string
}

// Vote represents a prompt selected by a player in a drawing
type Vote struct {
	Player         *Player
	SelectedPrompt *Prompt
}

// Drawing represents a drawing someone has made
type Drawing struct {
	ImageData      string
	Author         string
	DecoyPrompts   map[string]*Prompt
	OriginalPrompt *Prompt
	Votes          map[string]*Vote
	Scored         bool
}

// Game contains all data that represents the game at any point
type Game struct {
	GroupName    string
	Players      []*Player
	CurrentState GameState
	Prompts      []*Prompt
	Drawings     []*Drawing
}

// GetPromptWithIdentifier returns the prompt that has a given identifier in a drawing
func (drawing *Drawing) GetPromptWithIdentifier(identifier string) *Prompt {
	if identifierForPrompt(drawing.OriginalPrompt) == identifier {
		return drawing.OriginalPrompt
	}
	for _, prompt := range drawing.DecoyPrompts {
		if identifierForPrompt(prompt) == identifier {
			return prompt
		}
	}
	return nil
}

// AddPlayer adds a player to the game (if that player isn't there already)
func (game *Game) AddPlayer(player *Player) error {
	// First check if the player is already in the game and no-op if that's the case
	if !game.IsPlayerInGame(player.Name) {
		game.Players = append(game.Players, player)
	}
	return nil
}

// IsPlayerInGame  determines if a player is already in a game or not
func (game *Game) IsPlayerInGame(playerName string) bool {
	return game.GetPlayer(playerName) != nil
}

// GetPlayer returns a player object for the given player name
func (game *Game) GetPlayer(playerName string) *Player {
	for _, currentPlayer := range game.Players {
		if currentPlayer.Name == playerName {
			return currentPlayer
		}
	}
	return nil
}

// AddPrompt adds a player's prompt to the game
func (game *Game) AddPrompt(prompt *Prompt) error {
	game.Prompts = append(game.Prompts, prompt)
	return nil
}

// GetHostName Gets the name of the game's host, returns nil if there's no host yet
func (game *Game) GetHostName() *string {
	for _, currentPlayer := range game.Players {
		if currentPlayer.Host {
			return &currentPlayer.Name
		}
	}
	return nil
}

// GetActiveDrawing gets the drawing players are either entering prompts for or voting on prompts for
func (game *Game) GetActiveDrawing() *Drawing {
	for _, drawing := range game.Drawings {
		// Find the first drawing that has not been scored
		if !drawing.Scored {
			return drawing
		}
	}
	return nil
}

// BuildPrompt creates a prompt object with the right internal properties
func BuildPrompt(noun string, adjectives []string, author string) *Prompt {
	prompt := &Prompt{Noun: noun, Adjectives: adjectives, Author: author}
	prompt.Identifier = identifierForPrompt(prompt)
	return prompt
}

func identifierForPrompt(prompt *Prompt) string {
	hashFunction := fnv.New64()
	hashFunction.Write([]byte(
		fmt.Sprintf("%s-%s-%s", prompt.Adjectives[0], prompt.Adjectives[1], prompt.Noun),
	))
	return fmt.Sprintf("%d", hashFunction.Sum64())
}
