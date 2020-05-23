import React from 'react';
import PropTypes from 'prop-types';
import axios from 'axios';
import { formatServerError } from '../../utils/errorFormatting';
import './VotingScreen.css';
import UpdateGameState from '../../utils/updateGameState';

class VotingScreen extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      timerId: null,
      error: null,
      selectedPromptId: null,
    };
    this.updateGameState = this.updateGameState.bind(this);
    this.handleOptionChange = this.handleOptionChange.bind(this);
    this.onVoteClicked = this.onVoteClicked.bind(this);
    this.componentDidMount = this.componentDidMount.bind(this);
  }

  componentDidMount() {
    // Start listening for game state updates
    const timerId = setInterval(this.updateGameState, 3000);
    this.setState({ timerId });
  }

  componentWillUnmount() {
    const { timerId } = this.state;
    if (timerId !== null) {
      clearInterval(timerId);
    }
  }

  async onVoteClicked() {
    const { gameState, onGameStateChanged } = this.props;
    const { groupName, currentPlayer } = gameState;
    const { selectedPromptId } = this.state;
    const data = {
      playerName: currentPlayer.name, groupName, selectedPromptId,
    };

    try {
      const response = await axios.post('/api/cast-vote', data);
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

  handleOptionChange(voteEvent) {
    this.setState({ selectedPromptId: voteEvent.target.value });
  }

  render() {
    const { error, selectedPromptId } = this.state;
    const { gameState } = this.props;
    const { players, currentDrawing, currentPlayer } = gameState;
    const votingElements = (
      <form className="votingForm">
        <h4 className="votingQuestion">What was the original prompt for this drawing?</h4>
        {currentDrawing.prompts.map((prompt) => (
          <div className="votingOption" key={prompt.identifier}>
            <label htmlFor={prompt.identifier}>
              <input
                id={prompt.identifier}
                className="votingRadio"
                type="radio"
                value={prompt.identifier}
                checked={selectedPromptId === prompt.identifier}
                onChange={this.handleOptionChange}
              />
              {`${prompt.adjectives.join(', ')} ${prompt.noun}`}
            </label>
          </div>
        ))}
        <button className="buttonTypeA" type="button" onClick={this.onVoteClicked}>Vote</button>
      </form>
    );
    const waitingElements = (
      <div>
        <h3>Waiting for other players to vote on a prompt...</h3>
        <ul>
          {
            players.map((player) => (
              player.name === currentPlayer.name ? null : (
                <li key={player.name}>
                  {player.name}
                  {player.hasPendingAction ? ' is still voting' : ' is done'}
                </li>
              )
            ))
          }
        </ul>
      </div>
    );
    return (
      <div className="screen voteSelection">
        <img className="promptImage" src={currentDrawing.imageData} alt="a drawing" />
        {currentPlayer.hasCompletedAction ? waitingElements : votingElements}
        <h3 className="error">{error}</h3>
      </div>
    );
  }
}

VotingScreen.propTypes = {
  gameState: PropTypes.shape({
    currentPlayer: PropTypes.shape({
      name: PropTypes.string.isRequired,
      isHost: PropTypes.bool.isRequired,
      hasCompletedAction: PropTypes.bool.isRequired,
    }).isRequired,
    currentDrawing: PropTypes.shape({
      imageData: PropTypes.string.isRequired,
      prompts: PropTypes.arrayOf(PropTypes.shape({
        identifier: PropTypes.string.isRequired,
        noun: PropTypes.string.isRequired,
        adjectives: PropTypes.arrayOf(PropTypes.string).isRequired,
      })).isRequired,
    }),
    players: PropTypes.arrayOf(PropTypes.shape({
      name: PropTypes.string.isRequired,
    })),
    groupName: PropTypes.string.isRequired,
  }).isRequired,
  onGameStateChanged: PropTypes.func.isRequired,
};

export default VotingScreen;
