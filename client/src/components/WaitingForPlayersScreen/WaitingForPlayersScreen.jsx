import React from 'react';
import axios from 'axios';
import PropTypes from 'prop-types';
import { formatServerError } from '../../utils/errorFormatting';


class WaitingForPlayersScreen extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
    };
    this.updateGameState = this.updateGameState.bind(this);
  }

  componentDidMount() {
    setInterval(this.updateGameState, 3000);
  }

  // Todo: Probably move to a helper since it's going to be used in other screens
  async updateGameState() {
    const { gameState, onGameStateChanged } = this.props;
    const { groupName, currentPlayer } = gameState;
    const { name: playerName } = currentPlayer;
    try {
      const response = await axios.get(`/api/get-game-status/${groupName}?playerName=${playerName}`);
      onGameStateChanged(response.data);
    } catch (error) {
      this.setState({ error: formatServerError(error) });
    }
  }

  render() {
    const { error } = this.state;
    const { gameState } = this.props;
    const { players, groupName, currentPlayer } = gameState;
    const { name: currentPlayerName, isHost } = currentPlayer;
    const playerList = players.map((player) => (
      <li key={player.name}>
        {player.name}
        {player.name === currentPlayerName ? '*' : null}
      </li>
    ));
    return (
      <div className="waitingForPlayersScreen">
        <h1>
          Group name:
          {groupName}
        </h1>
        <ul>{playerList}</ul>
        { isHost ? <button type="button">Start game</button> : <h3>Waiting for the host to start the game...</h3> }
        <h3 className="error">{error}</h3>
      </div>
    );
  }
}

WaitingForPlayersScreen.propTypes = {
  gameState: PropTypes.shape({
    currentPlayer: PropTypes.shape({
      name: PropTypes.string.isRequired,
      isHost: PropTypes.bool.isRequired,
    }).isRequired,
    players: PropTypes.arrayOf(PropTypes.shape({
      name: PropTypes.string.isRequired,
    })),
    groupName: PropTypes.string.isRequired,
  }).isRequired,
  onGameStateChanged: PropTypes.func.isRequired,
};

export default WaitingForPlayersScreen;
