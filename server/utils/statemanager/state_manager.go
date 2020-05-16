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

// GameStatusResponse describes a response informing clients of the game status
type GameStatusResponse = map[string]interface{}

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

	if isHost &&
		stateManager.game.HostPlayer != "" &&
		playerName != stateManager.game.HostPlayer {
		return nil, fmt.Errorf("failed to add player %s as host - %s is already host", playerName, stateManager.game.HostPlayer)
	}

	// Add the group creator as the first player
	player := models.Player{Name: playerName, Host: isHost}
	err = stateManager.currentState.addPlayer(&player)
	if err != nil {
		return nil, err
	}
	models.SaveGame(stateManager.game)
	formattedState, err := formatGameStateForPlayer(stateManager.game, playerName)
	if err != nil {
		return nil, err
	}
	return formattedState, nil
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
		Group:      groupName,
		Noun:       noun,
		Adjective1: adjective1,
		Adjective2: adjective2,
	}

	err = stateManager.currentState.addPrompt(&newPrompt)
	if err != nil {
		return nil, err
	}
	models.SaveGame(stateManager.game)
	formattedState, err := formatGameStateForPlayer(stateManager.game, playerName)
	if err != nil {
		return nil, err
	}
	return formattedState, nil
}

// GetGameState gets the current state for a given game and player
func GetGameState(groupName string, playerName string) (*GameStatusResponse, error) {
	stateManager, err := getManagerForGroup(groupName)
	if err != nil {
		return nil, err
	}
	formattedState, err := formatGameStateForPlayer(stateManager.game, playerName)
	if err != nil {
		return nil, err
	}
	return formattedState, nil
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
	formattedState, err := formatGameStateForPlayer(stateManager.game, playerName)
	if err != nil {
		return nil, err
	}
	return formattedState, nil
}

// Todo: This should probably be moved to individual states since there will be different relevant data at each state
func formatGameStateForPlayer(game *models.Game, playerName string) (*GameStatusResponse, error) {
	gameHost, err := game.GetHostName()
	if err != nil {
		return nil, err
	}

	//TODO : Need to adjust this so we don't send all the prompt back to the client, someone can cheat by inspecting
	statusResponse := map[string]interface{}{
		"groupName": game.GroupName,
		"players":   game.Players,
		"currentPlayer": map[string]interface{}{
			"name":   playerName,
			"isHost": gameHost == playerName,
		},
		"currentState": game.CurrentState,
	}
	return &statusResponse, nil
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
		return createGameState("chats", []string{"graisseux", "frere jacques", "pepe le pew"}, gameState)
	case models.WaitingForPlayers:
		return createGameState("not cats", []string{"dog", "cat", "other dog"}, gameState)
	case models.InitialPromptCreation:
		return createGameState("fat cats", []string{"chubbs", "chonk", "beefcake"}, gameState)
	case models.DrawingsInProgress:
		return createGameState("human cats", []string{"sharon", "grandpa", "j. ralphio"}, gameState)
	default:
		return nil, fmt.Errorf("failed to set game to state %s", gameState)
	}
}

func createGameState(groupName string, players []string, gameState models.GameState) (*GameStatusResponse, error) {
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

	formattedState, err := formatGameStateForPlayer(game, hostName)
	if err != nil {
		return nil, err
	}
	return formattedState, nil
}
