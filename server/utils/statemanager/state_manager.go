package statemanager

import (
	"drawydraw/models"
	"errors"
)

// StateManager handles the different states and actions throughout the game
type StateManager struct {
	currentState state
	game         *models.Game
}

// GameStatusResponse describes a response informing clients of the game status
type GameStatusResponse = map[string]interface{}

// CreateGroup Handles creating a group other players can join
func CreateGroup(groupName string) (*GameStatusResponse, error) {
	if len(groupName) < 1 {
		return nil, errors.New("Group name too short.")
	}
	// See if there's already a game for that group name and error out if ther eis
	gameState := models.LoadGame(groupName)
	if gameState != nil {
		return nil, errors.New("A group with that name already exists")
	}
	// Games start in the waiting for players stage
	gameState = &models.Game{GroupName: groupName, CurrentState: models.WaitingForPlayers}
	models.SaveGame(gameState)
	formattedState, err := formatGameStateForPlayer(gameState, playerName)
	if err != nil {
		return nil, err
	}
	return formattedState, nil
}

// AddPlayer Handles adding a player to a game
func AddPlayer(playerName string, groupName string, isHost bool) (*GameStatusResponse, error) {

	if len(playerName) < 1 {
		return nil, errors.New("Player name too short.")
	}
	stateManager, err := getManagerForGroup(groupName)
	if err != nil {
		return nil, err
	}

	if isPlayerInGroup(playerName, stateManager.game.Players) {
			return nil, errors.New("A player with that name already exists in that group")
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

// Todo: This should probably be moved to individual states since there will be different relevant data at each state
func formatGameStateForPlayer(game *models.Game, playerName string) (*GameStatusResponse, error) {
	gameHost, err := game.GetHostName()
	if err != nil {
		return nil, err
	}
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
	case models.WaitingForPlayers:
		return waitingForPlayersState{game: game}, nil
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
