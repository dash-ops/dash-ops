/**
 * Tests for instance adapter functions
 */

import { describe, it, expect } from 'vitest';
import {
  transformTags,
  transformInstanceToDomain,
  transformInstancesToDomain,
  filterInstancesByState,
  sortInstancesByLaunchTime,
} from '../adapters/instanceAdapter';
import { AWSTypes } from '@/types';

describe('instanceAdapter', () => {
  describe('transformTags', () => {
    it('should transform tags array to key-value object', () => {
      const tags = [
        { key: 'Name', value: 'test-instance' },
        { key: 'Environment', value: 'production' },
      ];

      const result = transformTags(tags);

      expect(result).toEqual({
        Name: 'test-instance',
        Environment: 'production',
      });
    });

    it('should return empty object for empty tags array', () => {
      const result = transformTags([]);
      expect(result).toEqual({});
    });

    it('should return empty object for undefined tags', () => {
      const result = transformTags(undefined);
      expect(result).toEqual({});
    });

    it('should skip tags without key or value', () => {
      const tags = [
        { key: 'Name', value: 'test' },
        { key: '', value: 'empty-key' },
        { key: 'EmptyValue', value: '' },
        { key: 'Valid', value: 'value' },
      ];

      const result = transformTags(tags);

      expect(result).toEqual({
        Name: 'test',
        Valid: 'value',
      });
    });
  });

  describe('transformInstanceToDomain', () => {
    it('should transform API instance to domain model', () => {
      const apiInstance = {
        instance_id: 'i-1234567890abcdef0',
        name: 'test-instance',
        state: { name: 'running', code: 16 },
        platform: 'linux',
        instance_type: 't3.micro',
        public_ip: '203.0.113.12',
        private_ip: '10.0.0.100',
        subnet_id: 'subnet-12345678',
        vpc_id: 'vpc-12345678',
        cpu: { vcpus: 2, utilization: 45.5 },
        memory: { size_gb: 1, utilization: 60.2 },
        tags: [{ key: 'Name', value: 'test-instance' }],
        launch_time: '2023-01-01T00:00:00Z',
        account: 'prod',
        region: 'us-east-1',
        security_groups: [{ group_id: 'sg-12345678', group_name: 'default' }],
        cost_estimate: 0.0116,
      };

      const result = transformInstanceToDomain(apiInstance);

      expect(result).toEqual({
        id: 'i-1234567890abcdef0',
        name: 'test-instance',
        instance_id: 'i-1234567890abcdef0',
        state: { name: 'running', code: 16 },
        platform: 'linux',
        instance_type: 't3.micro',
        public_ip: '203.0.113.12',
        private_ip: '10.0.0.100',
        subnet_id: 'subnet-12345678',
        vpc_id: 'vpc-12345678',
        cpu: { vcpus: 2, utilization: 45.5 },
        memory: { size_gb: 1, utilization: 60.2 },
        tags: [{ key: 'Name', value: 'test-instance' }],
        launch_time: '2023-01-01T00:00:00Z',
        account: 'prod',
        region: 'us-east-1',
        security_groups: [{ group_id: 'sg-12345678', group_name: 'default' }],
        cost_estimate: 0.0116,
      });
    });

    it('should handle missing optional fields', () => {
      const apiInstance = {
        instance_id: 'i-1234567890abcdef0',
        name: 'test-instance',
        state: { name: 'stopped', code: 80 },
        platform: 'linux',
        instance_type: 't3.micro',
        public_ip: '',
        private_ip: '10.0.0.100',
        cpu: { vcpus: 2 },
        memory: { size_gb: 1 },
        tags: [],
        launch_time: '2023-01-01T00:00:00Z',
        account: 'prod',
        region: 'us-east-1',
        cost_estimate: 0,
      };

      const result = transformInstanceToDomain(apiInstance);

      expect(result.subnet_id).toBeUndefined();
      expect(result.vpc_id).toBeUndefined();
      expect(result.security_groups).toBeUndefined();
    });
  });

  describe('transformInstancesToDomain', () => {
    it('should transform multiple instances', () => {
      const apiInstances = [
        {
          instance_id: 'i-1',
          name: 'instance-1',
          state: { name: 'running', code: 16 },
          platform: 'linux',
          instance_type: 't3.micro',
          public_ip: '1.1.1.1',
          private_ip: '10.0.0.1',
          cpu: { vcpus: 2 },
          memory: { size_gb: 1 },
          tags: [],
          launch_time: '2023-01-01T00:00:00Z',
          account: 'prod',
          region: 'us-east-1',
          cost_estimate: 0,
        },
        {
          instance_id: 'i-2',
          name: 'instance-2',
          state: { name: 'stopped', code: 80 },
          platform: 'windows',
          instance_type: 't3.small',
          public_ip: '2.2.2.2',
          private_ip: '10.0.0.2',
          cpu: { vcpus: 2 },
          memory: { size_gb: 2 },
          tags: [],
          launch_time: '2023-01-02T00:00:00Z',
          account: 'prod',
          region: 'us-east-1',
          cost_estimate: 0,
        },
      ];

      const result = transformInstancesToDomain(apiInstances);

      expect(result).toHaveLength(2);
      expect(result[0].instance_id).toBe('i-1');
      expect(result[1].instance_id).toBe('i-2');
    });
  });

  describe('filterInstancesByState', () => {
    it('should filter instances by state name', () => {
      const instances: AWSTypes.Instance[] = [
        {
          id: 'i-1',
          name: 'instance-1',
          instance_id: 'i-1',
          state: { name: 'running', code: 16 },
        } as AWSTypes.Instance,
        {
          id: 'i-2',
          name: 'instance-2',
          instance_id: 'i-2',
          state: { name: 'stopped', code: 80 },
        } as AWSTypes.Instance,
        {
          id: 'i-3',
          name: 'instance-3',
          instance_id: 'i-3',
          state: { name: 'running', code: 16 },
        } as AWSTypes.Instance,
      ];

      const result = filterInstancesByState(instances, 'running');

      expect(result).toHaveLength(2);
      expect(result[0].instance_id).toBe('i-1');
      expect(result[1].instance_id).toBe('i-3');
    });

    it('should return empty array when no instances match state', () => {
      const instances: AWSTypes.Instance[] = [
        {
          id: 'i-1',
          name: 'instance-1',
          instance_id: 'i-1',
          state: { name: 'running', code: 16 },
        } as AWSTypes.Instance,
      ];

      const result = filterInstancesByState(instances, 'stopped');

      expect(result).toHaveLength(0);
    });
  });

  describe('sortInstancesByLaunchTime', () => {
    it('should sort instances by launch time in descending order by default', () => {
      const instances: AWSTypes.Instance[] = [
        {
          id: 'i-1',
          name: 'instance-1',
          instance_id: 'i-1',
          launch_time: '2023-01-01T00:00:00Z',
        } as AWSTypes.Instance,
        {
          id: 'i-2',
          name: 'instance-2',
          instance_id: 'i-2',
          launch_time: '2023-01-03T00:00:00Z',
        } as AWSTypes.Instance,
        {
          id: 'i-3',
          name: 'instance-3',
          instance_id: 'i-3',
          launch_time: '2023-01-02T00:00:00Z',
        } as AWSTypes.Instance,
      ];

      const result = sortInstancesByLaunchTime(instances);

      expect(result[0].instance_id).toBe('i-2'); // Most recent
      expect(result[1].instance_id).toBe('i-3');
      expect(result[2].instance_id).toBe('i-1'); // Oldest
    });

    it('should sort instances by launch time in ascending order when specified', () => {
      const instances: AWSTypes.Instance[] = [
        {
          id: 'i-1',
          name: 'instance-1',
          instance_id: 'i-1',
          launch_time: '2023-01-01T00:00:00Z',
        } as AWSTypes.Instance,
        {
          id: 'i-2',
          name: 'instance-2',
          instance_id: 'i-2',
          launch_time: '2023-01-03T00:00:00Z',
        } as AWSTypes.Instance,
      ];

      const result = sortInstancesByLaunchTime(instances, 'asc');

      expect(result[0].instance_id).toBe('i-1'); // Oldest first
      expect(result[1].instance_id).toBe('i-2');
    });

    it('should not mutate original array', () => {
      const instances: AWSTypes.Instance[] = [
        {
          id: 'i-1',
          name: 'instance-1',
          instance_id: 'i-1',
          launch_time: '2023-01-01T00:00:00Z',
        } as AWSTypes.Instance,
        {
          id: 'i-2',
          name: 'instance-2',
          instance_id: 'i-2',
          launch_time: '2023-01-03T00:00:00Z',
        } as AWSTypes.Instance,
      ];

      const originalFirst = instances[0];
      const result = sortInstancesByLaunchTime(instances);

      expect(result).not.toBe(instances); // Different array reference
      expect(instances[0]).toBe(originalFirst); // Original unchanged
    });
  });
});
