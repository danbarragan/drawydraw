const GameStates = Object.freeze({
  GroupSelection: 'GroupSelection',
  WaitingForPlayers: 'WaitingForPlayers',
  InitialPromptCreation: 'InitialPromptCreation',
  DrawingsInProgress: 'DrawingsInProgress',
  DecoyPromptCreation: 'DecoyPromptCreation',
  Voting: 'Voting',
  Scoring: 'Scoring',
});

exports.GameStates = GameStates;
