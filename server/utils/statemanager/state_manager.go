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

// CreateGroup Handles creating a group other players can join
func CreateGroup(playerName string, groupName string) (*models.Game, error) {
	// See if there's already a game for that group name and error out if ther eis
	gameState := models.LoadGame(groupName)
	if gameState != nil {
		return nil, errors.New("A group with that name already exists")
	}
	// Games start in the waiting for players stage
	gameState = &models.Game{GroupName: groupName, CurrentState: models.WaitingForPlayers}
	currentStage, err := getCurrentState(gameState)
	if err != nil {
		return nil, err
	}
	// Add the group creator as the first player
	player := models.Player{Name: playerName, Host: true}
	err = currentStage.addPlayer(&player)
	if err != nil {
		return nil, err
	}
	models.SaveGame(gameState)
	return gameState, nil
}

// AddPlayer Handles adding a player to a game
func AddPlayer(playerName string, groupName string) (*models.Game, error) {
	stateManager, err := getManagerForGroup(groupName)
	if err != nil {
		return nil, err
	}
	player := models.Player{Name: playerName}
	err = stateManager.currentState.addPlayer(&player)
	if err != nil {
		return nil, err
	}
	models.SaveGame(stateManager.game)
	return stateManager.game, nil
}

// GetGameState gets the current state for a given game
func GetGameState(groupName string) (*models.Game, error) {
	stateManager, err := getManagerForGroup(groupName)
	if err != nil {
		return nil, err
	}
	return stateManager.game, nil
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
