/**
 * Instance Adapter - Data transformation functions for AWS instances
 * 
 * Pure functions for transforming data between API responses and domain models.
 * No classes, no state - just pure functional transformations.
 */

import { AWSTypes } from '@/types';

/**
 * Transform API tags array to key-value object
 */
export function transformTags(tags?: AWSTypes.Tag[]): Record<string, string> {
  if (!tags || !Array.isArray(tags)) {
    return {};
  }

  return tags.reduce((acc, tag) => {
    if (tag.key && tag.value) {
      acc[tag.key] = tag.value;
    }
    return acc;
  }, {} as Record<string, string>);
}

/**
 * Transform raw API instance data to domain model
 */
export function transformInstanceToDomain(apiInstance: any): AWSTypes.Instance {
  return {
    id: apiInstance.instance_id,
    name: apiInstance.name,
    instance_id: apiInstance.instance_id,
    state: apiInstance.state,
    platform: apiInstance.platform,
    instance_type: apiInstance.instance_type,
    public_ip: apiInstance.public_ip,
    private_ip: apiInstance.private_ip,
    subnet_id: apiInstance.subnet_id,
    vpc_id: apiInstance.vpc_id,
    cpu: apiInstance.cpu,
    memory: apiInstance.memory,
    tags: apiInstance.tags || [],
    launch_time: apiInstance.launch_time,
    account: apiInstance.account,
    region: apiInstance.region,
    security_groups: apiInstance.security_groups,
    cost_estimate: apiInstance.cost_estimate,
  };
}

/**
 * Transform domain model to API request format
 */
export function transformInstanceToApiRequest(instance: AWSTypes.Instance): any {
  return {
    instance_id: instance.instance_id,
    // Add any API-specific transformations here
  };
}

/**
 * Transform multiple instances from API response
 */
export function transformInstancesToDomain(apiInstances: any[]): AWSTypes.Instance[] {
  return apiInstances.map(transformInstanceToDomain);
}

/**
 * Filter instances by state
 */
export function filterInstancesByState(
  instances: AWSTypes.Instance[], 
  state: string
): AWSTypes.Instance[] {
  return instances.filter(instance => instance.state.name === state);
}

/**
 * Sort instances by launch time
 */
export function sortInstancesByLaunchTime(
  instances: AWSTypes.Instance[], 
  order: 'asc' | 'desc' = 'desc'
): AWSTypes.Instance[] {
  return [...instances].sort((a, b) => {
    const timeA = new Date(a.launch_time).getTime();
    const timeB = new Date(b.launch_time).getTime();
    return order === 'asc' ? timeA - timeB : timeB - timeA;
  });
}
