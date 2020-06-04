import React from 'react';
import axios from 'axios';
import { FormattedMessage } from 'react-intl';
import PropTypes from 'prop-types';
import { formatServerError } from '../../utils/errorFormatting';
import UpdateGameState from '../../utils/updateGameState';

class WaitingForPlayersScreen extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      timerId: null,
    };
    this.updateGameState = this.updateGameState.bind(this);
    this.onStartGameButtonClicked = this.onStartGameButtonClicked.bind(this);
  }

  componentDidMount() {
    const timerId = setInterval(this.updateGameState, 3000);
    this.setState({ timerId });
  }

  componentWillUnmount() {
    const { timerId } = this.state;
    if (timerId !== null) {
      clearInterval(timerId);
    }
  }

  async onStartGameButtonClicked() {
    const { gameState, onGameStateChanged } = this.props;
    const { groupName, currentPlayer } = gameState;
    const { name } = currentPlayer;
    const data = { playerName: name, groupName };
    try {
      const response = await axios.post('/api/start-game', data);
      onGameStateChanged(response.data);
    } catch (error) {
      this.setState({ error: formatServerError(error) });
    }
  }

  updateGameState() {
    const { gameState, onGameStateChanged } = this.props;
    const { groupName, currentPlayer } = gameState;
    const { name: playerName } = currentPlayer;
    UpdateGameState(
      groupName,
      playerName,
      onGameStateChanged,
      (error) => { this.setState({ error: formatServerError(error) }); },
    );
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
      <div className="screen waitingForPlayersScreen">
        <h1>
          <FormattedMessage
            id="waitingForPlayersScreen.groupName"
            defaultMessage="Group name: {groupName}"
            values={{ groupName }}
          />
        </h1>
        <ul>{playerList}</ul>
        { isHost ? (
          <button type="button" className="buttonTypeA" onClick={this.onStartGameButtonClicked}>
            <FormattedMessage
              id="waitingForPlayersScreen.startGameButton"
              defaultMessage="Start"
            />
          </button>
        )
          : (
            <h3>
              <FormattedMessage
                id="waitingForPlayersScreen.waitingForHostMessage"
                defaultMessage="Waiting for the host to start the game..."
              />
            </h3>
          )}
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
