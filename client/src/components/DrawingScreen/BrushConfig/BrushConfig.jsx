/* eslint-disable jsx-a11y/click-events-have-key-events */
/* eslint-disable jsx-a11y/no-static-element-interactions */
import React from 'react';
import './BrushConfig.css';
import PropTypes from 'prop-types';
import { FormattedMessage } from 'react-intl';


const BrushColors = Object.freeze({
  // Map enum members to HTML colors
  Black: 'black',
  Grey: 'grey',
  White: 'white',
  Greenery: '#88B04B',
  Coral: '#FF6F61',
  Violet: '#6B5B95',
  Marsala: '#955251',
  Orchid: '#B565A7',
  Turquoise: '#45B8AC',
  Mimosa: '#EFC050',
});

const BrushSizes = Object.freeze({
  Small: { name: 'S', weight: 5 },
  Medium: { name: 'M', weight: 10 },
  Large: { name: 'L', weight: 20 },
  XL: { name: 'XL', weight: 40 },
});

const BrushConfig = (props) => {
  const {
    currentColor, onColorChange, currentSize, onWidthChange,
  } = props;
  return (
    <div className="colorPicker">
      <h3><FormattedMessage id="colorPicker.brushColor" defaultMessage="Color:" /></h3>
      {
          Object.values(BrushColors).map((color) => {
            const highlightClass = color === currentColor ? 'highlightedBrushColor' : '';
            return (
              <div key={`${color}Container`} className="brushColorContainer">
                <div
                  key={color}
                  className={`brushColor ${highlightClass}`}
                  style={{ backgroundColor: color }}
                  onClick={() => onColorChange(color)}
                />
              </div>
            );
          })
      }
      <h3><FormattedMessage id="colorPicker.brushSize" defaultMessage="Size:" /></h3>
      {
      Object.values(BrushSizes).map((size) => {
        const highlightClass = size === currentSize ? 'highlightedBrushSize' : '';
        return (
          <div key={`${size.name}Container`} className="brushSizeContainer">
            <div
              key={size.name}
              className={`brushSize ${highlightClass}`}
              onClick={() => onWidthChange(size)}
            >
              <p>{size.name}</p>
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
  currentSize: PropTypes.shape({
    name: PropTypes.string.isRequired,
    weight: PropTypes.number.isRequired,
  }).isRequired,
  onWidthChange: PropTypes.func.isRequired,
};

export { BrushConfig, BrushColors, BrushSizes };
