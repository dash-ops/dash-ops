/**
 * Tests for instance utility functions
 */

import { describe, it, expect } from 'vitest';
import {
  getInstanceDisplayName,
  getInstanceStateColor,
  canStartInstance,
  canStopInstance,
  formatLaunchTime,
  getInstanceTypeDisplay,
  hasPublicIp,
  getInstanceEnvironment,
} from '../utils/instanceUtils';
import { AWSTypes } from '@/types';

describe('instanceUtils', () => {
  const createMockInstance = (overrides: Partial<AWSTypes.Instance> = {}): AWSTypes.Instance => ({
    id: 'i-1234567890abcdef0',
    name: 'test-instance',
    instance_id: 'i-1234567890abcdef0',
    state: { name: 'running', code: 16 },
    platform: 'linux',
    instance_type: 't3.micro',
    public_ip: '203.0.113.12',
    private_ip: '10.0.0.100',
    cpu: { vcpus: 2 },
    memory: { size_gb: 1 },
    tags: [],
    launch_time: '2023-01-01T00:00:00Z',
    account: 'prod',
    region: 'us-east-1',
    cost_estimate: 0,
    ...overrides,
  });

  describe('getInstanceDisplayName', () => {
    it('should return name from Name tag if available', () => {
      const instance = createMockInstance({
        tags: [{ key: 'Name', value: 'tag-name' }],
        name: 'instance-name',
      });

      const result = getInstanceDisplayName(instance);

      expect(result).toBe('tag-name');
    });

    it('should return name from name tag (lowercase) if available', () => {
      const instance = createMockInstance({
        tags: [{ key: 'name', value: 'lowercase-name' }],
        name: 'instance-name',
      });

      const result = getInstanceDisplayName(instance);

      expect(result).toBe('lowercase-name');
    });

    it('should fallback to instance name if no name tag', () => {
      const instance = createMockInstance({
        tags: [{ key: 'Environment', value: 'prod' }],
        name: 'instance-name',
      });

      const result = getInstanceDisplayName(instance);

      expect(result).toBe('instance-name');
    });

    it('should fallback to instance_id if no name or tag', () => {
      const instance = createMockInstance({
        tags: [],
        name: '',
        instance_id: 'i-1234567890abcdef0',
      });

      const result = getInstanceDisplayName(instance);

      expect(result).toBe('i-1234567890abcdef0');
    });
  });

  describe('getInstanceStateColor', () => {
    it('should return green colors for running state', () => {
      const result = getInstanceStateColor('running');
      expect(result).toBe('text-green-600 bg-green-50 border-green-200');
    });

    it('should return red colors for stopped state', () => {
      const result = getInstanceStateColor('stopped');
      expect(result).toBe('text-red-600 bg-red-50 border-red-200');
    });

    it('should return yellow colors for stopping state', () => {
      const result = getInstanceStateColor('stopping');
      expect(result).toBe('text-yellow-600 bg-yellow-50 border-yellow-200');
    });

    it('should return yellow colors for pending state', () => {
      const result = getInstanceStateColor('pending');
      expect(result).toBe('text-yellow-600 bg-yellow-50 border-yellow-200');
    });

    it('should return gray colors for unknown state', () => {
      const result = getInstanceStateColor('unknown');
      expect(result).toBe('text-gray-600 bg-gray-50 border-gray-200');
    });

    it('should handle undefined state', () => {
      const result = getInstanceStateColor(undefined);
      expect(result).toBe('text-gray-600 bg-gray-50 border-gray-200');
    });

    it('should be case insensitive', () => {
      const result = getInstanceStateColor('RUNNING');
      expect(result).toBe('text-green-600 bg-green-50 border-green-200');
    });
  });

  describe('canStartInstance', () => {
    it('should return true for stopped instance', () => {
      const instance = createMockInstance({
        state: { name: 'stopped', code: 80 },
      });

      const result = canStartInstance(instance);

      expect(result).toBe(true);
    });

    it('should return false for running instance', () => {
      const instance = createMockInstance({
        state: { name: 'running', code: 16 },
      });

      const result = canStartInstance(instance);

      expect(result).toBe(false);
    });

    it('should be case insensitive', () => {
      const instance = createMockInstance({
        state: { name: 'STOPPED', code: 80 },
      });

      const result = canStartInstance(instance);

      expect(result).toBe(true);
    });

    it('should return false for undefined state name', () => {
      const instance = createMockInstance({
        state: { name: undefined as any, code: 80 },
      });

      const result = canStartInstance(instance);

      expect(result).toBe(false);
    });
  });

  describe('canStopInstance', () => {
    it('should return true for running instance', () => {
      const instance = createMockInstance({
        state: { name: 'running', code: 16 },
      });

      const result = canStopInstance(instance);

      expect(result).toBe(true);
    });

    it('should return false for stopped instance', () => {
      const instance = createMockInstance({
        state: { name: 'stopped', code: 80 },
      });

      const result = canStopInstance(instance);

      expect(result).toBe(false);
    });

    it('should be case insensitive', () => {
      const instance = createMockInstance({
        state: { name: 'RUNNING', code: 16 },
      });

      const result = canStopInstance(instance);

      expect(result).toBe(true);
    });
  });

  describe('formatLaunchTime', () => {
    it('should format valid ISO date string', () => {
      const result = formatLaunchTime('2023-01-15T14:30:00Z');

      expect(result).toMatch(/Jan 15, 2023/);
      // The exact time depends on timezone, so just check if it contains time info
      expect(result).toMatch(/\d{1,2}:\d{2}/);
    });

    it('should handle different timezone', () => {
      const result = formatLaunchTime('2023-12-25T09:15:30Z');

      expect(result).toMatch(/Dec 25, 2023/);
      // The exact time depends on timezone, so just check if it contains time info
      expect(result).toMatch(/\d{1,2}:\d{2}/);
    });

    it('should return "Invalid date" for invalid date string', () => {
      const result = formatLaunchTime('invalid-date');

      expect(result).toBe('Invalid Date');
    });

    it('should return "Invalid date" for empty string', () => {
      const result = formatLaunchTime('');

      expect(result).toBe('Invalid Date');
    });
  });

  describe('getInstanceTypeDisplay', () => {
    it('should return the instance type as-is', () => {
      const result = getInstanceTypeDisplay('t3.micro');

      expect(result).toBe('t3.micro');
    });

    it('should return "Unknown" for empty string', () => {
      const result = getInstanceTypeDisplay('');

      expect(result).toBe('Unknown');
    });

    it('should return "Unknown" for undefined', () => {
      const result = getInstanceTypeDisplay(undefined as any);

      expect(result).toBe('Unknown');
    });
  });

  describe('hasPublicIp', () => {
    it('should return true for instance with public IP', () => {
      const instance = createMockInstance({
        public_ip: '203.0.113.12',
      });

      const result = hasPublicIp(instance);

      expect(result).toBe(true);
    });

    it('should return false for instance with "None" public IP', () => {
      const instance = createMockInstance({
        public_ip: 'None',
      });

      const result = hasPublicIp(instance);

      expect(result).toBe(false);
    });

    it('should return false for instance with empty public IP', () => {
      const instance = createMockInstance({
        public_ip: '',
      });

      const result = hasPublicIp(instance);

      expect(result).toBe(false);
    });

    it('should return false for instance with undefined public IP', () => {
      const instance = createMockInstance({
        public_ip: undefined as any,
      });

      const result = hasPublicIp(instance);

      expect(result).toBe(false);
    });
  });

  describe('getInstanceEnvironment', () => {
    it('should return environment from Environment tag', () => {
      const instance = createMockInstance({
        tags: [
          { key: 'Environment', value: 'production' },
          { key: 'Name', value: 'test' },
        ],
      });

      const result = getInstanceEnvironment(instance);

      expect(result).toBe('production');
    });

    it('should return environment from environment tag (lowercase)', () => {
      const instance = createMockInstance({
        tags: [{ key: 'environment', value: 'staging' }],
      });

      const result = getInstanceEnvironment(instance);

      expect(result).toBe('staging');
    });

    it('should return environment from env tag', () => {
      const instance = createMockInstance({
        tags: [{ key: 'env', value: 'development' }],
      });

      const result = getInstanceEnvironment(instance);

      expect(result).toBe('development');
    });

    it('should prioritize Environment over other tags', () => {
      const instance = createMockInstance({
        tags: [
          { key: 'Environment', value: 'prod' },
          { key: 'env', value: 'dev' },
          { key: 'environment', value: 'staging' },
        ],
      });

      const result = getInstanceEnvironment(instance);

      expect(result).toBe('prod');
    });

    it('should return "Unknown" when no environment tags found', () => {
      const instance = createMockInstance({
        tags: [{ key: 'Name', value: 'test' }],
      });

      const result = getInstanceEnvironment(instance);

      expect(result).toBe('Unknown');
    });

    it('should return "Unknown" when no tags', () => {
      const instance = createMockInstance({
        tags: [],
      });

      const result = getInstanceEnvironment(instance);

      expect(result).toBe('Unknown');
    });
  });
});
