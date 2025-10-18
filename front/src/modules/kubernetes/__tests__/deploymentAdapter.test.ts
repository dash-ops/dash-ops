import { describe, it, expect } from 'vitest';
import {
  transformDeploymentToDomain,
  transformDeploymentsToDomain,
  transformDeploymentReplicas,
  transformDeploymentConditions,
  transformPodInfo,
  transformServiceContext,
} from '../adapters/deploymentAdapter';

describe('deploymentAdapter', () => {
  it('should transform API deployment to domain model', () => {
    const api = {
      id: 'dep-1',
      name: 'web',
      namespace: 'default',
      pod_count: 3,
      pod_info: { running: 2, pending: 1, failed: 0, total: 3 },
      replicas: { ready: 2, available: 2, current: 2, desired: 3 },
      age: '1d',
      created_at: '2023-01-01T00:00:00Z',
      conditions: [{ type: 'Available', status: 'True' }],
      service_context: { environment: 'prod', team: 'platform' },
    };

    const dep = transformDeploymentToDomain(api);
    expect(dep.id).toBe('dep-1');
    expect(dep.pod_info.total).toBe(3);
    expect(dep.replicas.desired).toBe(3);
    expect(dep.service_context?.environment).toBe('prod');
  });

  it('should transform array of deployments', () => {
    const arr = [{ name: 'a' }, { name: 'b' }];
    const res = transformDeploymentsToDomain(arr as any);
    expect(res.length).toBe(2);
    expect(res[1].name).toBe('b');
  });

  it('should transform replicas and podInfo with defaults', () => {
    const replicas = transformDeploymentReplicas({} as any);
    expect(replicas.ready).toBe(0);

    const info = transformPodInfo({} as any);
    expect(info.total).toBe(0);
  });

  it('should transform conditions and service context', () => {
    const cond = transformDeploymentConditions([{ type: 'Progressing', status: 'True' }] as any);
    expect(cond[0].type).toBe('Progressing');

    const ctx = transformServiceContext({ environment: 'staging' } as any);
    expect(ctx?.environment).toBe('staging');
    expect(transformServiceContext(undefined)).toBeUndefined();
  });
});
