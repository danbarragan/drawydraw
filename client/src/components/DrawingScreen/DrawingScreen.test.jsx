import React from 'react';
import { shallow } from 'enzyme';
import axios from 'axios';
import DrawingScreen, { Point, Stroke } from './DrawingScreen';
import '../../test/setupTests';

describe('DrawingScreen', () => {
  const mockGameState = { groupName: 'kitties4Life', currentPlayer: { name: 'baby cat' } };
  const mockOnGameStateChanged = jest.fn();
  const mockStrokes = [new Stroke([new Point(1, 5), new Point(1, 5.5), new Point(1, 6)], 'black', 3)];

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

  it('clears strokes when clear button is pressed', () => {
    const screen = shallow(
      <DrawingScreen onGameStateChanged={mockOnGameStateChanged} gameState={mockGameState} />,
    );
    screen.instance().setState({ strokes: mockStrokes });
    expect(screen.instance().state.strokes).toEqual(mockStrokes);
    screen.instance().onClearClick();
    expect(screen.instance().state.strokes).toEqual([]);
  });
});
