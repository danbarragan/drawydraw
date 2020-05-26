import React from 'react';
import axios from 'axios';
import PropTypes from 'prop-types';
import { formatServerError } from '../../utils/errorFormatting';
import UpdateGameState from '../../utils/updateGameState';

class ScoringScreen extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      timerId: null,
    };
    this.updateGameState = this.updateGameState.bind(this);
    this.onNextRoundButtonClicked = this.onNextRoundButtonClicked.bind(this);
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

  async onNextRoundButtonClicked() {
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
    const {
      players, currentPlayer, roundScores,
    } = gameState;
    const { name: currentPlayerName, isHost } = currentPlayer;
    const scoresBeforeRound = players.reduce(
      (dict, player) => ({ ...dict, [player.name]: player.points }), {},
    );
    const scores = Object.entries(roundScores).map(([player, score]) => (
      <li key={player}>
        {player === currentPlayerName ? '*' : null}
        {`${player}:a ${score + scoresBeforeRound[player]} points`}
        {` (+${score} points this round)`}
      </li>
    ));
    return (
      <div className="screen votingScreen">
        <ul>{scores}</ul>
        { isHost ? <button type="button" className="buttonTypeA" onClick={this.onNextRoundButtonClicked}>Next</button>
          : <h3>Waiting for the host to start the next round...</h3>}
        <h3 className="error">{error}</h3>
      </div>
    );
  }
}

ScoringScreen.propTypes = {
  gameState: PropTypes.shape({
    roundScores: PropTypes.objectOf(PropTypes.number).isRequired,
    currentPlayer: PropTypes.shape({
      name: PropTypes.string.isRequired,
      isHost: PropTypes.bool.isRequired,
    }).isRequired,
    players: PropTypes.arrayOf(PropTypes.shape({
      name: PropTypes.string.isRequired,
      points: PropTypes.number.isRequired,
    })),
    groupName: PropTypes.string.isRequired,
  }).isRequired,
  onGameStateChanged: PropTypes.func.isRequired,
};

export default ScoringScreen;
