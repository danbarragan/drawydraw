package statemanager

import (
	"drawydraw/models"
	"errors"
)

type initialPromptCreation struct {
	game *models.Game
}

func (state initialPromptCreation) addPlayer(player *models.Player) error {
	state.game.AddPlayer(player)
	return nil
}

func (state initialPromptCreation) startGame(groupName string, playerName string) error {
	return errors.New("startGame not supported for initial prompt creation state")
}
