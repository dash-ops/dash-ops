import { useState, useEffect, useCallback } from 'react';
import { getClusters } from '../resources/clusterResource';
import { transformClusterListResponseToDomain } from '../adapters/clusterAdapter';
import { KubernetesTypes } from '../types';

export interface UseClustersResult {
  clusters: KubernetesTypes.Cluster[];
  loading: boolean;
  error: string | null;
  fetchClusters: () => Promise<void>;
}

export const useClusters = (): UseClustersResult => {
  const [clusters, setClusters] = useState<KubernetesTypes.Cluster[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchClusters = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await getClusters();
      const transformedResponse = transformClusterListResponseToDomain(response.data);
      setClusters(transformedResponse.clusters);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to fetch clusters';
      setError(errorMessage);
      console.error('Error fetching clusters:', err);
      setClusters([]);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchClusters();
  }, [fetchClusters]);

  return {
    clusters,
    loading,
    error,
    fetchClusters,
  };
};
