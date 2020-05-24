package statemanager

import (
	"drawydraw/models"
	"errors"
	"sort"
)

type votingState struct {
	game *models.Game
}

func (state votingState) addPlayer(player *models.Player) error {
	// Only allow existing players to rejoin the game and in that case, no-op
	if state.game.IsPlayerInGame(player.Name) {
		return nil
	}
	return errors.New("Cannot add new players to a game in this state")
}

func (state votingState) startGame(groupName string, playerName string) error {
	return errors.New("startGame not supported for voting state")
}

func (state votingState) submitDrawing(playerName string, encodedImage string) error {
	return errors.New("Submitting drawings is not allowed in the voting state")
}

func (state votingState) addPrompt(prompts *models.Prompt) error {
	return errors.New("addprompts not supported for voting state")
}

func (state votingState) addGameStatusPropertiesForPlayer(player *models.Player, gameStatus *GameStatusResponse) error {
	activeDrawing := state.game.GetActiveDrawing()
	if activeDrawing == nil {
		return errors.New("There is no active drawing available for this state")
	}

	// Get all available prompts for the drawing
	prompts := make([]*Prompt, 0, len(state.game.Players))
	prompts = append(prompts, makeResponsePromptFromModelPrompt(activeDrawing.OriginalPrompt))
	for _, decoyPrompt := range activeDrawing.DecoyPrompts {
		prompts = append(prompts, makeResponsePromptFromModelPrompt(decoyPrompt))
	}
	// Sort the prompts by prompt id
	sort.Slice(prompts, func(i, j int) bool { return prompts[i].Identifier < prompts[j].Identifier })
	gameStatus.CurrentDrawing = &Drawing{
		ImageData: activeDrawing.ImageData,
		Prompts:   prompts,
	}
	// Mark players who haven't submitted their votes as having pending actions
	casterToVoteMap := map[string]*models.Vote{}
	for _, vote := range activeDrawing.Votes {
		casterToVoteMap[vote.Player.Name] = vote
	}
	for _, p := range gameStatus.Players {
		// The author of the drawing can't vote mark them  so they have no pending action
		if p.Name == activeDrawing.Author {
			p.HasPendingAction = false
		} else {
			_, hasVoted := casterToVoteMap[p.Name]
			p.HasPendingAction = !hasVoted
		}
	}
	// Current player has completed their action if they're the author of if they already voted
	if activeDrawing.Author == player.Name || casterToVoteMap[player.Name] != nil {
		gameStatus.CurrentPlayer.HasCompletedAction = true
	}
	return nil
}

func makeResponsePromptFromModelPrompt(prompt *models.Prompt) *Prompt {
	return &Prompt{Noun: prompt.Noun, Adjectives: prompt.Adjectives, Identifier: prompt.Identifier}
}

func (state votingState) castVote(player *models.Player, promptIdentifier string) error {
	activeDrawing := state.game.GetActiveDrawing()
	if activeDrawing == nil {
		return errors.New("There is no active drawing available for this state")
	}

	prompt := activeDrawing.GetPromptWithIdentifier(promptIdentifier)
	if prompt == nil {
		return errors.New("Could not find the chosen prompt in the active drawing")
	}

	activeDrawing.Votes[player.Name] = &models.Vote{Player: player, SelectedPrompt: prompt}
	// If all players have voted move to the scoring state
	if len(activeDrawing.Votes) == len(state.game.Players)-1 {
		state.game.CurrentState = models.Scoring
	}
	return nil
}
