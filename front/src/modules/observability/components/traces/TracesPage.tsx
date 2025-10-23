import { useMemo, useState } from 'react';
import { useTraces } from '../../hooks/useTraces';
import type { TracesQueryFilters } from '../../types';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';

export default function TracesPage(): JSX.Element {
  const [search, setSearch] = useState('');
  const [service, setService] = useState('all');
  const [status, setStatus] = useState<'all' | 'ok' | 'error'>('all');
  const [provider, setProvider] = useState('tempo-local');

  const initialFilters: TracesQueryFilters = useMemo(
    () => ({ status: 'all', limit: 50, page: 1, provider: 'tempo-local' }),
    []
  );

  const { data, spans, loading, error, updateFilters, statuses, fetchTraceSpans } = useTraces(initialFilters);

  const onApplyFilters = () => {
    const partial: Partial<TracesQueryFilters> = {
      status,
      provider,
      page: 1,
    };
    if (service !== 'all') partial.service = service;
    if (search) partial.search = search;
    updateFilters(partial);
  };

  return (
    <div className="flex-1 overflow-hidden p-6">
      <Card>
        <CardHeader>
          <CardTitle>Traces</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex gap-2 mb-4">
            <Input
              placeholder="Search (traceId, operation, tags)"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
            <Input
              placeholder="Service (optional)"
              value={service}
              onChange={(e) => setService(e.target.value)}
            />
            <select
              className="border rounded px-2"
              value={provider}
              onChange={(e) => setProvider(e.target.value)}
            >
              <option value="tempo-local">Tempo Local</option>
            </select>
            <select
              className="border rounded px-2"
              value={status}
              onChange={(e) => setStatus(e.target.value as any)}
            >
              {statuses.map((s) => (
                <option key={s} value={s}>
                  {s}
                </option>
              ))}
            </select>
            <button className="btn btn-primary" onClick={onApplyFilters}>
              Apply
            </button>
          </div>

          {error && (
            <div className="mb-2 text-red-600 text-sm">{error}</div>
          )}

          <div className="grid md:grid-cols-2 gap-4">
            <div className="border rounded p-2">
              <div className="font-medium mb-2">Trace List</div>
              <div className="space-y-2">
                {data.traces?.map((trace) => (
                  <div
                    key={trace.traceId}
                    className="p-2 border rounded cursor-pointer hover:bg-muted/50"
                    onClick={() => fetchTraceSpans(trace.traceId)}
                  >
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <Badge variant={trace.errors > 0 ? 'destructive' : 'secondary'}>
                          {trace.errors > 0 ? 'ERROR' : 'OK'}
                        </Badge>
                        <span className="font-mono text-xs">{trace.traceId}</span>
                      </div>
                      <div className="text-sm text-muted-foreground">
                        {Math.round(trace.totalDuration / 1000000)} ms
                      </div>
                    </div>
                    <div className="text-sm text-muted-foreground">{trace.rootOperation}</div>
                  </div>
                ))}
                {!loading && (!data.traces || data.traces.length === 0) && (
                  <div className="text-center py-6 text-muted-foreground">No traces</div>
                )}
              </div>
            </div>

            <div className="border rounded p-2">
              <div className="font-medium mb-2">Trace Spans</div>
              <div className="space-y-2">
                {spans.map((span) => (
                  <div key={span.id} className="p-2 border rounded">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center gap-2">
                        <Badge variant={span.status === 'error' ? 'destructive' : 'secondary'}>
                          {span.status.toUpperCase()}
                        </Badge>
                        <span className="font-mono text-xs">{span.operationName}</span>
                      </div>
                      <div className="text-sm text-muted-foreground">
                        {Math.round(span.duration)} ms
                      </div>
                    </div>
                    <div className="text-xs text-muted-foreground">{span.service}</div>
                  </div>
                ))}
                {!loading && spans.length === 0 && (
                  <div className="text-center py-6 text-muted-foreground">No spans</div>
                )}
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}


