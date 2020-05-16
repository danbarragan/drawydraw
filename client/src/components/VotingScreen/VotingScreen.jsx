import React from 'react';
import axios from 'axios';
import PropTypes from 'prop-types';
import { formatServerError } from '../../utils/errorFormatting';
import './VotingScreen.css';

class VotingScreen extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      voted: false,
      selectedOption: '',
    };
    this.updateGameState = this.updateGameState.bind(this);
    this.handleOptionChange = this.handleOptionChange.bind(this);
    this.castVote = this.castVote.bind(this);
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

  handleOptionChange(voteEvent) {
    const { voted } = this.state;
    if (voted) {
      return;
    }
    this.setState({
      selectedOption: voteEvent.target.value,
      error: '',
    });
  }

  castVote() {
    const { selectedOption } = this.state;
    if (!selectedOption) {
      this.setState({ error: 'select an option!   ' });
      return;
    }

    this.setState({
      voted: true,
    });
  }

  render() {
    const { error } = this.state;
    const { selectedOption } = this.state;
    const { voted } = this.state;
    const prompts = ['Tiny, Ugly, Duckling', 'Large, Fugly, Duckling', 'Awkard, Hairy, Duckling'];
    const options = [];
    for (let i = 0; i < prompts.length; i += 1) {
      const optionName = `option${i.toString()}`;
      const isSelected = selectedOption === optionName;
      let className = '';
      if (voted) {
        className = isSelected ? 'votedFor' : 'notVotedFor';
      } else {
        className = 'notVotedYet';
      }
      options.push(
        <div className={className} key={optionName}>
          <label htmlFor={optionName}>
            <input
              id={optionName}
              type="radio"
              value={optionName}
              checked={isSelected}
              onChange={this.handleOptionChange}
            />
            {prompts[i]}
          </label>
        </div>,
      );
    }
    return (
      <div className="voteSelection">
        <h1>Add Drawing HERE</h1>
        <form>
          {options}
        </form>
        <button disabled={voted} className="voteButton" type="button" onClick={this.castVote}>Vote</button>
        <h3 className="error">{error}</h3>
        <h2>X/Y Votes cast.</h2>
      </div>
    );
  }
}

VotingScreen.propTypes = {
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

export default VotingScreen;
