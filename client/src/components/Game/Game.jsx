import React from 'react';
import GroupSelectionScreen from '../GroupSelectionScreen/GroupSelectionScreen';
import WaitingForPlayersScreen from '../WaitingForPlayersScreen/WaitingForPlayersScreen';
import './Game.css';
import { GameStates } from '../../utils/constants';

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
      case GameStates.GroupSelection:
        return <GroupSelectionScreen onGameEntered={this.onGameEntered} />;
      case GameStates.WaitingForPlayers:
        return (
          <WaitingForPlayersScreen
            gameState={gameState}
            onGameStateChanged={this.onGameStateChanged}
          />
        );
      // Unknown group! Badness
      default:
        return <div><h1>We are sorry this is not implemented yet :(</h1></div>;
    }
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

  toggleConsole() {
    this.consoleEnabled = !this.consoleEnabled;
    // Re-render when the debug button is pressed.
    this.forceUpdate();
  }

  debugConsole() {
    return (
      <div className="debug">
        <button className="toggleConsole" type="button" onClick={this.toggleConsole}>debug</button>
        {this.consoleEnabled ? <div className="console"><pre> {JSON.stringify(this.state.gameState, null, 4)} </pre></div> : null}
      </div>
    );  
  }
}

export default Game;
