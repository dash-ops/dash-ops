/**
 * Utility functions for AWS instances
 * 
 * Pure functions for common instance operations and formatting.
 */

import { AWSTypes } from '@/types';

/**
 * Get display name for an instance
 */
export function getInstanceDisplayName(instance: AWSTypes.Instance): string {
  // Try to get name from tags first, then fallback to instance name
  const nameTag = instance.tags.find(tag => tag.key === 'Name' || tag.key === 'name');
  return nameTag?.value || instance.name || instance.instance_id;
}

/**
 * Get status color for instance state
 */
export function getInstanceStateColor(stateName?: string): string {
  switch (stateName?.toLowerCase()) {
    case 'running':
      return 'text-green-600 bg-green-50 border-green-200';
    case 'stopped':
      return 'text-red-600 bg-red-50 border-red-200';
    case 'stopping':
    case 'pending':
      return 'text-yellow-600 bg-yellow-50 border-yellow-200';
    default:
      return 'text-gray-600 bg-gray-50 border-gray-200';
  }
}

/**
 * Check if instance can be started
 */
export function canStartInstance(instance: AWSTypes.Instance): boolean {
  return instance.state.name?.toLowerCase() === 'stopped';
}

/**
 * Check if instance can be stopped
 */
export function canStopInstance(instance: AWSTypes.Instance): boolean {
  return instance.state.name?.toLowerCase() === 'running';
}

/**
 * Format launch time for display
 */
export function formatLaunchTime(launchTime: string): string {
  try {
    const date = new Date(launchTime);
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  } catch {
    return 'Invalid date';
  }
}

/**
 * Get instance type display name
 */
export function getInstanceTypeDisplay(type: string): string {
  return type || 'Unknown';
}

/**
 * Check if instance has public IP
 */
export function hasPublicIp(instance: AWSTypes.Instance): boolean {
  return Boolean(instance.public_ip && instance.public_ip !== 'None');
}

/**
 * Get environment from tags
 */
export function getInstanceEnvironment(instance: AWSTypes.Instance): string {
  const envTag = instance.tags.find(tag => 
    tag.key === 'Environment' || 
    tag.key === 'environment' || 
    tag.key === 'env'
  );
  return envTag?.value || 'Unknown';
}
