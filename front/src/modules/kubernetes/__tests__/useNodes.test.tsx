import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useNodes } from '../hooks/useNodes';
import { getNodes } from '../resources/nodesResource';
import { transformNodesToDomain } from '../adapters/nodeAdapter';

vi.mock('../resources/nodesResource');
vi.mock('../adapters/nodeAdapter');

describe('useNodes', () => {
  const mockContext = 'kind-kind';

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should fetch nodes successfully', async () => {
    const mockNodes = [{ name: 'node-1', status: 'Ready' }];
    const mockTransformedNodes = [{ name: 'node-1', status: 'Ready' }];
    
    vi.mocked(getNodes).mockResolvedValue({ data: mockNodes } as any);
    vi.mocked(transformNodesToDomain).mockReturnValue(mockTransformedNodes);

    const { result } = renderHook(() => useNodes(mockContext));

    expect(result.current.loading).toBe(true);
    expect(result.current.nodes).toEqual([]);

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(getNodes).toHaveBeenCalledWith(mockContext);
    expect(transformNodesToDomain).toHaveBeenCalledWith(mockNodes);
    expect(result.current.nodes).toEqual(mockTransformedNodes);
    expect(result.current.error).toBeNull();
  });

  it('should handle fetch error', async () => {
    const error = new Error('Network error');
    vi.mocked(getNodes).mockRejectedValue(error);

    const { result } = renderHook(() => useNodes(mockContext));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.error).toBe('Network error');
    expect(result.current.nodes).toEqual([]);
  });

  it('should not fetch when context is empty', async () => {
    const { result } = renderHook(() => useNodes(''));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(getNodes).not.toHaveBeenCalled();
    expect(result.current.nodes).toEqual([]);
  });

  it('should refetch when context changes', async () => {
    const { result, rerender } = renderHook(
      ({ context }) => useNodes(context),
      { initialProps: { context: mockContext } }
    );

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    const newContext = 'kind-kind-2';
    await act(async () => {
      rerender({ context: newContext });
    });

    await waitFor(() => {
      expect(getNodes).toHaveBeenCalledWith(newContext);
    });
  });
});
