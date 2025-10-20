import * as serviceResource from '../resources/serviceResource';
import * as serviceAdapter from '../adapters/serviceAdapter';
import type { Service, ServiceList } from '../types';

// Cache state
let servicesCache: Service[] | null = null;
let cacheTimestamp: number = 0;
const CACHE_TTL = 5 * 60 * 1000; // 5 minutes

/**
 * Get cached services or fetch from API
 */
export const getServicesCached = async (): Promise<Service[]> => {
  const now = Date.now();
  
  // Return cached data if still valid
  if (servicesCache && (now - cacheTimestamp) < CACHE_TTL) {
    return servicesCache;
  }
  
  try {
    const response = await serviceResource.getServices();
    servicesCache = serviceAdapter.transformServicesToDomain(response.services);
    cacheTimestamp = now;
    return servicesCache;
  } catch (error) {
    console.error('Failed to fetch services:', error);
    throw error;
  }
};

/**
 * Get cached service by name
 */
export const getServiceCached = async (name: string): Promise<Service | null> => {
  try {
    // First try to get from cache
    if (servicesCache) {
      const cached = servicesCache.find(s => s.metadata.name === name);
      if (cached) return cached;
    }
    
    // If not in cache, fetch from API
    const response = await serviceResource.getService(name);
    const service = serviceAdapter.transformServiceToDomain(response);
    
    // Update cache
    if (servicesCache) {
      const existingIndex = servicesCache.findIndex(s => s.metadata.name === name);
      if (existingIndex >= 0) {
        servicesCache[existingIndex] = service;
      } else {
        servicesCache.push(service);
      }
    }
    
    return service;
  } catch (error) {
    console.error(`Failed to fetch service ${name}:`, error);
    return null;
  }
};

/**
 * Invalidate services cache
 */
export const clearServicesCache = (): void => {
  servicesCache = null;
  cacheTimestamp = 0;
};

/**
 * Update service in cache
 */
export const updateServiceInCache = (service: Service): void => {
  if (!servicesCache) return;
  
  const index = servicesCache.findIndex(s => s.metadata.name === service.metadata.name);
  if (index >= 0) {
    servicesCache[index] = service;
  } else {
    servicesCache.push(service);
  }
  cacheTimestamp = Date.now();
};

/**
 * Remove service from cache
 */
export const removeServiceFromCache = (serviceName: string): void => {
  if (!servicesCache) return;
  
  servicesCache = servicesCache.filter(s => s.metadata.name !== serviceName);
  cacheTimestamp = Date.now();
};

/**
 * Check if cache is valid
 */
export const isCacheValid = (): boolean => {
  const now = Date.now();
  return !!(servicesCache && (now - cacheTimestamp) < CACHE_TTL);
};

/**
 * Get cache age in milliseconds
 */
export const getCacheAge = (): number => {
  if (!servicesCache) return Infinity;
  return Date.now() - cacheTimestamp;
};
