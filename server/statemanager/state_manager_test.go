package statemanager

import (
	"drawydraw/models"
	"drawydraw/test"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateGroup_NewGroup_Succeeds(t *testing.T) {
	test.SetupTestGameProvider(t)
	groupName := "group"
	err := CreateGroup(groupName)
	assert.Nil(t, err)
}

func TestCreateGroup_GroupExists_Fails(t *testing.T) {
	test.SetupTestGameProvider(t)
	groupName := "group"
	CreateGroup(groupName)
	err := CreateGroup(groupName)
	assert.NotNil(t, err)
}

func TestCreateGroup_ShortGroupName_Fails(t *testing.T) {
	test.SetupTestGameProvider(t)
	err := CreateGroup("")
	assert.NotNil(t, err)
}

func TestAddPlayer_AddHost_Succeeds(t *testing.T) {
	test.SetupTestGameProvider(t)
	groupName := "group"
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
	test.SetupTestGameProvider(t)
	groupName := "group"
	CreateGroup(groupName)
	AddPlayer("papa cat", groupName, true)
	gameState, _ := AddPlayer("mama cat", groupName, false)
	assert.NotNil(t, gameState)
}

func TestAddPlayer_AddToUnHostedGame_Fails(t *testing.T) {
	test.SetupTestGameProvider(t)
	groupName := "group"
	CreateGroup(groupName)
	gameState, _ := AddPlayer("mama cat", groupName, false)
	assert.Nil(t, gameState)
}

func TestAddPlayer_NoGroupCreated_Fails(t *testing.T) {
	test.SetupTestGameProvider(t)
	groupName := "group"
	_, err := AddPlayer("baby cat", groupName, false)
	assert.NotNil(t, err)
}

func TestAddPlayer_ShortPlayerName_Fails(t *testing.T) {
	test.SetupTestGameProvider(t)
	groupName := "group"
	_, err := AddPlayer("", groupName, false)
	assert.NotNil(t, err)
}

func TestAddPlayer_PlayerExistsInGroup_NoOps(t *testing.T) {
	test.SetupTestGameProvider(t)
	groupName := "group"
	playerName := "baby cat"
	CreateGroup(groupName)
	AddPlayer(playerName, groupName, true)
	gameStatus, err := AddPlayer(playerName, groupName, true)
	assert.Nil(t, err)
	assert.NotNil(t, gameStatus)
	expectedPlayers := []*Player{{Name: "baby cat", Host: true}}
	assert.EqualValues(t, gameStatus.Players, expectedPlayers)
}

func TestAddPlayer_AddSecondHost_Fails(t *testing.T) {
	test.SetupTestGameProvider(t)
	game := test.GameInWaitingForPlayersState()
	models.GetGameProvider().SaveGame(game)
	gameStatus, err := AddPlayer("extra cat", game.GroupName, true)
	assert.NotNil(t, err)
	assert.Nil(t, gameStatus)
}

func TestStartGame_Host_Succeeds(t *testing.T) {
	test.SetupTestGameProvider(t)
	game := test.GameInWaitingForPlayersState()
	models.GetGameProvider().SaveGame(game)
	gameStatus, err := StartGame(game.GroupName, game.Players[0].Name)
	assert.Nil(t, err)
	assert.NotNil(t, gameStatus)
	assert.EqualValues(t, gameStatus.CurrentState, models.InitialPromptCreation)
}

func TestStartGame_NonHost_Fails(t *testing.T) {
	test.SetupTestGameProvider(t)
	game := test.GameInWaitingForPlayersState()
	models.GetGameProvider().SaveGame(game)
	gameStatus, err := StartGame(game.GroupName, game.Players[1].Name)
	assert.NotNil(t, err)
	assert.Nil(t, gameStatus)
}

func TestAddPrompt_Succeeds(t *testing.T) {
	test.SetupTestGameProvider(t)
	game := test.GameInInitialPromptCreationState()
	models.GetGameProvider().SaveGame(game)
	// The game should only transition to the drawing state when all players submit their prompts
	for _, player := range game.Players[:2] {
		gameState, err := AddPrompt(player.Name, game.GroupName, "tuna", "stinky", "yummy")
		assert.Nil(t, err)
		assert.NotNil(t, gameState)
		assert.EqualValues(t, gameState.CurrentState, models.InitialPromptCreation)
	}
	// The game should only transition to the drawing state when all players submit their prompts
	gameState, err := AddPrompt(game.Players[2].Name, game.GroupName, "sardine", "small", "funny")
	assert.Nil(t, err)
	assert.NotNil(t, gameState)
	assert.EqualValues(t, gameState.CurrentState, models.DrawingsInProgress)
}

func TestGameStatusForPlayer_Fails_PlayerMissing(t *testing.T) {
	game := test.GameInInitialPromptCreationState()
	gameStatus, err := gameStatusForPlayer(game, "missing cat")
	assert.Nil(t, gameStatus)
	assert.NotNil(t, err)
}

func TestSubmitDrawing(t *testing.T) {
	test.SetupTestGameProvider(t)
	game := test.GameInDrawingsInProgressState()
	models.GetGameProvider().SaveGame(game)
	// The game should only transition to the decoy prompt phase when all players submit their drawings
	for _, player := range game.Players[:2] {
		gameState, err := SubmitDrawing(player.Name, game.GroupName, "mock data")
		assert.Nil(t, err)
		assert.NotNil(t, gameState)
		assert.EqualValues(t, gameState.CurrentState, models.DrawingsInProgress)
	}
	gameStatus, err := SubmitDrawing(game.Players[2].Name, game.GroupName, "mock data")
	assert.Nil(t, err)
	assert.NotNil(t, gameStatus)
	assert.EqualValues(t, gameStatus.CurrentState, models.DecoyPromptCreation)
}

func TestSubmitDrawing_Fails_PlayerMissing(t *testing.T) {
	test.SetupTestGameProvider(t)
	game := test.GameInDrawingsInProgressState()
	models.GetGameProvider().SaveGame(game)
	gameStatus, err := SubmitDrawing("Missing player", game.GroupName, "mock data")
	assert.NotNil(t, err)
	assert.Nil(t, gameStatus)
}

func TestAddDecoyPrompt_Success(t *testing.T) {
	test.SetupTestGameProvider(t)
	game := test.GameInDecoyPromptCreationState()
	models.GetGameProvider().SaveGame(game)
	// The game should only move to voting once all players submit decoy prompts
	gameStatus, err := AddPrompt(game.Players[0].Name, game.GroupName, "fish", "tasty", "red")
	assert.Nil(t, err)
	assert.NotNil(t, gameStatus)
	assert.EqualValues(t, gameStatus.CurrentState, models.DecoyPromptCreation)

	gameStatus, err = AddPrompt(game.Players[2].Name, game.GroupName, "salmon", "strange", "big")
	assert.Nil(t, err)
	assert.NotNil(t, gameStatus)
	assert.EqualValues(t, gameStatus.CurrentState, models.Voting)
}

func TestCastVote_Success(t *testing.T) {
	test.SetupTestGameProvider(t)
	game := test.GameInVotingState()
	models.GetGameProvider().SaveGame(game)
	activeDrawing := game.Drawings[0]
	gameStatus, err := CastVote(game.Players[0].Name, game.GroupName, activeDrawing.DecoyPrompts[game.Players[2].Name].Identifier)
	assert.Nil(t, err)
	assert.NotNil(t, gameStatus)
	assert.EqualValues(t, models.Voting, gameStatus.CurrentState)
	// Once all players vote we should move to scoring
	gameStatus, err = CastVote(game.Players[2].Name, game.GroupName, activeDrawing.OriginalPrompt.Identifier)
	assert.Nil(t, err)
	assert.NotNil(t, gameStatus)
	assert.EqualValues(t, models.Scoring, gameStatus.CurrentState)
}

func TestAddDecoyPrompt_Error_duplicatePromptEntry(t *testing.T) {
	test.SetupTestGameProvider(t)
	game := test.GameInDecoyPromptCreationState()
	models.GetGameProvider().SaveGame(game)
	AddPrompt(game.Players[0].Name, game.GroupName, "fish", "tasty", "red")
	gameStatus, err := AddPrompt(game.Players[0].Name, game.GroupName, "fish", "tasty", "red")
	assert.Nil(t, gameStatus)
	assert.NotNil(t, err)
}
