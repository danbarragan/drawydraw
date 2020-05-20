package models

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

// Prompt is a set of a noun and adjectives that describes a drawing someone will make or has made
type Prompt struct {
	Author     string
	Noun       string
	Adjectives []string
}

// Drawing represents a drawing someone has made
type Drawing struct {
	ImageData string
	Author    string
}

// Game contains all data that represents the game at any point
type Game struct {
	GroupName    string
	Players      []*Player
	CurrentState GameState
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
	game.Players = append(game.Players, player)
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
