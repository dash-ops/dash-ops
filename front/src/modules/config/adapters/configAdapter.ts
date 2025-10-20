import type { Plugin, PluginsResponse } from '../types';

export const transformPluginsToDomain = (apiPlugins: string[]): Plugin[] => {
  return apiPlugins.map(pluginName => ({
    id: pluginName,
    enabled: true, // Assume enabled if returned by API
  }));
};

export const transformPluginsResponseToDomain = (response: PluginsResponse): Plugin[] => {
  return transformPluginsToDomain(response.data);
};

export const transformPluginToApiRequest = (plugin: Plugin): string => {
  return plugin.id;
};

export const isPluginEnabled = (plugin: Plugin): boolean => {
  return plugin.enabled;
};

export const getPluginDisplayName = (plugin: Plugin): string => {
  return plugin.id.charAt(0).toUpperCase() + plugin.id.slice(1).replace(/-/g, ' ');
};

export const getPluginIcon = (pluginName: string): string => {
  const iconMap: Record<string, string> = {
    'aws': 'aws',
    'kubernetes': 'kubernetes',
    'service-catalog': 'layers-3',
    'oauth2': 'shield',
    'observability': 'activity',
    'config': 'settings',
  };
  
  return iconMap[pluginName] || 'package';
};

export const getPluginDescription = (pluginName: string): string => {
  const descriptions: Record<string, string> = {
    'aws': 'Amazon Web Services integration for managing EC2 instances and other AWS resources',
    'kubernetes': 'Kubernetes cluster management and monitoring capabilities',
    'service-catalog': 'Service catalog for managing and discovering application services',
    'oauth2': 'OAuth2 authentication and authorization management',
    'observability': 'Observability tools including metrics, logs, and traces',
    'config': 'Configuration management and system settings',
  };
  
  return descriptions[pluginName] || 'Plugin for managing system resources';
};

export const sortPluginsByName = (plugins: Plugin[]): Plugin[] => {
  return [...plugins].sort((a, b) => a.id.localeCompare(b.id));
};

export const filterEnabledPlugins = (plugins: Plugin[]): Plugin[] => {
  return plugins.filter(isPluginEnabled);
};

export const getPluginCategory = (pluginName: string): string => {
  const categories: Record<string, string> = {
    'aws': 'Cloud Provider',
    'kubernetes': 'Container Platform',
    'service-catalog': 'Service Management',
    'oauth2': 'Authentication',
    'observability': 'Monitoring',
    'config': 'System',
  };
  
  return categories[pluginName] || 'Other';
};

export const validatePluginName = (name: string): boolean => {
  // Plugin names should be lowercase, alphanumeric with hyphens
  const regex = /^[a-z0-9]([a-z0-9-]*[a-z0-9])?$/;
  return regex.test(name) && name.length >= 2 && name.length <= 50;
};

export const getPluginStatusColor = (enabled: boolean): string => {
  return enabled 
    ? 'bg-green-100 text-green-800 border-green-200'
    : 'bg-gray-100 text-gray-800 border-gray-200';
};

export const getPluginStatusLabel = (enabled: boolean): string => {
  return enabled ? 'Enabled' : 'Disabled';
};

export const formatPluginList = (plugins: Plugin[]): string => {
  if (plugins.length === 0) return 'No plugins available';
  if (plugins.length === 1) return plugins[0].id;
  if (plugins.length === 2) return `${plugins[0].id} and ${plugins[1].id}`;
  return `${plugins.slice(0, -1).map(p => p.id).join(', ')} and ${plugins[plugins.length - 1].id}`;
};
