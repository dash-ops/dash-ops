import { KubernetesTypes } from '../types';

/**
 * Get display name for node
 */
export const getNodeDisplayName = (node: KubernetesTypes.Node): string => {
  return node.name || node.id;
};

/**
 * Get node status color for UI
 */
export const getNodeStatusColor = (status: string): string => {
  switch (status.toLowerCase()) {
    case 'ready':
      return 'text-green-600 bg-green-100';
    case 'notready':
    case 'not ready':
      return 'text-red-600 bg-red-100';
    case 'unknown':
      return 'text-gray-600 bg-gray-100';
    default:
      return 'text-gray-600 bg-gray-100';
  }
};

/**
 * Check if node is ready
 */
export const isNodeReady = (node: KubernetesTypes.Node): boolean => {
  return node.status === 'Ready';
};

/**
 * Format node age for display
 */
export const formatNodeAge = (age: string): string => {
  return age || 'Unknown';
};

/**
 * Get node internal IP
 */
export const getNodeInternalIP = (node: KubernetesTypes.Node): string => {
  return node.internal_ip || 'No IP';
};

/**
 * Get node roles as string
 */
export const getNodeRoles = (node: KubernetesTypes.Node): string => {
  return node.roles?.join(', ') || 'No roles';
};

/**
 * Check if node is master/control plane
 */
export const isMasterNode = (node: KubernetesTypes.Node): boolean => {
  return node.roles?.some(role => 
    role.toLowerCase().includes('master') || 
    role.toLowerCase().includes('control-plane')
  ) || false;
};

/**
 * Check if node is worker
 */
export const isWorkerNode = (node: KubernetesTypes.Node): boolean => {
  return node.roles?.some(role => role.toLowerCase().includes('worker')) || false;
};

/**
 * Get CPU capacity
 */
export const getNodeCPUCapacity = (node: KubernetesTypes.Node): string => {
  return node.resources?.capacity?.cpu || '0';
};

/**
 * Get memory capacity
 */
export const getNodeMemoryCapacity = (node: KubernetesTypes.Node): string => {
  return node.resources?.capacity?.memory || '0';
};

/**
 * Get pod capacity
 */
export const getNodePodCapacity = (node: KubernetesTypes.Node): string => {
  return node.resources?.capacity?.pods || '0';
};

/**
 * Get CPU usage percentage
 */
export const getNodeCPUUsagePercentage = (node: KubernetesTypes.Node): number => {
  const capacity = parseFloat(getNodeCPUCapacity(node));
  const used = parseFloat(node.resources?.used?.cpu || '0');
  if (capacity === 0) return 0;
  return Math.round((used / capacity) * 100);
};

/**
 * Get memory usage percentage
 */
export const getNodeMemoryUsagePercentage = (node: KubernetesTypes.Node): number => {
  const capacity = parseFloat(getNodeMemoryCapacity(node));
  const used = parseFloat(node.resources?.used?.memory || '0');
  if (capacity === 0) return 0;
  return Math.round((used / capacity) * 100);
};

/**
 * Get node condition by type
 */
export const getNodeCondition = (node: KubernetesTypes.Node, conditionType: string): KubernetesTypes.NodeCondition | undefined => {
  return node.conditions?.find(condition => condition.type === conditionType);
};

/**
 * Check if node has specific condition status
 */
export const hasNodeConditionStatus = (node: KubernetesTypes.Node, conditionType: string, status: string): boolean => {
  const condition = getNodeCondition(node, conditionType);
  return condition?.status === status;
};
