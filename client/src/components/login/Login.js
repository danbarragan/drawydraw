import React from 'react';
import axios from "axios";

class Login extends React.Component {
    constructor(props) {
        super(props);
        this.state = {serverStatus: "TBD"};
        this.handleButtonClick = this.handleButtonClick.bind(this);
    }

    async handleButtonClick() {
        const data = {kittens: "Five million"}
        let response = null;
        try {
            response = await axios.post("api/echo", data)
            if (data.kittens == response.data.kittens) {
                this.setState({serverStatus: "happy"});
            }
        } catch (e) {
            this.setState({serverStatus: `Sad: ${e}`})
        }
    }

    render() {
        return(
            <div>
                <h1>Drawy draw</h1><br />
                <label>Your name</label><br />
                <input type="text"></input><br />
                <label>Room name</label><br />
                <input type="text"></input><br />
                <button onClick={this.handleButtonClick}>Let's draw</button><br />
                <label>Server status: {this.state.serverStatus}</label>
            </div>
        );
    }

}

export default Login;
