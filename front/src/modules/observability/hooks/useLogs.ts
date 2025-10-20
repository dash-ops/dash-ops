import { useCallback, useEffect, useState } from 'react';
import * as logsResource from '../resources/logsResource';
import * as logsAdapter from '../adapters/logsAdapter';
import type { LogEntry, LogsQueryFilters, LogsResponse } from '../types';

interface UseLogsState {
  data: LogsResponse;
  loading: boolean;
  error: string | null;
}

export const useLogs = (initialFilters: LogsQueryFilters) => {
  const [filters, setFilters] = useState<LogsQueryFilters>(initialFilters);
  const [state, setState] = useState<UseLogsState>({
    data: { items: [], total: 0, page: 1, pageSize: 0 },
    loading: false,
    error: null,
  });

  const fetchLogs = useCallback(async () => {
    try {
      setState((s) => ({ ...s, loading: true, error: null }));
      const resp = await logsResource.getLogs(filters);
      const domain = logsAdapter.transformLogsResponseToDomain(resp.data);
      setState({ data: domain, loading: false, error: null });
    } catch (e: unknown) {
      const message = e instanceof Error ? e.message : 'Failed to fetch logs';
      setState((s) => ({ ...s, loading: false, error: message }));
    }
  }, [filters]);

  const refresh = useCallback(async () => {
    await fetchLogs();
  }, [fetchLogs]);

  const updateFilters = useCallback((partial: Partial<LogsQueryFilters>) => {
    setFilters((f) => ({ ...f, ...partial }));
  }, []);

  useEffect(() => {
    fetchLogs();
  }, [fetchLogs]);

  const getLevels = useCallback((): Array<LogEntry['level'] | 'all'> => {
    return ['all', 'error', 'warn', 'info', 'debug'];
  }, []);

  return {
    ...state,
    filters,
    fetchLogs,
    refresh,
    updateFilters,
    getLevels,
  };
};


