import { useState, useEffect, useCallback } from 'react';
import { getPods } from '../resources/podsResource';
import { transformPodsToDomain } from '../adapters/podAdapter';
import { KubernetesTypes } from '../types';

export interface UsePodsResult {
  pods: KubernetesTypes.Pod[];
  loading: boolean;
  error: string | null;
  fetchPods: () => Promise<void>;
}

export const usePods = (filter: KubernetesTypes.PodFilter): UsePodsResult => {
  const [pods, setPods] = useState<KubernetesTypes.Pod[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchPods = useCallback(async () => {
    if (!filter.context || !filter.namespace) {
      setPods([]);
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await getPods(filter);
      const transformedPods = transformPodsToDomain(response.data || []);
      setPods(transformedPods);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to fetch pods';
      setError(errorMessage);
      console.error('Error fetching pods:', err);
      setPods([]);
    } finally {
      setLoading(false);
    }
  }, [filter.context, filter.namespace]);

  useEffect(() => {
    fetchPods();
  }, [fetchPods]);

  return {
    pods,
    loading,
    error,
    fetchPods,
  };
};
