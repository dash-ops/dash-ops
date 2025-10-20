import { useState, useCallback, useMemo } from 'react';
import * as serviceResource from '../resources/serviceResource';
import * as serviceAdapter from '../adapters/serviceAdapter';
import type { Service, ServiceFilters } from '../types';

interface UseServiceFiltersState {
  filters: ServiceFilters;
  loading: boolean;
  error: string | null;
}

interface UseServiceFiltersReturn extends UseServiceFiltersState {
  setSearch: (search: string) => void;
  setTier: (tier: string) => void;
  setTeam: (team: string) => void;
  setStatus: (status: string) => void;
  setSortBy: (sortBy: ServiceFilters['sortBy']) => void;
  clearFilters: () => void;
  applyFilters: () => Promise<Service[]>;
}

export const useServiceFilters = (
  initialFilters: Partial<ServiceFilters> = {}
): UseServiceFiltersReturn => {
  const [state, setState] = useState<UseServiceFiltersState>({
    filters: {
      search: '',
      tier: 'all',
      team: 'all',
      status: 'all',
      sortBy: 'name',
      ...initialFilters,
    },
    loading: false,
    error: null,
  });

  const setSearch = useCallback((search: string) => {
    setState(prev => ({
      ...prev,
      filters: { ...prev.filters, search },
    }));
  }, []);

  const setTier = useCallback((tier: string) => {
    setState(prev => ({
      ...prev,
      filters: { ...prev.filters, tier },
    }));
  }, []);

  const setTeam = useCallback((team: string) => {
    setState(prev => ({
      ...prev,
      filters: { ...prev.filters, team },
    }));
  }, []);

  const setStatus = useCallback((status: string) => {
    setState(prev => ({
      ...prev,
      filters: { ...prev.filters, status },
    }));
  }, []);

  const setSortBy = useCallback((sortBy: ServiceFilters['sortBy']) => {
    setState(prev => ({
      ...prev,
      filters: { ...prev.filters, sortBy },
    }));
  }, []);

  const clearFilters = useCallback(() => {
    setState(prev => ({
      ...prev,
      filters: {
        search: '',
        tier: 'all',
        team: 'all',
        status: 'all',
        sortBy: 'name',
      },
    }));
  }, []);

  const applyFilters = useCallback(async (): Promise<Service[]> => {
    try {
      setState(prev => ({ ...prev, loading: true, error: null }));

      let services: Service[] = [];

      // Use appropriate API endpoint based on filters
      if (state.filters.search) {
        const response = await serviceResource.searchServices(state.filters.search);
        services = serviceAdapter.transformServicesToDomain(response.services);
      } else if (state.filters.team && state.filters.team !== 'all') {
        const response = await serviceResource.getServicesByTeam(state.filters.team);
        services = serviceAdapter.transformServicesToDomain(response.services);
      } else if (state.filters.tier && state.filters.tier !== 'all') {
        const response = await serviceResource.getServicesByTier(state.filters.tier);
        services = serviceAdapter.transformServicesToDomain(response.services);
      } else {
        const response = await serviceResource.getServices();
        services = serviceAdapter.transformServicesToDomain(response.services);
      }

      // Apply client-side filtering and sorting
      const filteredServices = serviceAdapter.filterServices(services, {
        search: state.filters.search,
        tier: state.filters.tier,
        team: state.filters.team,
        status: state.filters.status,
      });

      const sortedServices = serviceAdapter.sortServices(filteredServices, state.filters.sortBy);

      setState(prev => ({ ...prev, loading: false }));
      return sortedServices;
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to apply filters';
      setState(prev => ({ ...prev, error: errorMessage, loading: false }));
      return [];
    }
  }, [state.filters]);

  return {
    ...state,
    setSearch,
    setTier,
    setTeam,
    setStatus,
    setSortBy,
    clearFilters,
    applyFilters,
  };
};
