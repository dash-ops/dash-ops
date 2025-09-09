/**
 * Auth Module specific types
 */

import { BaseEntity } from '../../types/api';

// OAuth2 User types
export interface UserData extends BaseEntity {
  id: string;
  login?: string;
  email?: string;
  avatar_url?: string;
  bio?: string;
  location?: string;
  company?: string;
  blog?: string;
  html_url?: string;
}

export interface Team extends BaseEntity {
  slug?: string;
}

export interface UserPermission extends BaseEntity {
  organization?: string;
  teams?: Team[];
  permissions?: {
    [plugin: string]: {
      [feature: string]: string | string[];
    };
  };
}

// Extended interfaces for ProfilePage
export interface PermissionData {
  name: string;
  resource: string;
  actions: string[];
}

export interface PluginPermissions {
  [feature: string]: string[];
}

export interface FeatureActions {
  [action: string]: unknown;
}
