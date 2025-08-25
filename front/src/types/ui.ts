/**
 * UI Component types
 */

import type { Page } from './common';

// Badge variants
export type BadgeVariant = 'default' | 'secondary' | 'outline' | 'destructive';

// Instance states for AWS
export type InstanceState =
  | 'pending'
  | 'running'
  | 'shutting-down'
  | 'terminated'
  | 'stopping'
  | 'stopped'
  | 'loading';

// Pod states for Kubernetes
export type PodStatus = 'Running' | 'Succeeded' | 'Pending' | 'Failed';

// Component props patterns
export interface RefreshProps {
  onReload: () => void | Promise<void>;
}

export interface ProgressDataProps {
  percent: number;
}

// Layout component interfaces
export interface ContentWithMenuProps {
  pages: Page[];
  paramName?: string;
  contextValue?: string;
}

export interface MenuItem {
  path: string;
  name: string;
}
