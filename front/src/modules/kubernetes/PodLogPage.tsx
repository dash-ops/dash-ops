import { useCallback, useEffect, useReducer, useState } from 'react';
import { useParams, useNavigate, useSearchParams } from 'react-router';
import { toast } from 'sonner';
import type { AxiosRequestConfig } from 'axios';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { ChevronLeft, Search, Copy, Maximize2, Minimize2 } from 'lucide-react';
import { cancelToken } from '../../helpers/http';
import { getPodLogs } from './podsResource';
import Refresh from '../../components/Refresh';
import { KubernetesTypes } from '@/types';

interface LogState {
  data: KubernetesTypes.PodLogsResponse | null;
  loading: boolean;
}

const INITIAL_STATE: LogState = { data: null, loading: false };
const LOADING = 'LOADING';
const SET_DATA = 'SET_DATA';

type LogAction =
  | { type: 'LOADING' }
  | { type: 'SET_DATA'; response: KubernetesTypes.PodLogsResponse };

function reducer(state: LogState, action: LogAction): LogState {
  switch (action.type) {
    case LOADING:
      return { ...state, loading: true, data: null };
    case SET_DATA:
      return { ...state, loading: false, data: action.response };
    default:
      return state;
  }
}

async function fetchData(
  dispatch: React.Dispatch<LogAction>,
  filter: KubernetesTypes.PodLogsFilter,
  config?: AxiosRequestConfig
): Promise<void> {
  try {
    const result = await getPodLogs(filter, config);
    dispatch({ type: SET_DATA, response: result.data });
  } catch (error) {
    // Ignore cancellation errors
    if (
      error &&
      typeof error === 'object' &&
      'message' in error &&
      (error as any).message === 'Request canceled'
    ) {
      return;
    }

    // Check if it's an axios error
    if (error && typeof error === 'object' && 'response' in error) {
      const axiosError = error as any;

      if (axiosError.response?.status === 404) {
        toast.error(
          'Pod logs not found. Check if the pod exists and is running.'
        );
      } else if (axiosError.response?.status === 401) {
        toast.error('Unauthorized. Please check your authentication.');
      } else {
        toast.error(
          `Failed to fetch pod logs: ${axiosError.response?.statusText || 'Unknown error'}`
        );
      }
    } else {
      toast.error('Ops... Failed to fetch API data');
    }

    dispatch({
      type: SET_DATA,
      response: { pod_name: '', namespace: '', logs: [], total_lines: 0 },
    });
  }
}

function copyToClipboard(text: string, containerName: string): void {
  navigator.clipboard
    .writeText(text)
    .then(() => {
      toast.success(`Logs copied for container: ${containerName}`);
    })
    .catch(() => {
      toast.error('Failed to copy logs to clipboard');
    });
}

function filterLogs(
  logs: KubernetesTypes.PodLogEntry[],
  searchTerm: string
): KubernetesTypes.PodLogEntry[] {
  if (!searchTerm.trim()) return logs;

  const lowerSearchTerm = searchTerm.toLowerCase();

  return logs.filter((log) =>
    log.message.toLowerCase().includes(lowerSearchTerm)
  );
}

