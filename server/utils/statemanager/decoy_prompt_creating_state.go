package statemanager

import (
	"drawydraw/models"
	"errors"
)

type decoyPromptCreatingState struct {
	game *models.Game
}

func (state decoyPromptCreatingState) addPlayer(player *models.Player) error {
	// Only allow existing players to rejoin the game and in that case, no-op
	if state.game.IsPlayerInGame(player) {
		return nil
	}
	return errors.New("Cannot add new players to a game in this state")
}

func (state decoyPromptCreatingState) startGame(groupName string, playerName string) error {
	return errors.New("startGame not supported for decoyPromptCreatingStage state")
}

func (state decoyPromptCreatingState) submitDrawing(playerName string, encodedImage string) error {
	return errors.New("submitDrawing not supported for decoyPromptCreatingStage state")

}

func (state decoyPromptCreatingState) addPrompt(prompt *models.Prompt) error {
	activeDrawing := getCurrentDrawing(state.game)
	if activeDrawing == nil {
		return errors.New("Cannot submit a prompt when there's no current drawing")
	}
	if _, hasPrompt := activeDrawing.DecoyPrompts[prompt.Author]; hasPrompt {
		return errors.New("Player has already submitted a prompt for this drawing")
	}
	activeDrawing.DecoyPrompts[prompt.Author] = prompt
	// If all players have added their prompts move to the voting state
	if len(activeDrawing.DecoyPrompts) == len(state.game.Players)-1 {
		state.game.CurrentState = models.Voting
	}
	return nil
}

func (state decoyPromptCreatingState) addGameStatusPropertiesForPlayer(player *models.Player, gameStatus *GameStatusResponse) error {
	activeDrawing := getCurrentDrawing(state.game)
	if activeDrawing == nil {
		return errors.New("There is no active drawing available for this state")
	}
	gameStatus.CurrentDrawing = &Drawing{
		ImageData: activeDrawing.ImageData,
	}
	authorToDecoyPromptMap := map[string]*models.Prompt{}
	for _, currentPrompt := range activeDrawing.DecoyPrompts {
		authorToDecoyPromptMap[currentPrompt.Author] = currentPrompt
	}
	// Mark players who haven't submitted their prompts as having pending actions
	for _, p := range gameStatus.Players {
		_, hasPrompt := authorToDecoyPromptMap[p.Name]
		p.HasPendingAction = !hasPrompt
		if p.Name == player.Name {
			gameStatus.CurrentPlayer.HasCompletedAction = hasPrompt
		}
	}
	// If the current player is the author of the active drawing they have nothing to do but wait
	if activeDrawing.Author == player.Name {
		gameStatus.CurrentPlayer.HasCompletedAction = true
	}
	return nil
}

func getCurrentDrawing(game *models.Game) *models.Drawing {
	var activeDrawing *models.Drawing = nil
	for _, drawing := range game.Drawings {
		// Find the first drawing that has not been voted on
		if len(drawing.Votes) < len(game.Players) {
			activeDrawing = drawing
			break
		}
	}
	return activeDrawing
}
