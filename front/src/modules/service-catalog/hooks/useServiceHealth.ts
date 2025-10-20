import { useState, useCallback } from 'react';
import * as serviceResource from '../resources/serviceResource';
import type { ServiceHealth, Service } from '../types';

interface UseServiceHealthState {
  healthData: Record<string, ServiceHealth>;
  loading: boolean;
  error: string | null;
}

interface UseServiceHealthReturn extends UseServiceHealthState {
  fetchServiceHealth: (serviceName: string) => Promise<ServiceHealth | null>;
  fetchBatchHealth: (serviceNames: string[]) => Promise<void>;
  getServiceHealth: (serviceName: string) => ServiceHealth | undefined;
  refreshHealth: (serviceName: string) => Promise<void>;
}

export const useServiceHealth = (): UseServiceHealthReturn => {
  const [state, setState] = useState<UseServiceHealthState>({
    healthData: {},
    loading: false,
    error: null,
  });

  const fetchServiceHealth = useCallback(async (serviceName: string): Promise<ServiceHealth | null> => {
    try {
      setState(prev => ({ ...prev, loading: true, error: null }));
      
      const health = await serviceResource.getServiceHealth(serviceName);
      
      setState(prev => ({
        ...prev,
        healthData: {
          ...prev.healthData,
          [serviceName]: health,
        },
        loading: false,
      }));
      
      return health;
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to fetch service health';
      setState(prev => ({ ...prev, error: errorMessage, loading: false }));
      
      // Return unknown health on error
      const unknownHealth: ServiceHealth = {
        service_name: serviceName,
        overall_status: 'unknown',
        environments: [],
        last_updated: new Date().toISOString(),
      };
      
      setState(prev => ({
        ...prev,
        healthData: {
          ...prev.healthData,
          [serviceName]: unknownHealth,
        },
      }));
      
      return unknownHealth;
    }
  }, []);

  const fetchBatchHealth = useCallback(async (serviceNames: string[]): Promise<void> => {
    try {
      setState(prev => ({ ...prev, loading: true, error: null }));
      
      const healthData = await serviceResource.getServiceHealthBatch(serviceNames);
      
      setState(prev => ({
        ...prev,
        healthData: {
          ...prev.healthData,
          ...healthData,
        },
        loading: false,
      }));
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to fetch batch health';
      setState(prev => ({ ...prev, error: errorMessage, loading: false }));
    }
  }, []);

  const getServiceHealth = useCallback((serviceName: string): ServiceHealth | undefined => {
    return state.healthData[serviceName];
  }, [state.healthData]);

  const refreshHealth = useCallback(async (serviceName: string): Promise<void> => {
    await fetchServiceHealth(serviceName);
  }, [fetchServiceHealth]);

  return {
    ...state,
    fetchServiceHealth,
    fetchBatchHealth,
    getServiceHealth,
    refreshHealth,
  };
};
