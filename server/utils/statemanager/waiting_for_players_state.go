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
	if playerName != state.game.HostPlayer {
		return errors.New("only the host can start a game")
	}
	state.game.CurrentState = models.InitialPromptCreation
	return nil
}
