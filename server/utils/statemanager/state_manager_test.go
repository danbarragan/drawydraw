package statemanager

import (
	"drawydraw/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Todo - add setup code or teardown helpers to clean up memcache values as they
// persist between tests causing false positive/negative results. For now just
// use unique group names per test
func randomGroupName() string {
	return uuid.New().String()
}

func TestCreateGroup_NewGroup_Succeeds(t *testing.T) {
	groupName := randomGroupName()
	err := CreateGroup(groupName)
	assert.Nil(t, err)
}

func TestCreateGroup_GroupExists_Fails(t *testing.T) {
	groupName := randomGroupName()
	CreateGroup(groupName)
	err := CreateGroup(groupName)
	assert.NotNil(t, err)
}

func TestCreateGroup_ShortGroupName_Fails(t *testing.T) {
	err := CreateGroup("")
	assert.NotNil(t, err)
}

func TestAddPlayer_AddHost_Succeeds(t *testing.T) {
	groupName := randomGroupName()
	CreateGroup(groupName)
	gameState, err := AddPlayer("mama cat", groupName, true)
	assert.Nil(t, err)
	assert.NotNil(t, gameState)
	expectedPlayers := []*models.Player{{Name: "mama cat", Host: true}}
	assert.EqualValues(t, (*gameState)["players"], expectedPlayers)
	currentPlayer := (*gameState)["currentPlayer"].(map[string]interface{})
	assert.Equal(t, currentPlayer["isHost"], true)
	assert.Equal(t, currentPlayer["name"], "mama cat")
}

//Should this be _Succeeds?
func TestAddPlayer_AddToHostedGame_Fails(t *testing.T) {
	groupName := randomGroupName()
	CreateGroup(groupName)
	AddPlayer("papa cat", groupName, true)
	gameState, _ := AddPlayer("mama cat", groupName, false)
	assert.NotNil(t, gameState)
}

func TestAddPlayer_AddToUnHostedGame_Fails(t *testing.T) {
	groupName := randomGroupName()
	CreateGroup(groupName)
	gameState, _ := AddPlayer("mama cat", groupName, false)
	assert.Nil(t, gameState)
}

func TestAddPlayer_NoGroupCreated_Fails(t *testing.T) {
	groupName := randomGroupName()
	_, err := AddPlayer("baby cat", groupName, false)
	assert.NotNil(t, err)
}

func TestAddPlayer_ShortPlayerName_Fails(t *testing.T) {
	groupName := randomGroupName()
	_, err := AddPlayer("", groupName, false)
	assert.NotNil(t, err)
}

func TestAddPlayer_PlayerExistsInGroup_NoOps(t *testing.T) {
	groupName := randomGroupName()
	playerName := "baby cat"
	CreateGroup(groupName)
	AddPlayer(playerName, groupName, true)
	statusResponse, err := AddPlayer(playerName, groupName, true)
	assert.Nil(t, err)
	assert.NotNil(t, statusResponse)
	expectedPlayers := []*models.Player{{Name: "baby cat", Host: true}}
	assert.EqualValues(t, (*statusResponse)["players"], expectedPlayers)
}

func TestAddPlayer_AddSecondHost_Fails(t *testing.T) {
	groupName := randomGroupName()
	CreateGroup(groupName)
	AddPlayer("old cat", groupName, true)
	_, err := AddPlayer("dead cat", groupName, true)
	assert.NotNil(t, err)
}

func TestStartGame_HostTriesandSucceeds(t *testing.T) {
	groupName := randomGroupName()
	CreateGroup(groupName)
	AddPlayer("host cat", groupName, true)
	AddPlayer("angry cat", groupName, false)
	addPlayerResponse, err := AddPlayer("annoyed cat", groupName, false)

	expectedState1 := models.WaitingForPlayers
	assert.EqualValues(t, (*addPlayerResponse)["currentState"], expectedState1)

	startResponse, err := StartGame(groupName, "host cat")
	assert.Nil(t, err)
	assert.NotNil(t, startResponse)

	expectedState2 := models.InitialPromptCreation
	assert.EqualValues(t, (*startResponse)["currentState"], expectedState2)
}

func TestStartGame_NonHost_Fails(t *testing.T) {
	groupName := randomGroupName()
	CreateGroup(groupName)
	AddPlayer("host cat", groupName, true)
	AddPlayer("angry cat", groupName, false)
	AddPlayer("annoyed cat", groupName, false)
	startResponse, err := StartGame(groupName, "angry cat")
	assert.NotNil(t, err)
	assert.Nil(t, startResponse)
}

func TestAddPrompts_Succeeds(t *testing.T) {
	//set up a group, add players, and start the game
	groupName := randomGroupName()
	CreateGroup(groupName)
	AddPlayer("host cat", groupName, true)
	AddPlayer("angry cat", groupName, false)
	AddPlayer("annoyed cat", groupName, false)
	startResponse, err := StartGame(groupName, "host cat")
	assert.NotNil(t, startResponse)

	//add prompt and check its in game state
	addPromptResponse, err := AddPrompts("friendly cat", groupName, "tuna", "stinky", "yummy")
	assert.Nil(t, err)
	assert.NotNil(t, addPromptResponse)

}
