package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
	Several problems with these tests:
	- They're really more like integration tests, we should improve abtractions sometime and have real unit tests
	- Asserting on a JSON string is error prone, we should probably have mock game states and deserialize JSON responses to compare against them
	- There's a lot of code repetition here and some shortcuts could be taken like setting the game state directly instead of going through different endpoints
*/

func TestCreateGameRoute(t *testing.T) {
	router := setupRouter("8080")
	w := httptest.NewRecorder()
	data := map[string]string{
		"groupName":  "Kitten Party",
		"playerName": "Baby Cat",
	}
	jsonData, err := json.Marshal(&data)
	assert.Nil(t, err)
	req, err := http.NewRequest("POST", "/api/create-game", bytes.NewBuffer(jsonData))
	assert.Nil(t, err)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	// Todo: Saner expected state
	expectedState := `{"currentPlayer":{"isHost":true,"name":"Baby Cat"},"currentState":"WaitingForPlayers","groupName":"Kitten Party","players":[{"name":"Baby Cat","host":true,"points":0}]}`
	assert.Equal(t, expectedState, w.Body.String())
}

func TestGetGameStateStatusRoute(t *testing.T) {
	// Create a game
	router := setupRouter("8080")
	w := httptest.NewRecorder()
	data := map[string]string{
		"groupName":  "somegame",
		"playerName": "Player",
	}
	jsonData, err := json.Marshal(&data)
	assert.Nil(t, err)
	req, err := http.NewRequest("POST", "/api/create-game", bytes.NewBuffer(jsonData))
	assert.Nil(t, err)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Get the game's status
	req, err = http.NewRequest("GET", "/api/get-game-status/somegame?playerName=Player", nil)
	assert.Nil(t, err)
	req.Header.Add("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	// Todo: Saner expected state
	expectedState := `{"currentPlayer":{"isHost":true,"name":"Player"},"currentState":"WaitingForPlayers","groupName":"somegame","players":[{"name":"Player","host":true,"points":0}]}`
	assert.Equal(t, expectedState, w.Body.String())
}

func TestCreateGameRoute__GameAlreadyExists(t *testing.T) {
	router := setupRouter("8080")
	w := httptest.NewRecorder()
	data := map[string]string{
		"groupName":  "magic group",
		"playerName": "some player",
	}
	jsonData, err := json.Marshal(&data)
	assert.Nil(t, err)
	// Create the game
	req, err := http.NewRequest("POST", "/api/create-game", bytes.NewBuffer(jsonData))
	assert.Nil(t, err)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Try to create the game again
	req, err = http.NewRequest("POST", "/api/create-game", bytes.NewBuffer(jsonData))
	assert.Nil(t, err)
	req.Header.Add("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddPlayerRoute(t *testing.T) {
	router := setupRouter("8080")
	w := httptest.NewRecorder()
	data := map[string]string{
		"groupName":  "group",
		"playerName": "player1",
	}
	jsonData, err := json.Marshal(&data)

	// Create the group
	req, err := http.NewRequest("POST", "/api/create-game", bytes.NewBuffer(jsonData))
	assert.Nil(t, err)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Add the player
	w = httptest.NewRecorder()
	data = map[string]string{
		"groupName":  "group",
		"playerName": "player2",
	}
	jsonData, err = json.Marshal(&data)
	req, err = http.NewRequest("POST", "/api/add-player", bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	expectedState := `{"currentPlayer":{"isHost":false,"name":"player2"},"currentState":"WaitingForPlayers","groupName":"group","players":[{"name":"player1","host":true,"points":0},{"name":"player2","host":false,"points":0}]}`
	assert.Equal(t, expectedState, w.Body.String())
}

func TestAddPlayerRoute_GroupNotSetup(t *testing.T) {
	router := setupRouter("8080")
	w := httptest.NewRecorder()
	data := map[string]string{
		"groupName":  "superGroup",
		"playerName": "player",
	}
	jsonData, err := json.Marshal(&data)
	assert.Nil(t, err)
	req, err := http.NewRequest("POST", "/api/add-player", bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStartGameRoute(t *testing.T) {
	// Test set up
	router := setupRouter("8080")
	w := httptest.NewRecorder()
	data := map[string]string{
		"groupName":  "startGameRoute",
		"playerName": "player1",
	}
	jsonData, err := json.Marshal(&data)

	// Create the group
	req, err := http.NewRequest("POST", "/api/create-game", bytes.NewBuffer(jsonData))
	assert.Nil(t, err)
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Actual test
	// Make a post to start game route
	data = map[string]string{
		"groupName":  "startGameRoute",
		"playerName": "player1",
	}
	w = httptest.NewRecorder()
	jsonData, err = json.Marshal(&data)
	req, err = http.NewRequest("POST", "/api/start-game", bytes.NewBuffer(jsonData))
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	expectedState := `{"currentPlayer":{"isHost":true,"name":"player1"},"currentState":"InitialPromptCreation","groupName":"startGameRoute","players":[{"name":"player1","host":true,"points":0}],"prompts":null}`
	actualResponse := w.Body.String()
	assert.Equal(t, expectedState, actualResponse)
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
	//prompt_uuid := actualResponse  //TODO : will need to get the real uuid from the the response and put it in the expected state
	expectedState := `{"currentPlayer":{"isHost":true,"name":"player1"},"currentState":"InitialPromptCreation","groupName":"addPromptRoute","players":[{"name":"player1","host":true,"points":0},{"name":"player2","host":false,"points":0}],"prompts":[{"Author":"player1","Group":"addPromptRoute","Noun":"chicken","Adjective1":"snazzy","Adjective2":"portly"}]}`
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
	//prompt_uuid := actualResponse  //TODO : will need to get the real uuid from the the response and put it in the expected state
	expectedState = `{"currentPlayer":{"isHost":false,"name":"player2"},"currentState":"Drawing","groupName":"addPromptRoute","players":[{"name":"player1","host":true,"points":0},{"name":"player2","host":false,"points":0}],"prompts":[{"Author":"player1","Group":"addPromptRoute","Noun":"chicken","Adjective1":"snazzy","Adjective2":"portly"},{"Author":"player2","Group":"addPromptRoute","Noun":"duck","Adjective1":"chilly","Adjective2":"sleepy"}]}`
	assert.Equal(t, expectedState, actualResponse)

}
