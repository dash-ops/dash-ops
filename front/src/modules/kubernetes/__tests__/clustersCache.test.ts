import { describe, it, expect, vi, beforeEach } from 'vitest';
import { getClustersCached, clearClustersCache } from '../utils/clustersCache';
import { getClusters } from '../resources/clusterResource';

vi.mock('../resources/clusterResource');

describe('clustersCache', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    clearClustersCache();
  });

  it('should fetch and cache clusters on first call', async () => {
    const mockClusters = [{ name: 'kind-kind' }];
    const mockResponse = { data: { clusters: mockClusters } } as any;
    vi.mocked(getClusters).mockResolvedValue(mockResponse as any);

    const result = await getClustersCached();

    expect(getClusters).toHaveBeenCalledTimes(1);
    expect(result).toEqual(mockClusters);
  });

  it('should return cached data on subsequent calls', async () => {
    const mockClusters = [{ name: 'kind-kind' }];
    const mockResponse = { data: { clusters: mockClusters } } as any;
    vi.mocked(getClusters).mockResolvedValue(mockResponse as any);

    // First call
    await getClustersCached();
    // Second call
    const result = await getClustersCached();

    expect(getClusters).toHaveBeenCalledTimes(1);
    expect(result).toEqual(mockClusters);
  });

  it('should handle concurrent calls with same promise', async () => {
    const mockClusters = [{ name: 'kind-kind' }];
    const mockResponse = { data: { clusters: mockClusters } } as any;
    vi.mocked(getClusters).mockResolvedValue(mockResponse as any);

    const [result1, result2] = await Promise.all([
      getClustersCached(),
      getClustersCached()
    ]);

    expect(getClusters).toHaveBeenCalledTimes(1);
    expect(result1).toEqual(mockClusters);
    expect(result2).toEqual(mockClusters);
  });

  it('should clear cache when clearClustersCache is called', async () => {
    const mockClusters = [{ name: 'kind-kind' }];
    const mockResponse = { data: { clusters: mockClusters } } as any;
    vi.mocked(getClusters).mockResolvedValue(mockResponse as any);

    // First call
    await getClustersCached();
    // Clear cache
    clearClustersCache();
    // Second call should fetch again
    await getClustersCached();

    expect(getClusters).toHaveBeenCalledTimes(2);
  });
});
