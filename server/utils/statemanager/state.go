package statemanager

import (
	"drawydraw/models"
)

// Interface implemented by each state's concrete handler
type state interface {
	addPlayer(player *models.Player) error
	addPrompt(prompt *models.Prompt) error
	startGame(groupName string, playerName string) error
}
