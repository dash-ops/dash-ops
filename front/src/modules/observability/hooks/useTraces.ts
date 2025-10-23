import { useCallback, useEffect, useMemo, useState } from 'react';
import * as tracesResource from '../resources/tracesResource';
import * as tracesAdapter from '../adapters/tracesAdapter';
import type { TracesQueryFilters, TracesResponse, TraceSpan } from '../types';

interface UseTracesState {
  data: TracesResponse;
  spans: TraceSpan[];
  loading: boolean;
  error: string | null;
}

export const useTraces = (initialFilters: TracesQueryFilters) => {
  const [filters, setFilters] = useState<TracesQueryFilters>(initialFilters);
  const [state, setState] = useState<UseTracesState>({
    data: { traces: [], total: 0, page: 1, pageSize: 0 },
    spans: [],
    loading: false,
    error: null,
  });

  const fetchTraces = useCallback(async () => {
    try {
      setState((s) => ({ ...s, loading: true, error: null }));
      const resp = await tracesResource.getTraces(filters);
      const traces = (resp.data.traces || []).map(tracesAdapter.transformTraceInfoToDomain);
      setState({ data: { ...resp.data, traces }, spans: [], loading: false, error: null });
    } catch (e: unknown) {
      const message = e instanceof Error ? e.message : 'Failed to fetch traces';
      setState((s) => ({ ...s, loading: false, error: message }));
    }
  }, [filters]);

  const fetchTraceSpans = useCallback(async (traceId: string) => {
    try {
      setState((s) => ({ ...s, loading: true, error: null }));
      const resp = await tracesResource.getTraceSpans(traceId);
      const spans = tracesAdapter.transformTraceSpansToDomain(resp.data);
      setState((s) => ({ ...s, spans, loading: false }));
    } catch (e: unknown) {
      const message = e instanceof Error ? e.message : 'Failed to fetch trace spans';
      setState((s) => ({ ...s, loading: false, error: message }));
    }
  }, []);

  const refresh = useCallback(async () => {
    await fetchTraces();
  }, [fetchTraces]);

  const updateFilters = useCallback((partial: Partial<TracesQueryFilters>) => {
    setFilters((f) => ({ ...f, ...partial }));
  }, []);

  useEffect(() => {
    fetchTraces();
  }, [fetchTraces]);

  const statuses = useMemo(() => ['all', 'ok', 'error'] as const, []);

  return {
    ...state,
    filters,
    fetchTraces,
    fetchTraceSpans,
    refresh,
    updateFilters,
    statuses,
  };
};


