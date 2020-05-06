package stagemanager

import (
	"drawydraw/models"
	"errors"
)

type loginStage struct{}

func (stage loginStage) addPlayer(player *models.Player, gameState *models.GameState) error {
	// To do: Error out if the player is already there (or if we have max players?)
	gameState.Players = append(gameState.Players, player)
	return nil
}

func (stage loginStage) enterPrompt(prompt *models.Prompt, gameState *models.GameState) error {
	return errors.New("Entering prompts is not allowed at this stage")
}
