import React from 'react';
// import axios from 'axios';
// import PropTypes from 'prop-types';
import Sketch from 'react-p5';
// import { formatServerError } from '../../utils/errorFormatting';
import './DrawingScreen.css';

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
  constructor(props) {
    super(props);
    this.state = {
      strokes: [],
    };
    this.mousePressed = this.mousePressed.bind(this);
    this.mouseDragged = this.mouseDragged.bind(this);
    this.renderCanvas = this.renderCanvas.bind(this);
  }

  setup(p5, canvasParentRef) {
    p5.createCanvas(500, 500).parent(canvasParentRef); // use parent to render canvas in this ref (without that p5 render this canvas outside your component)
  }

  mousePressed(event) {
    const { mouseX, mouseY } = event;
    // Add a new stroke to the set of strokes starting at the current mouse location
    let { strokes } = this.state;
    strokes = [...strokes, new Stroke([new Point(mouseX, mouseY)])];
    this.setState({ strokes });
  }

  mouseDragged(event) {
    const { mouseX, mouseY } = event;
    // Append the mouse's position to the most recent stroke
    const { strokes } = this.state;
    strokes[strokes.length - 1].addPoint(new Point(mouseX, mouseY));
    this.setState({ strokes });
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
    return (
      <div className="screen ">
        <h1>
          Draw some prompt
        </h1>
        <Sketch className="drawingCanvas" setup={this.setup} draw={this.renderCanvas} mouseDragged={this.mouseDragged} mousePressed={this.mousePressed} />
      </div>
    );
  }
}

DrawingScreen.propTypes = {
//   gameState: PropTypes.shape({
//     currentPlayer: PropTypes.shape({
//       name: PropTypes.string.isRequired,
//       isHost: PropTypes.bool.isRequired,
//     }).isRequired,
//     players: PropTypes.arrayOf(PropTypes.shape({
//       name: PropTypes.string.isRequired,
//     })),
//     groupName: PropTypes.string.isRequired,
//   }).isRequired,
//   onGameStateChanged: PropTypes.func.isRequired,
};

export default DrawingScreen;
