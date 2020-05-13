import React from 'react';
import axios from 'axios';
import PropTypes from 'prop-types';
import { formatServerError } from '../../utils/errorFormatting';


class InitialPromptCreationScreen extends React.Component {
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
    return (
      <div className="InitialPromptCreationScreen">
        Adjective, adjective, noun.
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
    }).isRequired,
    players: PropTypes.arrayOf(PropTypes.shape({
      name: PropTypes.string.isRequired,
    })),
    groupName: PropTypes.string.isRequired,
  }).isRequired,
  onGameStateChanged: PropTypes.func.isRequired,
};

export default InitialPromptCreationScreen;
