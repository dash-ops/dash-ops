import { KubernetesTypes } from '../types';

/**
 * Transform API response to domain model for Deployment
 */
export const transformDeploymentToDomain = (apiDeployment: any): KubernetesTypes.Deployment => {
  return {
    id: apiDeployment.id || apiDeployment.name,
    name: apiDeployment.name,
    namespace: apiDeployment.namespace,
    pod_count: apiDeployment.pod_count || 0,
    pod_info: transformPodInfo(apiDeployment.pod_info),
    replicas: transformDeploymentReplicas(apiDeployment.replicas),
    age: apiDeployment.age,
    created_at: apiDeployment.created_at,
    conditions: transformDeploymentConditions(apiDeployment.conditions || []),
    service_context: transformServiceContext(apiDeployment.service_context),
  };
};

/**
 * Transform array of API deployments to domain models
 */
export const transformDeploymentsToDomain = (apiDeployments: any[]): KubernetesTypes.Deployment[] => {
  return apiDeployments.map(transformDeploymentToDomain);
};

/**
 * Transform API pod info to domain model
 */
export const transformPodInfo = (apiPodInfo: any): KubernetesTypes.PodInfo => {
  return {
    running: apiPodInfo?.running || 0,
    pending: apiPodInfo?.pending || 0,
    failed: apiPodInfo?.failed || 0,
    total: apiPodInfo?.total || 0,
  };
};

/**
 * Transform API deployment replicas to domain model
 */
export const transformDeploymentReplicas = (apiReplicas: any): KubernetesTypes.DeploymentReplicas => {
  return {
    ready: apiReplicas?.ready || 0,
    available: apiReplicas?.available || 0,
    current: apiReplicas?.current || 0,
    desired: apiReplicas?.desired || 0,
  };
};

/**
 * Transform API deployment conditions to domain models
 */
export const transformDeploymentConditions = (apiConditions: any[]): KubernetesTypes.DeploymentCondition[] => {
  return apiConditions.map((condition) => ({
    type: condition.type,
    status: condition.status,
    reason: condition.reason,
    message: condition.message,
  }));
};

/**
 * Transform API service context to domain model
 */
export const transformServiceContext = (apiContext: any): KubernetesTypes.ServiceContext | undefined => {
  if (!apiContext) return undefined;
  
  return {
    service_name: apiContext.service_name,
    service_tier: apiContext.service_tier,
    environment: apiContext.environment,
    context: apiContext.context,
    team: apiContext.team,
    description: apiContext.description,
  };
};
