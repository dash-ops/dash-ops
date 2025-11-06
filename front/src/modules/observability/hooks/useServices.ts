import { useCallback, useEffect, useState } from 'react';
import * as servicesResource from '../resources/servicesResource';
import type { ServiceWithContext, ServicesQueryFilters, ServicesResponse } from '../types';

interface UseServicesState {
  data: ServicesResponse | null;
  loading: boolean;
  error: string | null;
}

export const useServices = (initialFilters?: ServicesQueryFilters) => {
  const [filters, setFilters] = useState<ServicesQueryFilters>(initialFilters || {});
  const [state, setState] = useState<UseServicesState>({
    data: null,
    loading: false,
    error: null,
  });

  const fetchServices = useCallback(async () => {
    try {
      setState((s) => ({ ...s, loading: true, error: null }));
      const resp = await servicesResource.getServices(filters);
      setState({ data: resp.data, loading: false, error: null });
    } catch (e: unknown) {
      const message = e instanceof Error ? e.message : 'Failed to fetch services';
      setState((s) => ({ ...s, loading: false, error: message }));
    }
  }, [filters]);

  const refresh = useCallback(async () => {
    await fetchServices();
  }, [fetchServices]);

  const updateFilters = useCallback((partial: Partial<ServicesQueryFilters>) => {
    setFilters((f) => ({ ...f, ...partial }));
  }, []);

  useEffect(() => {
    fetchServices();
  }, [fetchServices]);

  const getServiceNames = useCallback((): string[] => {
    if (!state.data?.services) return [];
    return state.data.services.map((service) => service.service_name);
  }, [state.data]);

  const findService = useCallback((serviceName: string): ServiceWithContext | undefined => {
    if (!state.data?.services) return undefined;
    return state.data.services.find((s) => s.service_name === serviceName);
  }, [state.data]);

  return {
    ...state,
    filters,
    fetchServices,
    refresh,
    updateFilters,
    getServiceNames,
    findService,
  };
};

