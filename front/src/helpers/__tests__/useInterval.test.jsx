import { render, cleanup } from '@testing-library/react';
import PropTypes from 'prop-types';
import useInterval from '../useInterval';

vi.useFakeTimers();

afterEach(cleanup);

it('should callback after delay', () => {
  const callback = vi.fn();

  const TestHookComponent = ({ callback }) => {
    useInterval(callback, 100);
    return null;
  };
  
  TestHookComponent.propTypes = {
    callback: PropTypes.func.isRequired,
  };

  render(<TestHookComponent callback={callback} />);

  expect(callback).toHaveBeenCalledTimes(0);
  vi.advanceTimersByTime(400);
  expect(callback).toHaveBeenCalledTimes(4);
});
