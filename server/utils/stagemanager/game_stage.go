package stagemanager

import (
	"drawydraw/models"
)

// Shared interface exposed by all of the game's stages
type gameStage interface {
	addPlayer(player *models.Player, gameState *models.GameState) error
}
