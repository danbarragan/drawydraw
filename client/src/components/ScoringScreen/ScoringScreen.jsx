import React from 'react';
import axios from 'axios';
import PropTypes from 'prop-types';
import { FormattedMessage } from 'react-intl';
import { formatServerError } from '../../utils/errorFormatting';
import UpdateGameState from '../../utils/updateGameState';
import './ScoringScreen.css';

class ScoringScreen extends React.Component {
  static formatBreakdownItem(breakdownItem) {
    const { reason, causingPlayer, amount } = breakdownItem;
    switch (reason) {
      case 'FooledPlayer':
        return (
          <FormattedMessage
            id="scoringScreen.fooledPlayerExplanation"
            defaultMessage="{amount} points because {causingPlayer} chose your decoy prompt"
            values={{ causingPlayer, amount }}
          />
        );
      case 'OtherChosePromptDrawn':
        return (
          <FormattedMessage
            id="scoringScreen.otherChosePromptDrawnExplanation"
            defaultMessage="{amount} points because {causingPlayer} chose the prompt you drew"
            values={{ causingPlayer, amount }}
          />
        );
      case 'ChoseCorrectPrompt':
        return (
          <FormattedMessage
            id="scoringScreen.choseCorrectPromptExplanation"
            defaultMessage="{amount} points because you chose the correct prompt"
            values={{ amount }}
          />
        );
      default:
        return (
          <FormattedMessage
            id="scoringScreen.unknownReasonExplanation"
            defaultMessage="{amount} points because ???"
            values={{ amount }}
          />
        );
    }
  }

  static formatPlayerScores(pointStandings, currentPlayerName) {
    const playerScores = [];
    // Sort standings by total score in descending order
    const sortedStandings = Object.values(pointStandings).sort(
      (a, b) => b.totalScore - a.totalScore,
    );
    sortedStandings.forEach((standing) => {
      const scoreItems = [];
      const {
        roundPointsBreakdown, totalScore, player,
      } = standing;
      roundPointsBreakdown.sort(
        (itemA, itemB) => (
          // Sort breakdown items first by score (desc) and then by reason
          itemA.amount === itemB.amount
            ? itemA.reason.localeCompare(itemB.reason)
            : itemB.amount - itemA.amount
        ),
      );
      let totalRoundScore = 0;
      roundPointsBreakdown.forEach((scoreItem) => {
        totalRoundScore += scoreItem.amount;
        scoreItems.push(
          <li key={`${player}-${scoreItem.amount}-${scoreItem.reason}`}>
            {ScoringScreen.formatBreakdownItem(scoreItem)}
          </li>,
        );
      });
      playerScores.push(
        <li key={player}>
          {player === currentPlayerName ? '*' : null}
          <FormattedMessage
            id="scoringScreen.pointSummary"
            defaultMessage="{player}: {totalScore} points ({totalRoundScore} points this round)"
            values={{ player, totalRoundScore, totalScore }}
          />
          <ul>
            {scoreItems}
          </ul>
        </li>,
      );
    });
    return playerScores;
  }

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
      currentPlayer, pointStandings, currentDrawing, pastDrawings,
    } = gameState;
    const { name: currentPlayerName, isHost } = currentPlayer;
    const playerScores = ScoringScreen.formatPlayerScores(pointStandings, currentPlayerName);
    const pastDrawingItems = pastDrawings.map((drawing) => (
      <div className="pastDrawingContainer" key={drawing.originalPrompt}>
        <img className="pastDrawing" src={drawing.imageData} alt="a drawing" />
        <span>
          <FormattedMessage
            id="scoringScreen.pastDrawingDescriptionFormat"
            defaultMessage="{adjective1} and {adjective2} {noun} by {author}"
            values={{
              author: drawing.author,
              adjective1: drawing.originalPrompt.adjectives[0],
              adjective2: drawing.originalPrompt.adjectives[1],
              noun: drawing.originalPrompt.noun,
            }}
          />
        </span>
      </div>
    ));
    return (
      <div className="screen votingScreen">
        <img className="promptImage" src={currentDrawing.imageData} alt="a drawing" />
        <p>
          <FormattedMessage
            id="scoringScreen.promptAuthorFormat"
            defaultMessage="The correct prompt for this image by {author} was:"
            values={{
              author: currentDrawing.author,
            }}
          />
          <br />
          <b>
            <FormattedMessage
              id="scoringScreen.promptFormat"
              defaultMessage="{adjective1} and {adjective2} {noun}"
              values={{
                author: currentDrawing.author,
                adjective1: currentDrawing.originalPrompt.adjectives[0],
                adjective2: currentDrawing.originalPrompt.adjectives[1],
                noun: currentDrawing.originalPrompt.noun,
              }}
            />
          </b>
        </p>
        <h4>
          <FormattedMessage
            id="scoringScreen.currentScoreHeader"
            defaultMessage="Current Scores:"
          />
        </h4>
        <ul>{playerScores}</ul>
        { isHost ? (
          <button type="button" className="buttonTypeA" onClick={this.onNextRoundButtonClicked}>
            <FormattedMessage
              id="scoringScreen.nextRoundButton"
              defaultMessage="Next round"
            />
          </button>
        )
          : (
            <h3>
              <FormattedMessage
                id="scoringScreen.waitingForNextRoundHeader"
                defaultMessage="Waiting for the host to start the next round..."
              />
            </h3>
          )}
        { pastDrawings.length > 0 ? (
          <div className="pastDrawings">
            <h3>
              <FormattedMessage
                id="scoringScreen.pastDrawingsLabel"
                defaultMessage="Past drawings from this round:"
              />
            </h3>
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
  originalPrompt: PropTypes.shape({
    noun: PropTypes.string.isRequired,
    adjectives: PropTypes.arrayOf(PropTypes.string).isRequired,
  }).isRequired,
  author: PropTypes.string.isRequired,
});

ScoringScreen.propTypes = {
  gameState: PropTypes.shape({
    pointStandings: PropTypes.objectOf(PropTypes.shape({
      totalScore: PropTypes.number.isRequired,
      player: PropTypes.string.isRequired,
      roundPointsBreakdown: PropTypes.arrayOf(PropTypes.shape({
        amount: PropTypes.number.isRequired,
        reason: PropTypes.string.isRequired,
      })).isRequired,
    })).isRequired,
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
