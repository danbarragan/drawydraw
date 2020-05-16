package statemanager

import (
	"drawydraw/models"
	"errors"
)

type voting struct {
	game *models.Game
}

func (state voting) addPlayer(player *models.Player) error {
	state.game.AddPlayer(player)
	return nil
}

func (state voting) startGame(groupName string, playerName string) error {
	return errors.New("startGame not supported for voting state")
}