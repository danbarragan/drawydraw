import React from 'react';
import axios from 'axios';
import PropTypes from 'prop-types';
import { formatServerError } from '../../utils/errorFormatting';
import UpdateGameState from '../../utils/updateGameState';
import './ScoringScreen.css';

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
      players, currentPlayer, roundScores, currentDrawing, pastDrawings,
    } = gameState;
    const { name: currentPlayerName, isHost } = currentPlayer;
    const scoresBeforeRound = players.reduce(
      (dict, player) => ({ ...dict, [player.name]: player.points }), {},
    );
    const playerScores = [];
    Object.entries(roundScores).forEach(([player, breakdown]) => {
      const scoreItems = [];
      let totalRoundScore = 0;
      breakdown.sort(
        (itemA, itemB) => (
          // Sort breakdown items first by score (desc) and then by reason
          itemA.amount === itemB.amount
            ? itemA.reason.localeCompare(itemB.reason)
            : itemB.amount - itemA.amount
        ),
      );
      breakdown.forEach((scoreItem) => {
        totalRoundScore += scoreItem.amount;
        scoreItems.push(
          <li key={`${player}-${scoreItem.amount}-${scoreItem.reason}`}>
            {`+${scoreItem.amount} because ${scoreItem.reason}`}
          </li>,
        );
      });
      playerScores.push(
        <li key={player}>
          {player === currentPlayerName ? '*' : null}
          {`${player}: ${totalRoundScore + scoresBeforeRound[player]} points`}
          {` (+${totalRoundScore} points this round)`}
          <ul>
            {scoreItems}
          </ul>
        </li>,
      );
    });
    const pastDrawingItems = pastDrawings.map((drawing) => (
      <div className="pastDrawingContainer" key={drawing.originalPrompt}>
        <img className="pastDrawing" src={drawing.imageData} alt="a drawing" />
        <span>{`${drawing.originalPrompt} by ${drawing.author}`}</span>
      </div>
    ));
    return (
      <div className="screen votingScreen">
        <img className="promptImage" src={currentDrawing.imageData} alt="a drawing" />
        <span>
          {`The correct prompt for this image by ${currentDrawing.author} was:`}
          <br />
          <b>{currentDrawing.originalPrompt}</b>
        </span>
        <h3>Current Scores:</h3>
        <ul>{playerScores}</ul>
        { isHost ? <button type="button" className="buttonTypeA" onClick={this.onNextRoundButtonClicked}>Next</button>
          : <h3>Waiting for the host to start the next round...</h3>}
        { pastDrawings.length > 0 ? (
          <div className="pastDrawings">
            <h3>Past drawings from this round</h3>
            {pastDrawingItems}
          </div>
        ) : null}
        <h3 className="error">{error}</h3>
      </div>
    );
  }
}

const drawingProptype = PropTypes.shape({
  imageData: PropTypes.string.isRequired,
  originalPrompt: PropTypes.string.isRequired,
  author: PropTypes.string.isRequired,
});

ScoringScreen.propTypes = {
  gameState: PropTypes.shape({
    roundScores: PropTypes.objectOf(PropTypes.arrayOf(PropTypes.shape({
      amount: PropTypes.number.isRequired,
      reason: PropTypes.string.isRequired,
    }))).isRequired,
    pastDrawings: PropTypes.arrayOf(drawingProptype).isRequired,
    currentDrawing: drawingProptype.isRequired,
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
