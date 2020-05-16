package statemanager

import (
	"drawydraw/models"
	"errors"
)

type drawingsInProgress struct {
	game *models.Game
}

func (state drawingsInProgress) addPlayer(player *models.Player) error {
	state.game.AddPlayer(player)
	return nil
}

func (state drawingsInProgress) startGame(groupName string, playerName string) error {
	return errors.New("startGame not supported for drawingsInProgress state")
}
