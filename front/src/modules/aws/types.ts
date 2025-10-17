/**
 * AWS Module specific types
 * 
 * These types mirror the exact structure returned by the backend API endpoints.
 * Any changes here should be synchronized with the backend wire types.
 */

import { BaseEntity } from '../../types/api';

// AWS Resources - mirroring backend wire/responses.go
export interface Account extends BaseEntity {
  name: string;
  key: string;
  region: string;
  status: string;
  error?: string;
}

export interface InstanceState {
  name: string;
  code: number;
}

export interface InstanceCPU {
  vcpus: number;
  utilization?: number;
}

export interface InstanceMemory {
  size_gb: number;
  utilization?: number;
}

export interface Tag {
  key: string;
  value: string;
}

export interface SecurityGroup {
  group_id: string;
  group_name: string;
}

export interface Instance extends BaseEntity {
  instance_id: string;
  name: string;
  state: InstanceState;
  platform: string;
  instance_type: string;
  public_ip: string;
  private_ip: string;
  subnet_id?: string;
  vpc_id?: string;
  cpu: InstanceCPU;
  memory: InstanceMemory;
  tags: Tag[];
  launch_time: string;
  account: string;
  region: string;
  security_groups?: SecurityGroup[];
  cost_estimate: number;
}

export interface AWSPermission extends BaseEntity {
  resource: string;
  actions: string[];
}

// State management types
export interface InstanceState {
  data: Instance[];
  loading: boolean;
}

// Action types for reducers
export type InstanceAction =
  | { type: 'LOADING' }
  | { type: 'SET_DATA'; response: Instance[] };

// AWS Component Props
export interface InstanceActionsProps {
  instance: Instance;
  toStart: () => void;
  toStop: () => void;
}

export interface InstanceTagProps {
  state: string | InstanceState;
}
