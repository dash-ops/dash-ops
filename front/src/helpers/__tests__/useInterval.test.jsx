import { render, cleanup } from '@testing-library/react';
import useInterval from '../useInterval';
import React from 'react';

vi.useFakeTimers();

afterEach(cleanup);

it('should callback after delay', () => {
  const callback = vi.fn();

  const TestHookComponent = ({ callback }) => {
    useInterval(callback, 100);
    return null;
  };

  render(<TestHookComponent callback={callback} />);

  expect(callback).toHaveBeenCalledTimes(0);
  vi.advanceTimersByTime(400);
  expect(callback).toHaveBeenCalledTimes(4);
});
