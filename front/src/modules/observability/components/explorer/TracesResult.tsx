import { useState, useMemo, useEffect } from 'react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { ScrollArea } from '@/components/ui/scroll-area';
import { 
  ChevronRight,
  Clock,
  Search
} from 'lucide-react';
import { cn } from '@/lib/utils';
import type { TraceSpan } from '../../types';
import { format } from 'date-fns';
import * as tracesResource from '../../resources/tracesResource';
import * as tracesAdapter from '../../adapters/tracesAdapter';

interface TracesResultProps {
  traces: TraceSpan[];
  provider?: string;
}

const COLOR_PALETTE = [
  '#8884d8',
  '#82ca9d',
  '#ffc658',
  '#ff7300',
  '#0088fe',
  '#00c49f',
  '#ff6b6b',
  '#4ecdc4',
  '#45b7d1',
  '#f7b731',
  '#5f27cd',
  '#00d2d3',
  '#ff9ff3',
  '#54a0ff',
  '#5f27cd',
  '#c44569'
];

const getServiceColor = (service: string, colorMap: Map<string, string>): string => {
  const cachedColor = colorMap.get(service);
  if (cachedColor) {
    return cachedColor;
  }
  
  let hash = 0;
  for (let i = 0; i < service.length; i++) {
    hash = service.charCodeAt(i) + ((hash << 5) - hash);
  }
  
  const colorIndex = Math.abs(hash) % COLOR_PALETTE.length;
  const selectedColor = COLOR_PALETTE[colorIndex];
  const color: string = selectedColor || COLOR_PALETTE[0] || '#8884d8';
  colorMap.set(service, color);
  return color;
};

