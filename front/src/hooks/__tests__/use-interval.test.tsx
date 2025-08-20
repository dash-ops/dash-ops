import { render, cleanup } from '@testing-library/react';
import { useInterval } from '../use-interval';

vi.useFakeTimers();

afterEach(cleanup);

it('should callback after delay', () => {
  const callback = vi.fn();

  const TestHookComponent = ({ callback }: { callback: () => void }) => {
    useInterval(callback, 100);
    return null;
  };

  render(<TestHookComponent callback={callback} />);

  expect(callback).toHaveBeenCalledTimes(0);
  vi.advanceTimersByTime(400);
  expect(callback).toHaveBeenCalledTimes(4);
});
