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

func (state promptCreatingState) addPrompts(prompts *models.Prompts) error {
	state.game.AddPrompts(prompts)
	return nil
}
