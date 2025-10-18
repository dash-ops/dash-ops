import { describe, it, expect, vi } from 'vitest';
import http from '../../../helpers/http';
import { getNamespaces } from '../resources/namespaceResource';

vi.mock('../../../helpers/http');

describe('namespaceResource', () => {
  const filter = { context: 'kind-kind' } as any;

  it('should call correct endpoint and return data', async () => {
    const mockResponse = { data: [{ name: 'default', status: 'Active' }] } as any;
    vi.mocked(http.get).mockResolvedValue(mockResponse as any);

    const resp = await getNamespaces(filter);
    expect(http.get).toHaveBeenCalledWith(`/v1/k8s/clusters/${filter.context}/namespaces`, undefined);
    expect(resp.data.length).toBe(1);
  });

  it('should fallback to empty array when data is falsy', async () => {
    const mockResponse = { data: null } as any;
    vi.mocked(http.get).mockResolvedValue(mockResponse as any);

    const resp = await getNamespaces(filter);
    expect(resp.data).toEqual([]);
  });
});
