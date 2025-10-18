import { describe, it, expect } from 'vitest';
import {
  transformNodeToDomain,
  transformNodesToDomain,
  transformNodeConditions,
  transformNodeResources,
  transformResourceSpec,
  transformAllocatedResources,
} from '../adapters/nodeAdapter';

describe('nodeAdapter', () => {
  it('should transform API node to domain model', () => {
    const apiNode = {
      id: 'node-1',
      name: 'ip-10-0-0-1',
      status: 'Ready',
      roles: ['worker'],
      age: '5d',
      version: 'v1.29.0',
      internal_ip: '10.0.0.1',
      conditions: [
        { type: 'Ready', status: 'True', reason: 'KubeletReady', message: 'ok', last_transition_time: '2023-01-01T00:00:00Z' },
      ],
      resources: {
        capacity: { cpu: '4', memory: '8Gi', pods: '110' },
        allocatable: { cpu: '4', memory: '7Gi', pods: '110' },
        used: { cpu: '2', memory: '3Gi', pods: '20' },
      },
      created_at: '2023-01-01T00:00:00Z',
    };

    const node = transformNodeToDomain(apiNode);
    expect(node.id).toBe('node-1');
    expect(node.name).toBe('ip-10-0-0-1');
    expect(node.status).toBe('Ready');
    expect(node.roles).toEqual(['worker']);
    expect(node.resources.capacity.cpu).toBe('4');
  });

  it('should transform array of API nodes', () => {
    const arr = [{ name: 'n1', status: 'Ready' }, { name: 'n2', status: 'NotReady' }];
    const result = transformNodesToDomain(arr as any);
    expect(result.length).toBe(2);
    expect(result[0].name).toBe('n1');
  });

  it('should transform node conditions', () => {
    const cond = transformNodeConditions([
      { type: 'Ready', status: 'True', reason: 'ok', message: 'fine', last_transition_time: 'now' },
    ]);
    expect(cond[0].type).toBe('Ready');
    expect(cond[0].status).toBe('True');
  });

  it('should transform resource specs and resources', () => {
    const spec = transformResourceSpec({ cpu: '1', memory: '2Gi', pods: '10' });
    expect(spec.memory).toBe('2Gi');

    const res = transformNodeResources({ capacity: spec, allocatable: spec, used: spec } as any);
    expect(res.allocatable.pods).toBe('10');
  });

  it('should transform allocated resources with defaults', () => {
    const alloc = transformAllocatedResources({ cpu_requests_fraction: 0.5, pod_capacity: 100 });
    expect(alloc.cpu_requests_fraction).toBe(0.5);
    expect(alloc.memory_limits_fraction).toBe(0);
  });
});
