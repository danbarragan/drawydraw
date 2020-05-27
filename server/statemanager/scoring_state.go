package statemanager

import (
	"drawydraw/models"
	"errors"
	"fmt"
	"strings"
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
	activeDrawing := state.game.GetActiveDrawing()
	if activeDrawing == nil {
		return errors.New("Could not find active drawing for game")
	}
	// Add all the round scores
	roundScores := state.calculateRoundScores(activeDrawing, state.game)
	for name, score := range *roundScores {
		player := state.game.GetPlayer(name)
		if player == nil {
			return fmt.Errorf("Could not find player %s in the game", name)
		}
		for _, points := range score {
			player.Points += points.Amount
		}
	}
	// Mark the active drawing as scored
	activeDrawing.Scored = true
	if state.game.GetActiveDrawing() != nil {
		// If there's another active drawing, go to the decoy prompts state
		state.game.CurrentState = models.DecoyPromptCreation
	} else {
		// If there isn't, then we need to reset prompts/drawings and go to prompts
		state.game.OriginalPrompts = []*models.Prompt{}
		state.game.GeneratedPrompts = []*models.Prompt{}
		state.game.Drawings = []*models.Drawing{}
		// Also clear out assigned prompts
		for _, player := range state.game.Players {
			player.AssignedPrompt = nil
		}
		state.game.CurrentState = models.InitialPromptCreation
	}
	return nil
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
	gameStatus.CurrentDrawing = &Drawing{
		ImageData: activeDrawing.ImageData,
		OriginalPrompt: fmt.Sprintf(
			"%s %s",
			strings.Join(activeDrawing.OriginalPrompt.Adjectives, ","),
			activeDrawing.OriginalPrompt.Noun,
		),
	}
	gameStatus.RoundScores = state.calculateRoundScores(activeDrawing, state.game)
	return nil
}

func (state scoringState) castVote(player *models.Player, promptIdentifier string) error {
	return errors.New("Casting votes is not allowed at this stage of the game")
}

func (state scoringState) calculateRoundScores(activeDrawing *models.Drawing, game *models.Game) *map[string][]*PointsBreakdown {
	scoreMap := map[string][]*PointsBreakdown{}
	// Initialize score map with empty score arrays
	for _, player := range game.Players {
		scoreMap[player.Name] = []*PointsBreakdown{}
	}
	for playerName, vote := range activeDrawing.Votes {
		if vote.SelectedPrompt == activeDrawing.OriginalPrompt {
			// Voter earns 3 points for picking the right prompt
			scoreMap[playerName] = append(
				scoreMap[playerName],
				&PointsBreakdown{Amount: 3, Reason: "you chose the correct prompt"},
			)
			// Author gets 1 point for someone picking the right prompt
			scoreMap[activeDrawing.Author] = append(
				scoreMap[activeDrawing.Author],
				&PointsBreakdown{Amount: 1, Reason: fmt.Sprintf("%s chose your original prompt", playerName)},
			)
		} else {
			// The person who fooled the voter earns 1 point
			scoreMap[vote.SelectedPrompt.Author] = append(
				scoreMap[vote.SelectedPrompt.Author],
				&PointsBreakdown{Amount: 1, Reason: fmt.Sprintf("%s chose your decoy prompt", playerName)},
			)
		}
	}
	return &scoreMap
}
