package statemanager

import (
	"drawydraw/models"
	"errors"
	"math/rand"
	"time"
)

type promptCreatingState struct {
	game *models.Game
}

func (state promptCreatingState) addPlayer(player *models.Player) error {
	// Only allow existing players to rejoin the game and in that case, no-op
	if state.game.IsPlayerInGame(player.Name) {
		return nil
	}
	return errors.New("Cannot add new players to a game in this state")
}

func (state promptCreatingState) startGame(groupName string, playerName string) error {
	return errors.New("startGame not supported for initial prompt creation state")
}

func (state promptCreatingState) addPrompt(prompt *models.Prompt) error {
	//check if the player had already entered a prompt (not sure if needed)
	for _, p := range state.game.OriginalPrompts {
		if prompt.Author == p.Author {
			return errors.New("The player has already entered their prompt")
		}
	}

	state.game.AddPrompt(prompt)
	//TODO better logic to change state when all players have added prompts
	if len(state.game.OriginalPrompts) == len(state.game.Players) {
		state.game.CurrentState = models.DrawingsInProgress
		generatePrompts(state.game)
	}
	return nil
}

func (state promptCreatingState) submitDrawing(playerName string, encodedImage string) error {
	return errors.New("Submitting drawings is not allowed in the initial prompt creation state")
}

func generatePrompts(game *models.Game) {
	playerCount := len(game.Players)
	// Create a pool of adjectives
	adjectives := make([]string, 0, playerCount*2)
	playerPromptMap := map[string]*models.Prompt{}
	for _, prompt := range game.OriginalPrompts {
		playerPromptMap[prompt.Author] = prompt
		adjectives = append(adjectives, prompt.Adjectives...)
	}
	rand.Seed(time.Now().UnixNano())
	for index, player := range game.Players {
		// Give each player the noun entered by the next player
		previousPlayerIndex := (index + 1) % playerCount
		assignedNounAuthor := game.Players[previousPlayerIndex].Name
		// Pick and remove two random adjectives from the list
		firstAdjectiveIndex := rand.Intn(len(adjectives))
		firstAdjective := adjectives[firstAdjectiveIndex]
		adjectives = append(adjectives[:firstAdjectiveIndex], adjectives[firstAdjectiveIndex+1:]...)
		secondAdjectiveIndex := rand.Intn(len(adjectives))
		secondAdjective := adjectives[secondAdjectiveIndex]
		adjectives = append(adjectives[:secondAdjectiveIndex], adjectives[secondAdjectiveIndex+1:]...)
		player.AssignedPrompt = models.BuildPrompt(
			playerPromptMap[assignedNounAuthor].Noun,
			[]string{firstAdjective, secondAdjective},
			"generated",
		)
		game.GeneratedPrompts = append(game.GeneratedPrompts, player.AssignedPrompt)
	}
}

func (state promptCreatingState) addGameStatusPropertiesForPlayer(player *models.Player, gameStatus *GameStatusResponse) error {

	authorToPromptMap := map[string]*models.Prompt{}
	for _, currentPrompt := range state.game.OriginalPrompts {
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

func (state promptCreatingState) castVote(player *models.Player, promptIdentifier string) error {
	return errors.New("Casting votes is not allowed at this stage of the game")
}
