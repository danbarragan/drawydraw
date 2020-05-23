package statemanager

import (
	"drawydraw/models"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// StateManager handles the different states and actions throughout the game
type StateManager struct {
	currentState state
	game         *models.Game
}

// Models used for describing the status of the game to clients
// Todo: Find a better home for these (or replace them with protos)

// Prompt represents the prompt currently assigned to a player for drawing
type Prompt struct {
	Identifier string   `json:"identifier"`
	Adjectives []string `json:"adjectives"`
	Noun       string   `json:"noun"`
}

// Player represents the status of a player other than the one making the request
type Player struct {
	Name             string `json:"name"`
	Host             bool   `json:"host"`
	Points           uint64 `json:"points"`
	HasPendingAction bool   `json:"hasPendingAction"`
}

// CurrentPlayer represents the status of the player making the request
type CurrentPlayer struct {
	AssignedPrompt     *Prompt `json:"assignedPrompt"`
	IsHost             bool    `json:"isHost"`
	Name               string  `json:"name"`
	HasCompletedAction bool    `json:"hasCompletedAction"`
}

// Drawing represents a drawing that players are either making prompts for or voting on prompts for it
type Drawing struct {
	ImageData string    `json:"imageData"`
	Prompts   []*Prompt `json:"prompts"`
}

// GameStatusResponse contains all the game status communicated to players
type GameStatusResponse struct {
	CurrentPlayer  *CurrentPlayer     `json:"currentPlayer"`
	CurrentState   string             `json:"currentState"`
	GroupName      string             `json:"groupName"`
	Players        []*Player          `json:"players"`
	CurrentDrawing *Drawing           `json:"currentDrawing"`
	RoundScores    *map[string]uint64 `json:"roundScores"`
}

// CreateGroup Handles creating a group other players can join
func CreateGroup(groupName string) error {
	if len(groupName) < 1 {
		return errors.New("no group name provided")
	}
	// See if there's already a game for that group name and error out if ther eis
	gameState := models.GetGameProvider().LoadGame(groupName)
	if gameState != nil {
		return fmt.Errorf("group '%s' already exists", groupName)
	}
	// Games start in the waiting for players stage
	gameState = &models.Game{
		GroupName: groupName, CurrentState: models.WaitingForPlayers,
	}
	models.GetGameProvider().SaveGame(gameState)
	return nil
}

// AddPlayer Handles adding a player to a game
func AddPlayer(playerName string, groupName string, isHost bool) (*GameStatusResponse, error) {
	if len(playerName) < 1 {
		return nil, errors.New("no player name provided")
	}
	stateManager, err := getManagerForGroup(groupName)
	if err != nil {
		return nil, err
	}

	hostName := stateManager.game.GetHostName()
	if isHost {
		if hostName != nil &&
			playerName != *hostName {
			return nil, fmt.Errorf("failed to add player %s as host - %s is already host", playerName, *hostName)
		}
	} else {
		// Non-host player is joining a game without a host - this should not be possible
		if hostName == nil {
			return nil, errors.New("cannot add a non-host player to a game without a host")
		}
	}

	// Add the group creator as the first player
	player := models.Player{Name: playerName, Host: isHost}
	err = stateManager.currentState.addPlayer(&player)
	if err != nil {
		return nil, err
	}
	models.GetGameProvider().SaveGame(stateManager.game)
	gameStatus, err := gameStatusForPlayer(stateManager.game, playerName)
	if err != nil {
		return nil, err
	}
	return gameStatus, nil
}

// AddPrompt handles adding the prompt a player created to the game state
func AddPrompt(playerName string, groupName string, noun string, adjective1 string, adjective2 string) (*GameStatusResponse, error) {
	//check if any of the prompt fields were empty
	if len(noun) < 1 ||
		len(adjective1) < 1 ||
		len(adjective2) < 1 {
		return nil, errors.New("Prompt is missing a field")
	}

	stateManager, err := getManagerForGroup(groupName)
	if err != nil {
		return nil, err
	}

	newPrompt := models.BuildPrompt(noun, []string{adjective1, adjective2}, playerName)

	err = stateManager.currentState.addPrompt(newPrompt)
	if err != nil {
		return nil, err
	}
	models.GetGameProvider().SaveGame(stateManager.game)
	gameStatus, err := gameStatusForPlayer(stateManager.game, playerName)
	if err != nil {
		return nil, err
	}
	return gameStatus, nil
}

// SubmitDrawing handles a player submitting a drawing
func SubmitDrawing(playerName string, groupName string, imageData string) (*GameStatusResponse, error) {
	//check if the image data is empty
	if len(imageData) < 1 {
		return nil, errors.New("Image data was not provided")
	}

	stateManager, err := getManagerForGroup(groupName)
	if err != nil {
		return nil, err
	}

	err = stateManager.currentState.submitDrawing(playerName, imageData)
	if err != nil {
		return nil, err
	}
	models.GetGameProvider().SaveGame(stateManager.game)
	gameStatus, err := gameStatusForPlayer(stateManager.game, playerName)
	if err != nil {
		return nil, err
	}
	return gameStatus, nil
}

// CastVote handles a player casting a vote for a prompt in a drawing
func CastVote(playerName string, groupName string, promptIdentifier string) (*GameStatusResponse, error) {
	stateManager, err := getManagerForGroup(groupName)
	if err != nil {
		return nil, err
	}

	player := stateManager.game.GetPlayer(playerName)
	if player == nil {
		return nil, errors.New("Player is not in the game")
	}

	err = stateManager.currentState.castVote(player, promptIdentifier)
	if err != nil {
		return nil, err
	}
	models.GetGameProvider().SaveGame(stateManager.game)
	gameStatus, err := gameStatusForPlayer(stateManager.game, playerName)
	if err != nil {
		return nil, err
	}
	return gameStatus, nil
}

// GetGameState gets the current state for a given game and player
func GetGameState(groupName string, playerName string) (*GameStatusResponse, error) {
	stateManager, err := getManagerForGroup(groupName)
	if err != nil {
		return nil, err
	}
	gameStatus, err := gameStatusForPlayer(stateManager.game, playerName)
	if err != nil {
		return nil, err
	}
	return gameStatus, nil
}

// StartGame starts the game with the current players
func StartGame(groupName string, playerName string) (*GameStatusResponse, error) {
	stateManager, err := getManagerForGroup(groupName)
	if err != nil {
		return nil, err
	}

	err = stateManager.currentState.startGame(groupName, playerName)
	if err != nil {
		return nil, err
	}
	models.GetGameProvider().SaveGame(stateManager.game)
	gameStatus, err := gameStatusForPlayer(stateManager.game, playerName)
	if err != nil {
		return nil, err
	}
	return gameStatus, nil
}

func gameStatusForPlayer(game *models.Game, playerName string) (*GameStatusResponse, error) {
	var currentPlayer *models.Player
	players := make([]*Player, len(game.Players))
	for i, player := range game.Players {
		players[i] = &Player{Name: player.Name, Points: player.Points, Host: player.Host}
		if player.Name == playerName {
			currentPlayer = player
		}
	}
	if currentPlayer == nil {
		return nil, errors.New("cannot find current player in game")
	}
	// Set base properties that do not depend on game state
	gameStatusResponse := &GameStatusResponse{
		GroupName:     game.GroupName,
		CurrentPlayer: &CurrentPlayer{Name: currentPlayer.Name, IsHost: currentPlayer.Host},
		CurrentState:  string(game.CurrentState),
		Players:       players,
	}
	// Add any state-dependent properties to the status
	currentState, err := getCurrentState(game)
	if err != nil {
		return nil, err
	}
	err = currentState.addGameStatusPropertiesForPlayer(currentPlayer, gameStatusResponse)
	if err != nil {
		return nil, err
	}
	return gameStatusResponse, nil
}

func getManagerForGroup(groupName string) (*StateManager, error) {
	gameState := models.GetGameProvider().LoadGame(groupName)
	if gameState == nil {
		return nil, errors.New("Could not find a group with that name")
	}
	stateHandler, err := getCurrentState(gameState)
	if err != nil {
		return nil, err
	}
	stateManager := StateManager{currentState: stateHandler, game: gameState}
	return &stateManager, nil
}

func getCurrentState(game *models.Game) (state, error) {
	switch currentState := game.CurrentState; currentState {
	case models.DecoyPromptCreation:
		return decoyPromptCreatingState{game: game}, nil
	case models.Voting:
		return votingState{game: game}, nil
	case models.WaitingForPlayers:
		return waitingForPlayersState{game: game}, nil
	case models.InitialPromptCreation:
		return promptCreatingState{game: game}, nil
	case models.DrawingsInProgress:
		return drawingsInProgressState{game: game}, nil
	case models.Scoring:
		return scoringState{game: game}, nil
	default:
		return nil, errors.New("Game is at an unknown state")
	}
}

func isPlayerInGroup(playerName string, playersInGroup []*models.Player) bool {
	for _, playerInGroup := range playersInGroup {
		if playerInGroup.Name == playerName {
			return true
		}
	}
	return false
}

// DEBUG CODE - dont keep this forever.

// SetGameState is a debug method for forcing the gamestate to make UI testing easier.
func SetGameState(gameStateName string) (*GameStatusResponse, error) {
	gameState := models.GameState(gameStateName)
	mockPrompts := [][]string{
		{"silly", "great", "beluga"},
		{"happy", "elderly", "unicorn"},
	}
	currentDir, _ := os.Getwd()
	mockImageDataPath := path.Join(currentDir, "test", "testImageData.txt")
	mockImageData, err := ioutil.ReadFile(mockImageDataPath)
	if err != nil {
		fmt.Printf("Could not find file at %s, current dir %s", mockImageDataPath, currentDir)
		return nil, err
	}
	switch currentState := gameState; currentState {
	case models.DecoyPromptCreation:
		mockDrawings := []*models.Drawing{
			{
				Votes:          map[string]*models.Vote{},
				Author:         "chair",
				ImageData:      string(mockImageData),
				DecoyPrompts:   map[string]*models.Prompt{},
				OriginalPrompt: models.BuildPrompt("lady", []string{"serious", "mysterious"}, "table"),
			},
			{
				Votes:          map[string]*models.Vote{},
				Author:         "table",
				ImageData:      string(mockImageData),
				DecoyPrompts:   map[string]*models.Prompt{},
				OriginalPrompt: models.BuildPrompt("woman", []string{"smiling", "kind"}, "chair"),
			},
		}
		return createGameState("furnitures", []string{"table", "chair"}, mockPrompts, mockDrawings, gameState)
	case models.Voting:
		mockDrawings := []*models.Drawing{
			{
				Votes:     map[string]*models.Vote{},
				Author:    "tablet",
				ImageData: string(mockImageData),
				DecoyPrompts: map[string]*models.Prompt{
					"phone": models.BuildPrompt("person", []string{"weird", "funky"}, "phone"),
				},
				OriginalPrompt: models.BuildPrompt("lady", []string{"serious", "mysterious"}, "tablet"),
			},
			{
				Votes:          map[string]*models.Vote{},
				Author:         "phone",
				ImageData:      string(mockImageData),
				DecoyPrompts:   map[string]*models.Prompt{},
				OriginalPrompt: models.BuildPrompt("woman", []string{"smiling", "kind"}, "phone"),
			},
		}
		return createGameState("electronics", []string{"phone", "tablet"}, mockPrompts, mockDrawings, gameState)
	case models.WaitingForPlayers:
		return createGameState("not cats", []string{"dog", "cat", "other dog"}, nil, nil, gameState)
	case models.InitialPromptCreation:
		return createGameState("fat cats", []string{"chubbs", "chonk", "beefcake"}, nil, nil, gameState)
	case models.DrawingsInProgress:
		return createGameState("human cats", []string{"sharon", "grandpa"}, mockPrompts, nil, gameState)
	default:
		return nil, fmt.Errorf("failed to set game to state %s", gameState)
	}
}

func createGameState(
	groupName string,
	players []string,
	prompts [][]string,
	drawings []*models.Drawing,
	gameState models.GameState,
) (*GameStatusResponse, error) {
	game := &models.Game{
		GroupName:    groupName,
		CurrentState: gameState,
		Players:      make([]*models.Player, len(players)),
	}
	for idx, playerName := range players {
		game.Players[idx] = &models.Player{Name: playerName, Host: idx == 0}
	}
	if prompts != nil {
		game.Prompts = make([]*models.Prompt, len(prompts))
		for index, prompt := range prompts {
			game.Prompts[index] = &models.Prompt{
				Adjectives: prompt[0:2],
				Noun:       prompt[2],
				Author:     players[index],
			}
		}
		assignPrompts(game)
	}
	game.Drawings = drawings
	models.GetGameProvider().SaveGame(game)
	gameStatus, err := gameStatusForPlayer(game, *game.GetHostName())
	if err != nil {
		return nil, err
	}
	return gameStatus, nil
}