export default function TracesResult({ traces, provider }: TracesResultProps): JSX.Element {
  const [selectedTrace, setSelectedTrace] = useState<string | null>(null);
  const [expandedSpans, setExpandedSpans] = useState<Set<string>>(new Set());
  const [traceDetails, setTraceDetails] = useState<Map<string, TraceSpan[]>>(new Map());
  const [loadingTrace, setLoadingTrace] = useState<string | null>(null);

  // Transform traces using adapter to handle backend format
  const transformedTraces = useMemo(() => {
    return traces.map(span => {
      // If span is already in correct format, use it; otherwise transform it
      if (span.id && span.service && typeof span.startTime === 'number') {
        return span;
      }
      return tracesAdapter.transformTraceSpanToDomain(span);
    });
  }, [traces]);

  // Group traces by traceId and build trace info list
  const traceInfoList = useMemo(() => {
    const groups = new Map<string, TraceSpan[]>();
    transformedTraces.forEach(span => {
      if (!groups.has(span.traceId)) {
        groups.set(span.traceId, []);
      }
      groups.get(span.traceId)!.push(span);
    });

    return Array.from(groups.entries()).map(([traceId, spans]) => {
      const sortedSpans = spans.sort((a, b) => a.startTime - b.startTime);
      const rootSpan = sortedSpans.find(s => !s.parentId) || sortedSpans[0];
      if (!rootSpan || sortedSpans.length === 0) {
        // Fallback if no spans
        return {
          traceId,
          rootOperation: 'Unknown',
          status: 'ok' as const,
          errors: 0,
          totalDuration: 0,
          spanCount: 0,
          serviceCount: 0,
          timestamp: new Date(),
          spans: []
        };
      }
      const maxTime = Math.max(...sortedSpans.map(s => s.startTime + s.duration));
      const minTime = Math.min(...sortedSpans.map(s => s.startTime));
      const totalDuration = maxTime - minTime;
      const errorSpans = sortedSpans.filter(s => s.status === 'error');
      const uniqueServices = new Set(sortedSpans.map(s => s.service));

      // startTime is already converted to milliseconds by the adapter
      const timestamp = new Date(rootSpan.startTime);

      return {
        traceId,
        rootOperation: rootSpan.operationName,
        status: errorSpans.length > 0 ? 'error' as const : 'ok' as const,
        errors: errorSpans.length,
        totalDuration,
        spanCount: sortedSpans.length,
        serviceCount: uniqueServices.size,
        timestamp,
        spans: sortedSpans
      };
    });
  }, [traces]);

  // Fetch trace details when a trace is selected
  useEffect(() => {
    if (!selectedTrace) {
      return;
    }

    // Check if we already have details for this trace
    if (traceDetails.has(selectedTrace)) {
      return;
    }

    // Find fallback spans from initial query
    const traceInfo = traceInfoList.find(t => t.traceId === selectedTrace);
    const fallbackSpans = traceInfo?.spans || [];

    // Fetch trace details from backend
    const fetchTraceDetails = async () => {
      if (!provider) {
        console.error('Provider is required to fetch trace details');
        return;
      }
      setLoadingTrace(selectedTrace);
      try {
        const response = await tracesResource.getTraceDetail(selectedTrace, provider);
        const spans = (response.data.spans || []).map(tracesAdapter.transformTraceSpanToDomain);
        setTraceDetails(prev => {
          const newMap = new Map(prev);
          newMap.set(selectedTrace, spans);
          return newMap;
        });
      } catch (error) {
        console.error('Failed to fetch trace details:', error);
        // Fallback to using spans from the initial query if available
        if (fallbackSpans.length > 0) {
          setTraceDetails(prev => {
            const newMap = new Map(prev);
            newMap.set(selectedTrace, fallbackSpans);
            return newMap;
          });
        }
      } finally {
        setLoadingTrace(null);
      }
    };

    fetchTraceDetails();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [selectedTrace, traceInfoList]);

  const selectedTraceData = traceInfoList.find(t => t.traceId === selectedTrace);
  const filteredTraces = selectedTrace && traceDetails.has(selectedTrace)
    ? traceDetails.get(selectedTrace)!
    : (selectedTraceData?.spans || []);

  const serviceColorMap = useMemo(() => {
    const map = new Map<string, string>();
    const allServices = new Set<string>();
    
    transformedTraces.forEach(span => {
      if (span.service) {
        allServices.add(span.service);
      }
    });
    
    filteredTraces.forEach(span => {
      if (span.service) {
        allServices.add(span.service);
      }
    });
    
    allServices.forEach(service => {
      getServiceColor(service, map);
    });
    
    return map;
  }, [transformedTraces, filteredTraces]);

  const formatDuration = (durationMs: number) => {
    if (durationMs < 1000) {
      return `${Math.round(durationMs)}ms`;
    }
    return `${(durationMs / 1000).toFixed(2)}s`;
  };

  const formatTime = (date: Date) => {
    if (!date || isNaN(date.getTime())) {
      return 'Invalid date';
    }
    return format(date, 'MMM d, HH:mm:ss');
  };

  const toggleSpanExpansion = (spanId: string) => {
    const newExpanded = new Set(expandedSpans);
    if (newExpanded.has(spanId)) {
      newExpanded.delete(spanId);
    } else {
      newExpanded.add(spanId);
    }
    setExpandedSpans(newExpanded);
  };

  const getSpanDepth = (span: TraceSpan, allSpans: TraceSpan[]): number => {
    if (!span.parentId) return 0;
    const parent = allSpans.find(s => s.id === span.parentId);
    return parent ? getSpanDepth(parent, allSpans) + 1 : 0;
  };

  const timelineData = useMemo(() => {
    if (filteredTraces.length === 0) {
      return { sortedSpans: [], minTime: 0, maxTime: 0, totalDuration: 0 };
    }
    const sortedSpans = [...filteredTraces].sort((a, b) => a.startTime - b.startTime);
    const minTime = Math.min(...sortedSpans.map(s => s.startTime));
    const maxTime = Math.max(...sortedSpans.map(s => s.startTime + s.duration));
    const totalDuration = maxTime - minTime;
    return { sortedSpans, minTime, maxTime, totalDuration };
  }, [filteredTraces]);

  return (
    <div className="flex flex-col h-full">
      {selectedTrace ? (
        // Timeline View
        <div className="space-y-4">
          <div className="flex items-center gap-2">
            <Button 
              variant="ghost" 
              size="sm"
              onClick={() => setSelectedTrace(null)}
            >
              ← Back to List
            </Button>
            <span className="text-sm text-muted-foreground">
              Trace ID: <span className="font-mono">{selectedTrace}</span>
            </span>
          </div>

          {/* Services List */}
                      <div className="border rounded-lg p-4 bg-muted/30">
            <h4 className="text-sm font-medium mb-3">Services</h4>
            <div className="flex flex-wrap gap-3">
              {Array.from(new Set(filteredTraces.map(s => s.service))).map(service => (
                <div key={service} className="flex items-center gap-2">
                  <div 
                    className="w-3 h-3 rounded"
                    style={{ backgroundColor: getServiceColor(service, serviceColorMap) }}
                  />
                  <span className="text-sm">{service}</span>
                </div>
              ))}
            </div>
          </div>

          {/* Timeline Visualization */}
          <div className="space-y-2">
            <h4 className="text-sm font-medium">Timeline</h4>
            <div className="border rounded-lg p-4 bg-background">
              {loadingTrace === selectedTrace ? (
                <div className="text-center py-8 text-muted-foreground">
                  <div className="animate-spin h-6 w-6 border-2 border-primary border-t-transparent rounded-full mx-auto mb-2"></div>
                  Loading trace details...
                </div>
              ) : filteredTraces.length === 0 ? (
                <div className="text-center py-8 text-muted-foreground">
                  No spans found for this trace
                </div>
              ) : (
                timelineData.sortedSpans.map((span) => {
                  const leftPercent = timelineData.totalDuration > 0 
                    ? ((span.startTime - timelineData.minTime) / timelineData.totalDuration) * 100 
                    : 0;
                  const widthPercent = timelineData.totalDuration > 0 
                    ? (span.duration / timelineData.totalDuration) * 100 
                    : 0;
                  const depth = getSpanDepth(span, timelineData.sortedSpans);
                
                  return (
                    <div key={span.id} className="relative mb-3">
                      <div className="flex items-center gap-4 mb-1">
                        <div className="w-48 text-sm truncate">
                          <div style={{ paddingLeft: `${depth * 20}px` }}>
                            {span.operationName}
                          </div>
                        </div>
                        <div className="flex-1 relative h-6 bg-muted rounded">
                          <div
                            className="absolute h-full rounded transition-all hover:opacity-80 cursor-pointer"
                            style={{
                              left: `${leftPercent}%`,
                              width: `${widthPercent}%`,
                              backgroundColor: getServiceColor(span.service, serviceColorMap),
                              opacity: span.status === 'error' ? 0.8 : 0.6
                            }}
                            onClick={() => toggleSpanExpansion(span.id)}
                          />
                          {span.status === 'error' && (
                            <div
                              className="absolute top-0 bottom-0 w-1 bg-red-500"
                              style={{ left: `${leftPercent}%` }}
                            />
                          )}
                        </div>
                        <div className="w-20 text-sm text-right font-mono">
                          {formatDuration(span.duration)}
                        </div>
                      </div>
                      
                      {expandedSpans.has(span.id) && (
                        <div className="ml-52 p-3 bg-muted/30 rounded border">
                          <div className="grid grid-cols-2 gap-4 text-sm">
                            <div>
                              <div><span className="text-muted-foreground">Service:</span> {span.service}</div>
                              <div><span className="text-muted-foreground">Operation:</span> {span.operationName}</div>
                              <div><span className="text-muted-foreground">Duration:</span> {formatDuration(span.duration)}</div>
                              <div><span className="text-muted-foreground">Status:</span> 
                                <Badge variant={span.status === 'ok' ? 'secondary' : 'destructive'} className="ml-2">
                                  {span.status.toUpperCase()}
                                </Badge>
                              </div>
                            </div>
                            <div>
                              {span.tags && Object.keys(span.tags).length > 0 ? (
                                <>
                                  <div className="text-muted-foreground mb-1">Tags:</div>
                                  <div className="flex flex-wrap gap-1">
                                    {Object.entries(span.tags).map(([key, value]) => (
                                      <Badge key={key} variant="outline" className="text-xs">
                                        {key}: {String(value)}
                                      </Badge>
                                    ))}
                                  </div>
                                </>
                              ) : (
                                <div className="text-muted-foreground">No tags</div>
                              )}
                            </div>
                          </div>
                        </div>
                      )}
                    </div>
                  );
                })
              )}
            </div>
          </div>
        </div>
      ) : (
        // Trace List View
        <div className="space-y-4">
          <div className="text-sm text-muted-foreground">
            {traceInfoList.length} traces found
          </div>

          {traceInfoList.length === 0 ? (
            <div className="text-center py-12">
              <Search className="h-12 w-12 mx-auto mb-4 text-muted-foreground opacity-50" />
              <h3 className="mb-2">No traces found</h3>
              <p className="text-muted-foreground">
                Try adjusting your query or time range
              </p>
            </div>
          ) : (
            <ScrollArea className="flex-1">
              <div className="space-y-2">
                {traceInfoList.map((trace) => (
                  <div
                    key={trace.traceId}
                    className="p-4 border rounded-lg cursor-pointer hover:bg-muted/50 transition-colors"
                    onClick={() => setSelectedTrace(trace.traceId)}
                  >
                    <div className="flex items-center justify-between mb-3">
                      <div className="flex items-center gap-3">
                        <div className={cn("w-3 h-3 rounded-full", 
                          trace.status === 'ok' ? 'bg-green-500' : 'bg-red-500'
                        )} />
                        <span className="font-medium">{trace.rootOperation}</span>
                        <Badge variant={trace.status === 'error' ? 'destructive' : 'secondary'}>
                          {trace.status === 'error' ? 'Error' : 'Success'}
                        </Badge>
                        {trace.errors > 0 && (
                          <Badge variant="outline" className="text-red-600">
                            {trace.errors} errors
                          </Badge>
                        )}
                      </div>
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <Clock className="h-4 w-4" />
                        {formatDuration(trace.totalDuration)}
                      </div>
                    </div>
                    
                    <div className="flex items-center justify-between text-sm">
                      <div className="flex items-center gap-4">
                        <span className="text-muted-foreground">
                          Trace ID: <span className="font-mono">{trace.traceId}</span>
                        </span>
                        <span className="text-muted-foreground">
                          {trace.spanCount} spans
                        </span>
                        <span className="text-muted-foreground">
                          {trace.serviceCount} services
                        </span>
                        <span className="text-muted-foreground">
                          {formatTime(trace.timestamp)}
                        </span>
                      </div>
                      <Button variant="ghost" size="sm" className="gap-1 h-7">
                        View Timeline
                        <ChevronRight className="h-3 w-3" />
                      </Button>
                    </div>
                  </div>
                ))}
              </div>
            </ScrollArea>
          )}
        </div>
      )}
    </div>
  );
}

