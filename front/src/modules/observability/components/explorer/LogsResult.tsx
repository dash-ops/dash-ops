import { useState } from 'react';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Separator } from '@/components/ui/separator';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { 
  ChevronDown, 
  ChevronUp, 
  Copy, 
  ExternalLink,
  List,
  Table as TableIcon
} from 'lucide-react';
import type { LogEntry } from '../../types';

interface LogsResultProps {
  logs: LogEntry[];
  onNavigateToTrace?: (traceId: string) => void;
}

export default function LogsResult({ logs, onNavigateToTrace }: LogsResultProps): JSX.Element {
  const [viewMode, setViewMode] = useState<'table' | 'list'>('table');
  const [expandedLogs, setExpandedLogs] = useState<Set<string>>(new Set());

  const getLevelClass = (level: LogEntry['level']): string => {
    switch (level) {
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
    const next = new Set(expandedLogs);
    if (next.has(id)) next.delete(id);
    else next.add(id);
    setExpandedLogs(next);
  };

  return (
    <div className="flex flex-col h-full">
      {/* View Mode Toggle */}
      <div className="flex items-center justify-between mb-4">
        <div className="text-sm text-muted-foreground">
          {logs.length} log entries
        </div>
        <div className="flex border rounded-lg">
          <Button 
            variant={viewMode === 'table' ? 'default' : 'ghost'}
            size="sm"
            onClick={() => setViewMode('table')}
            className="rounded-none rounded-l-lg"
          >
            <TableIcon className="h-4 w-4" />
          </Button>
          <Button 
            variant={viewMode === 'list' ? 'default' : 'ghost'}
            size="sm"
            onClick={() => setViewMode('list')}
            className="rounded-none rounded-r-lg"
          >
            <List className="h-4 w-4" />
          </Button>
        </div>
      </div>

      {/* Logs Display */}
      <ScrollArea className="flex-1">
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
              {logs.map((log) => (
                <>
                  <TableRow key={log.id} className="hover:bg-muted/50 cursor-pointer" onClick={() => toggleExpand(log.id)}>
                    <TableCell>
                      {expandedLogs.has(log.id) ? <ChevronUp className="h-3 w-3" /> : <ChevronDown className="h-3 w-3" />}
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
                        {log.traceId && onNavigateToTrace && (
                          <Button variant="ghost" size="sm" className="h-6 w-6 p-0" onClick={(e) => {
                            e.stopPropagation();
                            onNavigateToTrace(String(log.traceId));
                          }}>
                            <ExternalLink className="h-3 w-3" />
                          </Button>
                        )}
                        <Button variant="ghost" size="sm" className="h-6 w-6 p-0" onClick={(e) => e.stopPropagation()}>
                          <Copy className="h-3 w-3" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                  {expandedLogs.has(log.id) && (
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
              {logs.length === 0 && (
                <TableRow>
                  <TableCell colSpan={6} className="text-center py-6 text-muted-foreground">No logs found</TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        ) : (
          <div className="space-y-2">
            {logs.map((log) => (
              <div
                key={log.id}
                className="border rounded-lg p-3 hover:border-primary transition-colors cursor-pointer"
                onClick={() => toggleExpand(log.id)}
              >
                <div className="flex items-start gap-3">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <Badge variant="secondary" className="text-xs">
                        {log.service}
                      </Badge>
                      <Badge variant="outline" className={getLevelClass(log.level)}>
                        {log.level.toUpperCase()}
                      </Badge>
                      <span className="text-xs text-muted-foreground font-mono">
                        {new Date(log.timestamp).toLocaleTimeString()}
                      </span>
                    </div>
                    <p className={expandedLogs.has(log.id) ? '' : 'truncate'}>{log.message}</p>
                    {expandedLogs.has(log.id) && log.metadata && (
                      <div className="mt-2 pt-2 border-t text-sm">
                        {Object.entries(log.metadata).map(([key, value]) => (
                          <div key={key} className="text-muted-foreground">
                            <span>{key}:</span> {String(value)}
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                </div>
              </div>
            ))}
            {logs.length === 0 && (
              <div className="text-center py-6 text-muted-foreground">No logs found</div>
            )}
          </div>
        )}
      </ScrollArea>
    </div>
  );
}

