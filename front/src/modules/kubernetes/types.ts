/**
 * Kubernetes Module specific types
 */

import { BaseEntity, EntityWithStatus } from '../../types/api';

// Kubernetes Resources
export interface Cluster extends EntityWithStatus {}

export interface Namespace extends EntityWithStatus {}

export interface Node extends BaseEntity {
  ready: string;
  allocated_resources?: AllocatedResources;
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
  condition_status: ConditionStatus;
  restart_count: number;
  node_name: string;
}

export interface ConditionStatus {
  status: string;
}

export interface Deployment extends BaseEntity {
  namespace: string;
  pod_count: number;
  pod_info: PodInfo;
}

export interface PodInfo {
  current: number;
  desired: number;
}

export interface LogContainer {
  name: string;
  log: string;
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

// Component Props
export interface DeploymentActionsProps {
  context: string;
  deployment: Deployment;
  toUp: () => void;
  toDown: () => void;
}
