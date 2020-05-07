package stagemanager

import (
	"drawydraw/models"
)

type loginStage struct{}

func (stage loginStage) addPlayer(player *models.Player, gameState *models.GameState) error {
	// Todo: Error out if the player is already there (or if we have max players?)
	gameState.Players = append(gameState.Players, player)
	return nil
}
