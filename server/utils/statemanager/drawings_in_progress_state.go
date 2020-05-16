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

func (state drawingsInProgressState) addPrompts(prompts *models.Prompts) error {
	return errors.New("addprompts not supported for drawing state")
}
