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

// GameInVotingState describes a game where players are voting
func GameInVotingState() *models.Game {
	prompts := []*models.Prompt{
		models.BuildPrompt("chicken", []string{"snazzy", "portly"}, "player1"),
		models.BuildPrompt("tuna", []string{"big", "majestic"}, "player2"),
		models.BuildPrompt("boat", []string{"elegant", "sharp"}, "player3"),
	}
	drawings := []*models.Drawing{
		{
			ImageData: "mockImage",
			Author:    "player2",
			DecoyPrompts: map[string]*models.Prompt{
				"player1": models.BuildPrompt("toucan", []string{"happy", "big"}, "player1"),
				"player3": models.BuildPrompt("birb", []string{"jumpy", "edgy"}, "player1"),
			},
			OriginalPrompt: prompts[0],
			Votes:          map[string]*models.Vote{},
		},
		{
			ImageData: "mockImage",
			Author:    "player3",
			DecoyPrompts: map[string]*models.Prompt{
				"player1": models.BuildPrompt("raft", []string{"happy", "big"}, "player1"),
				"player2": models.BuildPrompt("kayak", []string{"jumpy", "edgy"}, "player1"),
			},
			OriginalPrompt: prompts[1],
			Votes:          map[string]*models.Vote{},
		},
		{
			ImageData: "mockImage",
			Author:    "player1",
			DecoyPrompts: map[string]*models.Prompt{
				"player2": models.BuildPrompt("fish", []string{"happy", "big"}, "player1"),
				"player3": models.BuildPrompt("sardine", []string{"jumpy", "edgy"}, "player1"),
			},
			OriginalPrompt: prompts[2],
			Votes:          map[string]*models.Vote{},
		},
	}
	return &models.Game{
		GroupName: "castVoteRoute",
		Players: []*models.Player{
			{Name: "player1", Host: true, AssignedPrompt: prompts[2]},
			{Name: "player2", AssignedPrompt: prompts[0]},
			{Name: "player3", AssignedPrompt: prompts[1]},
		},
		Prompts:      prompts,
		CurrentState: models.Voting,
		Drawings:     drawings,
	}
}
