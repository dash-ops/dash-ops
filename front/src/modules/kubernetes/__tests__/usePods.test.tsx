import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { usePods } from '../hooks/usePods';
import { getPods } from '../resources/podsResource';
import { transformPodsToDomain } from '../adapters/podAdapter';

vi.mock('../resources/podsResource');
vi.mock('../adapters/podAdapter');

describe('usePods', () => {
  const mockFilter = { context: 'kind-kind', namespace: 'default' };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should fetch pods successfully', async () => {
    const mockPods = [{ name: 'pod-1', namespace: 'default' }];
    const mockTransformedPods = [{ name: 'pod-1', namespace: 'default' }];
    
    vi.mocked(getPods).mockResolvedValue({ data: mockPods } as any);
    vi.mocked(transformPodsToDomain).mockReturnValue(mockTransformedPods);

    const { result } = renderHook(() => usePods(mockFilter));

    expect(result.current.loading).toBe(true);
    expect(result.current.pods).toEqual([]);

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(getPods).toHaveBeenCalledWith(mockFilter);
    expect(transformPodsToDomain).toHaveBeenCalledWith(mockPods);
    expect(result.current.pods).toEqual(mockTransformedPods);
    expect(result.current.error).toBeNull();
  });

  it('should handle fetch error', async () => {
    const error = new Error('Network error');
    vi.mocked(getPods).mockRejectedValue(error);

    const { result } = renderHook(() => usePods(mockFilter));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.error).toBe('Network error');
    expect(result.current.pods).toEqual([]);
  });

  it('should not fetch when context or namespace is missing', async () => {
    const { result } = renderHook(() => usePods({ context: '', namespace: 'default' }));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(getPods).not.toHaveBeenCalled();
    expect(result.current.pods).toEqual([]);
  });

  it('should refetch when filter changes', async () => {
    const { result, rerender } = renderHook(
      ({ filter }) => usePods(filter),
      { initialProps: { filter: mockFilter } }
    );

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    const newFilter = { context: 'kind-kind', namespace: 'kube-system' } as any;
    await act(async () => {
      rerender({ filter: newFilter });
    });

    await waitFor(() => {
      expect(getPods).toHaveBeenCalledWith(newFilter);
    });
  });
});
