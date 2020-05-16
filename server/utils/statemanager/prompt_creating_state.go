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

func (state promptCreatingState) addPrompts(prompts *models.Prompt) error {
	state.game.AddPrompts(prompts)
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
	// Really dumb prompt assignment, each player draws whatever the previous player entered
	playerPromptMap := map[string]*models.Prompt{}
	for _, prompt := range game.Prompts {
		playerPromptMap[prompt.Author] = prompt
	}
	playerCount := len(game.Players)
	for index, player := range game.Players {
		// Go's modulus operator is not as magic as python's
		previousPlayerIndex := (index-1)%playerCount
		assignedPromptAuthor := game.Players[].Name
		player.AssignedPrompt = playerPromptMap[assignedPromptAuthor]
	}
}
