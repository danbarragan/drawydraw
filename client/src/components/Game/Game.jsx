import React from 'react';
import axios from 'axios';
import { FormattedMessage, IntlProvider } from 'react-intl';
import DecoyPromptCreationScreen from '../DecoyPromptCreationScreen/DecoyPromptCreationScreen';
import DrawingScreen from '../DrawingScreen/DrawingScreen';
import GroupSelectionScreen from '../GroupSelectionScreen/GroupSelectionScreen';
import WaitingForPlayersScreen from '../WaitingForPlayersScreen/WaitingForPlayersScreen';
import InitialPromptCreationScreen from '../InitialPromptCreationScreen/InitialPromptCreationScreen';
import VotingScreen from '../VotingScreen/VotingScreen';
import ScoringScreen from '../ScoringScreen/ScoringScreen';
import './Game.css';
import { GameStates } from '../../utils/constants';
import { formatServerError } from '../../utils/errorFormatting';


import EnglishMessages from '../../translations/en.json';
import SpanishMessages from '../../translations/es.json';

const messages = {
  es: SpanishMessages,
  en: EnglishMessages,
};

class Game extends React.Component {
  constructor(props) {
    // Split the browser's locale string to get the language without the region
    const currentLanguage = navigator.language.split(/[-_]/)[0];
    super(props);
    this.state = {
      currentLanguage,
      consoleEnabled: false,
      gameState: {
        currentState: GameStates.GroupSelection, // Start in the group selection state
      },
    };
    this.getCurrentComponent = this.getCurrentComponent.bind(this);
    this.onGameEntered = this.onGameEntered.bind(this);
    this.onGameStateChanged = this.onGameStateChanged.bind(this);
    this.debugConsole = this.debugConsole.bind(this);
    this.toggleConsole = this.toggleConsole.bind(this);
    this.debugSetGameState = this.debugSetGameState.bind(this);
    this.onLocaleChanged = this.onLanguageSelected.bind(this);
  }

  onLanguageSelected(currentLanguage) {
    this.setState({ currentLanguage });
  }

  onGameEntered(gameState) {
    this.setState({ gameState });
  }

  onGameStateChanged(gameState) {
    this.setState({ gameState });
  }

  getCurrentComponent() {
    const { gameState } = this.state;
    const { currentState } = gameState;
    switch (currentState) {
      case GameStates.Scoring:
        return <ScoringScreen onGameStateChanged={this.onGameStateChanged} gameState={gameState} />;
      case GameStates.DecoyPromptCreation:
        return (
          <DecoyPromptCreationScreen
            onGameStateChanged={this.onGameStateChanged}
            gameState={gameState}
          />
        );
      case GameStates.Voting:
        return <VotingScreen onGameStateChanged={this.onGameStateChanged} gameState={gameState} />;
      case GameStates.DrawingsInProgress:
        return <DrawingScreen onGameStateChanged={this.onGameStateChanged} gameState={gameState} />;
      case GameStates.GroupSelection:
        return <GroupSelectionScreen onGameEntered={this.onGameEntered} />;
      case GameStates.WaitingForPlayers:
        return (
          <WaitingForPlayersScreen
            gameState={gameState}
            onGameStateChanged={this.onGameStateChanged}
          />
        );
      case GameStates.InitialPromptCreation:
        return (
          <InitialPromptCreationScreen
            gameState={gameState}
            onGameStateChanged={this.onGameStateChanged}
          />
        );
      // Unknown group! Badness
      default:
        return (
          <div>
            <h1>
              <FormattedMessage
                id="game.stateNotImplementedError"
                defaultMessage="We are sorry this is not implemented yet :("
              />
            </h1>
          </div>
        );
    }
  }

  toggleConsole() {
    const { consoleEnabled } = this.state;
    this.setState({ consoleEnabled: !consoleEnabled });
  }

  async debugSetGameState(gameStateName) {
    try {
      const response = await axios.post('/api/set-game-state', { gameStateName });
      this.setState({ gameState: response.data });
      this.setState({ error: '' });
      this.forceUpdate();
    } catch (error) {
      this.setState({ error: formatServerError(error) });
    }
  }

  debugConsole() {
    const { error } = this.state;
    const { gameState, consoleEnabled } = this.state;
    return (
      <div className="debug">
        <button className="debugButton buttonTypeA" type="button" onClick={this.toggleConsole}>debug</button>
        {consoleEnabled ? (
          <div className="console">
            <span>States:</span>
            {Object.values(GameStates).map((state) => (state !== GameStates.GroupSelection ? (
              <button key={state} className="debugButton buttonTypeB" type="button" onClick={(() => this.debugSetGameState(state))}>{state}</button>
            ) : null))}
            {error ? `Error:${error}` : null}
            <textarea className="debugGameState" value={JSON.stringify(gameState, null, 4)} disabled />
          </div>
        ) : null}
      </div>
    );
  }

  languagePicker() {
    const { currentLanguage } = this.state;
    const languageButtons = Object.keys(messages).map((language) => {
      const buttonStyle = language !== currentLanguage ? 'buttonTypeB' : 'buttonTypeA';
      return (
        <button
          key={language}
          className={`languageButton ${buttonStyle}`}
          type="button"
          onClick={(() => this.onLanguageSelected(language))}
        >
          {language}
        </button>
      );
    });
    return (
      <div className="languagePicker">
        <span>
          <FormattedMessage
            id="game.languageHeader"
            defaultMessage="Language:"
          />
        </span>
        {languageButtons}
      </div>
    );
  }

  render() {
    const { currentLanguage } = this.state;
    return (
      <IntlProvider locale={currentLanguage} messages={messages[currentLanguage]}>
        <div className="game">
          <div className="gameTitle">
            <h2>
              <FormattedMessage
                id="game.title"
                defaultMessage="Drawydraw"
              />
            </h2>
          </div>
          {this.getCurrentComponent()}
          {this.languagePicker()}
          {this.debugConsole()}
        </div>
      </IntlProvider>
    );
  }
}

export default Game;
