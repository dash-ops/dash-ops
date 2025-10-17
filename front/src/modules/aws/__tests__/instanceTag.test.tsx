import { render, screen, cleanup } from '@testing-library/react';
import { describe, it, expect, afterEach } from 'vitest';
import InstanceTag from '../components/instances/InstanceTag';

afterEach(cleanup);

describe('InstanceTag', () => {
  it('should render running state tag', () => {
    render(<InstanceTag state={{ name: 'running', code: 16 }} />);

    const tag = screen.getByText('running');
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveClass('bg-green-100', 'text-green-800');
  });

  it('should render stopped state tag', () => {
    render(<InstanceTag state={{ name: 'stopped', code: 80 }} />);

    const tag = screen.getByText('stopped');
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveClass('bg-red-100', 'text-red-800');
  });

  it('should render pending state tag', () => {
    render(<InstanceTag state={{ name: 'pending', code: 0 }} />);

    const tag = screen.getByText('pending');
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveClass('bg-yellow-100', 'text-yellow-800');
  });

  it('should render stopping state tag', () => {
    render(<InstanceTag state={{ name: 'stopping', code: 64 }} />);

    const tag = screen.getByText('stopping');
    expect(tag).toBeInTheDocument();
    expect(tag).toHaveClass('bg-purple-100', 'text-purple-800');
  });

  it('should render unknown state tag', () => {
    render(<InstanceTag state={{ name: 'unknown', code: 99 }} />);

    const tag = screen.getByText('unknown');
    expect(tag).toBeInTheDocument();
    // Unknown state doesn't have specific styles, just check if it renders
    expect(tag).toBeInTheDocument();
  });
});
