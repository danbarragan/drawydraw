import React from 'react';
import { mount } from 'enzyme';
import Game from './Game';
import '../../test/setupTests';
import { GameStates } from '../../utils/constants';

describe('Game Component', () => {
  const spyConsoleError = jest.spyOn(console, 'error');
  spyConsoleError.mockImplementation((message) => {
    if (message.lastIndexOf('Warning: Failed prop type:') === 0) {
      throw new Error(message);
    }
    return null;
  });

  // eslint-disable-next-line jest/expect-expect
  test.each(Object.values(GameStates))('can render state %p without crashing', (state) => {
    const game = mount(
      <Game />,
    );
    game.instance().setState({ currentState: state });
  });
});
