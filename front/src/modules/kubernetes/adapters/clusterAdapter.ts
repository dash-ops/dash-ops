import { KubernetesTypes } from '../types';

/**
 * Transform API response to domain model for Cluster
 */
export const transformClusterToDomain = (apiCluster: any): KubernetesTypes.Cluster => {
  return {
    id: apiCluster.id || apiCluster.name,
    name: apiCluster.name,
    context: apiCluster.context,
    version: apiCluster.version,
    status: apiCluster.status,
  };
};

/**
 * Transform array of API clusters to domain models
 */
export const transformClustersToDomain = (apiClusters: any[]): KubernetesTypes.Cluster[] => {
  return apiClusters.map(transformClusterToDomain);
};

/**
 * Transform API cluster list response to domain model
 */
export const transformClusterListResponseToDomain = (apiResponse: any): KubernetesTypes.ClusterListResponse => {
  return {
    clusters: transformClustersToDomain(apiResponse.clusters || []),
    total: apiResponse.total || 0,
  };
};

/**
 * Transform API namespace to domain model
 */
export const transformNamespaceToDomain = (apiNamespace: any): KubernetesTypes.Namespace => {
  return {
    id: apiNamespace.id || apiNamespace.name,
    name: apiNamespace.name,
    status: apiNamespace.status,
  };
};

/**
 * Transform array of API namespaces to domain models
 */
export const transformNamespacesToDomain = (apiNamespaces: any[]): KubernetesTypes.Namespace[] => {
  return apiNamespaces.map(transformNamespaceToDomain);
};
