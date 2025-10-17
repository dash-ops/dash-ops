import { useState, useEffect, useCallback } from 'react';
import { getNodes } from '../resources/nodesResource';
import { transformNodesToDomain } from '../adapters/nodeAdapter';
import { KubernetesTypes } from '../types';

export interface UseNodesResult {
  nodes: KubernetesTypes.Node[];
  loading: boolean;
  error: string | null;
  fetchNodes: () => Promise<void>;
}

export const useNodes = (context: string): UseNodesResult => {
  const [nodes, setNodes] = useState<KubernetesTypes.Node[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchNodes = useCallback(async () => {
    if (!context) {
      setNodes([]);
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await getNodes(context);
      const transformedNodes = transformNodesToDomain(response.data || []);
      setNodes(transformedNodes);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to fetch nodes';
      setError(errorMessage);
      console.error('Error fetching nodes:', err);
      setNodes([]);
    } finally {
      setLoading(false);
    }
  }, [context]);

  useEffect(() => {
    fetchNodes();
  }, [fetchNodes]);

  return {
    nodes,
    loading,
    error,
    fetchNodes,
  };
};
