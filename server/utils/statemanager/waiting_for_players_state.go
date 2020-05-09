package statemanager

import (
	"drawydraw/models"
)

type waitingForPlayersState struct {
	game *models.Game
}

func (state waitingForPlayersState) addPlayer(player *models.Player) error {
	// Todo: Error out if the player is already there (or if we have max players?)
	state.game.Players = append(state.game.Players, player)
	return nil
}
