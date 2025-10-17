import { KubernetesTypes } from '../types';

/**
 * Transform API response to domain model for Pod
 */
export const transformPodToDomain = (apiPod: any): KubernetesTypes.Pod => {
  return {
    id: apiPod.id || apiPod.name,
    name: apiPod.name,
    namespace: apiPod.namespace,
    status: apiPod.status,
    phase: apiPod.phase,
    node: apiPod.node,
    restarts: apiPod.restarts || 0,
    ready: apiPod.ready || '0/0',
    ip: apiPod.ip,
    age: apiPod.age,
    created_at: apiPod.created_at,
    containers: transformPodContainers(apiPod.containers || []),
    conditions: transformPodConditions(apiPod.conditions || []),
    qos_class: apiPod.qos_class,
  };
};

/**
 * Transform array of API pods to domain models
 */
export const transformPodsToDomain = (apiPods: any[]): KubernetesTypes.Pod[] => {
  return apiPods.map(transformPodToDomain);
};

/**
 * Transform API pod containers to domain models
 */
export const transformPodContainers = (apiContainers: any[]): KubernetesTypes.PodContainer[] => {
  return apiContainers.map((container) => ({
    name: container.name,
    image: container.image,
    ready: container.ready || false,
    restart_count: container.restart_count || 0,
    state: transformPodContainerState(container.state),
    resources: transformPodContainerResources(container.resources),
  }));
};

/**
 * Transform API pod container state to domain model
 */
export const transformPodContainerState = (apiState: any): KubernetesTypes.PodContainerState => {
  return {
    running: apiState.running ? {
      started_at: apiState.running.started_at,
    } : undefined,
    waiting: apiState.waiting ? {
      reason: apiState.waiting.reason,
      message: apiState.waiting.message,
    } : undefined,
    terminated: apiState.terminated ? {
      exit_code: apiState.terminated.exit_code,
      reason: apiState.terminated.reason,
      started_at: apiState.terminated.started_at,
      finished_at: apiState.terminated.finished_at,
    } : undefined,
  };
};

/**
 * Transform API pod container resources to domain model
 */
export const transformPodContainerResources = (apiResources: any): KubernetesTypes.PodContainerResources => {
  return {
    requests: {
      cpu: apiResources?.requests?.cpu || '0',
      memory: apiResources?.requests?.memory || '0',
    },
    limits: {
      cpu: apiResources?.limits?.cpu || '0',
      memory: apiResources?.limits?.memory || '0',
    },
  };
};

/**
 * Transform API pod conditions to domain models
 */
export const transformPodConditions = (apiConditions: any[]): KubernetesTypes.PodCondition[] => {
  return apiConditions.map((condition) => ({
    type: condition.type,
    status: condition.status,
    last_transition_time: condition.last_transition_time,
  }));
};

/**
 * Transform API pod logs response to domain model
 */
export const transformPodLogsToDomain = (apiResponse: any): KubernetesTypes.PodLogsResponse => {
  return {
    pod_name: apiResponse.pod_name,
    namespace: apiResponse.namespace,
    container_name: apiResponse.container_name,
    logs: transformPodLogEntries(apiResponse.logs || []),
    total_lines: apiResponse.total_lines || 0,
  };
};

/**
 * Transform API pod log entries to domain models
 */
export const transformPodLogEntries = (apiLogs: any[]): KubernetesTypes.PodLogEntry[] => {
  return apiLogs.map((log) => ({
    timestamp: log.timestamp,
    message: log.message,
  }));
};
