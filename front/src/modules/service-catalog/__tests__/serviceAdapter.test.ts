import { describe, it, expect } from 'vitest';
import * as serviceAdapter from '../adapters/serviceAdapter';
import type { Service, ServiceFormData } from '../types';

describe('serviceAdapter', () => {
  const mockApiService = {
    apiVersion: 'v1',
    kind: 'Service',
    metadata: {
      name: 'test-service',
      tier: 'TIER-1',
      created_at: '2023-01-01T00:00:00Z',
    },
    spec: {
      description: 'Test service',
      team: {
        github_team: 'test-team',
      },
      business: {
        impact: 'high',
        sla_target: '99.9%',
      },
      kubernetes: {
        environments: [{
          name: 'prod',
          context: 'prod-cluster',
          namespace: 'default',
          resources: {
            deployments: [{
              name: 'test-deployment',
              replicas: 3,
              resources: {
                requests: { cpu: '100m', memory: '128Mi' },
                limits: { cpu: '500m', memory: '512Mi' },
              },
            }],
          },
        }],
      },
    },
  };

  it('should transform API service to domain model', () => {
    const result = serviceAdapter.transformServiceToDomain(mockApiService);
    
    expect(result.apiVersion).toBe('v1');
    expect(result.kind).toBe('Service');
    expect(result.metadata.name).toBe('test-service');
    expect(result.metadata.tier).toBe('TIER-1');
    expect(result.spec.description).toBe('Test service');
    expect(result.spec.team.github_team).toBe('test-team');
    expect(result.spec.business?.impact).toBe('high');
    expect(result.spec.kubernetes?.environments).toHaveLength(1);
  });

  it('should transform services list to domain models', () => {
    const result = serviceAdapter.transformServicesToDomain([mockApiService]);
    
    expect(result).toHaveLength(1);
    expect(result[0].metadata.name).toBe('test-service');
  });

  it('should transform service to API request format', () => {
    const service: Service = {
      apiVersion: 'v1',
      kind: 'Service',
      metadata: { name: 'test', tier: 'TIER-1' },
      spec: { description: 'test', team: { github_team: 'team' } },
    };
    
    const result = serviceAdapter.transformServiceToApiRequest(service);
    
    expect(result.apiVersion).toBe('v1');
    expect(result.kind).toBe('Service');
    expect(result.metadata.name).toBe('test');
  });

  it('should transform form data to service', () => {
    const formData: ServiceFormData = {
      name: 'test-service',
      description: 'Test service',
      tier: 'TIER-1',
      github_team: 'test-team',
      impact: 'high',
      sla_target: '99.9%',
      language: 'typescript',
      framework: 'react',
      env_name: 'prod',
      env_context: 'prod-cluster',
      env_namespace: 'default',
      deployment_name: 'test-deployment',
      deployment_replicas: 3,
      cpu_request: '100m',
      memory_request: '128Mi',
      cpu_limit: '500m',
      memory_limit: '512Mi',
      metrics_url: 'http://metrics',
      logs_url: 'http://logs',
    };
    
    const result = serviceAdapter.transformFormDataToService(formData);
    
    expect(result.metadata.name).toBe('test-service');
    expect(result.metadata.tier).toBe('TIER-1');
    expect(result.spec.description).toBe('Test service');
    expect(result.spec.team.github_team).toBe('test-team');
    expect(result.spec.kubernetes?.environments).toHaveLength(1);
  });

  it('should filter services by criteria', () => {
    const services: Service[] = [
      {
        apiVersion: 'v1',
        kind: 'Service',
        metadata: { name: 'service-1', tier: 'TIER-1' },
        spec: { description: 'High tier service', team: { github_team: 'team-a' } },
      },
      {
        apiVersion: 'v1',
        kind: 'Service',
        metadata: { name: 'service-2', tier: 'TIER-2' },
        spec: { description: 'Medium tier service', team: { github_team: 'team-b' } },
      },
    ];
    
    const filtered = serviceAdapter.filterServices(services, { tier: 'TIER-1' });
    
    expect(filtered).toHaveLength(1);
    expect(filtered[0].metadata.name).toBe('service-1');
  });

  it('should sort services by criteria', () => {
    const services: Service[] = [
      {
        apiVersion: 'v1',
        kind: 'Service',
        metadata: { name: 'service-b', tier: 'TIER-2' },
        spec: { description: 'B service', team: { github_team: 'team' } },
      },
      {
        apiVersion: 'v1',
        kind: 'Service',
        metadata: { name: 'service-a', tier: 'TIER-1' },
        spec: { description: 'A service', team: { github_team: 'team' } },
      },
    ];
    
    const sorted = serviceAdapter.sortServices(services, 'name');
    
    expect(sorted[0].metadata.name).toBe('service-a');
    expect(sorted[1].metadata.name).toBe('service-b');
  });
});
