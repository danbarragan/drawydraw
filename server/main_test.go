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
	expectedState := `{"groupName":"Kitten Party","players":[{"name":"Baby Cat","host":true,"points":0}],"currentState":"WaitingForPlayers"}`
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
	req, err = http.NewRequest("GET", "/api/get-game-status/somegame", nil)
	assert.Nil(t, err)
	req.Header.Add("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	// Todo: Saner expected state
	expectedState := `{"groupName":"somegame","players":[{"name":"Player","host":true,"points":0}],"currentState":"WaitingForPlayers"}`
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
	expectedState := `{"groupName":"group","players":[{"name":"player1","host":true,"points":0},{"name":"player2","host":false,"points":0}],"currentState":"WaitingForPlayers"}`
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
