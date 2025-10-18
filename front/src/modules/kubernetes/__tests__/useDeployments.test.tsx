import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useDeployments } from '../hooks/useDeployments';
import { getDeployments } from '../resources/deploymentResource';
import { transformDeploymentsToDomain } from '../adapters/deploymentAdapter';

vi.mock('../resources/deploymentResource');
vi.mock('../adapters/deploymentAdapter');

describe('useDeployments', () => {
  const mockFilter = { context: 'kind-kind', namespace: 'default' };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should fetch deployments successfully', async () => {
    const mockDeployments = [{ name: 'deploy-1', namespace: 'default' }];
    const mockTransformedDeployments = [{ name: 'deploy-1', namespace: 'default' }];
    
    vi.mocked(getDeployments).mockResolvedValue({ data: mockDeployments } as any);
    vi.mocked(transformDeploymentsToDomain).mockReturnValue(mockTransformedDeployments);

    const { result } = renderHook(() => useDeployments(mockFilter));

    expect(result.current.loading).toBe(true);
    expect(result.current.deployments).toEqual([]);

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(getDeployments).toHaveBeenCalledWith(mockFilter);
    expect(transformDeploymentsToDomain).toHaveBeenCalledWith(mockDeployments);
    expect(result.current.deployments).toEqual(mockTransformedDeployments);
    expect(result.current.error).toBeNull();
  });

  it('should handle fetch error', async () => {
    const error = new Error('Network error');
    vi.mocked(getDeployments).mockRejectedValue(error);

    const { result } = renderHook(() => useDeployments(mockFilter));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.error).toBe('Network error');
    expect(result.current.deployments).toEqual([]);
  });

  it('should not fetch when context or namespace is missing', async () => {
    const { result } = renderHook(() => useDeployments({ context: '', namespace: 'default' }));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(getDeployments).not.toHaveBeenCalled();
    expect(result.current.deployments).toEqual([]);
  });

  it('should refetch when filter changes', async () => {
    const { result, rerender } = renderHook(
      ({ filter }) => useDeployments(filter),
      { initialProps: { filter: mockFilter } }
    );

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    const newFilter = { context: 'kind-kind', namespace: 'kube-system' };
    await act(async () => {
      rerender({ filter: newFilter });
    });

    await waitFor(() => {
      expect(getDeployments).toHaveBeenCalledWith(newFilter);
    });
  });
});
