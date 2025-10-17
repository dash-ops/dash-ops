import { useState, useEffect, useReducer, useCallback } from 'react';
import { useParams, useSearchParams } from 'react-router';
import { toast } from 'sonner';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { getNodes } from '../../resources/nodesResource';
import Refresh from '../../../../components/Refresh';
import ResourceProgressBar from '../shared/ResourceProgressBar';
import SimpleNodeStatus from './SimpleNodeStatus';
import { formatAge, calculateUsagePercentage } from '../../utils/helpers';
import { KubernetesTypes } from '@/types';

const INITIAL_STATE: KubernetesTypes.NodesState = { data: [], loading: false };
const LOADING = 'LOADING';
const SET_DATA = 'SET_DATA';

function reducer(
  state: KubernetesTypes.NodesState,
  action: KubernetesTypes.NodesAction
): KubernetesTypes.NodesState {
  switch (action.type) {
    case LOADING:
      return { ...state, loading: true, data: [] };
    case SET_DATA:
      return { ...state, loading: false, data: action.response };
    default:
      return state;
  }
}

export default function NodesPage(): JSX.Element {
  const { context } = useParams<{ context: string }>();
  const [searchParams] = useSearchParams();
  const [search, setSearch] = useState<string>(searchParams.get('node') ?? '');
  const [nodes, dispatch] = useReducer(reducer, INITIAL_STATE);

  const fetchData = useCallback(
    async (config?: { signal?: AbortSignal }) => {
      if (!context) return;

      try {
        dispatch({ type: LOADING });
        const result = await getNodes({ context }, config);
        dispatch({ type: SET_DATA, response: result.data });
      } catch (e: unknown) {
        if (e instanceof Error && e.message === 'Request canceled') {
          return;
        }
        console.error('Fetch error:', e);
        toast.error('Ops... Failed to fetch API data');
        dispatch({ type: SET_DATA, response: [] });
      }
    },
    [context]
  );

  useEffect(() => {
    const controller = new AbortController();
    const signal = controller.signal;
    fetchData({ signal });

    return () => {
      controller.abort();
    };
  }, [fetchData]);

  const onReload = useCallback(async () => {
    fetchData();
  }, [fetchData]);

  const filteredData = nodes.data.filter((node) =>
    node.name.toLowerCase().includes(search.toLowerCase())
  );

  if (!context) {
    return (
      <div className="text-center py-8">
        <span className="text-muted-foreground">
          No cluster context provided
        </span>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="md:col-span-1">
          <Input
            placeholder="Search by node name"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
        <div className="hidden md:block" />
        <div className="flex justify-end">
          <Refresh onReload={onReload} />
        </div>
      </div>

      <div className="border rounded-lg">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-[200px]">Name</TableHead>
              <TableHead className="w-[120px]">CPU Usage</TableHead>
              <TableHead className="w-[120px]">Memory Usage</TableHead>
              <TableHead className="w-[120px]">Pod Usage</TableHead>
              <TableHead className="w-[80px]">Status</TableHead>
              <TableHead className="w-[120px]">Roles</TableHead>
              <TableHead className="w-[100px]">Version</TableHead>
              <TableHead className="w-[80px]">Age</TableHead>
              <TableHead className="w-[100px]">Conditions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {nodes.loading ? (
              <TableRow>
                <TableCell colSpan={9} className="text-center py-8">
                  <div className="flex items-center justify-center">
                    <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-gray-900" />
                    <span className="ml-2">Loading...</span>
                  </div>
                </TableCell>
              </TableRow>
            ) : filteredData.length === 0 ? (
              <TableRow>
                <TableCell colSpan={9} className="text-center py-8">
                  <span className="text-muted-foreground">No nodes found</span>
                </TableCell>
              </TableRow>
            ) : (
              filteredData.map((node) => (
                <TableRow key={node.name} className="hover:bg-muted/50">
                  <TableCell className="font-medium text-foreground">
                    {node.name}
                  </TableCell>

                  <TableCell>
                    <ResourceProgressBar
                      label=""
                      percentage={calculateUsagePercentage(
                        node.resources.used.cpu,
                        node.resources.capacity.cpu
                      )}
                      color="cpu"
                      tooltip={`${node.resources.used.cpu} / ${node.resources.capacity.cpu}`}
                    />
                  </TableCell>

                  <TableCell>
                    <ResourceProgressBar
                      label=""
                      percentage={calculateUsagePercentage(
                        node.resources.used.memory,
                        node.resources.capacity.memory
                      )}
                      color="memory"
                      tooltip={`${node.resources.used.memory} / ${node.resources.capacity.memory}`}
                    />
                  </TableCell>

                  <TableCell>
                    <ResourceProgressBar
                      label=""
                      percentage={calculateUsagePercentage(
                        node.resources.used.pods,
                        node.resources.capacity.pods
                      )}
                      color="disk"
                      tooltip={`${node.resources.used.pods} / ${node.resources.capacity.pods} pods`}
                    />
                  </TableCell>

                  <TableCell className="text-center">
                    <Badge variant="outline" className="text-xs">
                      {node.status}
                    </Badge>
                  </TableCell>

                  <TableCell>
                    <div className="flex flex-wrap gap-1">
                      {(node.roles || []).map((role, index) => (
                        <Badge
                          key={index}
                          variant="secondary"
                          className="text-xs"
                        >
                          {role}
                        </Badge>
                      ))}
                    </div>
                  </TableCell>

                  <TableCell className="text-sm text-muted-foreground">
                    {node.version}
                  </TableCell>

                  <TableCell className="text-sm text-muted-foreground">
                    {formatAge(node.age)}
                  </TableCell>

                  <TableCell>
                    <SimpleNodeStatus conditions={node.conditions} />
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
