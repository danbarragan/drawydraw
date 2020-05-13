import React from 'react';
import axios from 'axios';
import PropTypes from 'prop-types';
import { formatServerError } from '../../utils/errorFormatting';


class InitialPromptCreation extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
    };
  }

  render() {
    return (
      <div className="InitialPromptCreation">
        Adjective, adjective, noun.
      </div>
    );
  }
}

InitialPromptCreation.propTypes = {
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

export default InitialPromptCreation;
