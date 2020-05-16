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

func (state promptCreatingState) addPrompt(prompt *models.Prompt) error {
	state.game.AddPrompt(prompt)

	//TODO better logic to change state when all players have added prompt
	if len(state.game.Prompts) == len(state.game.Players) {
		state.game.CurrentState = models.DrawingsInProgress
	}

	return nil
}
