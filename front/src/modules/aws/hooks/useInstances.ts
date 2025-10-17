/**
 * Custom hook for managing AWS instances state and operations
 * 
 * This hook encapsulates all the state management and side effects
 * related to AWS instances, following React best practices.
 */

import { useState, useEffect, useCallback } from 'react';
import { AWSTypes } from '@/types';
import { AccountFilter } from '@/types';
import { getInstances, startInstance, stopInstance } from '../resources/instanceResource';
import { transformInstancesToDomain } from '../adapters/instanceAdapter';

interface UseInstancesReturn {
  instances: AWSTypes.Instance[];
  loading: boolean;
  error: string | null;
  refresh: () => Promise<void>;
  startInstance: (instanceId: string) => Promise<void>;
  stopInstance: (instanceId: string) => Promise<void>;
}

export function useInstances(filter: AccountFilter): UseInstancesReturn {
  const [instances, setInstances] = useState<AWSTypes.Instance[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchInstances = useCallback(async () => {
    if (!filter.accountKey) return;

    setLoading(true);
    setError(null);

    try {
      const response = await getInstances(filter);
      const transformedInstances = transformInstancesToDomain(response.data);
      setInstances(transformedInstances);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to fetch instances';
      setError(errorMessage);
      console.error('Error fetching instances:', err);
    } finally {
      setLoading(false);
    }
  }, [filter.accountKey]);

  const handleStartInstance = useCallback(async (instanceId: string) => {
    try {
      setError(null);
      await startInstance(filter.accountKey, instanceId);
      
      // Refresh instances after successful operation
      await fetchInstances();
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to start instance';
      setError(errorMessage);
      console.error('Error starting instance:', err);
      // Don't call fetchInstances here as it would reset the error
    }
  }, [filter.accountKey, fetchInstances]);

  const handleStopInstance = useCallback(async (instanceId: string) => {
    try {
      setError(null);
      await stopInstance(filter.accountKey, instanceId);
      
      // Refresh instances after successful operation
      await fetchInstances();
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to stop instance';
      setError(errorMessage);
      console.error('Error stopping instance:', err);
      // Don't call fetchInstances here as it would reset the error
    }
  }, [filter.accountKey, fetchInstances]);

  // Fetch instances when filter changes
  useEffect(() => {
    fetchInstances();
  }, [fetchInstances]);

  return {
    instances,
    loading,
    error,
    refresh: fetchInstances,
    startInstance: handleStartInstance,
    stopInstance: handleStopInstance,
  };
}
