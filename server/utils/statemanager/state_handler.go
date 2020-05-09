package statemanager

import (
	"drawydraw/models"
)

// Shared interface implemented by each state's concrete handler
type stateHandler interface {
	addPlayer(player *models.Player) error
}
