import { describe, it, expect, vi } from 'vitest';
import http from '../../../helpers/http';
import { getPermissions } from '../resources/permissionResource';

vi.mock('../../../helpers/http');

describe('permissionResource', () => {
  const filter = { context: 'kind-kind' } as any;

  it('should call correct endpoint and return data', async () => {
    const mockResponse = { data: [{ id: 'perm-1', name: 'pods', resource: 'pods', actions: ['get'] }] } as any;
    vi.mocked(http.get).mockResolvedValue(mockResponse as any);

    const resp = await getPermissions(filter);
    expect(http.get).toHaveBeenCalledWith(`/v1/k8s/clusters/${filter.context}/permissions`, undefined);
    expect(resp.data.length).toBe(1);
  });

  it('should fallback to empty array when data is falsy', async () => {
    const mockResponse = { data: null } as any;
    vi.mocked(http.get).mockResolvedValue(mockResponse as any);

    const resp = await getPermissions(filter);
    expect(resp.data).toEqual([]);
  });
});
