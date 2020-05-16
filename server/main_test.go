package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"os"

	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"
)

/*
	Several problems with these tests:
	- They're really more like integration tests, we should improve abtractions sometime and have real unit tests
	- Asserting on a JSON string is error prone, we should probably have mock game states and deserialize JSON responses to compare against them
	- There's a lot of code repetition here and some shortcuts could be taken like setting the game state directly instead of going through different endpoints
*/

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
	data := map[string]string{
		"groupName":  "Kitten Party",
		"playerName": "Baby Cat",
	}
	req := createRequest(t, "POST", "/api/create-game", data)
	resp := testHTTPResponse(t, req, http.StatusOK)

	// Todo: Saner expected state
	expectedState := `{"currentPlayer":{"isHost":true,"name":"Baby Cat"},"currentState":"WaitingForPlayers","groupName":"Kitten Party","players":[{"name":"Baby Cat","host":true,"points":0}]}`
	assert.Equal(t, expectedState, resp)
}

func TestGetGameStateStatusRoute(t *testing.T) {
	// Create a game
	data := map[string]string{
		"groupName":  "somegame",
		"playerName": "Player",
	}
	req := createRequest(t, "POST", "/api/create-game", data)
	testHTTPResponse(t, req, http.StatusOK)

	// Get the game's status
	req = createRequest(t, "GET", "/api/get-game-status/somegame?playerName=Player", nil)
	resp := testHTTPResponse(t, req, http.StatusOK)

	// Todo: Saner expected state
	expectedState := `{"currentPlayer":{"isHost":true,"name":"Player"},"currentState":"WaitingForPlayers","groupName":"somegame","players":[{"name":"Player","host":true,"points":0}]}`
	assert.Equal(t, expectedState, resp)
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
	data = map[string]string{
		"groupName":  "group",
		"playerName": "player2",
	}
	req = createRequest(t, "POST", "/api/add-player", data)
	resp := testHTTPResponse(t, req, http.StatusOK)

	expectedState := `{"currentPlayer":{"isHost":false,"name":"player2"},"currentState":"WaitingForPlayers","groupName":"group","players":[{"name":"player1","host":true,"points":0},{"name":"player2","host":false,"points":0}]}`
	assert.Equal(t, expectedState, resp)
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

	expectedState := `{"currentPlayer":{"isHost":true,"name":"player1"},"currentState":"InitialPromptCreation","groupName":"startGameRoute","players":[{"name":"player1","host":true,"points":0}]}`
	assert.Equal(t, expectedState, resp)
}

// Helper function to process a request and test its response
func testHTTPResponse(t *testing.T, req *http.Request, statusCode int) string {

  // Create a response recorder// Test set up
  w := httptest.NewRecorder()

  // Create the service and process the above request.
  r := setupRouter("8080")
  r.ServeHTTP(w, req)

  if w.Code != statusCode {
    t.Fail()
  }

  return w.Body.String()
}

func createRequest(t *testing.T, method string, route string, data map[string]string) *http.Request {
	jsonData, err := json.Marshal(&data)
	assert.Nil(t, err)
	req, err := http.NewRequest(method, route, bytes.NewBuffer(jsonData))
	assert.Nil(t, err)
	req.Header.Add("Content-Type", "application/json")
	return req
}
