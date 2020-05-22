package statemanager

import (
	"drawydraw/models"
	"errors"
)

type scoringState struct {
	game *models.Game
}

func (state scoringState) addPlayer(player *models.Player) error {
	// Only allow existing players to rejoin the game and in that case, no-op
	if state.game.IsPlayerInGame(player.Name) {
		return nil
	}
	return errors.New("Cannot add new players to a game in this state")
}

func (state scoringState) startGame(groupName string, playerName string) error {
	return errors.New("startGame not supported for voting state")
}

func (state scoringState) submitDrawing(playerName string, encodedImage string) error {
	return errors.New("Submitting drawings is not allowed in the voting state")
}

func (state scoringState) addPrompt(prompts *models.Prompt) error {
	return errors.New("addprompts not supported for voting state")
}

func (state scoringState) addGameStatusPropertiesForPlayer(player *models.Player, gameStatus *GameStatusResponse) error {
	activeDrawing := state.game.GetActiveDrawing()
	if activeDrawing == nil {
		return errors.New("Could not find active drawing for game")
	}
	gameStatus.RoundScores = state.calculateRoundScores(activeDrawing)
	return nil
}

func (state scoringState) castVote(player *models.Player, promptIdentifier string) error {
	return errors.New("Casting votes is not allowed at this stage of the game")
}

func (state scoringState) calculateRoundScores(activeDrawing *models.Drawing) *map[string]uint64 {
	scoreMap := map[string]uint64{}
	// Initialize score map with 0s
	for _, player := range state.game.Players {
		scoreMap[player.Name] = 0
	}

	for playerName, vote := range activeDrawing.Votes {
		if vote.SelectedPrompt == activeDrawing.OriginalPrompt {
			// Voter earns 3 points for picking the right prompt
			scoreMap[playerName] += 3
			// Author gets 1 point for someone picking the right prompt
			scoreMap[activeDrawing.Author]++
		} else {
			// The person who fooled the voter earns 1 point
			scoreMap[vote.SelectedPrompt.Author]++
		}
	}
	return &scoreMap
}
