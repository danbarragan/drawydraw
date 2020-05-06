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

// CreateForGroup creates a stage manager instance for a given group's game
func CreateForGroup(groupName string) (*StageManager, error) {
	gameState := models.GetGameState(groupName)
	currentStage, err := getCurrentStage(gameState)
	if err != nil {
		return nil, err
	}
	stageManager := StageManager{currentStage: currentStage, gameState: gameState}
	return &stageManager, nil
}

// AddPlayer Handles adding a player to a game
func (stageManager *StageManager) AddPlayer(playerName string) error {
	player := models.Player{}
	err := stageManager.currentStage.addPlayer(&player, stageManager.gameState)
	if err != nil {
		return err
	}
	stageManager.gameState.Save()
	return nil
}

func getCurrentStage(gameState *models.GameState) (gameStage, error) {
	switch stage := gameState.CurrentStage; stage {
	case models.WaitingForPlayers:
		return loginStage{}, nil
	default:
		return nil, errors.New("Game is at an unknown stage")
	}
}
