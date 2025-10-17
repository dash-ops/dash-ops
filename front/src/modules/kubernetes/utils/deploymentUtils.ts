import { KubernetesTypes } from '../types';

/**
 * Get display name for deployment
 */
export const getDeploymentDisplayName = (deployment: KubernetesTypes.Deployment): string => {
  return deployment.name || deployment.id;
};

/**
 * Get deployment status color for UI
 */
export const getDeploymentStatusColor = (deployment: KubernetesTypes.Deployment): string => {
  const { ready, desired } = deployment.replicas;
  
  if (ready === desired && ready > 0) {
    return 'text-green-600 bg-green-100';
  } else if (ready > 0) {
    return 'text-yellow-600 bg-yellow-100';
  } else {
    return 'text-red-600 bg-red-100';
  }
};

/**
 * Check if deployment is healthy
 */
export const isDeploymentHealthy = (deployment: KubernetesTypes.Deployment): boolean => {
  const { ready, desired } = deployment.replicas;
  return ready === desired && ready > 0;
};

/**
 * Get deployment readiness percentage
 */
export const getDeploymentReadinessPercentage = (deployment: KubernetesTypes.Deployment): number => {
  const { ready, desired } = deployment.replicas;
  if (desired === 0) return 0;
  return Math.round((ready / desired) * 100);
};

/**
 * Format deployment age for display
 */
export const formatDeploymentAge = (age: string): string => {
  return age || 'Unknown';
};

/**
 * Get total pod count for deployment
 */
export const getDeploymentPodCount = (deployment: KubernetesTypes.Deployment): number => {
  return deployment.pod_info?.total || 0;
};

/**
 * Get running pod count for deployment
 */
export const getDeploymentRunningPodCount = (deployment: KubernetesTypes.Deployment): number => {
  return deployment.pod_info?.running || 0;
};

/**
 * Get pending pod count for deployment
 */
export const getDeploymentPendingPodCount = (deployment: KubernetesTypes.Deployment): number => {
  return deployment.pod_info?.pending || 0;
};

/**
 * Get failed pod count for deployment
 */
export const getDeploymentFailedPodCount = (deployment: KubernetesTypes.Deployment): number => {
  return deployment.pod_info?.failed || 0;
};

/**
 * Check if deployment has failed pods
 */
export const hasDeploymentFailedPods = (deployment: KubernetesTypes.Deployment): boolean => {
  return getDeploymentFailedPodCount(deployment) > 0;
};

/**
 * Get desired replica count
 */
export const getDeploymentDesiredReplicas = (deployment: KubernetesTypes.Deployment): number => {
  return deployment.replicas?.desired || 0;
};

/**
 * Get current replica count
 */
export const getDeploymentCurrentReplicas = (deployment: KubernetesTypes.Deployment): number => {
  return deployment.replicas?.current || 0;
};

/**
 * Get available replica count
 */
export const getDeploymentAvailableReplicas = (deployment: KubernetesTypes.Deployment): number => {
  return deployment.replicas?.available || 0;
};

/**
 * Check if deployment is scaling up
 */
export const isDeploymentScalingUp = (deployment: KubernetesTypes.Deployment): boolean => {
  const { current, desired } = deployment.replicas;
  return current < desired;
};

/**
 * Check if deployment is scaling down
 */
export const isDeploymentScalingDown = (deployment: KubernetesTypes.Deployment): boolean => {
  const { current, desired } = deployment.replicas;
  return current > desired;
};

/**
 * Get deployment service context information
 */
export const getDeploymentServiceContext = (deployment: KubernetesTypes.Deployment): KubernetesTypes.ServiceContext | undefined => {
  return deployment.service_context;
};

/**
 * Get deployment environment
 */
export const getDeploymentEnvironment = (deployment: KubernetesTypes.Deployment): string | undefined => {
  return deployment.service_context?.environment;
};

/**
 * Get deployment team
 */
export const getDeploymentTeam = (deployment: KubernetesTypes.Deployment): string | undefined => {
  return deployment.service_context?.team;
};

/**
 * Get deployment condition by type
 */
export const getDeploymentCondition = (deployment: KubernetesTypes.Deployment, conditionType: string): KubernetesTypes.DeploymentCondition | undefined => {
  return deployment.conditions?.find(condition => condition.type === conditionType);
};

/**
 * Check if deployment has specific condition status
 */
export const hasDeploymentConditionStatus = (deployment: KubernetesTypes.Deployment, conditionType: string, status: string): boolean => {
  const condition = getDeploymentCondition(deployment, conditionType);
  return condition?.status === status;
};
