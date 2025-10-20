import * as configResource from '../resources/configResource';
import * as configAdapter from '../adapters/configAdapter';
import type { Plugin } from '../types';

// Cache state
let pluginsCache: Plugin[] | null = null;
let cacheTimestamp: number = 0;
const CACHE_TTL = 10 * 60 * 1000; // 10 minutes

export const getPluginsCached = async (): Promise<Plugin[]> => {
  const now = Date.now();
  
  // Return cached data if still valid
  if (pluginsCache && (now - cacheTimestamp) < CACHE_TTL) {
    return pluginsCache;
  }
  
  try {
    const response = await configResource.getPlugins();
    pluginsCache = configAdapter.transformPluginsToDomain(response.data);
    cacheTimestamp = now;
    return pluginsCache;
  } catch (error) {
    console.error('Failed to fetch plugins:', error);
    throw error;
  }
};

export const getPluginCached = async (name: string): Promise<Plugin | null> => {
  try {
    // First try to get from cache
    if (pluginsCache) {
      const cached = pluginsCache.find(p => p.id === name);
      if (cached) return cached;
    }
    
    // If not in cache, fetch from API
    const plugins = await getPluginsCached();
    return plugins.find(p => p.id === name) || null;
  } catch (error) {
    console.error(`Failed to fetch plugin ${name}:`, error);
    return null;
  }
};

export const invalidatePluginsCache = (): void => {
  pluginsCache = null;
  cacheTimestamp = 0;
};

export const updatePluginInCache = (plugin: Plugin): void => {
  if (!pluginsCache) return;
  
  const index = pluginsCache.findIndex(p => p.id === plugin.id);
  if (index >= 0) {
    pluginsCache[index] = plugin;
  } else {
    pluginsCache.push(plugin);
  }
  cacheTimestamp = Date.now();
};

export const removePluginFromCache = (pluginName: string): void => {
  if (!pluginsCache) return;
  
  pluginsCache = pluginsCache.filter(p => p.id !== pluginName);
  cacheTimestamp = Date.now();
};

export const isCacheValid = (): boolean => {
  const now = Date.now();
  return !!(pluginsCache && (now - cacheTimestamp) < CACHE_TTL);
};

export const getCacheAge = (): number => {
  if (!pluginsCache) return Infinity;
  return Date.now() - cacheTimestamp;
};

export const clearPluginsCache = (): void => {
  pluginsCache = null;
  cacheTimestamp = 0;
};
