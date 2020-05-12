import React from 'react';
import axios from 'axios';
import GroupSelectionScreen from '../GroupSelectionScreen/GroupSelectionScreen';
import WaitingForPlayersScreen from '../WaitingForPlayersScreen/WaitingForPlayersScreen';
import './Game.css';
import { GameStates } from '../../utils/constants';
import { formatServerError } from '../../utils/errorFormatting';

class Game extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      gameState: {
        currentState: GameStates.GroupSelection, // Start in the group selection state
      },
    };
    this.getCurrentComponent = this.getCurrentComponent.bind(this);
    this.onGameEntered = this.onGameEntered.bind(this);
    this.onGameStateChanged = this.onGameStateChanged.bind(this);
    this.debugConsole = this.debugConsole.bind(this);
    this.toggleConsole = this.toggleConsole.bind(this);
    this.debugSetGameState = this.debugSetGameState.bind(this);
    this.consoleEnabled = true;
  }

  onGameEntered(gameState) {
    this.setState({ gameState });
  }

  onGameStateChanged(gameState) {
    this.setState({ gameState });
  }

  getCurrentComponent() {
    const { gameState } = this.state;
    const { currentState } = gameState;

    switch (currentState) {
      case GameStates.WaitingForPlayers:
        return (
          <WaitingForPlayersScreen
            gameState={gameState}
            onGameStateChanged={this.onGameStateChanged}
          />
        );
      // For now the group selection screen will be used as the default one
      default:
        return <GroupSelectionScreen onGameEntered={this.onGameEntered} />;
    }
  }

  toggleConsole() {
    this.consoleEnabled = !this.consoleEnabled;
    // Re-render when the debug button is pressed.
    this.forceUpdate();
  }

  async debugSetGameState(gameStateName) {
    try {
      const response = await axios.post('/api/set-game-state', { gameStateName });
      this.setState({ gameState: response.data });
      this.forceUpdate();
    } catch (error) {
      this.setState({ error: formatServerError(error) });
    }
  }

  debugConsole() {
    const { error } = this.state;
    const { gameState } = this.state;
    return (
      <div className="debug">
        <button className="toggleConsole" type="button" onClick={this.toggleConsole}>debug</button>
        <button className="setState" type="button" onClick={(() => this.debugSetGameState('WaitingForPlayers'))}>WaitingForPlayers</button>
        {this.consoleEnabled ? (
          <div className="console">
            <p>
              Errors:
              {error || ':)'}
            </p>
            <pre>
              {' '}
              {JSON.stringify(gameState, null, 4)}
              {' '}
            </pre>
          </div>
        ) : null}
      </div>
    );
  }

  render() {
    return (
      <div className="game">
        <div className="gameTitle"><h1>Drawydraw</h1></div>
        {this.getCurrentComponent()}
        {this.debugConsole()}
      </div>
    );
  }
}

export default Game;
