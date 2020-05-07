package stagemanager

import (
	"drawydraw/models"
	"errors"
)

// StageManager handles different actions throughout the game
type StageManager struct {
	currentStage gameStage
	gameState    *models.GameState
}

// CreateGroup Handles creating a group other players can join
func CreateGroup(playerName string, groupName string) (*models.GameState, error) {
	// See if there's already a game for that group name and error out if ther eis
	gameState := models.LoadGameState(groupName)
	if gameState != nil {
		return nil, errors.New("A group with that name already exists")
	}
	// Games start in the waiting for players stage
	gameState = &models.GameState{GroupName: groupName, CurrentStage: models.WaitingForPlayers}
	currentStage, err := getCurrentStage(gameState)
	if err != nil {
		return nil, err
	}
	// Add the group creator as the first player
	player := models.Player{Name: playerName, Host: true}
	err = currentStage.addPlayer(&player, gameState)
	if err != nil {
		return nil, err
	}
	models.SaveGameState(gameState)
	return gameState, nil
}

// AddPlayer Handles adding a player to a game
func AddPlayer(playerName string, groupName string) (*models.GameState, error) {
	stageManager, err := getManagerForGroup(groupName)
	if err != nil {
		return nil, err
	}
	player := models.Player{Name: playerName}
	err = stageManager.currentStage.addPlayer(&player, stageManager.gameState)
	if err != nil {
		return nil, err
	}
	models.SaveGameState(stageManager.gameState)
	return stageManager.gameState, nil
}

// GetGameState gets the current state for a given game
func GetGameState(groupName string) (*models.GameState, error) {
	stageManager, err := getManagerForGroup(groupName)
	if err != nil {
		return nil, err
	}
	return stageManager.gameState, nil
}

func getManagerForGroup(groupName string) (*StageManager, error) {
	gameState := models.LoadGameState(groupName)
	if gameState == nil {
		return nil, errors.New("Could not find a group with that name")
	}
	currentStage, err := getCurrentStage(gameState)
	if err != nil {
		return nil, err
	}
	stageManager := StageManager{currentStage: currentStage, gameState: gameState}
	return &stageManager, nil
}

func getCurrentStage(gameState *models.GameState) (gameStage, error) {
	switch stage := gameState.CurrentStage; stage {
	case models.WaitingForPlayers:
		return loginStage{}, nil
	default:
		return nil, errors.New("Game is at an unknown stage")
	}
}
