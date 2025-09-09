import { useState, useEffect, useReducer, useCallback } from 'react';
import {
  Link,
  useLocation,
  useNavigate,
  useParams,
  useSearchParams,
} from 'react-router';
import { toast } from 'sonner';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { getPods } from './podsResource';
import { getNamespacesCached } from './namespacesCache';
import Refresh from '../../components/Refresh';
import PodContainers from './PodContainers';
import PodQoSBadge from './PodQoSBadge';
import { formatAge } from './helpers';
import { FileText } from 'lucide-react';
import { KubernetesTypes, BadgeVariant } from '@/types';

const INITIAL_STATE: KubernetesTypes.PodState = { data: [], loading: false };
const LOADING = 'LOADING';
const SET_DATA = 'SET_DATA';

function reducer(
  state: KubernetesTypes.PodState,
  action: KubernetesTypes.PodAction
): KubernetesTypes.PodState {
  switch (action.type) {
    case LOADING:
      return { ...state, loading: true, data: [] };
    case SET_DATA:
      return { ...state, loading: false, data: action.response };
    default:
      return state;
  }
}

async function fetchData(
  dispatch: React.Dispatch<KubernetesTypes.PodAction>,
  filter: KubernetesTypes.PodFilter,
  config?: { signal?: AbortSignal }
): Promise<void> {
  try {
    const result = await getPods(filter, config);
    dispatch({ type: SET_DATA, response: result.data });
  } catch (e: unknown) {
    if (e instanceof Error && e.message === 'Request canceled') {
      return;
    }
    console.error('Fetch error:', e);
    toast.error('Ops... Failed to fetch API data');
    dispatch({ type: SET_DATA, response: [] });
  }
}

