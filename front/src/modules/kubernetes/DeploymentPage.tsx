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
  restartDeployment,
  scaleDeployment,
} from './deploymentResource';
import { getNamespacesCached } from './namespacesCache';
import Refresh from '../../components/Refresh';
import ModernDeploymentActions from './ModernDeploymentActions';
import SimpleDeploymentStatus from './SimpleDeploymentStatus';
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

async function handleRestart(
  context: string,
  deployment: KubernetesTypes.Deployment,
  onReload: () => void
): Promise<void> {
  try {
    await restartDeployment(context, deployment.name, deployment.namespace);
    toast.success(`Deployment ${deployment.name} restarted successfully`);
    onReload();
  } catch (e: unknown) {
    toast.error(
      `Failed to restart deployment: ${
        (e as { data?: { error?: string }; message?: string }).data?.error ||
        (e as { message?: string }).message
      }`
    );
  }
}

async function handleScale(
  context: string,
  deployment: KubernetesTypes.Deployment,
  replicas: number,
  onReload: () => void
): Promise<void> {
  try {
    await scaleDeployment(
      context,
      deployment.name,
      deployment.namespace,
      replicas
    );
    toast.success(
      `Deployment ${deployment.name} scaled to ${replicas} replicas`
    );
    onReload();
  } catch (e: unknown) {
    toast.error(
      `Failed to scale deployment: ${
        (e as { data?: { error?: string }; message?: string }).data?.error ||
        (e as { message?: string }).message
      }`
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
                <TableHead className="w-[200px]">Name</TableHead>
                <TableHead className="w-[120px]">Namespace</TableHead>
                <TableHead className="w-[100px]">Pods</TableHead>
                <TableHead className="w-[100px]">Replicas</TableHead>
                <TableHead className="w-[80px]">Age</TableHead>
                <TableHead className="w-[120px]">Conditions</TableHead>
                <TableHead className="w-[100px] text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {deployments.loading ? (
                <TableRow>
                  <TableCell colSpan={7} className="text-center py-8">
                    <div className="flex items-center justify-center">
                      <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-gray-900" />
                      <span className="ml-2">Loading...</span>
                    </div>
                  </TableCell>
                </TableRow>
              ) : filteredData.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} className="text-center py-8">
                    <span className="text-muted-foreground">
                      No deployments found
                    </span>
                  </TableCell>
                </TableRow>
              ) : (
                filteredData.map((deployment) => (
                  <TableRow key={deployment.name} className="hover:bg-muted/50">
                    <TableCell className="font-medium text-foreground">
                      {deployment.name}
                    </TableCell>

                    <TableCell>
                      <Badge variant="outline" className="text-xs">
                        {deployment.namespace}
                      </Badge>
                    </TableCell>

                    <TableCell>
                      <Badge
                        variant={
                          deployment.pod_info.current > 0
                            ? 'default'
                            : 'destructive'
                        }
                        className="text-xs"
                      >
                        {deployment.pod_info.current}/
                        {deployment.pod_info.desired}
                      </Badge>
                    </TableCell>

                    <TableCell>
                      <div className="text-sm">
                        <span className="text-foreground">
                          {deployment.replicas.ready}/
                          {deployment.replicas.desired}
                        </span>
                        {deployment.replicas.updated !==
                          deployment.replicas.desired && (
                          <span className="text-muted-foreground ml-1">
                            ({deployment.replicas.updated} updated)
                          </span>
                        )}
                      </div>
                    </TableCell>

                    <TableCell className="text-sm text-muted-foreground">
                      {deployment.age}
                    </TableCell>

                    <TableCell>
                      <SimpleDeploymentStatus
                        conditions={deployment.conditions}
                        replicas={deployment.replicas}
                      />
                    </TableCell>

                    <TableCell className="text-right">
                      <ModernDeploymentActions
                        context={context}
                        deployment={deployment}
                        onRestart={() =>
                          handleRestart(context, deployment, onReload)
                        }
                        onScale={(replicas) =>
                          handleScale(context, deployment, replicas, onReload)
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
