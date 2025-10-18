import { describe, it, expect, vi } from 'vitest';
import http from '../../../helpers/http';
import { getClusters } from '../resources/clusterResource';

vi.mock('../../../helpers/http');

describe('clusterResource', () => {
  it('should call correct endpoint and return data', async () => {
    const mockResponse = { data: { clusters: [{ name: 'kind-kind' }], total: 1 } } as any;
    vi.mocked(http.get).mockResolvedValue(mockResponse);

    const resp = await getClusters();
    expect(http.get).toHaveBeenCalledWith('/v1/k8s/clusters', undefined);
    expect(resp.data.total).toBe(1);
  });
});
