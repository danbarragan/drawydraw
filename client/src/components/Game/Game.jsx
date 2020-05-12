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
        <div className="gameTitle"><h1>Some game or something</h1></div>
        {this.getCurrentComponent()}
      </div>
    );
  }
}

export default Game;
