package main

import (
	"bytes"
	"drawydraw/models"
	"drawydraw/statemanager"
	"drawydraw/test"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// This function is used for setup before executing the test functions
func TestMain(m *testing.M) {
	//Set Gin to Test Mode
	gin.SetMode(gin.TestMode)
	// Setup here
	status := m.Run()
	// Cleanup here
	os.Exit(status)
}

func TestCreateGameRoute(t *testing.T) {
	test.SetupTestGameProvider(t)
	// Create the game
	hostName := "Baby Cat"
	groupName := "Kitten Party"
	data := map[string]string{
		"groupName":  groupName,
		"playerName": hostName,
	}
	req := createRequest(t, "POST", "/api/create-game", data)
	actualGameState := sendRequest(t, req, http.StatusOK)
	expectedGameState := &statemanager.GameStatusResponse{
		GroupName:     "Kitten Party",
		CurrentPlayer: &statemanager.CurrentPlayer{Name: "Baby Cat", IsHost: true},
		CurrentState:  string(models.WaitingForPlayers),
		Players: []*statemanager.Player{
			{Name: "Baby Cat", Host: true},
		},
	}
	assert.EqualValues(t, expectedGameState, actualGameState)
}

func TestGetGameStateStatusRoute(t *testing.T) {
	test.SetupTestGameProvider(t)
	models.GetGameProvider().SaveGame(test.GameInWaitingForPlayersState())
	req := createRequest(t, "GET", "/api/get-game-status/somegame?playerName=player1", nil)
	actualGameState := sendRequest(t, req, http.StatusOK)
	expectedGameState := &statemanager.GameStatusResponse{
		GroupName:     "somegame",
		CurrentPlayer: &statemanager.CurrentPlayer{Name: "player1", IsHost: true},
		CurrentState:  string(models.WaitingForPlayers),
		Players: []*statemanager.Player{
			{Name: "player1", Host: true},
			{Name: "player2"},
			{Name: "player3"},
		},
	}
	assert.EqualValues(t, expectedGameState, actualGameState)
}

func TestCreateGameRoute__GameAlreadyExists(t *testing.T) {
	test.SetupTestGameProvider(t)
	models.GetGameProvider().SaveGame(test.GameInWaitingForPlayersState())
	data := map[string]string{
		"groupName":  "somegame",
		"playerName": "some player",
	}
	req := createRequest(t, "POST", "/api/create-game", data)
	sendRequest(t, req, http.StatusBadRequest)
}

func TestAddPlayerRoute(t *testing.T) {
	test.SetupTestGameProvider(t)
	game := test.GameInWaitingForPlayersState()
	models.GetGameProvider().SaveGame(game)
	// Add the player
	data := map[string]string{
		"groupName":  game.GroupName,
		"playerName": "player4",
	}
	req := createRequest(t, "POST", "/api/add-player", data)
	actualGameState := sendRequest(t, req, http.StatusOK)
	expectedGameState := &statemanager.GameStatusResponse{
		GroupName:     "somegame",
		CurrentPlayer: &statemanager.CurrentPlayer{Name: "player4"},
		CurrentState:  string(models.WaitingForPlayers),
		Players: []*statemanager.Player{
			{Name: "player1", Host: true},
			{Name: "player2"},
			{Name: "player3"},
			{Name: "player4"},
		},
	}
	assert.EqualValues(t, expectedGameState, actualGameState)
}

func TestAddPlayerRoute_GroupNotSetup(t *testing.T) {
	test.SetupTestGameProvider(t)
	data := map[string]string{
		"groupName":  "superGroup",
		"playerName": "player",
	}

	req := createRequest(t, "POST", "/api/add-player", data)
	sendRequest(t, req, http.StatusBadRequest)
}

func TestStartGameRoute(t *testing.T) {
	test.SetupTestGameProvider(t)
	game := test.GameInWaitingForPlayersState()
	models.GetGameProvider().SaveGame(game)
	data := map[string]string{
		"groupName":  game.GroupName,
		"playerName": "player1",
	}
	req := createRequest(t, "POST", "/api/start-game", data)
	actualGameState := sendRequest(t, req, http.StatusOK)
	expectedGameState := &statemanager.GameStatusResponse{
		GroupName:     game.GroupName,
		CurrentPlayer: &statemanager.CurrentPlayer{IsHost: true, Name: "player1"},
		CurrentState:  string(models.InitialPromptCreation),
		Players: []*statemanager.Player{
			{Name: "player1", Host: true, HasPendingAction: true},
			{Name: "player2", HasPendingAction: true},
			{Name: "player3", HasPendingAction: true},
		},
	}
	assert.EqualValues(t, expectedGameState, actualGameState)
}

func TestAddPromptRoute(t *testing.T) {
	test.SetupTestGameProvider(t)
	game := test.GameInInitialPromptCreationState()
	models.GetGameProvider().SaveGame(game)
	//Make a post to the add prompts route from player 1, confirm state stays at "Initial Prompt Creation"
	data := map[string]string{
		"groupName":  game.GroupName,
		"playerName": "player1",
		"noun":       "chicken",
		"adjective1": "snazzy",
		"adjective2": "portly",
	}
	req := createRequest(t, "POST", "/api/add-prompt", data)
	actualGameState := sendRequest(t, req, http.StatusOK)
	expectedGameState := &statemanager.GameStatusResponse{
		GroupName: game.GroupName,
		CurrentPlayer: &statemanager.CurrentPlayer{
			IsHost:             true,
			Name:               "player1",
			HasCompletedAction: true,
		},
		CurrentState: string(models.InitialPromptCreation),
		Players: []*statemanager.Player{
			{Name: "player1", Host: true, HasPendingAction: false},
			{Name: "player2", HasPendingAction: true},
			{Name: "player3", HasPendingAction: true},
		},
	}
	assert.EqualValues(t, expectedGameState, actualGameState)
}

func TestAddPromptRoute_AssignsPrompts(t *testing.T) {
	test.SetupTestGameProvider(t)
	game := test.GameInInitialPromptCreationState()
	// Add the prompt from player 1
	game.OriginalPrompts = []*models.Prompt{
		models.BuildPrompt("chicken", []string{"snazzy", "portly"}, "player1"),
		models.BuildPrompt("tuna", []string{"big", "majestic"}, "player2"),
	}
	models.GetGameProvider().SaveGame(game)
	//Make a post to the add prompts route from player 2, state should transition to drawings in progress
	data := map[string]string{
		"groupName":  game.GroupName,
		"playerName": "player3",
		"noun":       "orangutan",
		"adjective1": "fiery",
		"adjective2": "friendly",
	}
	req := createRequest(t, "POST", "/api/add-prompt", data)
	actualGameState := sendRequest(t, req, http.StatusOK)
	expectedGameState := &statemanager.GameStatusResponse{
		GroupName: game.GroupName,
		CurrentPlayer: &statemanager.CurrentPlayer{
			Name: "player3",
			AssignedPrompt: &statemanager.Prompt{
				Noun:       "chicken",
				Adjectives: []string{"snazzy", "portly"},
			},
		},
		CurrentState: string(models.DrawingsInProgress),
		Players: []*statemanager.Player{
			{Name: "player1", Host: true, HasPendingAction: true},
			{Name: "player2", HasPendingAction: true},
			{Name: "player3", HasPendingAction: true},
		},
	}
	// Since prompt assignment is random, check that the adjectives assigned are in the original list
	allAdjectives := map[string]bool{"snazzy": true, "portly": true, "fiery": true, "friendly": true, "big": true, "majestic": true}
	for _, adjective := range actualGameState.CurrentPlayer.AssignedPrompt.Adjectives {
		assert.True(t, allAdjectives[adjective])
	}
	// The same adjective should not be picked twice
	assert.NotEqual(
		t,
		actualGameState.CurrentPlayer.AssignedPrompt.Adjectives[0],
		actualGameState.CurrentPlayer.AssignedPrompt.Adjectives[1],
	)
	// Since we checked the adjectives let's just ignore those in the comparison with the expected state
	expectedGameState.CurrentPlayer.AssignedPrompt.Adjectives = actualGameState.CurrentPlayer.AssignedPrompt.Adjectives
	assert.EqualValues(t, expectedGameState, actualGameState)
}

func TestSubmitDrawingRoute(t *testing.T) {
	test.SetupTestGameProvider(t)
	game := test.GameInDrawingsInProgressState()
	models.GetGameProvider().SaveGame(game)
	// Submit a drawing
	data := map[string]string{
		"groupName":  game.GroupName,
		"playerName": "player1",
		"imageData":  "someImageData",
	}
	req := createRequest(t, "POST", "/api/submit-drawing", data)
	actualGameState := sendRequest(t, req, http.StatusOK)
	expectedGameState := &statemanager.GameStatusResponse{
		GroupName: game.GroupName,
		CurrentPlayer: &statemanager.CurrentPlayer{
			IsHost:             true,
			Name:               "player1",
			HasCompletedAction: true,
			AssignedPrompt: &statemanager.Prompt{
				Noun:       "boat",
				Adjectives: []string{"elegant", "sharp"},
			},
		},
		CurrentState: string(models.DrawingsInProgress),
		Players: []*statemanager.Player{
			{Name: "player1", Host: true, HasPendingAction: false},
			{Name: "player2", HasPendingAction: true},
			{Name: "player3", HasPendingAction: true},
		},
	}
	assert.EqualValues(t, expectedGameState, actualGameState)
}

func TestCastVoteRoute(t *testing.T) {
	test.SetupTestGameProvider(t)
	game := test.GameInVotingState()
	models.GetGameProvider().SaveGame(game)
	// Cast a vote
	data := map[string]string{
		"groupName":        game.GroupName,
		"playerName":       "player1",
		"selectedPromptId": "7876445554424581103",
	}
	req := createRequest(t, "POST", "/api/cast-vote", data)
	actualGameState := sendRequest(t, req, http.StatusOK)
	expectedGameState := &statemanager.GameStatusResponse{
		GroupName: game.GroupName,
		CurrentPlayer: &statemanager.CurrentPlayer{
			IsHost:             true,
			Name:               "player1",
			HasCompletedAction: true,
		},
		CurrentState: string(models.Voting),
		Players: []*statemanager.Player{
			{Name: "player1", Host: true, HasPendingAction: false},
			{Name: "player2", HasPendingAction: false},
			{Name: "player3", HasPendingAction: true},
		},
		CurrentDrawing: &statemanager.Drawing{
			ImageData: "data:image/bmp;base64,Qk0eAAAAAAAAABoAAAAMAAAAAQABAAEAGAAAAP8A",
			Prompts: []*statemanager.Prompt{
				{Identifier: "2289583145965790902", Noun: "birb", Adjectives: []string{"jumpy", "edgy"}},
				{Identifier: "7876445554424581103", Noun: "chicken", Adjectives: []string{"snazzy", "portly"}},
				{Identifier: "9033667170926423839", Noun: "toucan", Adjectives: []string{"happy", "big"}},
			},
		},
	}
	assert.EqualValues(t, expectedGameState, actualGameState)
}

// Helper function to process a request and test its response
func sendRequest(t *testing.T, req *http.Request, statusCode int) *statemanager.GameStatusResponse {
	// Create a response recorder// Test set up
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	r := setupRouter("8080")
	r.ServeHTTP(w, req)

	if w.Code != statusCode {
		t.Fail()
	}

	actualGameState := &statemanager.GameStatusResponse{}
	json.Unmarshal([]byte(w.Body.String()), &actualGameState)
	return actualGameState
}

func createRequest(t *testing.T, method string, route string, data map[string]string) *http.Request {
	jsonData, err := json.Marshal(&data)
	assert.Nil(t, err)
	req, err := http.NewRequest(method, route, bytes.NewBuffer(jsonData))
	assert.Nil(t, err)
	req.Header.Add("Content-Type", "application/json")
	return req
}
