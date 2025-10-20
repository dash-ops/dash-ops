import { useState, useEffect, useCallback } from 'react';
import * as configResource from '../resources/configResource';
import * as configAdapter from '../adapters/configAdapter';
import type { Plugin } from '../types';

interface UsePluginsState {
  plugins: Plugin[];
  loading: boolean;
  error: string | null;
}

interface UsePluginsReturn extends UsePluginsState {
  fetchPlugins: () => Promise<void>;
  refresh: () => Promise<void>;
  getEnabledPlugins: () => Plugin[];
  getPluginByName: (name: string) => Plugin | undefined;
  isPluginEnabled: (name: string) => boolean;
}

export const usePlugins = (): UsePluginsReturn => {
  const [state, setState] = useState<UsePluginsState>({
    plugins: [],
    loading: false,
    error: null,
  });

  const fetchPlugins = useCallback(async () => {
    try {
      setState(prev => ({ ...prev, loading: true, error: null }));
      
      const response = await configResource.getPlugins();
      const plugins = configAdapter.transformPluginsToDomain(response.data);
      
      setState(prev => ({ 
        ...prev, 
        plugins, 
        loading: false 
      }));
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to fetch plugins';
      setState(prev => ({ 
        ...prev, 
        error: errorMessage, 
        loading: false 
      }));
    }
  }, []);

  const refresh = useCallback(async () => {
    await fetchPlugins();
  }, [fetchPlugins]);

  const getEnabledPlugins = useCallback((): Plugin[] => {
    return configAdapter.filterEnabledPlugins(state.plugins);
  }, [state.plugins]);

  const getPluginByName = useCallback((name: string): Plugin | undefined => {
    return state.plugins.find(plugin => plugin.id === name);
  }, [state.plugins]);

  const isPluginEnabled = useCallback((name: string): boolean => {
    const plugin = getPluginByName(name);
    return plugin ? configAdapter.isPluginEnabled(plugin) : false;
  }, [getPluginByName]);

  useEffect(() => {
    fetchPlugins();
  }, [fetchPlugins]);

  return {
    ...state,
    fetchPlugins,
    refresh,
    getEnabledPlugins,
    getPluginByName,
    isPluginEnabled,
  };
};
