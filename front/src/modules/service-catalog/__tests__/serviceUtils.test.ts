import { describe, it, expect } from 'vitest';
import * as serviceUtils from '../utils/serviceUtils';
import type { Service, ServiceHealth } from '../types';

describe('serviceUtils', () => {
  const mockService: Service = {
    apiVersion: 'v1',
    kind: 'Service',
    metadata: { name: 'test-service', tier: 'TIER-1' },
    spec: { description: 'Test service', team: { github_team: 'test-team' } },
  };

  it('should get service display name', () => {
    const result = serviceUtils.getServiceDisplayName(mockService);
    expect(result).toBe('test-service');
  });

  it('should get service tier color', () => {
    expect(serviceUtils.getServiceTierColor('TIER-1')).toContain('red');
    expect(serviceUtils.getServiceTierColor('TIER-2')).toContain('yellow');
    expect(serviceUtils.getServiceTierColor('TIER-3')).toContain('green');
  });

  it('should get service tier label', () => {
    expect(serviceUtils.getServiceTierLabel('TIER-1')).toBe('Tier 1');
    expect(serviceUtils.getServiceTierLabel('TIER-2')).toBe('Tier 2');
    expect(serviceUtils.getServiceTierLabel('TIER-3')).toBe('Tier 3');
  });

  it('should get health status color', () => {
    expect(serviceUtils.getHealthStatusColor('healthy')).toContain('green');
    expect(serviceUtils.getHealthStatusColor('degraded')).toContain('yellow');
    expect(serviceUtils.getHealthStatusColor('critical')).toContain('red');
    expect(serviceUtils.getHealthStatusColor('unknown')).toContain('gray');
  });

  it('should get health status icon', () => {
    expect(serviceUtils.getHealthStatusIcon('healthy')).toBe('check-circle');
    expect(serviceUtils.getHealthStatusIcon('degraded')).toBe('alert-triangle');
    expect(serviceUtils.getHealthStatusIcon('critical')).toBe('x-circle');
    expect(serviceUtils.getHealthStatusIcon('unknown')).toBe('help-circle');
  });

  it('should format service created date', () => {
    const result = serviceUtils.formatServiceCreatedDate('2023-01-01T00:00:00Z');
    expect(result).not.toBe('Unknown');
    expect(typeof result).toBe('string');
  });

  it('should format service updated date', () => {
    const recentDate = new Date(Date.now() - 1000 * 60 * 60).toISOString(); // 1 hour ago
    const result = serviceUtils.formatServiceUpdatedDate(recentDate);
    expect(result).toMatch(/h ago/);
  });

  it('should get service description preview', () => {
    const longDescription = 'This is a very long description that should be truncated';
    const result = serviceUtils.getServiceDescriptionPreview(longDescription, 20);
    expect(result.length).toBeLessThanOrEqual(23); // 20 + "..."
    expect(result.endsWith('...')).toBe(true);
  });

  it('should check if service has Kubernetes config', () => {
    const serviceWithK8s = {
      ...mockService,
      spec: {
        ...mockService.spec,
        kubernetes: { environments: [{ name: 'prod', context: 'prod', namespace: 'default', resources: { deployments: [] } }] },
      },
    };
    
    expect(serviceUtils.hasKubernetesConfig(mockService)).toBe(false);
    expect(serviceUtils.hasKubernetesConfig(serviceWithK8s)).toBe(true);
  });

  it('should get total deployments count', () => {
    const serviceWithDeployments = {
      ...mockService,
      spec: {
        ...mockService.spec,
        kubernetes: {
          environments: [
            { name: 'env1', context: 'ctx1', namespace: 'ns1', resources: { deployments: [{ name: 'dep1', replicas: 1, resources: { requests: { cpu: '100m', memory: '128Mi' }, limits: { cpu: '500m', memory: '512Mi' } } }] } },
            { name: 'env2', context: 'ctx2', namespace: 'ns2', resources: { deployments: [{ name: 'dep2', replicas: 1, resources: { requests: { cpu: '100m', memory: '128Mi' }, limits: { cpu: '500m', memory: '512Mi' } } }] } },
          ],
        },
      },
    };
    
    expect(serviceUtils.getTotalDeployments(mockService)).toBe(0);
    expect(serviceUtils.getTotalDeployments(serviceWithDeployments)).toBe(2);
  });

  it('should check if service is owned by team', () => {
    expect(serviceUtils.isServiceOwnedByTeam(mockService, 'test-team')).toBe(true);
    expect(serviceUtils.isServiceOwnedByTeam(mockService, 'other-team')).toBe(false);
  });

  it('should validate service name', () => {
    expect(serviceUtils.isValidServiceName('valid-service')).toBe(true);
    expect(serviceUtils.isValidServiceName('ValidService')).toBe(false);
    expect(serviceUtils.isValidServiceName('ab')).toBe(false);
    expect(serviceUtils.isValidServiceName('')).toBe(false);
  });

  it('should validate GitHub team name', () => {
    expect(serviceUtils.isValidGitHubTeam('valid-team')).toBe(true);
    expect(serviceUtils.isValidGitHubTeam('valid_team')).toBe(true);
    expect(serviceUtils.isValidGitHubTeam('ValidTeam')).toBe(true);
    expect(serviceUtils.isValidGitHubTeam('')).toBe(false);
  });

  it('should generate service URL slug', () => {
    expect(serviceUtils.generateServiceSlug('Test Service!')).toBe('test-service');
    expect(serviceUtils.generateServiceSlug('service@#$%name')).toBe('service-name');
  });
});
