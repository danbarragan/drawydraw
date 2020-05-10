package statemanager

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "github.com/google/uuid"
)

// Todo - add setup code or teardown helpers to clean up memcache values as they
// persist between tests causing false positive/negative results. For now just
// use unique group names per test
func randomGroupName() string {
  return uuid.New().String()
}

func TestCreateGroup_NewGroup_Succeeds(t *testing.T) {
  groupName := randomGroupName()
  gameState, err := CreateGroup(groupName)
  assert.Nil(t, err)
  assert.NotNil(t, gameState)
}

func TestCreateGroup_GroupExists_Fails(t *testing.T) {
  groupName := randomGroupName()
  CreateGroup(groupName)
  gameState, err := CreateGroup(groupName)
  assert.NotNil(t, err)
  assert.Nil(t, gameState)
}


func TestAddPlayer_AddHost_Succeeds(t *testing.T) {
  groupName := randomGroupName()
  CreateGroup(groupName)
  gameState, err := AddPlayer("mama cat", groupName, true)
  assert.Nil(t, err)
  assert.NotNil(t, gameState)
}

func TestAddPlayer_AddNonHost_Succeeds(t *testing.T) {
  groupName := randomGroupName()
  CreateGroup(groupName)
  gameState, err := AddPlayer("mama cat", groupName, false)
  assert.Nil(t, err)
  assert.NotNil(t, gameState)
}

func TestAddPlayer_NoGroupCreated_Fails(t *testing.T) {
  groupName := randomGroupName()
  gameState, err := AddPlayer("baby cat", groupName, false)
  assert.NotNil(t, err)
  assert.Nil(t, gameState)
}

func TestAddPlayer_PlayerExistsInGroup_Fails(t *testing.T) {
  groupName := randomGroupName()
  playerName := "baby cat"
  CreateGroup(groupName)
  AddPlayer(playerName, groupName, false)
  gameState, err := AddPlayer(playerName, groupName, false)
  assert.NotNil(t, err)
  assert.Nil(t, gameState)
}
