import type { ServicesResponse } from '../types';

/**
 * Transforms API response to domain format
 */
export const transformServicesResponseToDomain = (
  response: ServicesResponse
): ServicesResponse => {
  // API response already matches domain format, just return as is
  return response;
};

/**
 * Transforms domain format to API request format
 */
export const transformServicesToDomain = (
  services: ServicesResponse['services']
): ServicesResponse['services'] => {
  return services.map((service) => ({
    ...service,
    // Add any domain-specific transformations here if needed
  }));
};

