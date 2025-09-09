/**
 * Kubernetes Module specific types
 */

import { BaseEntity, EntityWithStatus } from '../../types/api';

// Kubernetes Resources
export interface Cluster extends EntityWithStatus {
  name: string;
  context: string;
  version: string;
  status: string;
}

export interface ClusterListResponse {
  clusters: Cluster[];
  total: number;
}

export interface Namespace extends EntityWithStatus {}

export interface Node extends BaseEntity {
  status: string;
  roles: string[];
  age: string;
  version: string;
  internal_ip: string;
  conditions: NodeCondition[];
  resources: NodeResources;
  created_at: string;
}

export interface NodeCondition {
  type: string;
  status: string;
  reason?: string;
  message?: string;
  last_transition_time?: string;
}

export interface NodeResources {
  capacity: ResourceSpec;
  allocatable: ResourceSpec;
  used: ResourceSpec;
}

export interface ResourceSpec {
  cpu: string;
  memory: string;
  pods: string;
}

export interface AllocatedResources {
  cpu_requests_fraction: number;
  cpu_limits_fraction: number;
  memory_requests_fraction: number;
  memory_limits_fraction: number;
  allocated_pods: number;
  pod_capacity: number;
}

export interface Pod extends BaseEntity {
  namespace: string;
  status: string;
  phase: string;
  node: string;
  restarts: number;
  ready: string;
  ip: string;
  age: string;
  created_at: string;
  containers: PodContainer[];
  conditions: PodCondition[];
  qos_class?: string;
}

export interface PodContainer {
  name: string;
  image: string;
  ready: boolean;
  restart_count: number;
  state: PodContainerState;
  resources: PodContainerResources;
}

export interface PodContainerState {
  running?: {
    started_at: string;
  };
  waiting?: {
    reason: string;
    message: string;
  };
  terminated?: {
    exit_code: number;
    reason: string;
    started_at: string;
    finished_at: string;
  };
}

export interface PodContainerResources {
  requests: {
    cpu: string;
    memory: string;
  };
  limits: {
    cpu: string;
    memory: string;
  };
}

export interface PodCondition {
  type: string;
  status: string;
  last_transition_time: string;
}

export interface ConditionStatus {
  status: string;
}

export interface ServiceContext {
  service_name?: string;
  service_tier?: string;
  environment?: string;
  context?: string;
  team?: string;
  description?: string;
}

export interface Deployment extends BaseEntity {
  namespace: string;
  pod_count: number;
  pod_info: PodInfo;
  replicas: DeploymentReplicas;
  age: string;
  created_at: string;
  conditions: DeploymentCondition[];
  service_context?: ServiceContext;
}

export interface PodInfo {
  running: number;
  pending: number;
  failed: number;
  total: number;
}

export interface DeploymentReplicas {
  ready: number;
  available: number;
  current: number;
  desired: number;
}

export interface DeploymentCondition {
  type: string;
  status: string;
  reason?: string;
  message?: string;
}

export interface LogContainer {
  name: string;
  log: string;
}

export interface PodLogEntry {
  timestamp: string;
  message: string;
}

export interface PodLogsResponse {
  pod_name: string;
  namespace: string;
  container_name?: string;
  logs: PodLogEntry[];
  total_lines: number;
}

export interface K8sPermission extends BaseEntity {
  resource: string;
  actions: string[];
}

// State management types
export interface PodState {
  data: Pod[];
  loading: boolean;
}

export interface DeploymentState {
  data: Deployment[];
  loading: boolean;
}

export interface NodesState {
  data: Node[];
  loading: boolean;
}

export interface LogState {
  data: LogContainer[];
  loading: boolean;
}

// Action types for reducers
export type PodAction =
  | { type: 'LOADING' }
  | { type: 'SET_DATA'; response: Pod[] };

export type DeploymentAction =
  | { type: 'LOADING' }
  | { type: 'SET_DATA'; response: Deployment[] };

export type NodesAction =
  | { type: 'LOADING' }
  | { type: 'SET_DATA'; response: Node[] };

export type LogAction =
  | { type: 'LOADING' }
  | { type: 'SET_DATA'; response: LogContainer[] };

// Filter types
export interface PodFilter {
  context: string;
  namespace: string;
}

export interface PodLogsFilter {
  context: string;
  name: string;
  namespace: string;
}

export interface DeploymentFilter {
  context: string;
  namespace: string;
}
