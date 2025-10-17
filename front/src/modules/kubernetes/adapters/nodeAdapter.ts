import { KubernetesTypes } from '../types';

/**
 * Transform API response to domain model for Node
 */
export const transformNodeToDomain = (apiNode: any): KubernetesTypes.Node => {
  return {
    id: apiNode.id || apiNode.name,
    name: apiNode.name,
    status: apiNode.status,
    roles: apiNode.roles || [],
    age: apiNode.age,
    version: apiNode.version,
    internal_ip: apiNode.internal_ip,
    conditions: transformNodeConditions(apiNode.conditions || []),
    resources: transformNodeResources(apiNode.resources),
    created_at: apiNode.created_at,
  };
};

/**
 * Transform array of API nodes to domain models
 */
export const transformNodesToDomain = (apiNodes: any[]): KubernetesTypes.Node[] => {
  return apiNodes.map(transformNodeToDomain);
};

/**
 * Transform API node conditions to domain models
 */
export const transformNodeConditions = (apiConditions: any[]): KubernetesTypes.NodeCondition[] => {
  return apiConditions.map((condition) => ({
    type: condition.type,
    status: condition.status,
    reason: condition.reason,
    message: condition.message,
    last_transition_time: condition.last_transition_time,
  }));
};

/**
 * Transform API node resources to domain model
 */
export const transformNodeResources = (apiResources: any): KubernetesTypes.NodeResources => {
  return {
    capacity: transformResourceSpec(apiResources?.capacity),
    allocatable: transformResourceSpec(apiResources?.allocatable),
    used: transformResourceSpec(apiResources?.used),
  };
};

/**
 * Transform API resource spec to domain model
 */
export const transformResourceSpec = (apiSpec: any): KubernetesTypes.ResourceSpec => {
  return {
    cpu: apiSpec?.cpu || '0',
    memory: apiSpec?.memory || '0',
    pods: apiSpec?.pods || '0',
  };
};

/**
 * Transform API allocated resources to domain model
 */
export const transformAllocatedResources = (apiAllocated: any): KubernetesTypes.AllocatedResources => {
  return {
    cpu_requests_fraction: apiAllocated?.cpu_requests_fraction || 0,
    cpu_limits_fraction: apiAllocated?.cpu_limits_fraction || 0,
    memory_requests_fraction: apiAllocated?.memory_requests_fraction || 0,
    memory_limits_fraction: apiAllocated?.memory_limits_fraction || 0,
    allocated_pods: apiAllocated?.allocated_pods || 0,
    pod_capacity: apiAllocated?.pod_capacity || 0,
  };
};
