import http from '../../../helpers/http';

// Get all services with optional filters
export function getServices(filters = {}) {
  const params = new URLSearchParams();

  if (filters.tier && filters.tier !== 'all') {
    params.append('tier', filters.tier);
  }
  if (filters.team) {
    params.append('team', filters.team);
  }
  if (filters.status) {
    params.append('status', filters.status);
  }
  if (filters.search) {
    params.append('search', filters.search);
  }

  const queryString = params.toString();
  const url = queryString
    ? `/v1/servicecatalog/services?${queryString}`
    : '/v1/servicecatalog/services';

  return http.get(url);
}

// Get service by ID
export function getService(id) {
  return http.get(`/v1/servicecatalog/services/${id}`);
}

// Create new service
export function createService(serviceData) {
  return http.post('/v1/servicecatalog/services', serviceData);
}

// Update service
export function updateService(id, serviceData) {
  return http.put(`/v1/servicecatalog/services/${id}`, serviceData);
}

// Delete service
export function deleteService(id) {
  return http.delete(`/v1/servicecatalog/services/${id}`);
}

// Get service catalog stats
export function getStats() {
  return http.get('/v1/servicecatalog/stats');
}
