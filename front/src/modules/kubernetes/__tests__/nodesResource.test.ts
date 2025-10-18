import { describe, it, expect, vi } from 'vitest';
import http from '../../../helpers/http';
import { getNodes } from '../resources/nodesResource';

vi.mock('../../../helpers/http');

describe('nodesResource', () => {
  const filter = { context: 'kind-kind' } as any;

  it('should call correct endpoint and return data', async () => {
    const mockResponse = { data: [{ name: 'node-1', status: 'Ready' }] } as any;
    vi.mocked(http.get).mockResolvedValue(mockResponse as any);

    const resp = await getNodes(filter);
    expect(http.get).toHaveBeenCalledWith(`/v1/k8s/clusters/${filter.context}/nodes`, undefined);
    expect(resp.data.length).toBe(1);
  });

  it('should fallback to empty array when data is falsy', async () => {
    const mockResponse = { data: null } as any;
    vi.mocked(http.get).mockResolvedValue(mockResponse as any);

    const resp = await getNodes(filter);
    expect(resp.data).toEqual([]);
  });
});
