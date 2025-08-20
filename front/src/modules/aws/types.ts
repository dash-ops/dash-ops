/**
 * AWS Module specific types
 */

import { BaseEntity } from '../../types/api';

// AWS Resources
export interface Account extends BaseEntity {
  key: string;
}

export interface Instance extends BaseEntity {
  instance_id: string;
  state: string;
  platform?: string;
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
  state: string;
}
