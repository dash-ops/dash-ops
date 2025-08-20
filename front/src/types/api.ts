/**
 * API and State management types
 */

import { Dispatch } from 'react';

// Generic API response patterns
export interface BaseEntity {
  name: string;
  id?: string;
}

export interface EntityWithStatus extends BaseEntity {
  status: string;
}

// State management patterns
export interface LoadingState<T> {
  data: T[];
  loading: boolean;
}

// Action patterns for reducers
export type LoadingAction<T> =
  | { type: 'LOADING' }
  | { type: 'SET_DATA'; response: T[] };

// Generic reducer type
export type StateReducer<T> = (
  state: LoadingState<T>,
  action: LoadingAction<T>
) => LoadingState<T>;

// Fetch data function signature
export type FetchDataFunction<T, F> = (
  dispatch: Dispatch<LoadingAction<T>>,
  filter: F,
  config?: { signal?: AbortSignal }
) => Promise<void>;

// Update function signature
export type UpdateFunction = (id: string, value: unknown) => void;
