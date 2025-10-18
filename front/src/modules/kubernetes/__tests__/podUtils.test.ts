import { describe, it, expect } from 'vitest';
import {
  getPodDisplayName,
  getPodStatusColor,
  isPodHealthy,
  getPodReadinessPercentage,
  formatPodAge,
  getPodRestartCount,
  hasRecentRestarts,
  getPodIP,
  isPodInError,
  getPodContainerCount,
  getRunningContainerCount,
} from '../utils/podUtils';
import type { Pod } from '../types';

describe('podUtils', () => {
  const mockPod: Pod = {
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
      {
        name: 'container-2',
        image: 'redis:latest',
        ready: true,
        restart_count: 0,
        state: { running: { started_at: '2023-01-01T00:00:00Z' } },
        resources: {
          requests: { cpu: '50m', memory: '64Mi' },
          limits: { cpu: '100m', memory: '128Mi' },
        },
      },
    ],
    conditions: [],
  };

  describe('getPodDisplayName', () => {
    it('should return pod name when available', () => {
      expect(getPodDisplayName(mockPod)).toBe('test-pod');
    });

    it('should return pod id when name is not available', () => {
      const podWithoutName = { ...mockPod, name: '' };
      expect(getPodDisplayName(podWithoutName)).toBe('pod-1');
    });
  });

  describe('getPodStatusColor', () => {
    it('should return green color for running pods', () => {
      expect(getPodStatusColor('Running')).toBe('text-green-600 bg-green-100');
    });

    it('should return yellow color for pending pods', () => {
      expect(getPodStatusColor('Pending')).toBe('text-yellow-600 bg-yellow-100');
    });

    it('should return red color for failed pods', () => {
      expect(getPodStatusColor('Failed')).toBe('text-red-600 bg-red-100');
    });

    it('should return blue color for succeeded pods', () => {
      expect(getPodStatusColor('Succeeded')).toBe('text-blue-600 bg-blue-100');
    });

    it('should return gray color for unknown pods', () => {
      expect(getPodStatusColor('Unknown')).toBe('text-gray-600 bg-gray-100');
    });

    it('should return gray color for unknown status', () => {
      expect(getPodStatusColor('SomeUnknownStatus')).toBe('text-gray-600 bg-gray-100');
    });
  });

  describe('isPodHealthy', () => {
    it('should return true for healthy running pods', () => {
      expect(isPodHealthy(mockPod)).toBe(true);
    });

    it('should return false for non-running pods', () => {
      const nonRunningPod = { ...mockPod, phase: 'Pending' };
      expect(isPodHealthy(nonRunningPod)).toBe(false);
    });

    it('should return false for pods with no ready containers', () => {
      const unreadyPod = { ...mockPod, phase: 'Running', ready: '0/0' };
      expect(isPodHealthy(unreadyPod)).toBe(false);
    });
  });

  describe('getPodReadinessPercentage', () => {
    it('should calculate readiness percentage correctly', () => {
      expect(getPodReadinessPercentage('2/2')).toBe(100);
      expect(getPodReadinessPercentage('1/2')).toBe(50);
      expect(getPodReadinessPercentage('0/2')).toBe(0);
    });

    it('should handle zero total containers', () => {
      expect(getPodReadinessPercentage('0/0')).toBe(0);
    });
  });

  describe('formatPodAge', () => {
    it('should return age when available', () => {
      expect(formatPodAge('1d')).toBe('1d');
    });

    it('should return Unknown when age is not available', () => {
      expect(formatPodAge('')).toBe('Unknown');
    });
  });

  describe('getPodRestartCount', () => {
    it('should return restart count when available', () => {
      expect(getPodRestartCount(mockPod)).toBe(2);
    });

    it('should return 0 when restart count is not available', () => {
      const podWithoutRestarts = { ...mockPod, restarts: undefined };
      expect(getPodRestartCount(podWithoutRestarts)).toBe(0);
    });
  });

  describe('hasRecentRestarts', () => {
    it('should return true when pod has restarts', () => {
      expect(hasRecentRestarts(mockPod)).toBe(true);
    });

    it('should return false when pod has no restarts', () => {
      const podWithoutRestarts = { ...mockPod, restarts: 0 };
      expect(hasRecentRestarts(podWithoutRestarts)).toBe(false);
    });
  });

  describe('getPodIP', () => {
    it('should return IP when available', () => {
      expect(getPodIP(mockPod)).toBe('10.0.0.1');
    });

    it('should return No IP when IP is not available', () => {
      const podWithoutIP = { ...mockPod, ip: '' };
      expect(getPodIP(podWithoutIP)).toBe('No IP');
    });
  });

  describe('isPodInError', () => {
    it('should return true for failed pods', () => {
      const failedPod = { ...mockPod, phase: 'Failed' };
      expect(isPodInError(failedPod)).toBe(true);
    });

    it('should return true for unknown pods', () => {
      const unknownPod = { ...mockPod, phase: 'Unknown' };
      expect(isPodInError(unknownPod)).toBe(true);
    });

    it('should return false for running pods', () => {
      expect(isPodInError(mockPod)).toBe(false);
    });
  });

  describe('getPodContainerCount', () => {
    it('should return container count', () => {
      expect(getPodContainerCount(mockPod)).toBe(2);
    });

    it('should return 0 when containers are not available', () => {
      const podWithoutContainers = { ...mockPod, containers: [] };
      expect(getPodContainerCount(podWithoutContainers)).toBe(0);
    });
  });

  describe('getRunningContainerCount', () => {
    it('should return count of ready containers', () => {
      expect(getRunningContainerCount(mockPod)).toBe(2);
    });

    it('should return 0 when no containers are ready', () => {
      const podWithUnreadyContainers = {
        ...mockPod,
        containers: [
          { ...mockPod.containers[0], ready: false },
          { ...mockPod.containers[1], ready: false },
        ],
      };
      expect(getRunningContainerCount(podWithUnreadyContainers)).toBe(0);
    });
  });
});
