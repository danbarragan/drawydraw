import React from 'react';
import Sketch from 'react-p5';
import './DrawingScreen.css';
import axios from 'axios';
import PropTypes from 'prop-types';
import { formatServerError } from '../../utils/errorFormatting';

// Tiny classes to make managing points easier
function Point(x, y) {
  this.x = x;
  this.y = y;
}

function Stroke(points) {
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
    };
    this.mousePressed = this.mousePressed.bind(this);
    this.mouseDragged = this.mouseDragged.bind(this);
    this.renderCanvas = this.renderCanvas.bind(this);
    this.onSubmitClick = this.onSubmitClick.bind(this);
    this.setupCanvas = this.setupCanvas.bind(this);
    this.renderStrokesAsDataURL = this.renderStrokesAsDataURL.bind(this);
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

  setupCanvas(p5, canvasParentRef) {
    const canvasContainer = p5.createCanvas(500, 500).parent(canvasParentRef);
    this.setState({ canvasContainer });
  }

  mousePressed(event) {
    const { mouseX, mouseY, canvas } = event;
    if (DrawingScreen.isPointInCanvas(mouseX, mouseY, canvas)) {
      // Add a new stroke to the set of strokes starting at the current mouse location
      let { strokes } = this.state;
      strokes = [...strokes, new Stroke([new Point(mouseX, mouseY)])];
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
    p5.stroke('black');
    p5.strokeWeight(3);
    const { strokes } = this.state;
    strokes.forEach((currentStroke) => {
      p5.beginShape();
      currentStroke.points.forEach((point) => {
        p5.curveVertex(point.x, point.y);
      });
      p5.endShape();
    });
  }

  render() {
    const { error } = this.state;
    return (
      <div className="screen">
        <h1>
          Draw some prompt
        </h1>
        <Sketch className="drawingCanvas" setup={this.setupCanvas} draw={this.renderCanvas} mouseDragged={this.mouseDragged} mousePressed={this.mousePressed} />
        <button type="button" className="button buttonTypeA" onClick={this.onSubmitClick}>Submit</button>
        <h3 className="error">{error}</h3>
      </div>
    );
  }
}

DrawingScreen.propTypes = {
  gameState: PropTypes.shape({
    currentPlayer: PropTypes.shape({
      name: PropTypes.string.isRequired,
    }).isRequired,
    groupName: PropTypes.string.isRequired,
  }).isRequired,
  onGameStateChanged: PropTypes.func.isRequired,
};

export { Stroke, Point, DrawingScreen as default };
