import { useMemo, useState, useEffect } from 'react';
import { useSearchParams } from 'react-router';
import { useLogs } from '../../hooks/useLogs';
import type { LogEntry, LogsQueryFilters } from '../../types';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Button } from '@/components/ui/button';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Separator } from '@/components/ui/separator';
import Refresh from '@/components/Refresh';
import { Database, BarChart3, Search, List, Table as TableIcon, ExternalLink, Copy, ChevronDown, ChevronUp } from 'lucide-react';

export default function LogsPage(): JSX.Element {
  const [searchParams, setSearchParams] = useSearchParams();
  const [logqlQuery, setLogqlQuery] = useState('');
  const [provider, setProvider] = useState('loki-local');
  const [viewMode, setViewMode] = useState<'table' | 'list'>('table');
  const [expanded, setExpanded] = useState<Set<string>>(new Set());

  // Check if URL has parameters (shared link)
  const hasUrlParams = searchParams.toString().length > 0;

  const initialFilters: LogsQueryFilters = useMemo(
    () => ({ level: 'all', limit: 50, page: 1, provider: 'loki-local' }),
    []
  );

  const { data, loading, error, filters, updateFilters, refresh, fetchLogs } = useLogs(initialFilters, hasUrlParams);

  // Load filters from URL on mount if present
  useEffect(() => {
    if (hasUrlParams) {
      const queryParam = searchParams.get('query');
      const providerParam = searchParams.get('provider');
      
      if (queryParam) setLogqlQuery(queryParam);
      if (providerParam) setProvider(providerParam);
      
      // Apply filters from URL
      const partial: Partial<LogsQueryFilters> = {
        page: 1,
        provider: providerParam || 'loki-local',
      };
      if (queryParam) {
        partial.query = queryParam;
      }
      updateFilters(partial);
    }
  }, []); // Only run on mount

  const onApplyFilters = () => {
    const partial: Partial<LogsQueryFilters> = {
      page: 1,
      provider,
    };
    // Use LogQL query if provided, otherwise let backend build from filters
    if (logqlQuery.trim()) {
      partial.query = logqlQuery;
    }
    // Update URL with current filters
    const params = new URLSearchParams();
    if (logqlQuery.trim()) {
      params.set('query', logqlQuery);
    }
    if (provider) {
      params.set('provider', provider);
    }
    setSearchParams(params);
    
    // Build complete filters and fetch directly
    const newFilters = { ...initialFilters, ...filters, ...partial };
    updateFilters(partial);
    // Manually fetch with new filters since autoFetch is false
    fetchLogs(newFilters);
  };

  const histogram = useMemo(() => {
    const buckets = new Map<string, { count: number; errors: number; warnings: number }>();
    data.items.forEach((log) => {
      const hour = new Date(log.timestamp).toISOString().slice(0, 13) + ':00:00.000Z';
      const entry = buckets.get(hour) || { count: 0, errors: 0, warnings: 0 };
      entry.count += 1;
      if (log.level === 'error') entry.errors += 1;
      if (log.level === 'warn') entry.warnings += 1;
      buckets.set(hour, entry);
    });
    return Array.from(buckets.entries())
      .map(([timestamp, v]) => ({ timestamp, ...v }))
      .sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime());
  }, [data.items]);

  const getLevelClass = (l: LogEntry['level']): string => {
    switch (l) {
      case 'error':
        return 'text-red-600 bg-red-50 border-red-200';
      case 'warn':
        return 'text-yellow-600 bg-yellow-50 border-yellow-200';
      case 'info':
        return 'text-blue-600 bg-blue-50 border-blue-200';
      case 'debug':
        return 'text-gray-600 bg-gray-50 border-gray-200';
      default:
        return 'text-gray-600 bg-gray-50 border-gray-200';
    }
  };

  const toggleExpand = (id: string) => {
    const next = new Set(expanded);
    if (next.has(id)) next.delete(id);
    else next.add(id);
    setExpanded(next);
  };

  return (
    <div className="flex-1 overflow-hidden p-6">
      <div className="grid grid-cols-1 gap-6">
        {/* Logs */}
        <Card className="flex-1 overflow-hidden">
          <CardHeader className="pb-4">
            <div className="flex items-center justify-between">
              <div>
                <CardTitle className="flex items-center gap-2">
                  <Database className="h-5 w-5" />
                  Application Logs
                </CardTitle>
              </div>
              <div className="flex items-center gap-2">
                <div className="flex border rounded-lg">
                  <Button variant={viewMode === 'table' ? 'default' : 'ghost'} size="sm" onClick={() => setViewMode('table')} className="rounded-none rounded-l-lg">
                    <TableIcon className="h-4 w-4" />
                  </Button>
                  <Button variant={viewMode === 'list' ? 'default' : 'ghost'} size="sm" onClick={() => setViewMode('list')} className="rounded-none rounded-r-lg">
                    <List className="h-4 w-4" />
                  </Button>
                </div>
                <Refresh onReload={refresh} />
              </div>
            </div>

            {/* Filters */}
            <div className="flex gap-2 mt-4">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground h-4 w-4" />
                <Input 
                  placeholder='LogQL query (e.g., {app="dashops-playground"} |= "error" or {cluster="kind-dashops-dev"})'
                  value={logqlQuery} 
                  onChange={(e) => setLogqlQuery(e.target.value)} 
                  className="pl-10 font-mono text-sm" 
                />
              </div>
              <Select value={provider} onValueChange={setProvider}>
                <SelectTrigger className="w-40">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="loki-local">Loki Local</SelectItem>
                </SelectContent>
              </Select>
              <Button onClick={onApplyFilters}>Apply</Button>
            </div>
          </CardHeader>
          
          {/* Histogram - Moved below filters */}
          <div className="px-6 pb-4">
            <div className="flex items-center gap-2 mb-3">
              <CardTitle className="flex items-center gap-2 text-sm">
                <BarChart3 className="h-4 w-4" />
                Log Distribution
              </CardTitle>
              <Badge variant="outline">{data.items.length} logs</Badge>
            </div>
            {/* lightweight bars without external chart lib */}
            <div className="flex items-end gap-1 h-24">
              {histogram.map((b) => (
                <div key={b.timestamp} className="flex-1 flex items-end gap-0.5">
                  <div title={`Total ${b.count}`} className="w-full bg-indigo-400/70" style={{ height: Math.min(96, b.count * 4) }} />
                </div>
              ))}
              {histogram.length === 0 && (
                <div className="text-sm text-muted-foreground">No data</div>
              )}
            </div>
          </div>

          <CardContent className="p-0">
            {error && (
              <div className="px-4 pt-4 text-sm text-red-600">{error}</div>
            )}
            <ScrollArea className="h-[calc(100vh-380px)]">
              {viewMode === 'table' ? (
                <Table>
                  <TableHeader className="sticky top-0 bg-background z-10">
                    <TableRow>
                      <TableHead className="w-10" />
                      <TableHead className="w-40">Date</TableHead>
                      <TableHead className="w-48">Service</TableHead>
                      <TableHead className="w-24">Level</TableHead>
                      <TableHead>Message</TableHead>
                      <TableHead className="w-16" />
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {data.items.map((log) => (
                      <>
                        <TableRow key={log.id} className="hover:bg-muted/50 cursor-pointer" onClick={() => toggleExpand(log.id)}>
                          <TableCell>
                            {expanded.has(log.id) ? <ChevronUp className="h-3 w-3" /> : <ChevronDown className="h-3 w-3" />}
                          </TableCell>
                          <TableCell className="font-mono text-xs">{new Date(log.timestamp).toLocaleString()}</TableCell>
                          <TableCell>
                            <Badge variant="secondary" className="text-xs">{log.service}</Badge>
                          </TableCell>
                          <TableCell>
                            <Badge variant="outline" className={getLevelClass(log.level)}>
                              {log.level.toUpperCase()}
                            </Badge>
                          </TableCell>
                          <TableCell className="truncate max-w-[70ch]">
                            <div className="max-w-[70ch]">
                              <p className="truncate">{log.message}</p>
                              <div className="flex items-center gap-2 mt-1 text-xs text-muted-foreground">
                                {log.source && <span className="font-mono">{log.source}</span>}
                                {log.traceId && <Badge variant="outline" className="text-xs font-mono">{String(log.traceId).slice(0, 8)}...</Badge>}
                              </div>
                            </div>
                          </TableCell>
                          <TableCell>
                            <div className="flex items-center gap-1">
                              {log.traceId && (
                                <Button variant="ghost" size="sm" className="h-6 w-6 p-0" onClick={(e) => e.stopPropagation()}>
                                  <ExternalLink className="h-3 w-3" />
                                </Button>
                              )}
                              <Button variant="ghost" size="sm" className="h-6 w-6 p-0" onClick={(e) => e.stopPropagation()}>
                                <Copy className="h-3 w-3" />
                              </Button>
                            </div>
                          </TableCell>
                        </TableRow>
                        {expanded.has(log.id) && (
                          <TableRow>
                            <TableCell colSpan={6} className="bg-muted/30">
                              <div className="p-4">
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4 text-sm">
                                  <div>
                                    <div><span className="text-muted-foreground">Timestamp:</span> {new Date(log.timestamp).toISOString()}</div>
                                    <div><span className="text-muted-foreground">Host:</span> {log.host ?? '-'}</div>
                                    <div><span className="text-muted-foreground">Source:</span> {log.source ?? '-'}</div>
                                    {log.traceId && <div><span className="text-muted-foreground">Trace ID:</span> {String(log.traceId)}</div>}
                                    {log.spanId && <div><span className="text-muted-foreground">Span ID:</span> {String(log.spanId)}</div>}
                                  </div>
                                  {log.metadata && (
                                    <div>
                                      <div className="text-muted-foreground mb-1">Metadata</div>
                                      <div className="space-y-1">
                                        {Object.entries(log.metadata).map(([k, v]) => (
                                          <div key={k}><span className="text-muted-foreground">{k}:</span> {String(v)}</div>
                                        ))}
                                      </div>
                                    </div>
                                  )}
                                </div>
                                <Separator className="my-2" />
                                <div className="bg-background border rounded p-3 text-sm font-mono">
                                  {log.message}
                                </div>
                              </div>
                            </TableCell>
                          </TableRow>
                        )}
                      </>
                    ))}
                    {!loading && data.items.length === 0 && (
                      <TableRow>
                        <TableCell colSpan={6} className="text-center py-6 text-muted-foreground">No logs</TableCell>
                      </TableRow>
                    )}
                  </TableBody>
                </Table>
              ) : (
                <div className="space-y-2 p-4">
                  {data.items.map((log) => (
                    <div key={log.id} className="p-4 rounded-lg border hover:bg-muted/50 transition-colors">
                      <div className="flex items-start justify-between gap-4 mb-3">
                        <div className="flex items-center gap-3">
                          <Badge variant="outline" className={getLevelClass(log.level)}>{log.level.toUpperCase()}</Badge>
                          <span className="text-sm font-mono text-muted-foreground">{new Date(log.timestamp).toLocaleTimeString()}</span>
                          <Badge variant="secondary" className="text-xs">{log.service}</Badge>
                          {log.host && <Badge variant="outline" className="text-xs font-mono">{log.host}</Badge>}
                        </div>
                        <div className="flex items-center gap-1">
                          {log.traceId && (
                            <Button variant="ghost" size="sm" className="h-6 w-6 p-0">
                              <ExternalLink className="h-3 w-3" />
                            </Button>
                          )}
                          <Button variant="ghost" size="sm" className="h-6 w-6 p-0">
                            <Copy className="h-3 w-3" />
                          </Button>
                        </div>
                      </div>
                      <p className="mb-3 leading-relaxed">{log.message}</p>
                      <div className="flex items-center justify-between text-sm">
                        {log.source && <span className="text-muted-foreground font-mono">{log.source}</span>}
                        <div className="flex items-center gap-2">
                          {log.traceId && (
                            <Badge variant="outline" className="text-xs font-mono">{String(log.traceId)}</Badge>
                          )}
                          {log.metadata && Object.keys(log.metadata).length > 0 && (
                            <Badge variant="secondary" className="text-xs">+{Object.keys(log.metadata).length} fields</Badge>
                          )}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </ScrollArea>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}


