import { useState, useEffect, useRef } from 'react';
import * as providersResource from '../resources/providersResource';
import type { ProvidersResponse } from '../resources/providersResource';

interface ProvidersData {
  logs: string[];
  traces: string[];
  metrics: string[];
}

// Singleton cache to avoid multiple requests
let cachedProviders: ProvidersData | null = null;
let loadingPromise: Promise<ProvidersResponse> | null = null;
const CACHE_TTL = 5 * 60 * 1000; // 5 minutes
let cacheTimestamp: number = 0;

const isCacheValid = (): boolean => {
  if (!cachedProviders || !cacheTimestamp) {
    return false;
  }
  return Date.now() - cacheTimestamp < CACHE_TTL;
};

const fetchProviders = async (): Promise<ProvidersData> => {
  // If cache is valid, return cached data
  if (isCacheValid() && cachedProviders) {
    return cachedProviders;
  }

  // If there's already a request in progress, wait for it
  if (loadingPromise) {
    const response = await loadingPromise;
    return {
      logs: response.data.logs_providers || [],
      traces: response.data.traces_providers || [],
      metrics: response.data.metrics_providers || [],
    };
  }

  // Start new request
  loadingPromise = providersResource.getProviders();
  
  try {
    const response = await loadingPromise;
    const data: ProvidersData = {
      logs: response.data.logs_providers || [],
      traces: response.data.traces_providers || [],
      metrics: response.data.metrics_providers || [],
    };
    
    // Update cache
    cachedProviders = data;
    cacheTimestamp = Date.now();
    
    return data;
  } finally {
    // Clear loading promise after request completes
    loadingPromise = null;
  }
};

export const useProviders = () => {
  const [providers, setProviders] = useState<ProvidersData>({
    logs: [],
    traces: [],
    metrics: [],
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  const mountedRef = useRef(true);

  useEffect(() => {
    mountedRef.current = true;
    
    const loadProviders = async () => {
      try {
        setLoading(true);
        setError(null);
        const data = await fetchProviders();
        
        // Only update state if component is still mounted
        if (mountedRef.current) {
          setProviders(data);
          setLoading(false);
        }
      } catch (err) {
        if (mountedRef.current) {
          setError(err instanceof Error ? err : new Error('Failed to fetch providers'));
          setLoading(false);
        }
      }
    };

    loadProviders();

    return () => {
      mountedRef.current = false;
    };
  }, []);

  return { providers, loading, error };
};

// Function to invalidate cache (useful for manual refresh)
export const invalidateProvidersCache = () => {
  cachedProviders = null;
  cacheTimestamp = 0;
  loadingPromise = null;
};
