package statemanager

import (
	"drawydraw/models"
	"errors"
)

type promptCreatingState struct {
	game *models.Game
}

func (state promptCreatingState) addPlayer(player *models.Player) error {
	state.game.AddPlayer(player)
	return nil
}

func (state promptCreatingState) startGame(groupName string, playerName string) error {
	return errors.New("startGame not supported for initial prompt creation state")
}

func (state promptCreatingState) addPrompt(prompt *models.Prompt) error {
	//check if the player had already entered a prompt (not sure if needed)
	for _, p := range state.game.Prompts {
		if prompt.Author == p.Author {
			return errors.New("The player has already entered their prompt")
		}
	}

	state.game.AddPrompt(prompt)
	//TODO better logic to change state when all players have added prompts
	if len(state.game.Prompts) == len(state.game.Players) {
		state.game.CurrentState = models.DrawingsInProgress
		assignPrompts(state.game)
	}
	return nil
}

func (state promptCreatingState) submitDrawing(playerName string, encodedImage string) error {
	return errors.New("Submitting drawings is not allowed in the initial prompt creation state")
}

func assignPrompts(game *models.Game) {
	// Really dumb prompt assignment, each player draws whatever the next player (in joining order) entered
	playerPromptMap := map[string]*models.Prompt{}
	for _, prompt := range game.Prompts {
		playerPromptMap[prompt.Author] = prompt
	}
	playerCount := len(game.Players)
	for index, player := range game.Players {
		previousPlayerIndex := (index + 1) % playerCount
		assignedPromptAuthor := game.Players[previousPlayerIndex].Name
		player.AssignedPrompt = playerPromptMap[assignedPromptAuthor]
	}
}

func (state promptCreatingState) addGameStatusPropertiesForPlayer(player *models.Player, gameStatus *GameStatusResponse) error {

	authorToPromptMap := map[string]*models.Prompt{}
	for _, currentPrompt := range state.game.Prompts {
		authorToPromptMap[currentPrompt.Author] = currentPrompt
	}

	// Mark players who haven't submitted their prompt as having pending actions
	for _, p := range gameStatus.Players {
		_, hasPrompt := authorToPromptMap[p.Name]
		p.HasPendingAction = !hasPrompt

		if p.Name == player.Name {
			gameStatus.CurrentPlayer.HasCompletedAction = hasPrompt
		}
	}

	return nil
}
