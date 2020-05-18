package statemanager

import (
	"drawydraw/models"
	"errors"
	"fmt"
)

// StateManager handles the different states and actions throughout the game
type StateManager struct {
	currentState state
	game         *models.Game
}

// Models used for describing the status of the game to clients
// Todo: Find a better home for these (or replace them with protos)
type AssignedPrompt struct {
	Adjectives []string `json:"adjectives"`
	Noun       string   `json:"noun"`
}

type Player struct {
	Name   string `json:"name"`
	Host   bool   `json:"host"`
	Points uint64 `json:"points"`
}

type CurrentPlayer struct {
	AssignedPrompt *AssignedPrompt `json:"assignedPrompt"`
	IsHost         bool            `json:"isHost"`
	Name           string          `json:"name"`
}

type GameStatusResponse struct {
	CurrentPlayer *CurrentPlayer `json:"currentPlayer"`
	CurrentState  string         `json:"currentState"`
	GroupName     string         `json:"groupName"`
	Players       []*Player      `json:"players"`
}

// CreateGroup Handles creating a group other players can join
func CreateGroup(groupName string) error {
	if len(groupName) < 1 {
		return errors.New("no group name provided")
	}
	// See if there's already a game for that group name and error out if ther eis
	gameState := models.LoadGame(groupName)
	if gameState != nil {
		return fmt.Errorf("group '%s' already exists", groupName)
	}
	// Games start in the waiting for players stage
	gameState = &models.Game{
		GroupName: groupName, CurrentState: models.WaitingForPlayers,
	}
	models.SaveGame(gameState)
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

	if isHost {
		if stateManager.game.HostPlayer != "" &&
			playerName != stateManager.game.HostPlayer {
			return nil, fmt.Errorf("failed to add player %s as host - %s is already host", playerName, stateManager.game.HostPlayer)
		}
	} else {
		// Non-host player is joining a game without a host - this should not be possible
		if stateManager.game.HostPlayer == "" {
			return nil, errors.New("cannot add a non-host player to a game without a host")
		}
	}

	// Add the group creator as the first player
	player := models.Player{Name: playerName, Host: isHost}
	err = stateManager.currentState.addPlayer(&player)
	if err != nil {
		return nil, err
	}
	models.SaveGame(stateManager.game)
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
		return nil, errors.New("One or more of the prompt was not provided")
	}

	stateManager, err := getManagerForGroup(groupName)
	if err != nil {
		return nil, err
	}

	//check if the player had already entered a prompt (not sure if needed)
	for _, prompt := range stateManager.game.Prompts {
		if playerName == prompt.Author {
			return nil, errors.New("The player has already entered their prompt")
		}
	}

	newPrompt := models.Prompt{
		Author:     playerName,
		Noun:       noun,
		Adjectives: []string{adjective1, adjective2}}

	err = stateManager.currentState.addPrompt(&newPrompt)
	if err != nil {
		return nil, err
	}
	models.SaveGame(stateManager.game)
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
	models.SaveGame(stateManager.game)
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
	gameState := models.LoadGame(groupName)
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
	case models.Voting:
		return voting{game: game}, nil
	case models.WaitingForPlayers:
		return waitingForPlayersState{game: game}, nil
	case models.InitialPromptCreation:
		return promptCreatingState{game: game}, nil
	case models.DrawingsInProgress:
		return drawingsInProgressState{game: game}, nil
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
	switch currentState := gameState; currentState {
	case models.Voting:
		return createGameState("chats", []string{"graisseux", "frere jacques", "pepe le pew"}, nil, gameState)
	case models.WaitingForPlayers:
		return createGameState("not cats", []string{"dog", "cat", "other dog"}, nil, gameState)
	case models.InitialPromptCreation:
		return createGameState("fat cats", []string{"chubbs", "chonk", "beefcake"}, nil, gameState)
	case models.DrawingsInProgress:
		mockPrompt := []string{"silly", "great", "beluga"}
		mockPrompts := [][]string{mockPrompt, mockPrompt, mockPrompt}
		return createGameState("human cats", []string{"sharon", "grandpa", "j. ralphio"}, mockPrompts, gameState)
	default:
		return nil, fmt.Errorf("failed to set game to state %s", gameState)
	}
}

func createGameState(groupName string, players []string, prompts [][]string, gameState models.GameState) (*GameStatusResponse, error) {
	hostName := players[0]

	game := &models.Game{
		GroupName: groupName, CurrentState: gameState, HostPlayer: hostName,
	}
	models.SaveGame(game)
	for idx, playerName := range players {
		isHost := false
		if idx == 0 {
			isHost = true
		}
		_, err := AddPlayer(playerName, groupName, isHost)
		if err != nil {
			return nil, err
		}
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
	gameStatus, err := gameStatusForPlayer(game, hostName)
	if err != nil {
		return nil, err
	}
	return gameStatus, nil
}
