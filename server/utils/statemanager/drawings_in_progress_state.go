package statemanager

import (
	"drawydraw/models"
	"errors"
)

type drawingsInProgressState struct {
	game *models.Game
}

func (state drawingsInProgressState) addPlayer(player *models.Player) error {
	// Only allow existing players to rejoin the game and in that case, no-op
	if state.game.IsPlayerInGame(player) {
		return nil
	}
	return errors.New("Cannot add new players to a game in this state")
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
	playerExists := false
	for _, player := range state.game.Players {
		if player.Name == playerName {
			playerExists = true
			break
		}
	}
	if !playerExists {
		return errors.New("player is not in the group")
	}
	drawing := models.Drawing{
		Author:       playerName,
		ImageData:    encodedImage,
		DecoyPrompts: map[string]*models.Prompt{},
	}
	state.game.Drawings = append(state.game.Drawings, &drawing)
	// If this is the last drawing, transition to the fake prompt creation state
	if len(state.game.Drawings) == len(state.game.Players) {
		state.game.CurrentState = models.DecoyPromptCreation
	}
	return nil
}

func (state drawingsInProgressState) addPrompt(prompts *models.Prompt) error {
	return errors.New("addprompts not supported for drawing state")
}

func (state drawingsInProgressState) addGameStatusPropertiesForPlayer(player *models.Player, gameStatus *GameStatusResponse) error {
	authorToDrawingMap := map[string]*models.Drawing{}
	for _, currentDrawing := range state.game.Drawings {
		authorToDrawingMap[currentDrawing.Author] = currentDrawing
	}
	// Mark players who haven't submitted their drawing as having pending actions
	for _, p := range gameStatus.Players {
		_, hasDrawing := authorToDrawingMap[p.Name]
		p.HasPendingAction = !hasDrawing
		if p.Name == player.Name {
			gameStatus.CurrentPlayer.HasCompletedAction = hasDrawing
		}
	}
	if player.AssignedPrompt != nil {
		gameStatus.CurrentPlayer.AssignedPrompt = &Prompt{
			Adjectives: player.AssignedPrompt.Adjectives,
			Noun:       player.AssignedPrompt.Noun,
		}
	}
	return nil
}
