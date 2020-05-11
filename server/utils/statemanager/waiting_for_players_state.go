package statemanager

import (
	"drawydraw/models"
)

type waitingForPlayersState struct {
	game *models.Game
}

func (state waitingForPlayersState) addPlayer(player *models.Player) error {
	state.game.AddPlayer(player)
	return nil
}
