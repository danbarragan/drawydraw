/* eslint-disable jsx-a11y/click-events-have-key-events */
/* eslint-disable jsx-a11y/no-static-element-interactions */
import React from 'react';
import './BrushConfig.css';
import PropTypes from 'prop-types';

const Colors = Object.freeze({
  // Map enum members to HTML colors
  Black: 'black',
  White: 'white',
  Greenery: '#88B04B',
  Coral: '#FF6F61',
  Violet: '#6B5B95',
  Marsala: '#955251',
  Orchid: '#B565A7',
  Turquoise: '#5B8AC',
  Mimosa: '#EFC050',
});

const Widths = Object.freeze({
  Small: { name: 'S', weight: 5 },
  Medium: { name: 'M', weight: 10 },
  Large: { name: 'L', weight: 15 },
  XL: { name: 'XL', weight: 20 },
});

const BrushConfig = (props) => {
  const {
    currentColor, onColorChange, currentWidth, onWidthChange,
  } = props;
  return (
    <div className="colorPicker">
      <h3>Color</h3>
      {
          Object.values(Colors).map((color) => {
            const highlightClass = color === currentColor ? 'highlightedColor' : '';
            return (
              <div key={`${color}Container`} className="colorContainer">
                <div
                  key={color}
                  className={`color ${highlightClass}`}
                  style={{ backgroundColor: color }}
                  onClick={() => onColorChange(color)}
                />
              </div>
            );
          })
      }
      <h3>Width</h3>
      {
      Object.values(Widths).map((width) => {
        const highlightClass = width === currentWidth ? 'highlightedWidth' : '';
        return (
          <div key={`${width.name}Container`} className="widthContainer">
            <div
              key={width.name}
              className={`width ${highlightClass}`}
              onClick={() => onWidthChange(width)}
            >
              <p>{width.name}</p>
            </div>
          </div>
        );
      })
      }
    </div>
  );
};

BrushConfig.propTypes = {
  currentColor: PropTypes.string.isRequired,
  onColorChange: PropTypes.func.isRequired,
  currentWidth: PropTypes.shape({
    name: PropTypes.string.isRequired,
    weight: PropTypes.number.isRequired,
  }).isRequired,
  onWidthChange: PropTypes.func.isRequired,
};

export { BrushConfig, Colors, Widths };
