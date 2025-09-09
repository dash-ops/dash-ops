/**
 * Centralized types export
 *
 * Import types from this file to use shared types across the application:
 *
 * @example
 * ```typescript
 * import { Menu, Router, LoadingState } from '@/types';
 * import { Instance } from '@/types/modules/aws';
 * ```
 */

// Re-export common types
export * from './common';
export * from './api';
export * from './ui';

// Re-export module-specific types with namespace
export * as AWSTypes from '../modules/aws/types';
export * as KubernetesTypes from '../modules/kubernetes/types';
export * as AuthTypes from '../modules/oauth2/types';
export * as ConfigTypes from '../modules/config/types';
