import * as pluginsResource from '../resources/pluginsResource';
import * as pluginsUtils from './plugins';
import type { Plugin } from '../types';

let pluginsCache: Plugin[] | null = null;
let cacheTimestamp = 0;
const CACHE_TTL = 10 * 60 * 1000; // 10 minutes

export const getPluginsCached = async (): Promise<Plugin[]> => {
  const now = Date.now();

  if (pluginsCache && now - cacheTimestamp < CACHE_TTL) {
    return pluginsCache;
  }

  try {
    const response = await pluginsResource.getPlugins();
    pluginsCache = pluginsUtils.transformPluginsToDomain(response.data);
    cacheTimestamp = now;
    return pluginsCache;
  } catch (error) {
    console.error('Failed to fetch plugins:', error);
    throw error;
  }
};

export const getPluginCached = async (name: string): Promise<Plugin | null> => {
  try {
    if (pluginsCache) {
      const cached = pluginsCache.find(p => p.id === name);
      if (cached) return cached;
    }

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
  return !!(pluginsCache && now - cacheTimestamp < CACHE_TTL);
};

export const getCacheAge = (): number => {
  if (!pluginsCache) return Infinity;
  return Date.now() - cacheTimestamp;
};

export const clearPluginsCache = (): void => {
  pluginsCache = null;
  cacheTimestamp = 0;
};

