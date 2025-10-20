import { useState, useEffect, useCallback } from 'react';
import * as serviceResource from '../resources/serviceResource';
import * as serviceAdapter from '../adapters/serviceAdapter';
import type { Service, ServiceFilters } from '../types';

interface UseServicesState {
  services: Service[];
  loading: boolean;
  error: string | null;
}

interface UseServicesReturn extends UseServicesState {
  fetchServices: () => Promise<void>;
  refresh: () => Promise<void>;
  createService: (service: Service) => Promise<Service | null>;
  updateService: (name: string, service: Service) => Promise<Service | null>;
  deleteService: (name: string) => Promise<boolean>;
}

export const useServices = (): UseServicesReturn => {
  const [state, setState] = useState<UseServicesState>({
    services: [],
    loading: false,
    error: null,
  });

  const fetchServices = useCallback(async () => {
    try {
      setState(prev => ({ ...prev, loading: true, error: null }));
      
      const response = await serviceResource.getServices();
      const services = serviceAdapter.transformServicesToDomain(response.services);
      
      setState(prev => ({ 
        ...prev, 
        services, 
        loading: false 
      }));
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to fetch services';
      setState(prev => ({ 
        ...prev, 
        error: errorMessage, 
        loading: false 
      }));
    }
  }, []);

  const refresh = useCallback(async () => {
    await fetchServices();
  }, [fetchServices]);

  const createService = useCallback(async (service: Service): Promise<Service | null> => {
    try {
      const apiService = serviceAdapter.transformServiceToApiRequest(service);
      const response = await serviceResource.createService(apiService);
      const newService = serviceAdapter.transformServiceToDomain(response);
      
      setState(prev => ({
        ...prev,
        services: [...prev.services, newService],
      }));
      
      return newService;
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to create service';
      setState(prev => ({ ...prev, error: errorMessage }));
      return null;
    }
  }, []);

  const updateService = useCallback(async (name: string, service: Service): Promise<Service | null> => {
    try {
      const apiService = serviceAdapter.transformServiceToApiRequest(service);
      const response = await serviceResource.updateService(name, apiService);
      const updatedService = serviceAdapter.transformServiceToDomain(response);
      
      setState(prev => ({
        ...prev,
        services: prev.services.map(s => 
          s.metadata.name === name ? updatedService : s
        ),
      }));
      
      return updatedService;
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to update service';
      setState(prev => ({ ...prev, error: errorMessage }));
      return null;
    }
  }, []);

  const deleteService = useCallback(async (name: string): Promise<boolean> => {
    try {
      await serviceResource.deleteService(name);
      
      setState(prev => ({
        ...prev,
        services: prev.services.filter(s => s.metadata.name !== name),
      }));
      
      return true;
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to delete service';
      setState(prev => ({ ...prev, error: errorMessage }));
      return false;
    }
  }, []);

  useEffect(() => {
    fetchServices();
  }, [fetchServices]);

  return {
    ...state,
    fetchServices,
    refresh,
    createService,
    updateService,
    deleteService,
  };
};
