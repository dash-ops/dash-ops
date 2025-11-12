import { describe, it, expect } from 'vitest';
import * as pluginsUtils from '../utils/plugins';
import type { Plugin } from '../types';

describe('pluginsUtils', () => {
  const mockPlugin: Plugin = {
    id: 'aws',
    enabled: true,
  };

  it('transforms API plugins to domain models', () => {
    const apiPlugins = ['aws', 'kubernetes'];
    expect(pluginsUtils.transformPluginsToDomain(apiPlugins)).toEqual([
      { id: 'aws', enabled: true },
      { id: 'kubernetes', enabled: true },
    ]);
  });

  it('transforms plugins response to domain', () => {
    const response = { data: ['aws', 'kubernetes'] };
    expect(pluginsUtils.transformPluginsResponseToDomain(response)).toEqual([
      { id: 'aws', enabled: true },
      { id: 'kubernetes', enabled: true },
    ]);
  });

  it('transforms plugin to API request payload', () => {
    expect(pluginsUtils.transformPluginToApiRequest(mockPlugin)).toBe('aws');
  });

  it('should get plugin display name', () => {
    expect(pluginsUtils.getPluginDisplayName(mockPlugin)).toBe('Aws');

    const pluginWithHyphens: Plugin = { id: 'service-catalog', enabled: true };
    expect(pluginsUtils.getPluginDisplayName(pluginWithHyphens)).toBe('Service catalog');
  });

  it('should get plugin icon', () => {
    expect(pluginsUtils.getPluginIcon('aws')).toBe('aws');
    expect(pluginsUtils.getPluginIcon('kubernetes')).toBe('kubernetes');
    expect(pluginsUtils.getPluginIcon('service-catalog')).toBe('layers-3');
    expect(pluginsUtils.getPluginIcon('oauth2')).toBe('shield');
    expect(pluginsUtils.getPluginIcon('unknown')).toBe('package');
  });

  it('should get plugin description', () => {
    expect(pluginsUtils.getPluginDescription('aws')).toContain('Amazon Web Services');
    expect(pluginsUtils.getPluginDescription('kubernetes')).toContain('Kubernetes cluster');
    expect(pluginsUtils.getPluginDescription('unknown')).toBe('Plugin for managing system resources');
  });

  it('should get plugin category', () => {
    expect(pluginsUtils.getPluginCategory('aws')).toBe('Cloud Provider');
    expect(pluginsUtils.getPluginCategory('kubernetes')).toBe('Container Platform');
    expect(pluginsUtils.getPluginCategory('oauth2')).toBe('Authentication');
    expect(pluginsUtils.getPluginCategory('unknown')).toBe('Other');
  });

  it('should get plugin status color', () => {
    expect(pluginsUtils.getPluginStatusColor(true)).toContain('green');
    expect(pluginsUtils.getPluginStatusColor(false)).toContain('gray');
  });

  it('should get plugin status label', () => {
    expect(pluginsUtils.getPluginStatusLabel(true)).toBe('Enabled');
    expect(pluginsUtils.getPluginStatusLabel(false)).toBe('Disabled');
  });

  it('should get plugin status icon', () => {
    expect(pluginsUtils.getPluginStatusIcon(true)).toBe('check-circle');
    expect(pluginsUtils.getPluginStatusIcon(false)).toBe('x-circle');
  });

  it('should sort plugins by name', () => {
    const plugins: Plugin[] = [
      { id: 'kubernetes', enabled: true },
      { id: 'aws', enabled: true },
      { id: 'service-catalog', enabled: true },
    ];

    const result = pluginsUtils.sortPluginsByName(plugins);
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

    const result = pluginsUtils.sortPluginsByCategory(plugins);
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

    const result = pluginsUtils.filterPluginsByCategory(plugins, 'Cloud Provider');
    expect(result).toHaveLength(1);
    expect(result[0].id).toBe('aws');
  });

  it('should filter enabled plugins', () => {
    const plugins: Plugin[] = [
      { id: 'aws', enabled: true },
      { id: 'kubernetes', enabled: false },
      { id: 'service-catalog', enabled: true },
    ];

    const result = pluginsUtils.filterEnabledPlugins(plugins);
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

    const result = pluginsUtils.filterDisabledPlugins(plugins);
    expect(result).toHaveLength(1);
    expect(result[0].id).toBe('kubernetes');
  });

  it('should search plugins', () => {
    const plugins: Plugin[] = [
      { id: 'aws', enabled: true },
      { id: 'kubernetes', enabled: true },
      { id: 'service-catalog', enabled: true },
    ];

    const result = pluginsUtils.searchPlugins(plugins, 'aws');
    expect(result).toHaveLength(1);
    expect(result[0].id).toBe('aws');

    const result2 = pluginsUtils.searchPlugins(plugins, 'cloud');
    expect(result2).toHaveLength(1);
    expect(result2[0].id).toBe('aws');
  });

  it('should get plugin counts', () => {
    const plugins: Plugin[] = [
      { id: 'aws', enabled: true },
      { id: 'kubernetes', enabled: false },
      { id: 'service-catalog', enabled: true },
    ];

    expect(pluginsUtils.getPluginCount(plugins)).toBe(3);
    expect(pluginsUtils.getEnabledPluginCount(plugins)).toBe(2);
    expect(pluginsUtils.getDisabledPluginCount(plugins)).toBe(1);
  });

  it('should get plugin categories', () => {
    const plugins: Plugin[] = [
      { id: 'aws', enabled: true },
      { id: 'kubernetes', enabled: true },
      { id: 'oauth2', enabled: true },
    ];

    const result = pluginsUtils.getPluginCategories(plugins);
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

    const result = pluginsUtils.getPluginsByCategory(plugins);
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

    expect(pluginsUtils.formatPluginList([])).toBe('No plugins available');
    expect(pluginsUtils.formatPluginList([plugins[0]])).toBe('aws');
    expect(pluginsUtils.formatPluginList(plugins)).toBe('aws and kubernetes');
  });

  it('should validate plugin name', () => {
    expect(pluginsUtils.validatePluginName('aws')).toBe(true);
    expect(pluginsUtils.validatePluginName('service-catalog')).toBe(true);
    expect(pluginsUtils.validatePluginName('AWS')).toBe(false);
    expect(pluginsUtils.validatePluginName('')).toBe(false);
    expect(pluginsUtils.validatePluginName('a')).toBe(false);
  });

  it('should check if plugin is core', () => {
    expect(pluginsUtils.isPluginCore('config')).toBe(true);
    expect(pluginsUtils.isPluginCore('oauth2')).toBe(true);
    expect(pluginsUtils.isPluginCore('aws')).toBe(false);
  });

  it('should check if plugin is optional', () => {
    expect(pluginsUtils.isPluginOptional('observability')).toBe(true);
    expect(pluginsUtils.isPluginOptional('service-catalog')).toBe(true);
    expect(pluginsUtils.isPluginOptional('aws')).toBe(false);
  });

  it('should get plugin priority', () => {
    expect(pluginsUtils.getPluginPriority('config')).toBe(1);
    expect(pluginsUtils.getPluginPriority('oauth2')).toBe(2);
    expect(pluginsUtils.getPluginPriority('aws')).toBe(3);
    expect(pluginsUtils.getPluginPriority('unknown')).toBe(999);
  });

  it('should sort plugins by priority', () => {
    const plugins: Plugin[] = [
      { id: 'aws', enabled: true },
      { id: 'config', enabled: true },
      { id: 'oauth2', enabled: true },
    ];

    const result = pluginsUtils.sortPluginsByPriority(plugins);
    expect(result[0].id).toBe('config');
    expect(result[1].id).toBe('oauth2');
    expect(result[2].id).toBe('aws');
  });
});

