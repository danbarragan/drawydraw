package test

import (
	"drawydraw/models"
)

// This file contains game models at different states to be used in testing

// GameInWaitingForPlayersState describes a newly created game waiting for players
func GameInWaitingForPlayersState() *models.Game {
	return &models.Game{
		GroupName:    "somegame",
		Players:      []*models.Player{{Name: "Player", Host: true}},
		CurrentState: models.WaitingForPlayers,
	}
}

// GameInInitialPromptCreationState describes a game where we're waiting for initial prompts
func GameInInitialPromptCreationState() *models.Game {
	return &models.Game{
		GroupName: "addPromptRoute",
		Players: []*models.Player{
			{Name: "player1", Host: true},
			{Name: "player2"},
		},
		CurrentState: models.InitialPromptCreation,
	}
}

// GameInDrawingsInProgressState describes a game where players are drawing
func GameInDrawingsInProgressState() *models.Game {
	player1Prompt := &models.Prompt{
		Noun:       "chicken",
		Adjectives: []string{"snazzy", "portly"}, Author: "player1",
	}
	player2Prompt := &models.Prompt{
		Noun:       "tuna",
		Adjectives: []string{"snazzy", "portly"}, Author: "player2",
	}
	return &models.Game{
		GroupName: "submitDrawingRoute",
		Players: []*models.Player{
			{Name: "player1", Host: true, AssignedPrompt: player2Prompt},
			{Name: "player2", AssignedPrompt: player1Prompt},
		},
		Prompts:      []*models.Prompt{player1Prompt, player2Prompt},
		CurrentState: models.DrawingsInProgress,
	}
}
