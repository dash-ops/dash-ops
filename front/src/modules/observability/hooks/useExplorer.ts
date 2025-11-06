import { useCallback, useState } from 'react';
import * as explorerResource from '../resources/explorerResource';
import * as explorerAdapter from '../adapters/explorerAdapter';
import type { LogEntry, TraceSpan } from '../types';

interface UseExplorerState {
  dataSource: 'logs' | 'traces' | 'metrics' | null;
  results: LogEntry[] | TraceSpan[] | any[];
  total: number;
  query: string;
  executionTimeMs: number;
  loading: boolean;
  error: string | null;
}

export const useExplorer = () => {
  const [state, setState] = useState<UseExplorerState>({
    dataSource: null,
    results: [],
    total: 0,
    query: '',
    executionTimeMs: 0,
    loading: false,
    error: null,
  });

  const executeQuery = useCallback(async (
    query: string,
    timeRange?: { from: Date; to: Date },
    provider?: string
  ) => {
    try {
      setState((s) => ({ ...s, loading: true, error: null }));
      
      const request = explorerAdapter.buildExplorerRequest(query, timeRange, provider);
      const resp = await explorerResource.executeQuery(request);
      const domain = explorerAdapter.transformExplorerResponseToDomain(resp.data);
      
      setState({
        dataSource: domain.dataSource,
        results: domain.results,
        total: domain.total,
        query: domain.query,
        executionTimeMs: domain.executionTimeMs,
        loading: false,
        error: null,
      });
      
      return domain;
    } catch (e: unknown) {
      const message = e instanceof Error ? e.message : 'Failed to execute query';
      setState((s) => ({ 
        ...s, 
        loading: false, 
        error: message,
        results: [],
        total: 0,
      }));
      throw e;
    }
  }, []);

  return {
    ...state,
    executeQuery,
  };
};

