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
	req := createRequest(t, "GET", "/api/get-game-status/somegame?playerName=Player", nil)
	actualGameState := sendRequest(t, req, http.StatusOK)
	expectedGameState := &statemanager.GameStatusResponse{
		GroupName:     "somegame",
		CurrentPlayer: &statemanager.CurrentPlayer{Name: "Player", IsHost: true},
		CurrentState:  string(models.WaitingForPlayers),
		Players: []*statemanager.Player{
			{Name: "Player", Host: true},
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
	models.GetGameProvider().SaveGame(test.GameInWaitingForPlayersState())
	// Add the player
	playerName := "player2"
	data := map[string]string{
		"groupName":  "somegame",
		"playerName": playerName,
	}
	req := createRequest(t, "POST", "/api/add-player", data)
	actualGameState := sendRequest(t, req, http.StatusOK)
	expectedGameState := &statemanager.GameStatusResponse{
		GroupName:     "somegame",
		CurrentPlayer: &statemanager.CurrentPlayer{Name: "player2"},
		CurrentState:  string(models.WaitingForPlayers),
		Players: []*statemanager.Player{
			{Name: "Player", Host: true},
			{Name: "player2"},
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
	models.GetGameProvider().SaveGame(test.GameInWaitingForPlayersState())
	data := map[string]string{
		"groupName":  "somegame",
		"playerName": "Player",
	}
	req := createRequest(t, "POST", "/api/start-game", data)
	actualGameState := sendRequest(t, req, http.StatusOK)
	expectedGameState := &statemanager.GameStatusResponse{
		GroupName:     "somegame",
		CurrentPlayer: &statemanager.CurrentPlayer{IsHost: true, Name: "Player"},
		CurrentState:  string(models.InitialPromptCreation),
		Players: []*statemanager.Player{
			{Name: "Player", Host: true, HasPendingAction: true},
		},
	}
	assert.EqualValues(t, expectedGameState, actualGameState)
}

func TestAddPromptRoute(t *testing.T) {
	test.SetupTestGameProvider(t)
	models.GetGameProvider().SaveGame(test.GameInInitialPromptCreationState())
	//Make a post to the add prompts route from player 1, confirm state stays at "Initial Prompt Creation"
	data := map[string]string{
		"groupName":  "addPromptRoute",
		"playerName": "player1",
		"noun":       "chicken",
		"adjective1": "snazzy",
		"adjective2": "portly",
	}
	req := createRequest(t, "POST", "/api/add-prompt", data)
	actualGameState := sendRequest(t, req, http.StatusOK)
	expectedGameState := &statemanager.GameStatusResponse{
		GroupName: "addPromptRoute",
		CurrentPlayer: &statemanager.CurrentPlayer{
			IsHost:             true,
			Name:               "player1",
			HasCompletedAction: true,
		},
		CurrentState: string(models.InitialPromptCreation),
		Players: []*statemanager.Player{
			{Name: "player1", Host: true, HasPendingAction: false},
			{Name: "player2", HasPendingAction: true},
		},
	}
	assert.EqualValues(t, expectedGameState, actualGameState)
}

func TestAddPromptRoute_AssignsPrompts(t *testing.T) {
	test.SetupTestGameProvider(t)
	gameState := test.GameInInitialPromptCreationState()
	// Add the prompt from player 1
	gameState.Prompts = []*models.Prompt{{
		Noun:       "chicken",
		Adjectives: []string{"snazzy", "portly"}, Author: "player1",
	}}
	models.GetGameProvider().SaveGame(gameState)
	//Make a post to the add prompts route from player 2, state should transition to drawings in progress
	data := map[string]string{
		"groupName":  "addPromptRoute",
		"playerName": "player2",
		"noun":       "orangutan",
		"adjective1": "fiery",
		"adjective2": "friendly",
	}
	req := createRequest(t, "POST", "/api/add-prompt", data)
	actualGameState := sendRequest(t, req, http.StatusOK)
	expectedGameState := &statemanager.GameStatusResponse{
		GroupName: "addPromptRoute",
		CurrentPlayer: &statemanager.CurrentPlayer{
			Name: "player2",
			AssignedPrompt: &statemanager.Prompt{
				Noun:       "chicken",
				Adjectives: []string{"snazzy", "portly"},
			},
		},
		CurrentState: string(models.DrawingsInProgress),
		Players: []*statemanager.Player{
			{Name: "player1", Host: true, HasPendingAction: true},
			{Name: "player2", HasPendingAction: true},
		},
	}
	assert.EqualValues(t, expectedGameState, actualGameState)
}

func TestSubmitDrawingRoute(t *testing.T) {
	test.SetupTestGameProvider(t)
	models.GetGameProvider().SaveGame(test.GameInDrawingsInProgressState())
	// Submit a drawing
	data := map[string]string{
		"groupName":  "submitDrawingRoute",
		"playerName": "player1",
		"imageData":  "someImageData",
	}
	req := createRequest(t, "POST", "/api/submit-drawing", data)
	actualGameState := sendRequest(t, req, http.StatusOK)
	expectedGameState := &statemanager.GameStatusResponse{
		GroupName: "submitDrawingRoute",
		CurrentPlayer: &statemanager.CurrentPlayer{
			IsHost:             true,
			Name:               "player1",
			HasCompletedAction: true,
			AssignedPrompt: &statemanager.Prompt{
				Noun:       "tuna",
				Adjectives: []string{"snazzy", "portly"},
			},
		},
		CurrentState: string(models.DrawingsInProgress),
		Players: []*statemanager.Player{
			{Name: "player1", Host: true, HasPendingAction: false},
			{Name: "player2", HasPendingAction: true},
		},
	}
	assert.EqualValues(t, expectedGameState, actualGameState)
}

func TestCastVoteRoute(t *testing.T) {
	test.SetupTestGameProvider(t)
	models.GetGameProvider().SaveGame(test.GameInVotingState())
	// Cast a vote
	data := map[string]string{
		"groupName":        "castVoteRoute",
		"playerName":       "player1",
		"selectedPromptId": "7876445554424581103",
	}
	req := createRequest(t, "POST", "/api/cast-vote", data)
	actualGameState := sendRequest(t, req, http.StatusOK)
	expectedGameState := &statemanager.GameStatusResponse{
		GroupName: "castVoteRoute",
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
