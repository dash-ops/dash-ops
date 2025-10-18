import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useClusters } from '../hooks/useClusters';
import { getClusters } from '../resources/clusterResource';
import { transformClusterListResponseToDomain } from '../adapters/clusterAdapter';

vi.mock('../resources/clusterResource');
vi.mock('../adapters/clusterAdapter');

describe('useClusters', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should fetch clusters successfully', async () => {
    const mockResponse = { clusters: [{ name: 'kind-kind' }], total: 1 };
    const mockTransformedResponse = { clusters: [{ name: 'kind-kind' }], total: 1 };
    
    vi.mocked(getClusters).mockResolvedValue({ data: mockResponse } as any);
    vi.mocked(transformClusterListResponseToDomain).mockReturnValue(mockTransformedResponse);

    const { result } = renderHook(() => useClusters());

    expect(result.current.loading).toBe(true);
    expect(result.current.clusters).toEqual([]);

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(getClusters).toHaveBeenCalled();
    expect(transformClusterListResponseToDomain).toHaveBeenCalledWith(mockResponse);
    expect(result.current.clusters).toEqual(mockTransformedResponse.clusters);
    expect(result.current.error).toBeNull();
  });

  it('should handle fetch error', async () => {
    const error = new Error('Network error');
    vi.mocked(getClusters).mockRejectedValue(error);

    const { result } = renderHook(() => useClusters());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.error).toBe('Network error');
    expect(result.current.clusters).toEqual([]);
  });

  it('should refetch when fetchClusters is called', async () => {
    const { result } = renderHook(() => useClusters());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    const initialCallCount = vi.mocked(getClusters).mock.calls.length;

    await act(async () => {
      await result.current.fetchClusters();
    });

    expect(vi.mocked(getClusters).mock.calls.length).toBe(initialCallCount + 1);
  });
});
