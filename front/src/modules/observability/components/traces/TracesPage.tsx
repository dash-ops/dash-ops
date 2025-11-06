import { useMemo, useState, useEffect } from 'react';
import { useSearchParams } from 'react-router';
import { useTraces } from '../../hooks/useTraces';
import type { TracesQueryFilters } from '../../types';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { useServices } from '../../hooks/useServices';
import { Button } from '@/components/ui/button';
import Refresh from '@/components/Refresh';

export default function TracesPage(): JSX.Element {
  const [searchParams, setSearchParams] = useSearchParams();
  const [search, setSearch] = useState('');
  const [service, setService] = useState('all');
  const [status, setStatus] = useState<'all' | 'ok' | 'error'>('all');
  const [provider, setProvider] = useState('tempo-local');

  // Check if URL has parameters (shared link)
  const hasUrlParams = searchParams.toString().length > 0;

  const initialFilters: TracesQueryFilters = useMemo(
    () => ({ status: 'all', limit: 50, page: 1, provider: 'tempo-local' }),
    []
  );

  const { data, spans, loading, error, filters, updateFilters, statuses, fetchTraceSpans, refresh, fetchTraces } = useTraces(initialFilters, hasUrlParams);
  
  // Fetch services from Service Catalog for dropdown
  const { data: servicesData, loading: servicesLoading } = useServices({ limit: 100 });

  // Load filters from URL on mount if present
  useEffect(() => {
    if (hasUrlParams) {
      const searchParam = searchParams.get('search');
      const serviceParam = searchParams.get('service');
      const statusParam = searchParams.get('status') as 'all' | 'ok' | 'error' | null;
      const providerParam = searchParams.get('provider');
      
      if (searchParam) setSearch(searchParam);
      if (serviceParam) setService(serviceParam);
      if (statusParam) setStatus(statusParam);
      if (providerParam) setProvider(providerParam);
      
      // Apply filters from URL
      const partial: Partial<TracesQueryFilters> = {
        status: statusParam || 'all',
        provider: providerParam || 'tempo-local',
        page: 1,
      };
      if (serviceParam && serviceParam !== 'all') partial.service = serviceParam;
      if (searchParam) partial.search = searchParam;
      updateFilters(partial);
    }
  }, []); // Only run on mount

  const onApplyFilters = () => {
    const partial: Partial<TracesQueryFilters> = {
      status,
      provider,
      page: 1,
    };
    if (service !== 'all') partial.service = service;
    if (search) partial.search = search;
    
    // Update URL with current filters
    const params = new URLSearchParams();
    if (search) {
      params.set('search', search);
    }
    if (service && service !== 'all') {
      params.set('service', service);
    }
    if (status && status !== 'all') {
      params.set('status', status);
    }
    if (provider) {
      params.set('provider', provider);
    }
    setSearchParams(params);
    
    const newFilters = { ...initialFilters, ...filters, ...partial };
    updateFilters(partial);
    // Manually fetch with new filters
    fetchTraces(newFilters);
  };

  return (
    <div className="flex-1 overflow-hidden p-6">
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle>Traces</CardTitle>
            <Refresh onReload={refresh} />
          </div>
        </CardHeader>
        <CardContent>
          <div className="flex gap-2 mb-4">
            <Input
              placeholder="Search (traceId, operation, tags)"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
            <Select value={service} onValueChange={setService} disabled={servicesLoading}>
              <SelectTrigger className="w-56">
                <SelectValue placeholder={servicesLoading ? "Loading services..." : "All Services"} />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Services</SelectItem>
                {servicesData?.services.map((svc) => (
                  <SelectItem key={svc.service_name} value={svc.service_name}>
                    {svc.service_name}
                    {svc.namespace && ` (${svc.namespace})`}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Select value={provider} onValueChange={setProvider}>
              <SelectTrigger className="w-40">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="tempo-local">Tempo Local</SelectItem>
              </SelectContent>
            </Select>
            <Select value={status} onValueChange={(v) => setStatus(v as any)}>
              <SelectTrigger className="w-32">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {statuses.map((s) => (
                  <SelectItem key={s} value={s}>
                    {s}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Button onClick={onApplyFilters}>Apply</Button>
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


