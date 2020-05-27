import React from 'react';
import axios from 'axios';
import PropTypes from 'prop-types';
import './InitialPromptCreationScreen.css';
import { formatServerError } from '../../utils/errorFormatting';
import UpdateGameState from '../../utils/updateGameState';


class InitialPromptCreationScreen extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      noun: '',
      adjective1: '',
      adjective2: '',
      error: null,
    };

    this.onSubmitPromptButtonClicked = this.onSubmitPromptButtonClicked.bind(this);
    this.updateGameState = this.updateGameState.bind(this);
    this.onNounChange = this.onNounChange.bind(this);
    this.onAdjective1Change = this.onAdjective1Change.bind(this);
    this.onAdjective2Change = this.onAdjective2Change.bind(this);
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

  onNounChange(event) {
    this.setState({ noun: event.target.value });
  }

  onAdjective1Change(event) {
    this.setState({ adjective1: event.target.value });
  }

  onAdjective2Change(event) {
    this.setState({ adjective2: event.target.value });
  }

  async onSubmitPromptButtonClicked() {
    const { gameState, onGameStateChanged } = this.props;
    const { groupName, currentPlayer } = gameState;
    const { name } = currentPlayer;
    const { noun, adjective1, adjective2 } = this.state;
    const data = {
      playerName: name, groupName, noun, adjective1, adjective2,
    };

    try {
      const response = await axios.post('/api/add-prompt', data);
      onGameStateChanged(response.data);
    } catch (error) {
      this.setState({ error: formatServerError(error) });
    }
  }

  // Todo: Probably move to a helper since it's going to be used in other screens
  async updateGameState() {
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
    const {
      noun, adjective1, adjective2, error,
    } = this.state;
    const { gameState } = this.props;
    const { players, currentPlayer } = gameState;


    const promptCreatingElements = (
      <div>
        <h3>Please provide a few words to generate drawing prompts</h3>
        <label htmlFor="adjective1">
          First Adjective
          <input id="adj1" type="text" value={adjective1} onChange={this.onAdjective1Change} autoComplete="off" />
        </label>
        <label htmlFor="adjective2">
          Second Adjective
          <input id="adj2" type="text" value={adjective2} onChange={this.onAdjective2Change} autoComplete="off" />
        </label>
        <label htmlFor="noun">
          Noun
          <input id="noun" type="text" value={noun} onChange={this.onNounChange} autoComplete="off" />
        </label>
        <button className="buttonTypeA" type="button" onClick={this.onSubmitPromptButtonClicked}>Submit</button>
      </div>
    );

    const filteredPlayers = players.filter((player) => player.hasPendingAction === true);
    const playersPendingAction = filteredPlayers.map((player) => (
      <li key={player.name}>
        {player.name}
      </li>
    ));

    const waitingElements = (
      <div>
        <h3 className="initialPromptCreationScreen">Thanks, waiting for these players to finish submitting their prompts:</h3>
        <ul>{playersPendingAction}</ul>
      </div>
    );


    return (
      <div className="screen initialPromptCreationScreen">
        {currentPlayer.hasCompletedAction ? waitingElements : promptCreatingElements}
        <h3 className="error">{error}</h3>
      </div>
    );
  }
}

InitialPromptCreationScreen.propTypes = {
  gameState: PropTypes.shape({
    currentPlayer: PropTypes.shape({
      name: PropTypes.string.isRequired,
      isHost: PropTypes.bool.isRequired,
      hasCompletedAction: PropTypes.bool.isRequired,
    }).isRequired,
    players: PropTypes.arrayOf(PropTypes.shape({
      name: PropTypes.string.isRequired,
    })),
    groupName: PropTypes.string.isRequired,
  }).isRequired,
  onGameStateChanged: PropTypes.func.isRequired,
};


export default InitialPromptCreationScreen;
