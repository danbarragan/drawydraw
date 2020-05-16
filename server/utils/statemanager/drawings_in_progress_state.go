package statemanager

import (
	"drawydraw/models"
	"errors"
)

type drawingsInProgressState struct {
	game *models.Game
}

func (state drawingsInProgressState) addPlayer(player *models.Player) error {
	state.game.AddPlayer(player)
	return nil
}

func (state drawingsInProgressState) startGame(groupName string, playerName string) error {
	return errors.New("startGame not supported for drawingsInProgress state")
}

func (state drawingsInProgressState) addPrompt(prompt *models.Prompt) error {
	return errors.New("addprompt not supported for drawing state")
}
