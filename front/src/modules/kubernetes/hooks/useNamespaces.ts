import { useState, useEffect, useCallback } from 'react';
import { getNamespaces } from '../resources/namespaceResource';
import { transformNamespacesToDomain } from '../adapters/clusterAdapter';
import { KubernetesTypes } from '../types';

export interface UseNamespacesResult {
  namespaces: KubernetesTypes.Namespace[];
  loading: boolean;
  error: string | null;
  fetchNamespaces: () => Promise<void>;
}

export const useNamespaces = (context: string): UseNamespacesResult => {
  const [namespaces, setNamespaces] = useState<KubernetesTypes.Namespace[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchNamespaces = useCallback(async () => {
    if (!context) {
      setNamespaces([]);
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await getNamespaces(context);
      const transformedNamespaces = transformNamespacesToDomain(response.data || []);
      setNamespaces(transformedNamespaces);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to fetch namespaces';
      setError(errorMessage);
      console.error('Error fetching namespaces:', err);
      setNamespaces([]);
    } finally {
      setLoading(false);
    }
  }, [context]);

  useEffect(() => {
    fetchNamespaces();
  }, [fetchNamespaces]);

  return {
    namespaces,
    loading,
    error,
    fetchNamespaces,
  };
};
