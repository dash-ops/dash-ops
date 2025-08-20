/**
 * Config Module specific types
 */

import { BaseEntity } from '../../types/api';

export interface Plugin extends BaseEntity {
  enabled: boolean;
}

export interface PluginsResponse {
  data: string[]; // API returns array of plugin names as strings
}
