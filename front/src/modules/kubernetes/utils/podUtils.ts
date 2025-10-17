import { KubernetesTypes } from '../types';

/**
 * Get display name for pod
 */
export const getPodDisplayName = (pod: KubernetesTypes.Pod): string => {
  return pod.name || pod.id;
};

/**
 * Get pod status color for UI
 */
export const getPodStatusColor = (phase: string): string => {
  switch (phase.toLowerCase()) {
    case 'running':
      return 'text-green-600 bg-green-100';
    case 'pending':
      return 'text-yellow-600 bg-yellow-100';
    case 'failed':
    case 'error':
      return 'text-red-600 bg-red-100';
    case 'succeeded':
      return 'text-blue-600 bg-blue-100';
    case 'unknown':
      return 'text-gray-600 bg-gray-100';
    default:
      return 'text-gray-600 bg-gray-100';
  }
};

/**
 * Check if pod is healthy
 */
export const isPodHealthy = (pod: KubernetesTypes.Pod): boolean => {
  return pod.phase === 'Running' && pod.ready !== '0/0';
};

/**
 * Get pod readiness percentage
 */
export const getPodReadinessPercentage = (ready: string): number => {
  const [readyCount, totalCount] = ready.split('/').map(Number);
  if (totalCount === 0) return 0;
  return Math.round((readyCount / totalCount) * 100);
};

/**
 * Format pod age for display
 */
export const formatPodAge = (age: string): string => {
  return age || 'Unknown';
};

/**
 * Get pod restart count
 */
export const getPodRestartCount = (pod: KubernetesTypes.Pod): number => {
  return pod.restarts || 0;
};

/**
 * Check if pod has restarted recently
 */
export const hasRecentRestarts = (pod: KubernetesTypes.Pod): boolean => {
  return getPodRestartCount(pod) > 0;
};

/**
 * Get pod IP address
 */
export const getPodIP = (pod: KubernetesTypes.Pod): string => {
  return pod.ip || 'No IP';
};

/**
 * Check if pod is in error state
 */
export const isPodInError = (pod: KubernetesTypes.Pod): boolean => {
  return pod.phase === 'Failed' || pod.phase === 'Unknown';
};

/**
 * Get pod container count
 */
export const getPodContainerCount = (pod: KubernetesTypes.Pod): number => {
  return pod.containers?.length || 0;
};

/**
 * Get running container count
 */
export const getRunningContainerCount = (pod: KubernetesTypes.Pod): number => {
  return pod.containers?.filter(container => container.ready).length || 0;
};
