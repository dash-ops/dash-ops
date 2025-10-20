/**
 * Service Catalog Adapters
 * Pure functions for transforming service data between API and domain models
 */

import type {
  Service,
  ServiceHealth,
  ServiceCardData,
  ServiceFormData,
  ServiceStats,
} from '../types';

/**
 * Transform API service to domain model
 */
export const transformServiceToDomain = (apiService: any): Service => {
  return {
    apiVersion: apiService.apiVersion || 'v1',
    kind: apiService.kind || 'Service',
    metadata: {
      name: apiService.metadata?.name || '',
      tier: apiService.metadata?.tier || 'TIER-3',
      created_at: apiService.metadata?.created_at,
      created_by: apiService.metadata?.created_by,
      updated_at: apiService.metadata?.updated_at,
      updated_by: apiService.metadata?.updated_by,
      version: apiService.metadata?.version,
    },
    spec: {
      description: apiService.spec?.description || '',
      team: {
        github_team: apiService.spec?.team?.github_team || '',
      },
      business: {
        sla_target: apiService.spec?.business?.sla_target,
        dependencies: apiService.spec?.business?.dependencies,
        impact: apiService.spec?.business?.impact,
      },
      technology: apiService.spec?.technology ? {
        language: apiService.spec.technology.language,
        framework: apiService.spec.technology.framework,
      } : undefined,
      kubernetes: apiService.spec?.kubernetes ? {
        environments: apiService.spec.kubernetes.environments?.map((env: any) => ({
          name: env.name || '',
          context: env.context || '',
          namespace: env.namespace || '',
          resources: {
            deployments: env.resources?.deployments?.map((dep: any) => ({
              name: dep.name || '',
              replicas: dep.replicas || 1,
              resources: {
                requests: {
                  cpu: dep.resources?.requests?.cpu || '100m',
                  memory: dep.resources?.requests?.memory || '128Mi',
                },
                limits: {
                  cpu: dep.resources?.limits?.cpu || '500m',
                  memory: dep.resources?.limits?.memory || '512Mi',
                },
              },
            })) || [],
            services: env.resources?.services || [],
            configmaps: env.resources?.configmaps || [],
          },
        })) || [],
      } : undefined,
      observability: apiService.spec?.observability ? {
        metrics: apiService.spec.observability.metrics,
        logs: apiService.spec.observability.logs,
        traces: apiService.spec.observability.traces,
      } : undefined,
      runbooks: apiService.spec?.runbooks?.map((runbook: any) => ({
        name: runbook.name || '',
        url: runbook.url || '',
      })) || [],
    },
  };
};

/**
 * Transform services list to domain models
 */
export const transformServicesToDomain = (apiServices: any[]): Service[] => {
  return apiServices.map(transformServiceToDomain);
};

/**
 * Transform service to API request format
 */
export const transformServiceToApiRequest = (service: Service): any => {
  return {
    apiVersion: service.apiVersion,
    kind: service.kind,
    metadata: service.metadata,
    spec: service.spec,
  };
};

/**
 * Transform form data to service
 */
export const transformFormDataToService = (formData: ServiceFormData): Service => {
  return {
    apiVersion: 'v1',
    kind: 'Service',
    metadata: {
      name: formData.name,
      tier: formData.tier,
    },
    spec: {
      description: formData.description,
      team: {
        github_team: formData.github_team,
      },
      business: {
        impact: formData.impact,
        sla_target: formData.sla_target || undefined,
      },
      technology: {
        language: formData.language || undefined,
        framework: formData.framework || undefined,
      },
      kubernetes: {
        environments: [{
          name: formData.env_name,
          context: formData.env_context,
          namespace: formData.env_namespace,
          resources: {
            deployments: [{
              name: formData.deployment_name,
              replicas: formData.deployment_replicas,
              resources: {
                requests: {
                  cpu: formData.cpu_request,
                  memory: formData.memory_request,
                },
                limits: {
                  cpu: formData.cpu_limit,
                  memory: formData.memory_limit,
                },
              },
            }],
          },
        }],
      },
      observability: {
        metrics: formData.metrics_url || undefined,
        logs: formData.logs_url || undefined,
      },
    },
  };
};

