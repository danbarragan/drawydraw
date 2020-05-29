package test

import (
	"drawydraw/models"
)

// This file contains game models at different states to be used in testing

// GameInWaitingForPlayersState describes a newly created game waiting for players
func GameInWaitingForPlayersState() *models.Game {
	return &models.Game{
		GroupName: "somegame",
		Players: []*models.Player{
			{Name: "player1", Host: true},
			{Name: "player2"},
			{Name: "player3"},
		},
		CurrentState: models.WaitingForPlayers,
	}
}

// GameInInitialPromptCreationState describes a game where we're waiting for initial prompts
func GameInInitialPromptCreationState() *models.Game {
	// Just take the game in waiting for players state and move it to the next state
	game := GameInWaitingForPlayersState()
	game.CurrentState = models.InitialPromptCreation
	return game
}

// GameInDrawingsInProgressState describes a game where players are drawing
func GameInDrawingsInProgressState() *models.Game {
	game := GameInInitialPromptCreationState()
	playerPrompts := []*models.Prompt{
		models.BuildPrompt("chicken", []string{"snazzy", "portly"}, "player1"),
		models.BuildPrompt("tuna", []string{"big", "majestic"}, "player2"),
		models.BuildPrompt("boat", []string{"elegant", "sharp"}, "player3"),
	}
	for index, player := range game.Players {
		player.AssignedPrompt = playerPrompts[(index+2)%len(game.Players)]
	}
	game.OriginalPrompts = playerPrompts
	game.GeneratedPrompts = playerPrompts
	game.CurrentState = models.DrawingsInProgress
	return game
}

// DecoyPromptCreation describes a game where players are creating decoy prompts
func GameInDecoyPromptCreationState() *models.Game {
	game := GameInDrawingsInProgressState()
	mockImageData := "data:image/bmp;base64,Qk0eAAAAAAAAABoAAAAMAAAAAQABAAEAGAAAAP8A"
	drawings := []*models.Drawing{
		{
			ImageData:      mockImageData,
			Author:         "player2",
			DecoyPrompts:   map[string]*models.Prompt{},
			OriginalPrompt: game.GeneratedPrompts[0],
			Votes:          map[string]*models.Vote{},
		},
		{
			ImageData:      mockImageData,
			Author:         "player3",
			DecoyPrompts:   map[string]*models.Prompt{},
			OriginalPrompt: game.GeneratedPrompts[1],
			Votes:          map[string]*models.Vote{},
		},
		{
			ImageData:      mockImageData,
			Author:         "player1",
			DecoyPrompts:   map[string]*models.Prompt{},
			OriginalPrompt: game.GeneratedPrompts[2],
			Votes:          map[string]*models.Vote{},
		},
	}
	game.Drawings = drawings
	game.CurrentState = models.DecoyPromptCreation
	return game
}

// GameInVotingState describes a game where players are voting
func GameInVotingState() *models.Game {
	// Start from the decoy prompt creation state and just add decoy prompts to the active drawing
	game := GameInDecoyPromptCreationState()
	activeDrawing := game.GetActiveDrawing()
	activeDrawing.DecoyPrompts = map[string]*models.Prompt{
		"player1": models.BuildPrompt("toucan", []string{"happy", "big"}, "player1"),
		"player3": models.BuildPrompt("birb", []string{"jumpy", "edgy"}, "player1"),
	}
	game.CurrentState = models.Voting
	return game
}

// GameInScoringState describes a game where the round is being scored
func GameInScoringState() *models.Game {
	// Start off from the voting state and just add some votes to it
	game := GameInVotingState()
	game.CurrentState = models.Scoring
	activeDrawing := game.GetActiveDrawing()
	activeDrawing.Votes = map[string]*models.Vote{
		"player1": {Player: game.GetPlayer("player1"), SelectedPrompt: activeDrawing.OriginalPrompt},
		"player3": {Player: game.GetPlayer("player3"), SelectedPrompt: activeDrawing.DecoyPrompts["player1"]},
	}
	return game
}
