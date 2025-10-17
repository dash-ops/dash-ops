import { describe, it, expect } from 'vitest';
import {
  transformPodToDomain,
  transformPodsToDomain,
  transformPodContainers,
  transformPodContainerState,
  transformPodContainerResources,
  transformPodConditions,
  transformPodLogsToDomain,
  transformPodLogEntries,
} from '../adapters/podAdapter';

describe('podAdapter', () => {
  describe('transformPodToDomain', () => {
    it('should transform API pod to domain model', () => {
      const apiPod = {
        id: 'pod-1',
        name: 'test-pod',
        namespace: 'default',
        status: 'Running',
        phase: 'Running',
        node: 'node-1',
        restarts: 2,
        ready: '2/2',
        ip: '10.0.0.1',
        age: '1d',
        created_at: '2023-01-01T00:00:00Z',
        containers: [
          {
            name: 'container-1',
            image: 'nginx:latest',
            ready: true,
            restart_count: 1,
            state: { running: { started_at: '2023-01-01T00:00:00Z' } },
            resources: {
              requests: { cpu: '100m', memory: '128Mi' },
              limits: { cpu: '200m', memory: '256Mi' },
            },
          },
        ],
        conditions: [
          {
            type: 'Ready',
            status: 'True',
            last_transition_time: '2023-01-01T00:00:00Z',
          },
        ],
        qos_class: 'Burstable',
      };

      const result = transformPodToDomain(apiPod);

      expect(result.id).toBe('pod-1');
      expect(result.name).toBe('test-pod');
      expect(result.namespace).toBe('default');
      expect(result.status).toBe('Running');
      expect(result.phase).toBe('Running');
      expect(result.node).toBe('node-1');
      expect(result.restarts).toBe(2);
      expect(result.ready).toBe('2/2');
      expect(result.ip).toBe('10.0.0.1');
      expect(result.age).toBe('1d');
      expect(result.created_at).toBe('2023-01-01T00:00:00Z');
      expect(result.containers).toHaveLength(1);
      expect(result.conditions).toHaveLength(1);
      expect(result.qos_class).toBe('Burstable');
    });

    it('should handle missing optional fields', () => {
      const apiPod = {
        name: 'test-pod',
        namespace: 'default',
        status: 'Running',
        phase: 'Running',
        node: 'node-1',
        ready: '1/1',
        ip: '10.0.0.1',
        age: '1d',
        created_at: '2023-01-01T00:00:00Z',
        containers: [],
        conditions: [],
      };

      const result = transformPodToDomain(apiPod);

      expect(result.id).toBe('test-pod');
      expect(result.restarts).toBe(0);
      expect(result.containers).toHaveLength(0);
      expect(result.conditions).toHaveLength(0);
      expect(result.qos_class).toBeUndefined();
    });
  });

  describe('transformPodsToDomain', () => {
    it('should transform array of API pods to domain models', () => {
      const apiPods = [
        { name: 'pod-1', namespace: 'default', status: 'Running' },
        { name: 'pod-2', namespace: 'default', status: 'Pending' },
      ];

      const result = transformPodsToDomain(apiPods);

      expect(result).toHaveLength(2);
      expect(result[0].name).toBe('pod-1');
      expect(result[1].name).toBe('pod-2');
    });
  });

  describe('transformPodContainerState', () => {
    it('should transform running state', () => {
      const apiState = {
        running: {
          started_at: '2023-01-01T00:00:00Z',
        },
      };

      const result = transformPodContainerState(apiState);

      expect(result.running).toEqual({
        started_at: '2023-01-01T00:00:00Z',
      });
      expect(result.waiting).toBeUndefined();
      expect(result.terminated).toBeUndefined();
    });

    it('should transform waiting state', () => {
      const apiState = {
        waiting: {
          reason: 'ImagePullBackOff',
          message: 'Failed to pull image',
        },
      };

      const result = transformPodContainerState(apiState);

      expect(result.waiting).toEqual({
        reason: 'ImagePullBackOff',
        message: 'Failed to pull image',
      });
      expect(result.running).toBeUndefined();
      expect(result.terminated).toBeUndefined();
    });

    it('should transform terminated state', () => {
      const apiState = {
        terminated: {
          exit_code: 1,
          reason: 'Error',
          started_at: '2023-01-01T00:00:00Z',
          finished_at: '2023-01-01T01:00:00Z',
        },
      };

      const result = transformPodContainerState(apiState);

      expect(result.terminated).toEqual({
        exit_code: 1,
        reason: 'Error',
        started_at: '2023-01-01T00:00:00Z',
        finished_at: '2023-01-01T01:00:00Z',
      });
      expect(result.running).toBeUndefined();
      expect(result.waiting).toBeUndefined();
    });
  });

  describe('transformPodContainerResources', () => {
    it('should transform container resources with defaults', () => {
      const apiResources = {
        requests: { cpu: '100m', memory: '128Mi' },
        limits: { cpu: '200m', memory: '256Mi' },
      };

      const result = transformPodContainerResources(apiResources);

      expect(result.requests).toEqual({ cpu: '100m', memory: '128Mi' });
      expect(result.limits).toEqual({ cpu: '200m', memory: '256Mi' });
    });

    it('should handle missing resources with defaults', () => {
      const apiResources = {};

      const result = transformPodContainerResources(apiResources);

      expect(result.requests).toEqual({ cpu: '0', memory: '0' });
      expect(result.limits).toEqual({ cpu: '0', memory: '0' });
    });
  });

  describe('transformPodLogsToDomain', () => {
    it('should transform API pod logs response to domain model', () => {
      const apiResponse = {
        pod_name: 'test-pod',
        namespace: 'default',
        container_name: 'container-1',
        logs: [
          { timestamp: '2023-01-01T00:00:00Z', message: 'Log message 1' },
          { timestamp: '2023-01-01T00:01:00Z', message: 'Log message 2' },
        ],
        total_lines: 2,
      };

      const result = transformPodLogsToDomain(apiResponse);

      expect(result.pod_name).toBe('test-pod');
      expect(result.namespace).toBe('default');
      expect(result.container_name).toBe('container-1');
      expect(result.logs).toHaveLength(2);
      expect(result.total_lines).toBe(2);
      expect(result.logs[0].message).toBe('Log message 1');
    });
  });
});
