import { describe, it, expect } from 'vitest';
import * as configAdapter from '../adapters/configAdapter';
import type { Plugin } from '../types';

describe('configAdapter', () => {
  const mockApiPlugins = ['aws', 'kubernetes', 'service-catalog'];

  it('should transform API plugins to domain models', () => {
    const result = configAdapter.transformPluginsToDomain(mockApiPlugins);
    
    expect(result).toHaveLength(3);
    expect(result[0]).toEqual({
      id: 'aws',
      enabled: true,
    });
    expect(result[1]).toEqual({
      id: 'kubernetes',
      enabled: true,
    });
    expect(result[2]).toEqual({
      id: 'service-catalog',
      enabled: true,
    });
  });

  it('should transform plugins response to domain models', () => {
    const response = { data: mockApiPlugins };
    const result = configAdapter.transformPluginsResponseToDomain(response);
    
    expect(result).toHaveLength(3);
    expect(result[0].id).toBe('aws');
  });

  it('should transform plugin to API request format', () => {
    const plugin: Plugin = {
      id: 'aws',
      enabled: true,
    };
    
    const result = configAdapter.transformPluginToApiRequest(plugin);
    expect(result).toBe('aws');
  });

  it('should check if plugin is enabled', () => {
    const enabledPlugin: Plugin = { id: 'aws', enabled: true };
    const disabledPlugin: Plugin = { id: 'kubernetes', enabled: false };
    
    expect(configAdapter.isPluginEnabled(enabledPlugin)).toBe(true);
    expect(configAdapter.isPluginEnabled(disabledPlugin)).toBe(false);
  });

  it('should get plugin display name', () => {
    const plugin: Plugin = { id: 'service-catalog', enabled: true };
    const result = configAdapter.getPluginDisplayName(plugin);
    expect(result).toBe('Service catalog');
  });

  it('should get plugin icon', () => {
    expect(configAdapter.getPluginIcon('aws')).toBe('aws');
    expect(configAdapter.getPluginIcon('kubernetes')).toBe('kubernetes');
    expect(configAdapter.getPluginIcon('unknown')).toBe('package');
  });

  it('should get plugin description', () => {
    const result = configAdapter.getPluginDescription('aws');
    expect(result).toContain('Amazon Web Services');
  });

  it('should get plugin category', () => {
    expect(configAdapter.getPluginCategory('aws')).toBe('Cloud Provider');
    expect(configAdapter.getPluginCategory('oauth2')).toBe('Authentication');
    expect(configAdapter.getPluginCategory('unknown')).toBe('Other');
  });

  it('should sort plugins by name', () => {
    const plugins: Plugin[] = [
      { id: 'kubernetes', enabled: true },
      { id: 'aws', enabled: true },
      { id: 'service-catalog', enabled: true },
    ];
    
    const result = configAdapter.sortPluginsByName(plugins);
    expect(result[0].id).toBe('aws');
    expect(result[1].id).toBe('kubernetes');
    expect(result[2].id).toBe('service-catalog');
  });

  it('should filter enabled plugins', () => {
    const plugins: Plugin[] = [
      { id: 'aws', enabled: true },
      { id: 'kubernetes', enabled: false },
      { id: 'service-catalog', enabled: true },
    ];
    
    const result = configAdapter.filterEnabledPlugins(plugins);
    expect(result).toHaveLength(2);
    expect(result[0].id).toBe('aws');
    expect(result[1].id).toBe('service-catalog');
  });

  it('should validate plugin name', () => {
    expect(configAdapter.validatePluginName('aws')).toBe(true);
    expect(configAdapter.validatePluginName('service-catalog')).toBe(true);
    expect(configAdapter.validatePluginName('AWS')).toBe(false);
    expect(configAdapter.validatePluginName('')).toBe(false);
    expect(configAdapter.validatePluginName('a')).toBe(false);
  });

  it('should get plugin status color', () => {
    expect(configAdapter.getPluginStatusColor(true)).toContain('green');
    expect(configAdapter.getPluginStatusColor(false)).toContain('gray');
  });

  it('should get plugin status label', () => {
    expect(configAdapter.getPluginStatusLabel(true)).toBe('Enabled');
    expect(configAdapter.getPluginStatusLabel(false)).toBe('Disabled');
  });

  it('should format plugin list', () => {
    const plugins: Plugin[] = [
      { id: 'aws', enabled: true },
      { id: 'kubernetes', enabled: true },
    ];
    
    expect(configAdapter.formatPluginList([])).toBe('No plugins available');
    expect(configAdapter.formatPluginList([plugins[0]])).toBe('aws');
    expect(configAdapter.formatPluginList(plugins)).toBe('aws and kubernetes');
  });
});
