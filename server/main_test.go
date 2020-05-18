package main

import (
	"bytes"
	"drawydraw/models"
	"drawydraw/utils/statemanager"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var groupNameKey = "groupName"
var playersKey = "players"
var currentPlayerKey = "currentPlayer"
var currentStateKey = "currentState"
var waitingForPlayers = "WaitingForPlayers"
var initialPromptCreation = "InitialPromptCreation"
var drawingsInProgress = "DrawingsInProgress"
var nameKey = "name"
var isHostKey = "isHost"

type Response = map[string]interface{}

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
	// Create a game
	hostName := "Player"
	groupName := "somegame"
	data := map[string]string{
		"groupName":  groupName,
		"playerName": hostName,
	}
	req := createRequest(t, "POST", "/api/create-game", data)
	sendRequest(t, req, http.StatusOK)

	// Get the game's status
	req = createRequest(t, "GET", "/api/get-game-status/somegame?playerName=Player", nil)
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
	// Create the game
	data := map[string]string{
		"groupName":  "magic group",
		"playerName": "some player",
	}
	req := createRequest(t, "POST", "/api/create-game", data)
	sendRequest(t, req, http.StatusOK)

	// Try to create the game again
	req = createRequest(t, "POST", "/api/create-game", data)
	sendRequest(t, req, http.StatusBadRequest)
}

func TestAddPlayerRoute(t *testing.T) {
	// Create the game
	data := map[string]string{
		"groupName":  "group",
		"playerName": "player1",
	}
	req := createRequest(t, "POST", "/api/create-game", data)
	sendRequest(t, req, http.StatusOK)

	// Add the player
	playerName := "player2"
	data = map[string]string{
		"groupName":  "group",
		"playerName": playerName,
	}
	req = createRequest(t, "POST", "/api/add-player", data)
	actualGameState := sendRequest(t, req, http.StatusOK)
	expectedGameState := &statemanager.GameStatusResponse{
		GroupName:     "group",
		CurrentPlayer: &statemanager.CurrentPlayer{Name: "player2"},
		CurrentState:  string(models.WaitingForPlayers),
		Players: []*statemanager.Player{
			{Name: "player1", Host: true},
			{Name: "player2"},
		},
	}
	assert.EqualValues(t, expectedGameState, actualGameState)
}

func TestAddPlayerRoute_GroupNotSetup(t *testing.T) {
	data := map[string]string{
		"groupName":  "superGroup",
		"playerName": "player",
	}

	req := createRequest(t, "POST", "/api/add-player", data)
	sendRequest(t, req, http.StatusBadRequest)
}

func TestStartGameRoute(t *testing.T) {
	// Create the game
	data := map[string]string{
		"groupName":  "startGameRoute",
		"playerName": "player1",
	}
	req := createRequest(t, "POST", "/api/create-game", data)
	sendRequest(t, req, http.StatusOK)

	// Actual test
	// Make a post to start game route
	data = map[string]string{
		"groupName":  "startGameRoute",
		"playerName": "player1",
	}
	req = createRequest(t, "POST", "/api/start-game", data)
	actualGameState := sendRequest(t, req, http.StatusOK)
	expectedGameState := &statemanager.GameStatusResponse{
		GroupName:     "startGameRoute",
		CurrentPlayer: &statemanager.CurrentPlayer{IsHost: true, Name: "player1"},
		CurrentState:  string(models.InitialPromptCreation),
		Players: []*statemanager.Player{
			{Name: "player1", Host: true},
		},
	}
	assert.EqualValues(t, expectedGameState, actualGameState)
}

func TestAddPromptRoute(t *testing.T) {
	// Test set up
	data := map[string]string{
		"groupName":  "addPromptRoute",
		"playerName": "player1",
	}
	req := createRequest(t, "POST", "/api/create-game", data)
	sendRequest(t, req, http.StatusOK)

	// Add another player
	data = map[string]string{
		"groupName":  "addPromptRoute",
		"playerName": "player2",
	}
	req = createRequest(t, "POST", "/api/add-player", data)
	sendRequest(t, req, http.StatusOK)

	// Start the game
	data = map[string]string{
		"groupName":  "addPromptRoute",
		"playerName": "player1",
	}
	req = createRequest(t, "POST", "/api/start-game", data)
	sendRequest(t, req, http.StatusOK)

	//Make a post to the add prompts route from player 1, confirm state stays at "Initial Prompt Creation"
	data = map[string]string{
		"groupName":  "addPromptRoute",
		"playerName": "player1",
		"noun":       "chicken",
		"adjective1": "snazzy",
		"adjective2": "portly",
	}
	req = createRequest(t, "POST", "/api/add-prompt", data)
	actualGameState := sendRequest(t, req, http.StatusOK)
	expectedGameState := &statemanager.GameStatusResponse{
		GroupName:     "addPromptRoute",
		CurrentPlayer: &statemanager.CurrentPlayer{IsHost: true, Name: "player1"},
		CurrentState:  string(models.InitialPromptCreation),
		Players: []*statemanager.Player{
			{Name: "player1", Host: true},
			{Name: "player2"},
		},
	}
	assert.EqualValues(t, expectedGameState, actualGameState)
}

func TestSubmitDrawingRoute(t *testing.T) {
	// Test set up
	data := map[string]string{
		"groupName":  "submitDrawingRoute",
		"playerName": "player1",
	}
	req := createRequest(t, "POST", "/api/create-game", data)
	sendRequest(t, req, http.StatusOK)

	// Add another player
	data = map[string]string{
		"groupName":  "submitDrawingRoute",
		"playerName": "player2",
	}
	req = createRequest(t, "POST", "/api/add-player", data)
	sendRequest(t, req, http.StatusOK)

	// Start the game
	data = map[string]string{
		"groupName":  "submitDrawingRoute",
		"playerName": "player1",
	}
	req = createRequest(t, "POST", "/api/start-game", data)
	sendRequest(t, req, http.StatusOK)

	//Make a post to the add prompts route from player 1, confirm state stays at "Initial Prompt Creation"
	data = map[string]string{
		"groupName":  "submitDrawingRoute",
		"playerName": "player1",
		"noun":       "chicken",
		"adjective1": "snazzy",
		"adjective2": "portly",
	}
	req = createRequest(t, "POST", "/api/add-prompt", data)
	sendRequest(t, req, http.StatusOK)
	data["playerName"] = "player2"
	data["noun"] = "tuna"
	req = createRequest(t, "POST", "/api/add-prompt", data)
	sendRequest(t, req, http.StatusOK)

	// Submit a drawing
	data = map[string]string{
		"groupName":  "submitDrawingRoute",
		"playerName": "player1",
		"imageData":  "someImageData",
	}
	req = createRequest(t, "POST", "/api/submit-drawing", data)
	actualGameState := sendRequest(t, req, http.StatusOK)
	expectedGameState := &statemanager.GameStatusResponse{
		GroupName: "submitDrawingRoute",
		CurrentPlayer: &statemanager.CurrentPlayer{
			IsHost:             true,
			Name:               "player1",
			HasCompletedAction: true,
			AssignedPrompt: &statemanager.AssignedPrompt{
				Noun:       "tuna",
				Adjectives: []string{"snazzy", "portly"},
			},
		},
		CurrentState: string(models.DrawingsInProgress),
		Players: []*statemanager.Player{
			{Name: "player1", Host: true, HasPendingActions: false},
			{Name: "player2", HasPendingActions: true},
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
