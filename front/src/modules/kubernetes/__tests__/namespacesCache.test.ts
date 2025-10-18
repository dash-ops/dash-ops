import { describe, it, expect, vi, beforeEach } from 'vitest';
import { getNamespacesCached, clearNamespacesCache } from '../utils/namespacesCache';
import { getNamespaces } from '../resources/namespaceResource';

vi.mock('../resources/namespaceResource');

describe('namespacesCache', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    clearNamespacesCache();
  });

  it('should fetch and cache namespaces for a context on first call', async () => {
    const mockNamespaces = [{ name: 'default', status: 'Active' }];
    const mockResponse = { data: mockNamespaces };
    vi.mocked(getNamespaces).mockResolvedValue(mockResponse);

    const result = await getNamespacesCached('kind-kind');

    expect(getNamespaces).toHaveBeenCalledWith({ context: 'kind-kind' });
    expect(result).toEqual(mockNamespaces);
  });

  it('should return cached data for same context on subsequent calls', async () => {
    const mockNamespaces = [{ name: 'default', status: 'Active' }];
    const mockResponse = { data: mockNamespaces };
    vi.mocked(getNamespaces).mockResolvedValue(mockResponse);

    // First call
    await getNamespacesCached('kind-kind');
    // Second call
    const result = await getNamespacesCached('kind-kind');

    expect(getNamespaces).toHaveBeenCalledTimes(1);
    expect(result).toEqual(mockNamespaces);
  });

  it('should handle concurrent calls for same context with same promise', async () => {
    const mockNamespaces = [{ name: 'default', status: 'Active' }];
    const mockResponse = { data: mockNamespaces };
    vi.mocked(getNamespaces).mockResolvedValue(mockResponse);

    const [result1, result2] = await Promise.all([
      getNamespacesCached('kind-kind'),
      getNamespacesCached('kind-kind')
    ]);

    expect(getNamespaces).toHaveBeenCalledTimes(1);
    expect(result1).toEqual(mockNamespaces);
    expect(result2).toEqual(mockNamespaces);
  });

  it('should fetch separately for different contexts', async () => {
    const mockNamespaces1 = [{ name: 'default', status: 'Active' }];
    const mockNamespaces2 = [{ name: 'kube-system', status: 'Active' }];
    vi.mocked(getNamespaces)
      .mockResolvedValueOnce({ data: mockNamespaces1 })
      .mockResolvedValueOnce({ data: mockNamespaces2 });

    await Promise.all([
      getNamespacesCached('kind-kind'),
      getNamespacesCached('kind-kind-2')
    ]);

    expect(getNamespaces).toHaveBeenCalledTimes(2);
    expect(getNamespaces).toHaveBeenCalledWith({ context: 'kind-kind' });
    expect(getNamespaces).toHaveBeenCalledWith({ context: 'kind-kind-2' });
  });

  it('should clear cache for specific context', async () => {
    const mockNamespaces = [{ name: 'default', status: 'Active' }];
    const mockResponse = { data: mockNamespaces };
    vi.mocked(getNamespaces).mockResolvedValue(mockResponse);

    // First call
    await getNamespacesCached('kind-kind');
    // Clear cache for specific context
    clearNamespacesCache('kind-kind');
    // Second call should fetch again
    await getNamespacesCached('kind-kind');

    expect(getNamespaces).toHaveBeenCalledTimes(2);
  });

  it('should clear all cache when no context provided', async () => {
    const mockNamespaces = [{ name: 'default', status: 'Active' }];
    const mockResponse = { data: mockNamespaces };
    vi.mocked(getNamespaces).mockResolvedValue(mockResponse);

    // First call
    await getNamespacesCached('kind-kind');
    // Clear all cache
    clearNamespacesCache();
    // Second call should fetch again
    await getNamespacesCached('kind-kind');

    expect(getNamespaces).toHaveBeenCalledTimes(2);
  });

  it('should handle fetch errors and clean up loading promise', async () => {
    const error = new Error('Network error');
    vi.mocked(getNamespaces).mockRejectedValue(error);

    await expect(getNamespacesCached('kind-kind')).rejects.toThrow('Network error');

    // Should allow retry after error
    const mockNamespaces = [{ name: 'default', status: 'Active' }];
    const mockResponse = { data: mockNamespaces };
    vi.mocked(getNamespaces).mockResolvedValue(mockResponse);

    const result = await getNamespacesCached('kind-kind');
    expect(result).toEqual(mockNamespaces);
  });
});
