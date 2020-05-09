import React from 'react';
import GroupSelector from '../GroupSelector/GroupSelector';
import './Game.css';
import { GameStates } from '../../utils/constants';

class Game extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      currentGameState: GameStates.GroupSelection,
    };
    this.getCurrentComponent = this.getCurrentComponent.bind(this);
  }


  getCurrentComponent() {
    const { currentGameState } = this.state;
    switch (currentGameState) {
      default:
        return <GroupSelector />;
    }
  }

  render() {

    return (
      <div className="Game">
        {this.getCurrentComponent()}
      </div>
    );
  }
}

export default Game;