export default function PodLogPage(): JSX.Element {
  const { context } = useParams<{ context: string }>();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const name = searchParams.get('name') ?? '';
  const namespace = searchParams.get('namespace') ?? 'default';
  const [logs, dispatch] = useReducer(reducer, INITIAL_STATE);
  const [searchTerm, setSearchTerm] = useState('');
  const [isExpanded, setIsExpanded] = useState(false);

  useEffect(() => {
    if (!context || !name) return;

    const source = cancelToken.source();
    dispatch({ type: LOADING });
    fetchData(
      dispatch,
      { context, name, namespace },
      { cancelToken: source.token }
    );
    return () => {
      source.cancel();
    };
  }, [context, name, namespace]);

  const onReload = useCallback(async () => {
    if (!context || !name) return;
    fetchData(dispatch, { context, name, namespace });
  }, [context, name, namespace]);

  if (!context || !name) {
    return (
      <div className="text-center py-8">
        <span className="text-muted-foreground">
          Missing required parameters (context or pod name)
        </span>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Header with navigation and controls */}
      <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
        <div className="flex items-center gap-4">
          <Button
            onClick={() => navigate(-1)}
            variant="outline"
            size="sm"
            className="gap-2"
          >
            <ChevronLeft className="h-4 w-4" />
            Go Back
          </Button>
          <div className="text-sm text-muted-foreground">
            Pod: <span className="font-medium">{name}</span> | Namespace:{' '}
            <span className="font-medium">{namespace}</span>
          </div>
        </div>
        <Refresh onReload={onReload} />
      </div>

      {/* Search bar */}
      {logs.data && logs.data.logs.length > 0 && (
        <div className="flex items-center gap-2">
          <div className="relative flex-1 max-w-md">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search in logs..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="pl-9"
            />
          </div>
          <Button
            onClick={() => setIsExpanded(!isExpanded)}
            variant="outline"
            size="sm"
            className="gap-2"
          >
            {isExpanded ? (
              <Minimize2 className="h-4 w-4" />
            ) : (
              <Maximize2 className="h-4 w-4" />
            )}
            {isExpanded ? 'Collapse' : 'Expand'}
          </Button>
        </div>
      )}

      {/* Log display */}
      {logs.data && logs.data.logs.length > 0 && (
        <div className="space-y-4">
          <div className="flex items-center justify-between p-4 border rounded-lg bg-muted/50">
            <div className="flex items-center gap-3">
              <span className="font-semibold text-base">
                Pod: {logs.data.pod_name}
              </span>
              <Badge variant="outline" className="text-xs">
                {logs.data.namespace}
              </Badge>
              <Badge variant="secondary" className="text-xs">
                {logs.data.total_lines} lines
              </Badge>
            </div>
            <div className="flex items-center gap-2">
              <Button
                onClick={() => {
                  const allLogs = logs
                    .data!.logs.map((log) => log.message)
                    .join('\n');
                  copyToClipboard(allLogs, logs.data!.pod_name);
                }}
                variant="ghost"
                size="sm"
                className="gap-1"
              >
                <Copy className="h-3 w-3" />
                Copy All
              </Button>
            </div>
          </div>

          <div className="border rounded-lg overflow-hidden">
            <div
              className={`bg-slate-950 text-slate-100 font-mono text-sm ${
                isExpanded ? 'h-[80vh]' : 'max-h-96'
              } overflow-auto`}
            >
              <div className="p-4">
                {(() => {
                  const filteredLogs = filterLogs(logs.data!.logs, searchTerm);
                  const matchCount = searchTerm.trim()
                    ? filteredLogs.length
                    : 0;

                  if (searchTerm.trim()) {
                    return (
                      <div className="mb-2 text-yellow-400 text-sm">
                        {matchCount} matches found
                      </div>
                    );
                  }

                  return null;
                })()}

                {(() => {
                  const filteredLogs = filterLogs(logs.data!.logs, searchTerm);

                  if (filteredLogs.length > 0) {
                    return (
                      <div className="space-y-0">
                        {filteredLogs.map((log, index) => {
                          const shouldHighlight =
                            searchTerm.trim() &&
                            log.message
                              .toLowerCase()
                              .includes(searchTerm.toLowerCase());

                          return (
                            <div key={index} className="flex items-start group">
                              <div className="select-none text-slate-500 text-xs min-w-[3rem] pr-3 py-0.5 text-right border-r border-slate-700 mr-3">
                                {index + 1}
                              </div>
                              <div className="flex-1 py-0.5">
                                <div className="text-slate-400 text-xs mb-1">
                                  {new Date(log.timestamp).toLocaleString()}
                                </div>
                                <div
                                  className={`whitespace-pre-wrap break-words ${
                                    shouldHighlight
                                      ? 'bg-yellow-900/50 text-yellow-200'
                                      : ''
                                  }`}
                                >
                                  {log.message || '\u00A0'}
                                </div>
                              </div>
                            </div>
                          );
                        })}
                      </div>
                    );
                  } else {
                    return (
                      <div className="text-slate-400 text-center py-8">
                        No logs match your search criteria
                      </div>
                    );
                  }
                })()}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Loading state */}
      {logs.loading && (
        <div className="flex flex-col items-center justify-center py-12 space-y-3">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary" />
          <div className="text-sm text-muted-foreground">
            Loading pod logs...
          </div>
        </div>
      )}

      {/* Empty state */}
      {!logs.loading && (!logs.data || logs.data.logs.length === 0) && (
        <div className="text-center py-12">
          <div className="text-muted-foreground">
            No logs found for this pod
          </div>
          <div className="text-sm text-muted-foreground mt-1">
            Try refreshing or check if the pod is running
          </div>
        </div>
      )}
    </div>
  );
}
