import { describe, it, expect } from 'vitest';
import {
  transformClusterToDomain,
  transformClustersToDomain,
  transformClusterListResponseToDomain,
  transformNamespaceToDomain,
  transformNamespacesToDomain,
} from '../adapters/clusterAdapter';

describe('clusterAdapter', () => {
  it('should transform API cluster to domain model', () => {
    const api = { id: 'c1', name: 'kind-kind', context: 'kind-kind', version: '1.29', status: 'Ready' };
    const c = transformClusterToDomain(api);
    expect(c.id).toBe('c1');
    expect(c.name).toBe('kind-kind');
  });

  it('should transform array of clusters and list response', () => {
    const arr = transformClustersToDomain([{ name: 'a' }, { name: 'b' }] as any);
    expect(arr.length).toBe(2);

    const list = transformClusterListResponseToDomain({ clusters: [{ name: 'x' }], total: 1 } as any);
    expect(list.total).toBe(1);
    expect(list.clusters[0].name).toBe('x');
  });

  it('should transform namespaces', () => {
    const ns = transformNamespaceToDomain({ name: 'default', status: 'Active' } as any);
    expect(ns.name).toBe('default');

    const nss = transformNamespacesToDomain([{ name: 'kube-system' }] as any);
    expect(nss[0].name).toBe('kube-system');
  });
});
