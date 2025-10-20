import { describe, it, expect } from 'vitest';
import * as configUtils from '../utils/configUtils';
import type { Plugin } from '../types';

describe('configUtils', () => {
  const mockPlugin: Plugin = {
    id: 'aws',
    enabled: true,
  };

  it('should get plugin display name', () => {
    expect(configUtils.getPluginDisplayName(mockPlugin)).toBe('Aws');
    
    const pluginWithHyphens: Plugin = { id: 'service-catalog', enabled: true };
    expect(configUtils.getPluginDisplayName(pluginWithHyphens)).toBe('Service catalog');
  });

  it('should get plugin icon', () => {
    expect(configUtils.getPluginIcon('aws')).toBe('aws');
    expect(configUtils.getPluginIcon('kubernetes')).toBe('kubernetes');
    expect(configUtils.getPluginIcon('service-catalog')).toBe('layers-3');
    expect(configUtils.getPluginIcon('oauth2')).toBe('shield');
    expect(configUtils.getPluginIcon('unknown')).toBe('package');
  });

  it('should get plugin description', () => {
    expect(configUtils.getPluginDescription('aws')).toContain('Amazon Web Services');
    expect(configUtils.getPluginDescription('kubernetes')).toContain('Kubernetes cluster');
    expect(configUtils.getPluginDescription('unknown')).toBe('Plugin for managing system resources');
  });

  it('should get plugin category', () => {
    expect(configUtils.getPluginCategory('aws')).toBe('Cloud Provider');
    expect(configUtils.getPluginCategory('kubernetes')).toBe('Container Platform');
    expect(configUtils.getPluginCategory('oauth2')).toBe('Authentication');
    expect(configUtils.getPluginCategory('unknown')).toBe('Other');
  });

  it('should get plugin status color', () => {
    expect(configUtils.getPluginStatusColor(true)).toContain('green');
    expect(configUtils.getPluginStatusColor(false)).toContain('gray');
  });

  it('should get plugin status label', () => {
    expect(configUtils.getPluginStatusLabel(true)).toBe('Enabled');
    expect(configUtils.getPluginStatusLabel(false)).toBe('Disabled');
  });

  it('should get plugin status icon', () => {
    expect(configUtils.getPluginStatusIcon(true)).toBe('check-circle');
    expect(configUtils.getPluginStatusIcon(false)).toBe('x-circle');
  });

  it('should sort plugins by name', () => {
    const plugins: Plugin[] = [
      { id: 'kubernetes', enabled: true },
      { id: 'aws', enabled: true },
      { id: 'service-catalog', enabled: true },
    ];
    
    const result = configUtils.sortPluginsByName(plugins);
    expect(result[0].id).toBe('aws');
    expect(result[1].id).toBe('kubernetes');
    expect(result[2].id).toBe('service-catalog');
  });

  it('should sort plugins by category', () => {
    const plugins: Plugin[] = [
      { id: 'kubernetes', enabled: true },
      { id: 'aws', enabled: true },
      { id: 'oauth2', enabled: true },
    ];
    
    const result = configUtils.sortPluginsByCategory(plugins);
    expect(result[0].id).toBe('oauth2'); // Authentication (A comes first alphabetically)
    expect(result[1].id).toBe('aws'); // Cloud Provider
    expect(result[2].id).toBe('kubernetes'); // Container Platform
  });

  it('should filter plugins by category', () => {
    const plugins: Plugin[] = [
      { id: 'aws', enabled: true },
      { id: 'kubernetes', enabled: true },
      { id: 'oauth2', enabled: true },
    ];
    
    const result = configUtils.filterPluginsByCategory(plugins, 'Cloud Provider');
    expect(result).toHaveLength(1);
    expect(result[0].id).toBe('aws');
  });

  it('should filter enabled plugins', () => {
    const plugins: Plugin[] = [
      { id: 'aws', enabled: true },
      { id: 'kubernetes', enabled: false },
      { id: 'service-catalog', enabled: true },
    ];
    
    const result = configUtils.filterEnabledPlugins(plugins);
    expect(result).toHaveLength(2);
    expect(result[0].id).toBe('aws');
    expect(result[1].id).toBe('service-catalog');
  });

  it('should filter disabled plugins', () => {
    const plugins: Plugin[] = [
      { id: 'aws', enabled: true },
      { id: 'kubernetes', enabled: false },
      { id: 'service-catalog', enabled: true },
    ];
    
    const result = configUtils.filterDisabledPlugins(plugins);
    expect(result).toHaveLength(1);
    expect(result[0].id).toBe('kubernetes');
  });

  it('should search plugins', () => {
    const plugins: Plugin[] = [
      { id: 'aws', enabled: true },
      { id: 'kubernetes', enabled: true },
      { id: 'service-catalog', enabled: true },
    ];
    
    const result = configUtils.searchPlugins(plugins, 'aws');
    expect(result).toHaveLength(1);
    expect(result[0].id).toBe('aws');
    
    const result2 = configUtils.searchPlugins(plugins, 'cloud');
    expect(result2).toHaveLength(1);
    expect(result2[0].id).toBe('aws');
  });

  it('should get plugin counts', () => {
    const plugins: Plugin[] = [
      { id: 'aws', enabled: true },
      { id: 'kubernetes', enabled: false },
      { id: 'service-catalog', enabled: true },
    ];
    
    expect(configUtils.getPluginCount(plugins)).toBe(3);
    expect(configUtils.getEnabledPluginCount(plugins)).toBe(2);
    expect(configUtils.getDisabledPluginCount(plugins)).toBe(1);
  });

  it('should get plugin categories', () => {
    const plugins: Plugin[] = [
      { id: 'aws', enabled: true },
      { id: 'kubernetes', enabled: true },
      { id: 'oauth2', enabled: true },
    ];
    
    const result = configUtils.getPluginCategories(plugins);
    expect(result).toContain('Cloud Provider');
    expect(result).toContain('Container Platform');
    expect(result).toContain('Authentication');
  });

  it('should get plugins by category', () => {
    const plugins: Plugin[] = [
      { id: 'aws', enabled: true },
      { id: 'kubernetes', enabled: true },
      { id: 'oauth2', enabled: true },
    ];
    
    const result = configUtils.getPluginsByCategory(plugins);
    expect(result['Cloud Provider']).toHaveLength(1);
    expect(result['Cloud Provider'][0].id).toBe('aws');
    expect(result['Container Platform']).toHaveLength(1);
    expect(result['Container Platform'][0].id).toBe('kubernetes');
  });

  it('should format plugin list', () => {
    const plugins: Plugin[] = [
      { id: 'aws', enabled: true },
      { id: 'kubernetes', enabled: true },
    ];
    
    expect(configUtils.formatPluginList([])).toBe('No plugins available');
    expect(configUtils.formatPluginList([plugins[0]])).toBe('aws');
    expect(configUtils.formatPluginList(plugins)).toBe('aws and kubernetes');
  });

  it('should validate plugin name', () => {
    expect(configUtils.validatePluginName('aws')).toBe(true);
    expect(configUtils.validatePluginName('service-catalog')).toBe(true);
    expect(configUtils.validatePluginName('AWS')).toBe(false);
    expect(configUtils.validatePluginName('')).toBe(false);
    expect(configUtils.validatePluginName('a')).toBe(false);
  });

  it('should check if plugin is core', () => {
    expect(configUtils.isPluginCore('config')).toBe(true);
    expect(configUtils.isPluginCore('oauth2')).toBe(true);
    expect(configUtils.isPluginCore('aws')).toBe(false);
  });

  it('should check if plugin is optional', () => {
    expect(configUtils.isPluginOptional('observability')).toBe(true);
    expect(configUtils.isPluginOptional('service-catalog')).toBe(true);
    expect(configUtils.isPluginOptional('aws')).toBe(false);
  });

  it('should get plugin priority', () => {
    expect(configUtils.getPluginPriority('config')).toBe(1);
    expect(configUtils.getPluginPriority('oauth2')).toBe(2);
    expect(configUtils.getPluginPriority('aws')).toBe(3);
    expect(configUtils.getPluginPriority('unknown')).toBe(999);
  });

  it('should sort plugins by priority', () => {
    const plugins: Plugin[] = [
      { id: 'aws', enabled: true },
      { id: 'config', enabled: true },
      { id: 'oauth2', enabled: true },
    ];
    
    const result = configUtils.sortPluginsByPriority(plugins);
    expect(result[0].id).toBe('config');
    expect(result[1].id).toBe('oauth2');
    expect(result[2].id).toBe('aws');
  });
});
