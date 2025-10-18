import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useNamespaces } from '../hooks/useNamespaces';
import { getNamespaces } from '../resources/namespaceResource';
import { transformNamespacesToDomain } from '../adapters/clusterAdapter';

vi.mock('../resources/namespaceResource');
vi.mock('../adapters/clusterAdapter');

describe('useNamespaces', () => {
  const mockContext = 'kind-kind';

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should fetch namespaces successfully', async () => {
    const mockNamespaces = [{ name: 'default', status: 'Active' }];
    const mockTransformedNamespaces = [{ name: 'default', status: 'Active' }];
    
    vi.mocked(getNamespaces).mockResolvedValue({ data: mockNamespaces } as any);
    vi.mocked(transformNamespacesToDomain).mockReturnValue(mockTransformedNamespaces);

    const { result } = renderHook(() => useNamespaces(mockContext));

    expect(result.current.loading).toBe(true);
    expect(result.current.namespaces).toEqual([]);

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(getNamespaces).toHaveBeenCalledWith(mockContext);
    expect(transformNamespacesToDomain).toHaveBeenCalledWith(mockNamespaces);
    expect(result.current.namespaces).toEqual(mockTransformedNamespaces);
    expect(result.current.error).toBeNull();
  });

  it('should handle fetch error', async () => {
    const error = new Error('Network error');
    vi.mocked(getNamespaces).mockRejectedValue(error);

    const { result } = renderHook(() => useNamespaces(mockContext));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.error).toBe('Network error');
    expect(result.current.namespaces).toEqual([]);
  });

  it('should not fetch when context is empty', async () => {
    const { result } = renderHook(() => useNamespaces(''));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(getNamespaces).not.toHaveBeenCalled();
    expect(result.current.namespaces).toEqual([]);
  });

  it('should refetch when context changes', async () => {
    const { result, rerender } = renderHook(
      ({ context }) => useNamespaces(context),
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
      expect(getNamespaces).toHaveBeenCalledWith(newContext);
    });
  });
});
