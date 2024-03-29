import React from 'react';
import { FormattedMessage } from 'react-intl';
import Sketch from 'react-p5';
import './DrawingScreen.css';
import axios from 'axios';
import PropTypes from 'prop-types';
import { formatServerError } from '../../utils/errorFormatting';
import { BrushColors, BrushConfig, BrushSizes } from './BrushConfig/BrushConfig';
import UpdateGameState from '../../utils/updateGameState';

// Tiny classes to make managing points easier
function Point(x, y) {
  this.x = x;
  this.y = y;
}

function Stroke(points, color, weight) {
  this.weight = weight;
  this.color = color;
  this.points = points || [];
  this.addPoint = (point) => {
    this.points.push(point);
  };
}

class DrawingScreen extends React.Component {
  // It turns out that the p5 listeners also listen to mouse events outside the canvas D:
  // so we need to figure out whether each event happened inside or outside the canvas
  static isPointInCanvas(x, y, canvas) {
    // Coords are relative to the canvas
    return (canvas.width >= x && x >= 0) && (canvas.height >= y && y >= 0);
  }

  constructor(props) {
    super(props);
    this.state = {
      strokes: [],
      currentBrushColor: BrushColors.Black,
      currentBrushSize: BrushSizes.Small,
      timerId: null,
    };
    this.mousePressed = this.mousePressed.bind(this);
    this.mouseDragged = this.mouseDragged.bind(this);
    this.renderCanvas = this.renderCanvas.bind(this);
    this.onSubmitClick = this.onSubmitClick.bind(this);
    this.onClearClick = this.onClearClick.bind(this);
    this.setupCanvas = this.setupCanvas.bind(this);
    this.renderStrokesAsDataURL = this.renderStrokesAsDataURL.bind(this);
    this.onBrushColorChange = this.onBrushColorChange.bind(this);
    this.onBrushSizeChange = this.onBrushSizeChange.bind(this);
    this.updateGameState = this.updateGameState.bind(this);
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

  async onSubmitClick() {
    const { onGameStateChanged, gameState } = this.props;
    const { groupName, currentPlayer } = gameState;
    const { name: playerName } = currentPlayer;
    // We might want to consider lossier compression if images are too chunky
    const imageData = this.renderStrokesAsDataURL();
    const data = { playerName, groupName, imageData };
    try {
      const response = await axios.post('api/submit-drawing', data);
      onGameStateChanged(response.data);
    } catch (error) {
      this.setState({ error: formatServerError(error) });
    }
  }

  onClearClick() {
    this.setState({ strokes: [] });
  }

  onBrushColorChange(currentBrushColor) {
    this.setState({ currentBrushColor });
  }

  onBrushSizeChange(currentBrushSize) {
    this.setState({ currentBrushSize });
  }

  setupCanvas(p5, canvasParentRef) {
    const canvasContainer = p5.createCanvas(700, 700).parent(canvasParentRef);
    this.setState({ canvasContainer });
  }

  mousePressed(event) {
    const { mouseX, mouseY, canvas } = event;
    const { currentBrushColor, currentBrushSize } = this.state;
    if (DrawingScreen.isPointInCanvas(mouseX, mouseY, canvas)) {
      // Add a new stroke to the set of strokes starting at the current mouse location
      let { strokes } = this.state;
      strokes = [
        ...strokes,
        new Stroke([new Point(mouseX, mouseY)], currentBrushColor, currentBrushSize.weight),
      ];
      this.setState({ strokes });
    }
  }

  mouseDragged(event) {
    const { mouseX, mouseY, canvas } = event;
    if (DrawingScreen.isPointInCanvas(mouseX, mouseY, canvas)) {
      // Append the mouse's position to the most recent stroke
      const { strokes } = this.state;
      strokes[strokes.length - 1].addPoint(new Point(mouseX, mouseY));
      this.setState({ strokes });
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

  renderStrokesAsDataURL() {
    const { canvasContainer } = this.state;
    return canvasContainer.canvas.toDataURL();
  }

  renderCanvas(p5) {
    p5.background('white');
    p5.noFill();
    const { strokes } = this.state;
    strokes.forEach((currentStroke) => {
      p5.stroke(currentStroke.color);
      p5.strokeWeight(currentStroke.weight);
      // Render strokes differently based on how many points they have
      switch (currentStroke.points.length) {
        // 1 point - draw a point
        case 1:
          p5.point(currentStroke.points[0].x, currentStroke.points[0].y);
          break;
        // 2 or 3 points - draw a line that passes through the first two points
        case 2:
        case 3: {
          const coordinates = currentStroke.points.reduce(
            (coords, point) => ([...coords, point.x, point.y]),
            [],
          );
          p5.line(...coordinates.slice(0, 4));
          break;
        }
        // 4+ points - create a shape with one vertex per point
        default:
          p5.beginShape();
          currentStroke.points.forEach((point) => {
            p5.curveVertex(point.x, point.y);
          });
          p5.endShape();
      }
    });
  }

  render() {
    const { error, currentBrushColor, currentBrushSize } = this.state;
    const { gameState } = this.props;
    const { currentPlayer, players } = gameState;
    const { noun, adjectives } = currentPlayer.assignedPrompt;
    const drawingElements = (
      <div>
        <h2>
          <FormattedMessage
            id="drawingScreen.drawHeader"
            defaultMessage="Draw: {adjective1} and {adjective2} {noun}"
            values={{ adjective1: adjectives[0], adjective2: adjectives[1], noun }}
          />
        </h2>
        <BrushConfig
          onColorChange={this.onBrushColorChange}
          currentColor={currentBrushColor}
          onWidthChange={this.onBrushSizeChange}
          currentSize={currentBrushSize}
        />
        <Sketch
          className="drawingCanvas"
          setup={this.setupCanvas}
          draw={this.renderCanvas}
          touchMoved={this.mouseDragged}
          touchStarted={this.mousePressed}
        />
        <button type="button" className="buttonTypeA" onClick={this.onSubmitClick}>
          <FormattedMessage
            id="common.submitButton"
            defaultMessage="Submit"
          />
        </button>
        <button type="button" className="buttonTypeB" onClick={this.onClearClick}>
          <FormattedMessage
            id="drawingScreen.clearButton"
            defaultMessage="Clear"
          />
        </button>
      </div>
    );
    const waitingElements = (
      <div>
        <h3>
          <FormattedMessage
            id="drawingScreen.waitingForOthersDrawingsMessage"
            defaultMessage="Waiting for these players to finish their drawings..."
          />
        </h3>
        <ul>
          {
            players.map((player) => (
              player.hasPendingAction ? <li key={player.name}>{player.name}</li> : null
            ))
          }
        </ul>
      </div>
    );
    return (
      <div className="screen">
        {currentPlayer.hasCompletedAction ? waitingElements : drawingElements}
        <h3 className="error">{error}</h3>
      </div>
    );
  }
}

DrawingScreen.propTypes = {
  gameState: PropTypes.shape({
    players: PropTypes.arrayOf(PropTypes.shape({
      name: PropTypes.string.isRequired,
      hasPendingAction: PropTypes.bool.isRequired,
    })),
    currentPlayer: PropTypes.shape({
      name: PropTypes.string.isRequired,
      assignedPrompt: PropTypes.shape({
        adjectives: PropTypes.arrayOf(PropTypes.string).isRequired,
        noun: PropTypes.string,
      }),
      hasCompletedAction: PropTypes.bool.isRequired,
    }).isRequired,
    groupName: PropTypes.string.isRequired,
  }).isRequired,
  onGameStateChanged: PropTypes.func.isRequired,
};

export { Stroke, Point, DrawingScreen as default };
