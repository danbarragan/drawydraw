package statemanager

import (
	"drawydraw/models"
	"errors"
)

type waitingForPlayersState struct {
	game *models.Game
}

func (state waitingForPlayersState) addPlayer(player *models.Player) error {
	state.game.AddPlayer(player)
	return nil
}

func (state waitingForPlayersState) startGame(groupName string, playerName string) error {
	if playerName != *state.game.GetHostName() {
		return errors.New("only the host can start a game")
	}
	state.game.CurrentState = models.InitialPromptCreation
	return nil
}

func (state waitingForPlayersState) submitDrawing(playerName string, encodedImage string) error {
	return errors.New("Submitting drawings is not allowed in the wating for players state")
}

func (state waitingForPlayersState) addPrompt(prompts *models.Prompt) error {
	return errors.New("addPrompt not supported for waiting for players state")
}

func (state waitingForPlayersState) addGameStatusPropertiesForPlayer(player *models.Player, gameStatus *GameStatusResponse) error {
	return nil
}

func (state waitingForPlayersState) castVote(player *models.Player, promptIdentifier string) error {
	return errors.New("Casting votes is not allowed at this stage of the game")
}