/**
 * Transform service to card data
 */
export const transformServiceToCardData = (
  service: Service,
  health?: ServiceHealth,
  hasWriteAccess: boolean = false,
  isMyTeam: boolean = false
): ServiceCardData => {
  return {
    service,
    health,
    hasWriteAccess,
    isMyTeam,
  };
};

/**
 * Transform services to card data array
 */
export const transformServicesToCardData = (
  services: Service[],
  healthData: Record<string, ServiceHealth> = {},
  userTeam?: string,
  permissions: Record<string, boolean> = {}
): ServiceCardData[] => {
  return services.map(service => {
    const health = healthData[service.metadata.name];
    const hasWriteAccess = permissions[service.metadata.name] || false;
    const isMyTeam = userTeam ? service.spec.team.github_team === userTeam : false;
    
    return transformServiceToCardData(service, health, hasWriteAccess, isMyTeam);
  });
};

/**
 * Calculate service statistics
 */
export const calculateServiceStats = (
  services: Service[],
  healthData: Record<string, ServiceHealth> = {},
  userTeam?: string
): ServiceStats => {
  const total = services.length;
  const myTeam = userTeam ? services.filter(s => s.spec.team.github_team === userTeam).length : 0;
  const tier1 = services.filter(s => s.metadata.tier === 'TIER-1').length;
  const tier2 = services.filter(s => s.metadata.tier === 'TIER-2').length;
  const tier3 = services.filter(s => s.metadata.tier === 'TIER-3').length;
  
  const critical = Object.values(healthData).filter(h => 
    h.overall_status === 'critical' || h.overall_status === 'down'
  ).length;
  
  return {
    total,
    myTeam,
    tier1,
    tier2,
    tier3,
    critical,
    editable: total, // Assuming all services are editable for now
  };
};

/**
 * Filter services by criteria
 */
export const filterServices = (
  services: Service[],
  filters: {
    search?: string;
    tier?: string;
    team?: string;
    status?: string;
  }
): Service[] => {
  return services.filter(service => {
    // Search filter
    if (filters.search) {
      const searchLower = filters.search.toLowerCase();
      const matchesSearch = 
        service.metadata.name.toLowerCase().includes(searchLower) ||
        service.spec.description.toLowerCase().includes(searchLower) ||
        service.spec.team.github_team.toLowerCase().includes(searchLower);
      
      if (!matchesSearch) return false;
    }
    
    // Tier filter
    if (filters.tier && filters.tier !== 'all') {
      if (service.metadata.tier !== filters.tier) return false;
    }
    
    // Team filter
    if (filters.team && filters.team !== 'all') {
      if (service.spec.team.github_team !== filters.team) return false;
    }
    
    return true;
  });
};

/**
 * Sort services by criteria
 */
export const sortServices = (
  services: Service[],
  sortBy: 'name' | 'tier' | 'team' | 'updated_at'
): Service[] => {
  return [...services].sort((a, b) => {
    switch (sortBy) {
      case 'name':
        return a.metadata.name.localeCompare(b.metadata.name);
      case 'tier':
        const tierOrder = { 'TIER-1': 1, 'TIER-2': 2, 'TIER-3': 3 };
        return (tierOrder[a.metadata.tier] || 3) - (tierOrder[b.metadata.tier] || 3);
      case 'team':
        return a.spec.team.github_team.localeCompare(b.spec.team.github_team);
      case 'updated_at':
        const aTime = new Date(a.metadata.updated_at || a.metadata.created_at || 0).getTime();
        const bTime = new Date(b.metadata.updated_at || b.metadata.created_at || 0).getTime();
        return bTime - aTime; // Most recent first
      default:
        return 0;
    }
  });
};
