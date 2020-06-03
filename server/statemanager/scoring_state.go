package statemanager

import (
	"drawydraw/models"
	"errors"
	"fmt"
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
	// Calculate standings and update total points
	standings := state.calculateStandings(activeDrawing, state.game)
	for name, standing := range *standings {
		player := state.game.GetPlayer(name)
		if player == nil {
			return fmt.Errorf("Could not find player %s in the game", name)
		}
		player.Points = standing.TotalScore
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
	gameStatus.CurrentDrawing = gameStatusDrawingFromDrawing(activeDrawing)
	// Add drawings that have been scored to past drawings
	gameStatus.PastDrawings = make([]*Drawing, 0, len(state.game.Drawings))
	for _, drawing := range state.game.Drawings {
		if drawing.Scored {
			gameStatus.PastDrawings = append(
				gameStatus.PastDrawings,
				gameStatusDrawingFromDrawing(drawing),
			)
		}
	}
	gameStatus.PointStandings = state.calculateStandings(activeDrawing, state.game)
	return nil
}

func (state scoringState) castVote(player *models.Player, promptIdentifier string) error {
	return errors.New("Casting votes is not allowed at this stage of the game")
}

func gameStatusDrawingFromDrawing(drawing *models.Drawing) *Drawing {
	return &Drawing{
		Author:         drawing.Author,
		ImageData:      drawing.ImageData,
		OriginalPrompt: makeResponsePromptFromModelPrompt(drawing.OriginalPrompt),
	}
}

func (state scoringState) calculateStandings(activeDrawing *models.Drawing, game *models.Game) *map[string]*PointStanding {
	pointStandings := map[string]*PointStanding{}
	// Initialize standings with the point totals before this round
	for _, player := range game.Players {
		pointStandings[player.Name] = &PointStanding{
			Player:               player.Name,
			RoundPointsBreakdown: []*PointsBreakdown{},
			TotalScore:           player.Points,
		}
	}
	for playerName, vote := range activeDrawing.Votes {
		if vote.SelectedPrompt == activeDrawing.OriginalPrompt {
			// Voter earns 3 points for picking the right prompt
			pointStandings[playerName].TotalScore += 3
			pointStandings[playerName].RoundPointsBreakdown = append(
				pointStandings[playerName].RoundPointsBreakdown,
				&PointsBreakdown{Amount: 3, Reason: ChoseCorrectPrompt, CausingPlayer: playerName},
			)
			// Author gets 1 point for someone picking the right prompt
			pointStandings[activeDrawing.Author].TotalScore += 1
			pointStandings[activeDrawing.Author].RoundPointsBreakdown = append(
				pointStandings[activeDrawing.Author].RoundPointsBreakdown,
				&PointsBreakdown{Amount: 1, Reason: OtherChosePromptDrawn, CausingPlayer: playerName},
			)
		} else {
			// The person who fooled the voter earns 1 point
			pointStandings[vote.SelectedPrompt.Author].TotalScore += 1
			pointStandings[vote.SelectedPrompt.Author].RoundPointsBreakdown = append(
				pointStandings[vote.SelectedPrompt.Author].RoundPointsBreakdown,
				&PointsBreakdown{Amount: 1, Reason: FooledPlayer, CausingPlayer: playerName},
			)
		}
	}
	return &pointStandings
}
