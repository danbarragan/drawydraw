package statemanager

import (
	"drawydraw/models"
	"errors"
)

type promptCreatingState struct {
	game *models.Game
}

func (state promptCreatingState) addPlayer(player *models.Player) error {
	return errors.New("Add Player cannot be performed in the Prompt Creating State")
}

func (state promptCreatingState) startGame(groupName string, playerName string) error {
	return errors.New("startGame not supported for initial prompt creation state")
}

func (state promptCreatingState) addPrompts(prompts *models.Prompts) error {
	state.game.AddPrompts(prompts)

	//TODO better logic to change state when all players have added prompts
	if len(state.game.Prompts) == len(state.game.Players) {
		state.game.CurrentState = models.DrawingsInProgress
	}

	return nil
}

func (state promptCreatingState) submitDrawing(groupName string, playerName string, encodedImage string) error {
	return errors.New("Submitting drawings is not allowed in the initial prompt creation state")
}
