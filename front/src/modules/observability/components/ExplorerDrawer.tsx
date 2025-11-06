import { useState, useMemo, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Textarea } from '@/components/ui/textarea';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { 
  X, 
  Plus,
  Play,
  ChevronDown,
  Database,
  GitBranch,
  BarChart3,
  Sparkles
} from 'lucide-react';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { cn } from '@/lib/utils';
import TimeRangePicker, { type TimeRange } from '@/components/TimeRangePicker';
import { useExplorer } from '../hooks/useExplorer';
import { useProviders } from '../hooks/useProviders';
import type { QueryTab, DataSource } from '../types/explorer';
import type { LogEntry, TraceSpan } from '../types';
import LogsResult from './explorer/LogsResult';
import TracesResult from './explorer/TracesResult';
import MetricsResult from './explorer/MetricsResult';

interface ExplorerDrawerProps {
  isOpen: boolean;
  onClose: () => void;
}

export default function ExplorerDrawer({ isOpen, onClose }: ExplorerDrawerProps): JSX.Element | null {
  const { executeQuery: apiExecuteQuery, loading, error } = useExplorer();
  const [tabs, setTabs] = useState<QueryTab[]>([
    {
      id: 'tab-1',
      title: 'Query 1',
      query: '',
      dataSource: null,
      results: null,
      isActive: false,
      timestamp: Date.now()
    }
  ]);
  const [activeTabId, setActiveTabId] = useState('tab-1');
  const [provider, setProvider] = useState<string>('');
  const [timeRange, setTimeRange] = useState<TimeRange>({ 
    value: '1h', 
    label: 'Last hour',
    from: new Date(Date.now() - 60 * 60 * 1000),
    to: new Date(),
  });

  const activeTab = tabs.find(t => t.id === activeTabId);

  // Fetch providers using hook with cache and request deduplication
  const { providers } = useProviders();

  // Detect data source from query and return available providers
  const availableProviders = useMemo(() => {
    if (!activeTab?.query) {
      return [];
    }
    const queryUpper = activeTab.query.toUpperCase();
    if (queryUpper.includes('FROM LOGS') || queryUpper.includes('FROM LOG')) {
      return providers.logs;
    } else if (queryUpper.includes('FROM TRACES') || queryUpper.includes('FROM TRACE')) {
      return providers.traces;
    } else if (queryUpper.includes('FROM METRICS') || queryUpper.includes('FROM METRIC')) {
      return providers.metrics;
    }
    // If query starts with { or contains =", it's likely LogQL (logs)
    if (activeTab.query.trim().startsWith('{') || activeTab.query.includes('="')) {
      return providers.logs;
    }
    return [];
  }, [activeTab?.query, providers]);

  // Auto-select first provider when data source changes
  useEffect(() => {
    if (availableProviders.length > 0 && !availableProviders.includes(provider)) {
      const firstProvider = availableProviders[0];
      if (firstProvider) {
        setProvider(firstProvider);
      }
    }
  }, [availableProviders, provider]);

  const createNewTab = () => {
    const newTab: QueryTab = {
      id: `tab-${Date.now()}`,
      title: `Query ${tabs.length + 1}`,
      query: '',
      dataSource: null,
      results: null,
      isActive: false,
      timestamp: Date.now()
    };
    setTabs([...tabs, newTab]);
    setActiveTabId(newTab.id);
  };

  const closeTab = (tabId: string) => {
    const newTabs = tabs.filter(t => t.id !== tabId);
    if (newTabs.length === 0) {
      onClose();
      return;
    }
    setTabs(newTabs);
    if (activeTabId === tabId && newTabs.length > 0) {
      setActiveTabId(newTabs[0]!.id);
    }
  };

  const updateTabQuery = (tabId: string, query: string) => {
    setTabs(tabs.map(t => t.id === tabId ? { ...t, query } : t));
  };

  const executeQuery = async (tabId: string) => {
    const tab = tabs.find(t => t.id === tabId);
    if (!tab || !tab.query) return;

    if (!provider) {
      console.error('Provider is required to execute query');
      return;
    }

    try {
      // Call backend API with query, timeRange, and provider
      const result = await apiExecuteQuery(
        tab.query,
        timeRange.from && timeRange.to ? { from: timeRange.from, to: timeRange.to } : undefined,
        provider
      );

      // Update tab with results
      const dataSourceMap: Record<string, DataSource> = {
        'logs': 'Logs',
        'traces': 'Traces',
        'metrics': 'Metrics',
      };

      setTabs(tabs.map(t => 
        t.id === tabId 
          ? { 
              ...t, 
              dataSource: dataSourceMap[result.dataSource] || null, 
              results: result.results, 
              isActive: true 
            }
          : t
      ));
    } catch (err) {
      console.error('Failed to execute query:', err);
      // Keep tab in error state with empty results
      setTabs(tabs.map(t => 
        t.id === tabId 
          ? { ...t, results: [], isActive: false }
          : t
      ));
    }
  };

  const insertDataSource = (tabId: string, source: DataSource) => {
    const tab = tabs.find(t => t.id === tabId);
    if (!tab) return;
    
    const newQuery = `FROM ${source} `;
    updateTabQuery(tabId, tab.query + newQuery);
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-x-0 bottom-0 z-50 bg-background border-t shadow-2xl flex flex-col"
         style={{ height: '70vh' }}>
      {/* Tabs Header */}
      <div className="flex items-center gap-2 px-4 py-2 border-b bg-muted/30">
        <div className="flex items-center gap-1 flex-1 overflow-x-auto">
          {tabs.map((tab) => (
            <div
              key={tab.id}
              className={cn(
                "flex items-center gap-2 px-3 py-1.5 rounded-t border border-b-0 cursor-pointer transition-colors min-w-fit",
                activeTabId === tab.id 
                  ? "bg-background border-border" 
                  : "bg-muted/50 hover:bg-muted border-transparent"
              )}
              onClick={() => setActiveTabId(tab.id)}
            >
              <span className="text-sm">{tab.title}</span>
              {tab.dataSource && (
                <Badge variant="outline" className="text-xs">
                  {tab.dataSource}
                </Badge>
              )}
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  closeTab(tab.id);
                }}
                className="hover:bg-muted rounded p-0.5"
              >
                <X className="h-3 w-3" />
              </button>
            </div>
          ))}
          <Button
            variant="ghost"
            size="sm"
            onClick={createNewTab}
            className="h-8 px-2"
          >
            <Plus className="h-4 w-4" />
          </Button>
        </div>
        <div className="flex items-center gap-2">
          {availableProviders.length > 0 && (
            <Select value={provider} onValueChange={setProvider}>
              <SelectTrigger className="w-[180px]">
                <SelectValue placeholder="Select provider" />
              </SelectTrigger>
              <SelectContent>
                {availableProviders.map((p: string) => (
                  <SelectItem key={p} value={p}>
                    {p}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          )}
          <TimeRangePicker value={timeRange} onChange={setTimeRange} />
          <Button
            variant="ghost"
            size="sm"
            onClick={onClose}
            className="flex-shrink-0"
          >
            <X className="h-4 w-4" />
          </Button>
        </div>
      </div>

      {/* Active Tab Content */}
      {activeTab && (
        <div className="flex-1 flex flex-col overflow-hidden">
          {/* Query Editor */}
          <div className="flex-none p-4 border-b space-y-3 bg-muted/20">
            <div className="flex items-start gap-2">
              <div className="flex-1 relative">
                <Textarea
                  value={activeTab.query}
                  onChange={(e) => updateTabQuery(activeTab.id, e.target.value)}
                  placeholder='Write your query... Example: FROM Logs WHERE level = "error"'
                  className="font-mono text-sm min-h-[80px] resize-none"
                  onKeyDown={(e) => {
                    if (e.key === 'Enter' && (e.metaKey || e.ctrlKey)) {
                      executeQuery(activeTab.id);
                    }
                  }}
                />
                <div className="absolute bottom-2 left-2 flex gap-1">
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="outline" size="sm" className="h-6 text-xs gap-1">
                        FROM
                        <ChevronDown className="h-3 w-3" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent>
                      <DropdownMenuItem 
                        onClick={() => activeTab && insertDataSource(activeTab.id, 'Logs')}
                        className="gap-2"
                      >
                        <Database className="h-4 w-4" />
                        Logs
                      </DropdownMenuItem>
                      <DropdownMenuItem 
                        onClick={() => activeTab && insertDataSource(activeTab.id, 'Traces')}
                        className="gap-2"
                      >
                        <GitBranch className="h-4 w-4" />
                        Traces
                      </DropdownMenuItem>
                      <DropdownMenuItem 
                        onClick={() => activeTab && insertDataSource(activeTab.id, 'Metrics')}
                        className="gap-2"
                      >
                        <BarChart3 className="h-4 w-4" />
                        Metrics
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              </div>
              <div className="flex flex-col gap-2">
                <Button
                  onClick={() => executeQuery(activeTab.id)}
                  className="gap-2"
                  disabled={!activeTab.query}
                >
                  <Play className="h-4 w-4" />
                  Run
                </Button>
                <Button
                  variant="outline"
                  className="gap-2"
                  size="sm"
                >
                  <Sparkles className="h-4 w-4" />
                  AI
                </Button>
              </div>
            </div>
            
            <div className="text-xs text-muted-foreground">
              Press <kbd className="px-1.5 py-0.5 bg-muted border rounded">⌘</kbd> + <kbd className="px-1.5 py-0.5 bg-muted border rounded">Enter</kbd> to run query
            </div>
          </div>

          {/* Results */}
          <div className="flex-1 overflow-hidden p-4">
            {error && (
              <div className="mb-4 p-3 bg-red-50 dark:bg-red-950 border border-red-200 dark:border-red-800 rounded text-red-600 dark:text-red-400 text-sm">
                {error}
              </div>
            )}
            
            {loading ? (
              <div className="flex flex-col items-center justify-center h-full text-center">
                <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center mb-4 animate-pulse">
                  <Database className="h-8 w-8 text-primary" />
                </div>
                <h3 className="text-lg font-medium mb-2">Executing query...</h3>
                <p className="text-muted-foreground">Please wait</p>
              </div>
            ) : !activeTab.results ? (
              <div className="flex flex-col items-center justify-center h-full text-center">
                <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center mb-4">
                  <Database className="h-8 w-8 text-primary" />
                </div>
                <h3 className="text-lg font-medium mb-2">Ready to Explore</h3>
                <p className="text-muted-foreground max-w-md">
                  Write a query and click Run to start exploring your observability data
                </p>
                <div className="mt-4 text-sm text-muted-foreground space-y-1">
                  <p>Examples:</p>
                  <code className="block bg-muted px-3 py-1 rounded">FROM Logs WHERE level = "error"</code>
                  <code className="block bg-muted px-3 py-1 rounded">FROM Traces WHERE status = "error"</code>
                  <code className="block bg-muted px-3 py-1 rounded">FROM Metrics</code>
                </div>
              </div>
            ) : (
              <>
                {activeTab.dataSource === 'Logs' && (
                  <LogsResult logs={activeTab.results as LogEntry[]} />
                )}
                                    {activeTab.dataSource === 'Traces' && (
                      <TracesResult traces={activeTab.results as TraceSpan[]} provider={provider} />
                    )}
                {activeTab.dataSource === 'Metrics' && (
                  <MetricsResult metrics={activeTab.results} />
                )}
              </>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
