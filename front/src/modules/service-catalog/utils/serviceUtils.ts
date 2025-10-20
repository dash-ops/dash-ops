import type { Service, ServiceHealth, ServiceCardData } from '../types';

/**
 * Get service display name
 */
export const getServiceDisplayName = (service: Service): string => {
  return service.metadata.name || 'Unknown Service';
};

/**
 * Get service tier badge color
 */
export const getServiceTierColor = (tier: Service['metadata']['tier']): string => {
  switch (tier) {
    case 'TIER-1':
      return 'bg-red-100 text-red-800 border-red-200';
    case 'TIER-2':
      return 'bg-yellow-100 text-yellow-800 border-yellow-200';
    case 'TIER-3':
      return 'bg-green-100 text-green-800 border-green-200';
    default:
      return 'bg-gray-100 text-gray-800 border-gray-200';
  }
};

/**
 * Get service tier label
 */
export const getServiceTierLabel = (tier: Service['metadata']['tier']): string => {
  switch (tier) {
    case 'TIER-1':
      return 'Tier 1';
    case 'TIER-2':
      return 'Tier 2';
    case 'TIER-3':
      return 'Tier 3';
    default:
      return 'Unknown';
  }
};

/**
 * Get health status color
 */
export const getHealthStatusColor = (status: ServiceHealth['overall_status']): string => {
  switch (status) {
    case 'healthy':
      return 'bg-green-100 text-green-800 border-green-200';
    case 'degraded':
      return 'bg-yellow-100 text-yellow-800 border-yellow-200';
    case 'down':
    case 'critical':
      return 'bg-red-100 text-red-800 border-red-200';
    case 'drift':
      return 'bg-orange-100 text-orange-800 border-orange-200';
    case 'unknown':
    default:
      return 'bg-gray-100 text-gray-800 border-gray-200';
  }
};

/**
 * Get health status icon
 */
export const getHealthStatusIcon = (status: ServiceHealth['overall_status']): string => {
  switch (status) {
    case 'healthy':
      return 'check-circle';
    case 'degraded':
      return 'alert-triangle';
    case 'down':
    case 'critical':
      return 'x-circle';
    case 'drift':
      return 'trending-up';
    case 'unknown':
    default:
      return 'help-circle';
  }
};

/**
 * Format service creation date
 */
export const formatServiceCreatedDate = (dateString?: string): string => {
  if (!dateString) return 'Unknown';
  
  try {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  } catch {
    return 'Invalid Date';
  }
};

/**
 * Format service updated date
 */
export const formatServiceUpdatedDate = (dateString?: string): string => {
  if (!dateString) return 'Never';
  
  try {
    const date = new Date(dateString);
    const now = new Date();
    const diffInHours = Math.floor((now.getTime() - date.getTime()) / (1000 * 60 * 60));
    
    if (diffInHours < 1) return 'Just now';
    if (diffInHours < 24) return `${diffInHours}h ago`;
    if (diffInHours < 168) return `${Math.floor(diffInHours / 24)}d ago`;
    
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  } catch {
    return 'Invalid Date';
  }
};

/**
 * Get service description preview
 */
export const getServiceDescriptionPreview = (description: string, maxLength: number = 100): string => {
  if (description.length <= maxLength) return description;
  return description.substring(0, maxLength).trim() + '...';
};

/**
 * Check if service has Kubernetes configuration
 */
export const hasKubernetesConfig = (service: Service): boolean => {
  return !!(service.spec.kubernetes?.environments?.length);
};

/**
 * Get total deployments count
 */
export const getTotalDeployments = (service: Service): number => {
  if (!service.spec.kubernetes?.environments) return 0;
  
  return service.spec.kubernetes.environments.reduce((total, env) => {
    return total + (env.resources.deployments?.length || 0);
  }, 0);
};

/**
 * Get service environments
 */
export const getServiceEnvironments = (service: Service): string[] => {
  if (!service.spec.kubernetes?.environments) return [];
  return service.spec.kubernetes.environments.map(env => env.name);
};

/**
 * Check if service is owned by user team
 */
export const isServiceOwnedByTeam = (service: Service, userTeam?: string): boolean => {
  if (!userTeam) return false;
  return service.spec.team.github_team === userTeam;
};

/**
 * Get service impact level
 */
export const getServiceImpact = (service: Service): string => {
  return service.spec.business?.impact || 'low';
};

/**
 * Get service impact color
 */
export const getServiceImpactColor = (impact: string): string => {
  switch (impact) {
    case 'high':
      return 'bg-red-100 text-red-800 border-red-200';
    case 'medium':
      return 'bg-yellow-100 text-yellow-800 border-yellow-200';
    case 'low':
    default:
      return 'bg-green-100 text-green-800 border-green-200';
  }
};

/**
 * Get service technology stack
 */
export const getServiceTechnology = (service: Service): string => {
  const tech = service.spec.technology;
  if (!tech) return 'Unknown';
  
  const parts = [];
  if (tech.language) parts.push(tech.language);
  if (tech.framework) parts.push(tech.framework);
  
  return parts.join(' + ') || 'Unknown';
};

/**
 * Check if service has observability configured
 */
export const hasObservabilityConfig = (service: Service): boolean => {
  const obs = service.spec.observability;
  return !!(obs?.metrics || obs?.logs || obs?.traces);
};

/**
 * Get service runbooks count
 */
export const getRunbooksCount = (service: Service): number => {
  return service.spec.runbooks?.length || 0;
};

/**
 * Validate service name
 */
export const isValidServiceName = (name: string): boolean => {
  // Service names should be lowercase, alphanumeric with hyphens
  const regex = /^[a-z0-9]([a-z0-9-]*[a-z0-9])?$/;
  return regex.test(name) && name.length >= 3 && name.length <= 63;
};

/**
 * Validate GitHub team name
 */
export const isValidGitHubTeam = (team: string): boolean => {
  // GitHub team names should be alphanumeric with hyphens and underscores
  const regex = /^[a-zA-Z0-9]([a-zA-Z0-9_-]*[a-zA-Z0-9])?$/;
  return regex.test(team) && team.length >= 1 && team.length <= 39;
};

/**
 * Generate service URL slug
 */
export const generateServiceSlug = (name: string): string => {
  return name.toLowerCase().replace(/[^a-z0-9-]/g, '-').replace(/-+/g, '-').replace(/^-|-$/g, '');
};

/**
 * Check if service card data has write access
 */
export const canEditService = (cardData: ServiceCardData): boolean => {
  return cardData.hasWriteAccess;
};

/**
 * Check if service is critical based on health
 */
export const isServiceCritical = (health?: ServiceHealth): boolean => {
  if (!health) return false;
  return health.overall_status === 'critical' || health.overall_status === 'down';
};

/**
 * Get service priority score for sorting
 */
export const getServicePriority = (service: Service, health?: ServiceHealth): number => {
  let score = 0;
  
  // Tier priority
  switch (service.metadata.tier) {
    case 'TIER-1':
      score += 100;
      break;
    case 'TIER-2':
      score += 50;
      break;
    case 'TIER-3':
      score += 10;
      break;
  }
  
  // Health priority
  if (health) {
    switch (health.overall_status) {
      case 'critical':
        score += 1000;
        break;
      case 'down':
        score += 500;
        break;
      case 'degraded':
        score += 100;
        break;
      case 'drift':
        score += 50;
        break;
      case 'unknown':
        score += 25;
        break;
      case 'healthy':
        score += 0;
        break;
    }
  }
  
  return score;
};
