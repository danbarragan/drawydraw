import React from 'react';
import Sketch from 'react-p5';
import './DrawingScreen.css';
import axios from 'axios';
import PropTypes from 'prop-types';
import { formatServerError } from '../../utils/errorFormatting';
import { BrushColors, BrushConfig, BrushSizes } from './BrushConfig/BrushConfig';

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
    const canvasContainer = p5.createCanvas(900, 900).parent(canvasParentRef);
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
      p5.beginShape();
      // Draw an individual point if the stroke only has one point
      if (currentStroke.points.length === 1) {
        const point = currentStroke.points[0];
        p5.point(point.x, point.y);
      } else {
        currentStroke.points.forEach((point) => {
          p5.curveVertex(point.x, point.y);
        });
      }
      p5.endShape();
    });
  }

  render() {
    const { error, currentBrushColor, currentBrushSize } = this.state;
    const { gameState } = this.props;
    const {currentPlayer} = gameState;
    const { noun, adjectives } = currentPlayer.assignedPrompt;
    const { players } = gameState;
    const drawingElements = (
      <div>
        <h1>
          Draw
          {' '}
          {adjectives[0]}
          ,
          {' '}
          {adjectives[1]}
          {' '}
          {noun}
        </h1>
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
        <button type="button" className="button buttonTypeA" onClick={this.onSubmitClick}>Submit</button>
        <button type="button" className="button buttonTypeB" onClick={this.onClearClick}>Clear</button>
      </div>
    );
    const waitingElements = (
      <div>
        <h3>Thank you for your drawing, waiting for other players...</h3>
        <ul>
          {
            players.filter((player) => player.name !== currentPlayer.name).map((player) => (
              <li key={player.name}>
                {player.name}
                {player.hasPendingActions ? ' is still drawing' : ' is done'}
              </li>
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
      hasPendingActions: PropTypes.bool.isRequired,
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
