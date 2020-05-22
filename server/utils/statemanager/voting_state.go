package statemanager

import (
	"drawydraw/models"
	"errors"
)

type voting struct {
	game *models.Game
}

func (state voting) addPlayer(player *models.Player) error {
	// Only allow existing players to rejoin the game and in that case, no-op
	if state.game.IsPlayerInGame(player) {
		return nil
	}
	return errors.New("Cannot add new players to a game in this state")
}

func (state voting) startGame(groupName string, playerName string) error {
	return errors.New("startGame not supported for voting state")
}

func (state voting) submitDrawing(playerName string, encodedImage string) error {
	return errors.New("Submitting drawings is not allowed in the voting state")
}

func (state voting) addPrompt(prompts *models.Prompt) error {
	return errors.New("addprompts not supported for voting state")
}

func (state voting) addGameStatusPropertiesForPlayer(player *models.Player, gameStatus *GameStatusResponse) error {
	return nil
}
