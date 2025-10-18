import { describe, it, expect, vi } from 'vitest';
import http from '../../../helpers/http';
import { getPods, getPodLogs } from '../resources/podsResource';

vi.mock('../../../helpers/http');

describe('podsResource', () => {
  const filter = { context: 'kind-kind', namespace: 'default' } as any;

  it('should call correct endpoint and return pods', async () => {
    const mockResponse = { data: [{ name: 'pod-1', namespace: 'default' }] } as any;
    vi.mocked(http.get).mockResolvedValue(mockResponse as any);

    const resp = await getPods(filter);
    expect(http.get).toHaveBeenCalledWith(`/v1/k8s/clusters/${filter.context}/pods?namespace=${filter.namespace}`, undefined);
    expect(resp.data.length).toBe(1);
  });

  it('should fallback to empty array when data is falsy', async () => {
    const mockResponse = { data: null } as any;
    vi.mocked(http.get).mockResolvedValue(mockResponse as any);

    const resp = await getPods(filter);
    expect(resp.data).toEqual([]);
  });

  it('should call getPodLogs with correct URL', async () => {
    const mockResponse = { data: { logs: [] } } as any;
    vi.mocked(http.get).mockResolvedValue(mockResponse as any);

    const logFilter = { context: 'kind-kind', namespace: 'default', name: 'pod-1' } as any;
    const resp = await getPodLogs(logFilter);
    expect(http.get).toHaveBeenCalledWith(`/v1/k8s/clusters/${logFilter.context}/namespaces/${logFilter.namespace}/pods/${logFilter.name}/logs`, undefined);
    expect(resp.data.logs).toEqual([]);
  });
});
