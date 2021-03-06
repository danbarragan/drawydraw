package statemanager

import (
	"drawydraw/models"
)

// Interface implemented by each state's concrete handler
type state interface {
	addPlayer(player *models.Player) error
	addPrompt(prompt *models.Prompt) error
	startGame(groupName string, playerName string) error
	submitDrawing(playerName string, encodedImage string) error
	castVote(player *models.Player, promptIdentifier string) error
	addGameStatusPropertiesForPlayer(player *models.Player, gameStatus *GameStatusResponse) error
}
