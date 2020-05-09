import React from 'react';
import axios from 'axios';
import './Login.css';

class Login extends React.Component {
  constructor(props) {
    super(props);
    this.state = { loginError: null };
    this.handleJoinGroupClick = this.handleJoinGroupClick.bind(this);
    this.handleCreateGroupClick = this.handleCreateGroupClick.bind(this);
  }

  async handleJoinGroupClick() {
    const data = { kittens: 'Five million' };
    let response = null;
    try {
      response = await axios.post('api/echo', data);
      if (data.kittens === response.data.kittens) {
        this.setState({ loginError: 'happy' });
      }
    } catch (e) {
      this.setState({ loginError: `Sad: ${e}` });
    }
  }

  async handleCreateGroupClick() {
    const data = { kittens: 'Five million' };
    let response = null;
    try {
      response = await axios.post('api/echo', data);
      if (data.kittens === response.data.kittens) {
        this.setState({ loginError: 'happy' });
      }
    } catch (e) {
      this.setState({ loginError: `Sad: ${e}` });
    }
  }

  render() {
    const { loginError } = this.state;
    return (
      <div className="login">
        <h1>Drawy draw</h1>
        <label>Your name</label>
        <input type="text" />
        <label>Room name</label>
        <input type="text" />
        <button type="button" onClick={this.handleJoinGroupClick}>Join group</button>
        <button type="button" onClick={this.handleCreateGroupClick}>Create group</button>
        <label className="error">
          {loginError}
        </label>
      </div>
    );
  }
}

export default Login;
