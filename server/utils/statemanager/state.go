package statemanager

import (
	"drawydraw/models"
)

// Interface implemented by each state's concrete handler
type state interface {
	addPlayer(player *models.Player) error
	addPrompts(prompts *models.Prompts) error
}
