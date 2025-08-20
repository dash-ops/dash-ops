/**
 * Config Module specific types
 */

import { BaseEntity } from '../../types/api';

export interface Plugin extends BaseEntity {
  enabled: boolean;
}
