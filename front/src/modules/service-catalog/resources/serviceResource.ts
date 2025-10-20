import http from '../../../helpers/http';
import type {
  Service,
  ServiceList,
  ServiceHealth,
  ServiceFilters,
} from './types';

const BASE_URL = '/v1/service-catalog';

// Service CRUD operations
export const getServices = async (): Promise<ServiceList> => {
  const response = await http.get<ServiceList>(`${BASE_URL}/services`);
  return response.data;
};

export const getService = async (name: string): Promise<Service> => {
  const response = await http.get<Service>(`${BASE_URL}/services/${name}`);
  return response.data;
};

export const createService = async (service: Service): Promise<Service> => {
  const response = await http.post<Service>(`${BASE_URL}/services`, service);
  return response.data;
};

export const updateService = async (
  name: string,
  service: Service
): Promise<Service> => {
  const response = await http.put<Service>(
    `${BASE_URL}/services/${name}`,
    service
  );
  return response.data;
};

export const deleteService = async (
  name: string
): Promise<{ message: string }> => {
  const response = await http.delete<{ message: string }>(
    `${BASE_URL}/services/${name}`
  );
  return response.data;
};

// Service filtering and search
export const searchServices = async (query: string): Promise<ServiceList> => {
  const response = await http.get<ServiceList>(`${BASE_URL}/services/search`, {
    params: { q: query },
  });
  return response.data;
};

export const getServicesByTeam = async (team: string): Promise<ServiceList> => {
  const response = await http.get<ServiceList>(
    `${BASE_URL}/services/by-team/${encodeURIComponent(team)}`
  );
  return response.data;
};

export const getServicesByTier = async (tier: string): Promise<ServiceList> => {
  const response = await http.get<ServiceList>(
    `${BASE_URL}/services/by-tier/${tier}`
  );
  return response.data;
};

// Service health and monitoring
export const getServiceHealth = async (
  name: string
): Promise<ServiceHealth> => {
  const response = await http.get<ServiceHealth>(
    `${BASE_URL}/services/${name}/health`
  );
  return response.data;
};

export const getServiceHistory = async (name: string) => {
  const response = await http.get(`${BASE_URL}/services/${name}/history`);
  return response.data;
};

// System information
export const getSystemStatus = async () => {
  const response = await http.get(`${BASE_URL}/system/status`);
  return response.data;
};

export const getAllHistory = async (limit?: number) => {
  const response = await http.get(`${BASE_URL}/system/history`, {
    params: limit ? { limit } : {},
  });
  return response.data;
};

// Helper function to build complex queries
export const getFilteredServices = async (
  filters: ServiceFilters
): Promise<ServiceList> => {
  // If there's a search term, use search endpoint
  if (filters.search) {
    return await searchServices(filters.search);
  }

  // If filtering by specific team, use team endpoint
  if (filters.team && filters.team !== 'all') {
    return await getServicesByTeam(filters.team);
  }

  // If filtering by specific tier, use tier endpoint
  if (filters.tier && filters.tier !== 'all') {
    return await getServicesByTier(filters.tier);
  }

  // Otherwise, get all services
  return await getServices();
};

// Health status helpers
export const getServiceHealthBatch = async (
  serviceNames: string[]
): Promise<Record<string, ServiceHealth>> => {
  // Make parallel requests for health data
  const healthPromises = serviceNames.map(async (name) => {
    try {
      const health = await getServiceHealth(name);
      return { name, health };
    } catch (_error) {
      // Return unknown health if API fails
      return {
        name,
        health: {
          service_name: name,
          overall_status: 'unknown' as const,
          environments: [],
          last_updated: new Date().toISOString(),
        },
      };
    }
  });

  const healthResults = await Promise.all(healthPromises);

  // Convert to record format
  return healthResults.reduce(
    (acc, { name, health }) => {
      acc[name] = health;
      return acc;
    },
    {} as Record<string, ServiceHealth>
  );
};