export default function PodPage(): JSX.Element {
  const { context } = useParams<{ context: string }>();
  const navigate = useNavigate();
  const location = useLocation();
  const [searchParams] = useSearchParams();
  const [search, setSearch] = useState<string>(searchParams.get('name') ?? '');
  const [namespace, setNamespace] = useState<string>(
    searchParams.get('namespace') ?? 'All'
  );
  const [namespaces, setNamespaces] = useState<KubernetesTypes.Namespace[]>([]);
  const [pods, dispatch] = useReducer(reducer, INITIAL_STATE);

  useEffect(() => {
    if (!context) return;

    getNamespacesCached(context)
      .then((namespaces) => {
        setNamespaces([{ name: 'All', status: 'Active' }, ...namespaces]);
      })
      .catch((e: unknown) => {
        console.error('Error fetching namespaces:', e);
      });
  }, [context]);

  useEffect(() => {
    if (!context) return;

    const controller = new AbortController();
    const signal = controller.signal;
    dispatch({ type: LOADING });

    const namespaceFilter = namespace === 'All' ? '' : namespace;
    fetchData(dispatch, { context, namespace: namespaceFilter }, { signal });

    return () => {
      controller.abort();
    };
  }, [context, namespace]);

  const onReload = useCallback(async () => {
    if (!context) return;
    const namespaceFilter = namespace === 'All' ? '' : namespace;
    fetchData(dispatch, { context, namespace: namespaceFilter });
  }, [context, namespace]);

  const searchHandler = (value: string): void => {
    setSearch(value);
    navigate(`${location.pathname}?name=${value}&namespace=${namespace}`);
  };

  const handleNamespaceChange = (newNamespace: string): void => {
    setNamespace(newNamespace);
    navigate(`${location.pathname}?name=${search}&namespace=${newNamespace}`);
  };

  const filteredData = pods.data.filter(
    (p) => search === '' || p.name.includes(search)
  );

  const getStatusColor = (status: string): BadgeVariant => {
    switch (status) {
      case 'Running':
        return 'default';
      case 'Succeeded':
        return 'secondary';
      case 'Pending':
        return 'outline';
      default:
        return 'destructive';
    }
  };

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
      <div className="grid grid-cols-1 md:grid-cols-12 gap-4">
        <div className="md:col-span-3">
          <Input
            placeholder="Search pods..."
            onChange={(e) => searchHandler(e.target.value)}
            value={search}
          />
        </div>
        <div className="md:col-span-1">
          <Button
            variant="outline"
            onClick={() => searchHandler('')}
            className="w-full"
          >
            Clear
          </Button>
        </div>
        <div className="md:col-span-3">
          <div className="space-y-1">
            <Select value={namespace} onValueChange={handleNamespaceChange}>
              <SelectTrigger>
                <SelectValue placeholder="Select namespace" />
              </SelectTrigger>
              <SelectContent>
                {namespaces.map((ns) => (
                  <SelectItem key={ns.name} value={ns.name}>
                    {ns.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>
        <div className="hidden md:block md:col-span-2" />
        <div className="md:col-span-3 flex justify-end">
          <Refresh onReload={onReload} />
        </div>
      </div>

      {pods.data.length > 0 && (
        <div className="border rounded-lg overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-[250px]">Name</TableHead>
                <TableHead className="w-[120px]">Namespace</TableHead>
                <TableHead className="w-[120px]">Containers</TableHead>
                <TableHead className="w-[80px]">Restarts</TableHead>
                <TableHead className="w-[150px]">Controlled By</TableHead>
                <TableHead className="w-[150px]">Node</TableHead>
                <TableHead className="w-[100px]">QoS</TableHead>
                <TableHead className="w-[80px]">Age</TableHead>
                <TableHead className="w-[100px]">Status</TableHead>
                <TableHead className="w-[60px] text-center">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {pods.loading ? (
                <TableRow>
                  <TableCell colSpan={10} className="text-center py-8">
                    <div className="flex items-center justify-center">
                      <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-gray-900" />
                      <span className="ml-2">Loading...</span>
                    </div>
                  </TableCell>
                </TableRow>
              ) : filteredData.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={10} className="text-center py-8">
                    <span className="text-muted-foreground">No pods found</span>
                  </TableCell>
                </TableRow>
              ) : (
                filteredData.map((pod) => (
                  <TableRow key={pod.name} className="hover:bg-muted/50">
                    <TableCell className="font-medium text-foreground">
                      {pod.name}
                    </TableCell>

                    <TableCell>
                      <Badge variant="outline" className="text-xs">
                        {pod.namespace}
                      </Badge>
                    </TableCell>

                    <TableCell>
                      <PodContainers containers={pod.containers} />
                    </TableCell>

                    <TableCell className="text-center">
                      <Badge
                        variant={
                          pod.restarts > 0 ? 'destructive' : 'secondary'
                        }
                        className="text-xs"
                      >
                        {pod.restarts}
                      </Badge>
                    </TableCell>

                    <TableCell>
                      <div className="text-sm text-muted-foreground truncate">
                        {pod.containers && pod.containers.length > 0 && pod.containers[0] ? pod.containers[0].name : '-'}
                      </div>
                    </TableCell>

                    <TableCell>
                      <div className="text-sm text-blue-600 hover:text-blue-800 cursor-pointer truncate">
                        {pod.node}
                      </div>
                    </TableCell>

                    <TableCell>
                      <PodQoSBadge qosClass={pod.qos_class || "BestEffort"} />
                    </TableCell>

                    <TableCell className="text-sm text-muted-foreground">
                      {formatAge(pod.age)}
                    </TableCell>

                    <TableCell>
                      <Badge
                        variant={getStatusColor(pod.status)}
                        className="text-xs"
                      >
                        {pod.status}
                      </Badge>
                    </TableCell>

                    <TableCell className="text-center">
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <Button
                              variant="ghost"
                              size="sm"
                              className="h-8 w-8 p-0"
                              asChild
                            >
                              <Link
                                to={`/k8s/${context}/pod/logs?name=${pod.name}&namespace=${namespace === 'All' ? pod.namespace : namespace}`}
                              >
                                <FileText className="h-4 w-4" />
                              </Link>
                            </Button>
                          </TooltipTrigger>
                          <TooltipContent>
                            <p>View container logs</p>
                          </TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </div>
      )}
    </div>
  );
}
