package statemanager

import (
	"drawydraw/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateGroup_NewGroup_Succeeds(t *testing.T) {
	groupName := "group"
	models.SetGameProvider(models.NewTestGameProvider())
	err := CreateGroup(groupName)
	assert.Nil(t, err)
}

func TestCreateGroup_GroupExists_Fails(t *testing.T) {
	groupName := "group"
	models.SetGameProvider(models.NewTestGameProvider())
	CreateGroup(groupName)
	err := CreateGroup(groupName)
	assert.NotNil(t, err)
}

func TestCreateGroup_ShortGroupName_Fails(t *testing.T) {
	err := CreateGroup("")
	assert.NotNil(t, err)
}

func TestAddPlayer_AddHost_Succeeds(t *testing.T) {
	groupName := "group"
	models.SetGameProvider(models.NewTestGameProvider())
	CreateGroup(groupName)
	gameStatus, err := AddPlayer("mama cat", groupName, true)
	assert.Nil(t, err)
	assert.NotNil(t, gameStatus)
	expectedPlayers := []*Player{{Name: "mama cat", Host: true}}
	assert.EqualValues(t, gameStatus.Players, expectedPlayers)
	expectedCurrentPlayer := &CurrentPlayer{Name: "mama cat", IsHost: true}
	assert.EqualValues(t, expectedCurrentPlayer, gameStatus.CurrentPlayer)
}

func TestAddPlayer_AddToHostedGame_Succeeds(t *testing.T) {
	groupName := "group"
	models.SetGameProvider(models.NewTestGameProvider())
	CreateGroup(groupName)
	AddPlayer("papa cat", groupName, true)
	gameState, _ := AddPlayer("mama cat", groupName, false)
	assert.NotNil(t, gameState)
}

func TestAddPlayer_AddToUnHostedGame_Fails(t *testing.T) {
	groupName := "group"
	models.SetGameProvider(models.NewTestGameProvider())
	CreateGroup(groupName)
	gameState, _ := AddPlayer("mama cat", groupName, false)
	assert.Nil(t, gameState)
}

func TestAddPlayer_NoGroupCreated_Fails(t *testing.T) {
	groupName := "group"
	models.SetGameProvider(models.NewTestGameProvider())
	_, err := AddPlayer("baby cat", groupName, false)
	assert.NotNil(t, err)
}

func TestAddPlayer_ShortPlayerName_Fails(t *testing.T) {
	groupName := "group"
	models.SetGameProvider(models.NewTestGameProvider())
	_, err := AddPlayer("", groupName, false)
	assert.NotNil(t, err)
}

func TestAddPlayer_PlayerExistsInGroup_NoOps(t *testing.T) {
	groupName := "group"
	playerName := "baby cat"
	models.SetGameProvider(models.NewTestGameProvider())
	CreateGroup(groupName)
	AddPlayer(playerName, groupName, true)
	gameStatus, err := AddPlayer(playerName, groupName, true)
	assert.Nil(t, err)
	assert.NotNil(t, gameStatus)
	expectedPlayers := []*Player{{Name: "baby cat", Host: true}}
	assert.EqualValues(t, gameStatus.Players, expectedPlayers)
}

func TestAddPlayer_AddSecondHost_Fails(t *testing.T) {
	groupName := "group"
	models.SetGameProvider(models.NewTestGameProvider())
	CreateGroup(groupName)
	AddPlayer("old cat", groupName, true)
	_, err := AddPlayer("dead cat", groupName, true)
	assert.NotNil(t, err)
}

func TestStartGame_Host_Succeeds(t *testing.T) {
	groupName := "group"
	models.SetGameProvider(models.NewTestGameProvider())
	CreateGroup(groupName)
	AddPlayer("host cat", groupName, true)
	AddPlayer("angry cat", groupName, false)
	statusAfterAddPlayer, err := AddPlayer("annoyed cat", groupName, false)
	assert.EqualValues(t, statusAfterAddPlayer.CurrentState, models.WaitingForPlayers)

	statusAfterStart, err := StartGame(groupName, "host cat")
	assert.Nil(t, err)
	assert.NotNil(t, statusAfterStart)
	assert.EqualValues(t, statusAfterStart.CurrentState, models.InitialPromptCreation)
}

func TestStartGame_NonHost_Fails(t *testing.T) {
	groupName := "group"
	models.SetGameProvider(models.NewTestGameProvider())
	CreateGroup(groupName)
	AddPlayer("host cat", groupName, true)
	AddPlayer("angry cat", groupName, false)
	AddPlayer("annoyed cat", groupName, false)
	startResponse, err := StartGame(groupName, "angry cat")
	assert.NotNil(t, err)
	assert.Nil(t, startResponse)
}

func TestAddPrompt_Succeeds(t *testing.T) {
	//set up a group, add players, and start the game
	groupName := "group"
	models.SetGameProvider(models.NewTestGameProvider())
	CreateGroup(groupName)
	AddPlayer("host cat", groupName, true)
	AddPlayer("angry cat", groupName, false)
	AddPlayer("annoyed cat", groupName, false)
	startResponse, err := StartGame(groupName, "host cat")
	assert.NotNil(t, startResponse)

	//add prompt and check its in game state
	addPromptResponse, err := AddPrompt("annoyed cat", groupName, "tuna", "stinky", "yummy")
	assert.Nil(t, err)
	assert.NotNil(t, addPromptResponse)
}

func TestGameStatusForPlayer_Fails_PlayerMissing(t *testing.T) {
	//set up a group, add players, and start the game
	groupName := "group"
	models.SetGameProvider(models.NewTestGameProvider())
	CreateGroup(groupName)
	AddPlayer("host cat", groupName, true)
	game := models.GetGameProvider().LoadGame(groupName)
	gameStatus, err := gameStatusForPlayer(game, "missing cat")
	assert.Nil(t, gameStatus)
	assert.NotNil(t, err)
}

func TestSubmitDrawing(t *testing.T) {
	//set up a group, add players, add a prompt
	groupName := "group"
	models.SetGameProvider(models.NewTestGameProvider())
	CreateGroup(groupName)
	AddPlayer("host cat", groupName, true)
	AddPlayer("annoyed cat", groupName, false)
	StartGame(groupName, "host cat")
	AddPrompt("annoyed cat", groupName, "tuna", "stinky", "yummy")
	AddPrompt("host cat", groupName, "big", "handsome", "can")
	gameStatus, err := SubmitDrawing("annoyed cat", groupName, "mock data")
	assert.Nil(t, err)
	assert.NotNil(t, gameStatus)
}

func TestSubmitDrawing_Fails_PlayerMissing(t *testing.T) {
	//set up a group, add players, add a prompt
	groupName := "group"
	models.SetGameProvider(models.NewTestGameProvider())
	CreateGroup(groupName)
	AddPlayer("host cat", groupName, true)
	AddPlayer("annoyed cat", groupName, false)
	StartGame(groupName, "host cat")
	AddPrompt("annoyed cat", groupName, "tuna", "stinky", "yummy")
	AddPrompt("host cat", groupName, "big", "handsome", "can")
	gameStatus, err := SubmitDrawing("ninja cat", groupName, "mock data")
	assert.Nil(t, gameStatus)
	assert.NotNil(t, err)
}
