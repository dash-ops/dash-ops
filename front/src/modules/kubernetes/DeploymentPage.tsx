import { useState, useEffect, useReducer, useCallback } from 'react';
import { useParams } from 'react-router';
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
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import {
  getDeployments,
  upDeployment,
  downDeployment,
} from './deploymentResource';
import { getNamespacesCached } from './namespacesCache';
import Refresh from '../../components/Refresh';
import DeploymentActions from './DeploymentActions';
import { KubernetesTypes } from '@/types';

const INITIAL_STATE: KubernetesTypes.DeploymentState = {
  data: [],
  loading: false,
};
const LOADING = 'LOADING';
const SET_DATA = 'SET_DATA';

function reducer(
  state: KubernetesTypes.DeploymentState,
  action: KubernetesTypes.DeploymentAction
): KubernetesTypes.DeploymentState {
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
  dispatch: React.Dispatch<KubernetesTypes.DeploymentAction>,
  filter: KubernetesTypes.DeploymentFilter,
  config?: { signal?: AbortSignal }
): Promise<void> {
  try {
    const result = await getDeployments(filter, config);
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

async function toUp(
  context: string,
  deployment: KubernetesTypes.Deployment,
  setNewPodCount: (name: string, podCount: number) => void
): Promise<void> {
  try {
    setNewPodCount(deployment.name, 1);
    await upDeployment(context, deployment.name, deployment.namespace);
  } catch (e: any) {
    setNewPodCount(deployment.name, 0);
    toast.error(
      `Failed to try to up deployment: ${e.data?.error || e.message}`
    );
  }
}

async function toDown(
  context: string,
  deployment: KubernetesTypes.Deployment,
  setNewPodCount: (name: string, podCount: number) => void
): Promise<void> {
  try {
    setNewPodCount(deployment.name, 0);
    await downDeployment(context, deployment.name, deployment.namespace);
  } catch (e: any) {
    setNewPodCount(deployment.name, 1);
    toast.error(
      `Failed to try to down deployment: ${e.data?.error || e.message}`
    );
  }
}

export default function DeploymentPage(): JSX.Element {
  const { context } = useParams<{ context: string }>();
  const [search, setSearch] = useState<string>('');
  const [namespace, setNamespace] = useState<string>('All');
  const [namespaces, setNamespaces] = useState<KubernetesTypes.Namespace[]>([]);
  const [deployments, dispatch] = useReducer(reducer, INITIAL_STATE);

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

  const handleNamespaceChange = (newNamespace: string): void => {
    setNamespace(newNamespace);
  };

  const updatePodCount = (name: string, podCount: number): void => {
    const newDeployments = deployments.data.map((dep) =>
      dep.name === name ? { ...dep, pod_count: podCount } : dep
    );
    dispatch({ type: SET_DATA, response: newDeployments });
  };

  const filteredData = deployments.data.filter(
    (deployment) => search === '' || deployment.name.includes(search)
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
      <div className="grid grid-cols-1 md:grid-cols-12 gap-4">
        <div className="md:col-span-3">
          <Input
            placeholder="Search deployments..."
            onChange={(e) => setSearch(e.target.value)}
            value={search}
          />
        </div>
        <div className="md:col-span-1">
          <Button
            variant="outline"
            onClick={() => setSearch('')}
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

      {deployments.data.length > 0 && (
        <div className="border rounded-lg overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-[300px]">Name</TableHead>
                <TableHead>Pods Info</TableHead>
                <TableHead className="w-[140px] text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {deployments.loading ? (
                <TableRow>
                  <TableCell colSpan={3} className="text-center py-8">
                    <div className="flex items-center justify-center">
                      <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-gray-900" />
                      <span className="ml-2">Loading...</span>
                    </div>
                  </TableCell>
                </TableRow>
              ) : filteredData.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={3} className="text-center py-8">
                    <span className="text-muted-foreground">
                      No deployments found
                    </span>
                  </TableCell>
                </TableRow>
              ) : (
                filteredData.map((deployment) => (
                  <TableRow key={deployment.name}>
                    <TableCell className="font-medium">
                      {deployment.name}
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant={
                          deployment.pod_info.current > 0
                            ? 'default'
                            : 'destructive'
                        }
                      >
                        {deployment.pod_info.current}/
                        {deployment.pod_info.desired}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-right">
                      <DeploymentActions
                        context={context}
                        deployment={deployment}
                        toUp={() => toUp(context, deployment, updatePodCount)}
                        toDown={() =>
                          toDown(context, deployment, updatePodCount)
                        }
                      />
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
