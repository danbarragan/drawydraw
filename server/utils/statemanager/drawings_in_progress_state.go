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

func (state drawingsInProgressState) submitDrawing(playerName string, encodedImage string) error {
	for _, currentDrawing := range state.game.Drawings {
		if currentDrawing.Author == playerName {
			return errors.New("player has already submitted a drawing")
		}
	}
	drawing := models.Drawing{Author: playerName, ImageData: encodedImage}
	state.game.Drawings = append(state.game.Drawings, &drawing)
	return nil
}

func (state drawingsInProgressState) addPrompt(prompts *models.Prompt) error {
	return errors.New("addprompts not supported for drawing state")
}

func (state drawingsInProgressState) addGameStatusPropertiesForPlayer(player *models.Player, gameStatus *GameStatusResponse) error {
	gameStatus.CurrentPlayer.AssignedPrompt = &AssignedPrompt{
		Adjectives: player.AssignedPrompt.Adjectives,
		Noun:       player.AssignedPrompt.Noun,
	}
	return nil
}
