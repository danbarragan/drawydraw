package main

import (
	"bytes"
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
	resp := testHTTPResponse(t, req, http.StatusOK)

	assert.Equal(t, resp[groupNameKey], groupName)
	assert.Equal(t, resp[currentStateKey], waitingForPlayers)
	assert.Equal(t, resp[currentPlayerKey], Response{nameKey: hostName, isHostKey: true})
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
	testHTTPResponse(t, req, http.StatusOK)

	// Get the game's status
	req = createRequest(t, "GET", "/api/get-game-status/somegame?playerName=Player", nil)
	resp := testHTTPResponse(t, req, http.StatusOK)

	// Todo: Saner expected state
	assert.Equal(t, resp[groupNameKey], groupName)
	assert.Equal(t, resp[currentStateKey], waitingForPlayers)
	assert.Equal(t, resp[currentPlayerKey], Response{nameKey: hostName, isHostKey: true})
}

func TestCreateGameRoute__GameAlreadyExists(t *testing.T) {
	// Create the game
	data := map[string]string{
		"groupName":  "magic group",
		"playerName": "some player",
	}
	req := createRequest(t, "POST", "/api/create-game", data)
	testHTTPResponse(t, req, http.StatusOK)

	// Try to create the game again
	req = createRequest(t, "POST", "/api/create-game", data)
	testHTTPResponse(t, req, http.StatusBadRequest)
}

func TestAddPlayerRoute(t *testing.T) {
	// Create the game
	data := map[string]string{
		"groupName":  "group",
		"playerName": "player1",
	}
	req := createRequest(t, "POST", "/api/create-game", data)
	testHTTPResponse(t, req, http.StatusOK)

	// Add the player
	playerName := "player2"
	data = map[string]string{
		"groupName":  "group",
		"playerName": playerName,
	}
	req = createRequest(t, "POST", "/api/add-player", data)
	resp := testHTTPResponse(t, req, http.StatusOK)

	assert.Equal(t, resp[currentPlayerKey], Response{nameKey: playerName, isHostKey: false})
}

func TestAddPlayerRoute_GroupNotSetup(t *testing.T) {
	data := map[string]string{
		"groupName":  "superGroup",
		"playerName": "player",
	}

	req := createRequest(t, "POST", "/api/add-player", data)
	testHTTPResponse(t, req, http.StatusBadRequest)
}

func TestStartGameRoute(t *testing.T) {
	// Create the game
	data := map[string]string{
		"groupName":  "startGameRoute",
		"playerName": "player1",
	}
	req := createRequest(t, "POST", "/api/create-game", data)
	testHTTPResponse(t, req, http.StatusOK)

	// Actual test
	// Make a post to start game route
	data = map[string]string{
		"groupName":  "startGameRoute",
		"playerName": "player1",
	}
	req = createRequest(t, "POST", "/api/start-game", data)
	resp := testHTTPResponse(t, req, http.StatusOK)

	assert.Equal(t, resp[currentStateKey], initialPromptCreation)
}

// Helper function to process a request and test its response
func testHTTPResponse(t *testing.T, req *http.Request, statusCode int) Response {

	// Create a response recorder// Test set up
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	r := setupRouter("8080")
	r.ServeHTTP(w, req)

	if w.Code != statusCode {
		t.Fail()
	}

	resp := Response{}
	json.Unmarshal([]byte(w.Body.String()), &resp)
	return resp
}

func createRequest(t *testing.T, method string, route string, data map[string]string) *http.Request {
	jsonData, err := json.Marshal(&data)
	assert.Nil(t, err)
	req, err := http.NewRequest(method, route, bytes.NewBuffer(jsonData))
	assert.Nil(t, err)
	req.Header.Add("Content-Type", "application/json")
	return req
}

func TestAddPromptRoute(t *testing.T) {
	// Test set up
	router := setupRouter("8080")
	w := httptest.NewRecorder()
	data := map[string]string{
		"groupName":  "addPromptRoute",
		"playerName": "player1",
	}
	jsonData, err := json.Marshal(&data)

	// Create the group
	req, err := http.NewRequest("POST", "/api/create-game", bytes.NewBuffer(jsonData))
	assert.Nil(t, err)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Add another player
	w = httptest.NewRecorder()
	data = map[string]string{
		"groupName":  "addPromptRoute",
		"playerName": "player2",
	}
	jsonData, err = json.Marshal(&data)
	assert.Nil(t, err)
	req, err = http.NewRequest("POST", "/api/add-player", bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Start the game
	data = map[string]string{
		"groupName":  "addPromptRoute",
		"playerName": "player1",
	}
	w = httptest.NewRecorder()
	jsonData, err = json.Marshal(&data)
	req, err = http.NewRequest("POST", "/api/start-game", bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	//Make a post to the add prompts route from player 1, confirm state stays at "Initial Prompt Creation"
	data = map[string]string{
		"groupName":  "addPromptRoute",
		"playerName": "player1",
		"noun":       "chicken",
		"adjective1": "snazzy",
		"adjective2": "portly",
	}
	w = httptest.NewRecorder()
	jsonData, err = json.Marshal(&data)
	req, err = http.NewRequest("POST", "/api/add-prompts", bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	actualResponse := w.Body.String()
	expectedState := `{"currentPlayer":{"isHost":true,"name":"player1"},"currentState":"InitialPromptCreation","groupName":"addPromptRoute","players":[{"name":"player1","host":true,"points":0},{"name":"player2","host":false,"points":0}]}`
	assert.Equal(t, expectedState, actualResponse)

	//Make a post to the add prompts route from player 2, confirm state has changed to "Drawing"
	data = map[string]string{
		"groupName":  "addPromptRoute",
		"playerName": "player2",
		"noun":       "duck",
		"adjective1": "chilly",
		"adjective2": "sleepy",
	}
	w = httptest.NewRecorder()
	jsonData, err = json.Marshal(&data)
	req, err = http.NewRequest("POST", "/api/add-prompts", bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	actualResponse = w.Body.String()
	expectedState = `{"currentPlayer":{"isHost":false,"name":"player2"},"currentState":"DrawingsInProgress","groupName":"addPromptRoute","players":[{"name":"player1","host":true,"points":0},{"name":"player2","host":false,"points":0}]}`
	assert.Equal(t, expectedState, actualResponse)
}
