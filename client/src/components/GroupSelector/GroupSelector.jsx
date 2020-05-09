import React from 'react';
import axios from 'axios';
import './GroupSelector.css';
import PropTypes from 'prop-types';
import { formatServerError } from '../../utils/errorFormatting';

class RoomSelector extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      roomSelectionError: null,
      playerName: '',
      groupName: '',
    };
    this.onJoinGroupClick = this.onJoinGroupClick.bind(this);
    this.onCreateGroupClick = this.onCreateGroupClick.bind(this);
    this.onPlayerNameChange = this.onPlayerNameChange.bind(this);
    this.onGroupNameChange = this.onGroupNameChange.bind(this);
  }

  async onJoinGroupClick() {
    const { onGroupSelected } = this.props;
    const { playerName, groupName } = this.state;
    const data = { playerName, groupName };
    let response = null;
    try {
      // Todo: Probably worth renaming this endpoint join-group
      response = await axios.post('/api/add-player', data);
      onGroupSelected(response.data);
      if (data.kittens === response.data.kittens) {
        this.setState({ roomSelectionError: 'happy' });
      }
    } catch (error) {
      this.setState({ roomSelectionError: formatServerError(error) });
    }
  }

  async onCreateGroupClick() {
    const { onGroupSelected } = this.props;
    const { playerName, groupName } = this.state;
    const data = { playerName, groupName };
    let response = null;
    try {
      response = await axios.post('api/create-group', data);
      onGroupSelected(response.data);
    } catch (error) {
      this.setState({ roomSelectionError: formatServerError(error) });
    }
  }

  onPlayerNameChange(event) {
    this.setState({ playerName: event.target.value });
  }

  onGroupNameChange(event) {
    this.setState({ groupName: event.target.value });
  }

  render() {
    const { roomSelectionError, playerName, groupName } = this.state;
    return (
      <div className="login">
        <h1>Drawy draw</h1>
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
        <h3 className="error">{roomSelectionError}</h3>
      </div>
    );
  }
}

RoomSelector.propTypes = {
  onGroupSelected: PropTypes.func.isRequired,
};

export default RoomSelector;
