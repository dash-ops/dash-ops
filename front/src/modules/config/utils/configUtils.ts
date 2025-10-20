import type { Plugin } from '../types';

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

export const getPluginStatusColor = (enabled: boolean): string => {
  return enabled 
    ? 'bg-green-100 text-green-800 border-green-200'
    : 'bg-gray-100 text-gray-800 border-gray-200';
};

export const getPluginStatusLabel = (enabled: boolean): string => {
  return enabled ? 'Enabled' : 'Disabled';
};

export const getPluginStatusIcon = (enabled: boolean): string => {
  return enabled ? 'check-circle' : 'x-circle';
};

export const sortPluginsByName = (plugins: Plugin[]): Plugin[] => {
  return [...plugins].sort((a, b) => a.id.localeCompare(b.id));
};

export const sortPluginsByCategory = (plugins: Plugin[]): Plugin[] => {
  return [...plugins].sort((a, b) => {
    const categoryA = getPluginCategory(a.id);
    const categoryB = getPluginCategory(b.id);
    
    if (categoryA === categoryB) {
      return a.id.localeCompare(b.id);
    }
    
    return categoryA.localeCompare(categoryB);
  });
};

export const filterPluginsByCategory = (plugins: Plugin[], category: string): Plugin[] => {
  return plugins.filter(plugin => getPluginCategory(plugin.id) === category);
};

export const filterEnabledPlugins = (plugins: Plugin[]): Plugin[] => {
  return plugins.filter(plugin => plugin.enabled);
};

export const filterDisabledPlugins = (plugins: Plugin[]): Plugin[] => {
  return plugins.filter(plugin => !plugin.enabled);
};

export const searchPlugins = (plugins: Plugin[], query: string): Plugin[] => {
  const lowercaseQuery = query.toLowerCase();
  return plugins.filter(plugin => 
    plugin.id.toLowerCase().includes(lowercaseQuery) ||
    getPluginDisplayName(plugin).toLowerCase().includes(lowercaseQuery) ||
    getPluginDescription(plugin.id).toLowerCase().includes(lowercaseQuery) ||
    getPluginCategory(plugin.id).toLowerCase().includes(lowercaseQuery)
  );
};

export const getPluginCount = (plugins: Plugin[]): number => {
  return plugins.length;
};

export const getEnabledPluginCount = (plugins: Plugin[]): number => {
  return filterEnabledPlugins(plugins).length;
};

export const getDisabledPluginCount = (plugins: Plugin[]): number => {
  return filterDisabledPlugins(plugins).length;
};

export const getPluginCategories = (plugins: Plugin[]): string[] => {
  const categories = new Set(plugins.map(plugin => getPluginCategory(plugin.id)));
  return Array.from(categories).sort();
};

export const getPluginsByCategory = (plugins: Plugin[]): Record<string, Plugin[]> => {
  const categorized: Record<string, Plugin[]> = {};
  
  plugins.forEach(plugin => {
    const category = getPluginCategory(plugin.id);
    if (!categorized[category]) {
      categorized[category] = [];
    }
    categorized[category].push(plugin);
  });
  
  // Sort plugins within each category
  Object.keys(categorized).forEach(category => {
    categorized[category] = sortPluginsByName(categorized[category]);
  });
  
  return categorized;
};

export const formatPluginList = (plugins: Plugin[]): string => {
  if (plugins.length === 0) return 'No plugins available';
  if (plugins.length === 1) return plugins[0].id;
  if (plugins.length === 2) return `${plugins[0].id} and ${plugins[1].id}`;
  return `${plugins.slice(0, -1).map(p => p.id).join(', ')} and ${plugins[plugins.length - 1].id}`;
};

export const validatePluginName = (name: string): boolean => {
  // Plugin names should be lowercase, alphanumeric with hyphens
  const regex = /^[a-z0-9]([a-z0-9-]*[a-z0-9])?$/;
  return regex.test(name) && name.length >= 2 && name.length <= 50;
};

export const isPluginCore = (pluginName: string): boolean => {
  const corePlugins = ['config', 'oauth2'];
  return corePlugins.includes(pluginName);
};

export const isPluginOptional = (pluginName: string): boolean => {
  const optionalPlugins = ['observability', 'service-catalog'];
  return optionalPlugins.includes(pluginName);
};

export const getPluginPriority = (pluginName: string): number => {
  const priorities: Record<string, number> = {
    'config': 1,
    'oauth2': 2,
    'aws': 3,
    'kubernetes': 4,
    'service-catalog': 5,
    'observability': 6,
  };
  
  return priorities[pluginName] || 999;
};

export const sortPluginsByPriority = (plugins: Plugin[]): Plugin[] => {
  return [...plugins].sort((a, b) => {
    const priorityA = getPluginPriority(a.id);
    const priorityB = getPluginPriority(b.id);
    
    if (priorityA === priorityB) {
      return a.id.localeCompare(b.id);
    }
    
    return priorityA - priorityB;
  });
};
