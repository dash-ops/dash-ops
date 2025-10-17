import { useState, useEffect, useCallback } from 'react';
import { getDeployments } from '../resources/deploymentResource';
import { transformDeploymentsToDomain } from '../adapters/deploymentAdapter';
import { KubernetesTypes } from '../types';

export interface UseDeploymentsResult {
  deployments: KubernetesTypes.Deployment[];
  loading: boolean;
  error: string | null;
  fetchDeployments: () => Promise<void>;
}

export const useDeployments = (filter: KubernetesTypes.DeploymentFilter): UseDeploymentsResult => {
  const [deployments, setDeployments] = useState<KubernetesTypes.Deployment[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchDeployments = useCallback(async () => {
    if (!filter.context || !filter.namespace) {
      setDeployments([]);
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await getDeployments(filter);
      const transformedDeployments = transformDeploymentsToDomain(response.data || []);
      setDeployments(transformedDeployments);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to fetch deployments';
      setError(errorMessage);
      console.error('Error fetching deployments:', err);
      setDeployments([]);
    } finally {
      setLoading(false);
    }
  }, [filter.context, filter.namespace]);

  useEffect(() => {
    fetchDeployments();
  }, [fetchDeployments]);

  return {
    deployments,
    loading,
    error,
    fetchDeployments,
  };
};
