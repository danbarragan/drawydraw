// Link.react.test.js
import React from 'react';
import { mount, shallow } from 'enzyme';
import axios from 'axios';
import DrawingScreen, { Point, Stroke } from './DrawingScreen';
import '../../test/setupTests';

describe('DrawingScreen', () => {
  const mockGameState = { groupName: 'kitties4Life', currentPlayer: { name: 'baby cat' } };
  const mockOnGameStateChanged = jest.fn();

  it('renders strokes as a data URL', () => {
    const expectedDataURL = 'data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAfQAAAH0CAYAAADL1t+KAAAABmJLR0QA/wD/AP+gvaeTAAAD30lEQVR4nO3BAQEAAACCIP+vbkhAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA8GJFFQABGYPuoAAAAABJRU5ErkJggg==';
    const screen = mount(
      <DrawingScreen onGameStateChanged={mockOnGameStateChanged} gameState={mockGameState} />,
    );
    const strokes = [new Stroke([new Point(1, 5), new Point(1, 5.5), new Point(1, 6)])];
    screen.setState({ strokes });
    const imageData = screen.instance().renderStrokesAsDataURL();
    expect(imageData).toBe(expectedDataURL);
  });

  it('posts canvas image to server', async () => {
    const mockResponseData = { foo: 1 };
    axios.post = jest.fn(() => ({ data: mockResponseData }));
    const screen = shallow(
      <DrawingScreen onGameStateChanged={mockOnGameStateChanged} gameState={mockGameState} />,
    );
    // Mock out converting the canvas contents into an image
    const mockToDataURL = jest.fn(() => 'img/png SOME_DATA');
    const canvasContainer = { canvas: { toDataURL: mockToDataURL } };
    screen.setState({ canvasContainer });
    await screen.instance().onSubmitClick();
    const expectedData = { groupName: 'kitties4Life', playerName: 'baby cat', imageData: 'img/png SOME_DATA' };
    expect(axios.post).toHaveBeenCalledWith('api/submit-drawing', expectedData);
    expect(mockOnGameStateChanged).toHaveBeenCalledWith(mockResponseData);
  });
});
