import React from 'react';
import axios from 'axios';
import './GroupSelectionScreen.css';
import PropTypes from 'prop-types';
import { formatServerError } from '../../utils/errorFormatting';

class GroupSelectionScreen extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      playerName: '',
      groupName: '',
    };
    this.onJoinGroupClick = this.onJoinGroupClick.bind(this);
    this.onCreateGroupClick = this.onCreateGroupClick.bind(this);
    this.onPlayerNameChange = this.onPlayerNameChange.bind(this);
    this.onGroupNameChange = this.onGroupNameChange.bind(this);
  }

  async onJoinGroupClick() {
    const { onGameEntered } = this.props;
    const { playerName, groupName } = this.state;
    const data = { playerName, groupName };
    try {
      // Todo: Probably worth renaming this endpoint join-group
      const response = await axios.post('/api/add-player', data);
      onGameEntered(response.data);
    } catch (error) {
      this.setState({ error: formatServerError(error) });
    }
  }

  async onCreateGroupClick() {
    const { onGameEntered } = this.props;
    const { playerName, groupName } = this.state;
    const data = { playerName, groupName };
    try {
      const response = await axios.post('api/create-game', data);
      onGameEntered(response.data);
    } catch (error) {
      this.setState({ error: formatServerError(error) });
    }
  }

  onPlayerNameChange(event) {
    this.setState({ playerName: event.target.value });
  }

  onGroupNameChange(event) {
    this.setState({ groupName: event.target.value });
  }

  render() {
    const { error, playerName, groupName } = this.state;
    return (
      <div className="groupSelectionScreen">
        <h3>Join or create a group</h3>
        <label htmlFor="playerName">
          Your name
          <input id="playerName" type="text" value={playerName} onChange={this.onPlayerNameChange} />
        </label>
        <label htmlFor="groupName">
          Group name
          <input id="groupName" type="text" value={groupName} onChange={this.onGroupNameChange} />
        </label>
        <button type="button" onClick={this.onJoinGroupClick}>Join group</button>
        <button type="button" onClick={this.onCreateGroupClick}>Create group</button>
        <h3 className="error">{error}</h3>
      </div>
    );
  }
}

GroupSelectionScreen.propTypes = {
  onGameEntered: PropTypes.func.isRequired,
};

export default GroupSelectionScreen;