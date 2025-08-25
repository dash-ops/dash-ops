import { useState, useEffect } from 'react';
import { useLocation, useNavigate } from 'react-router';
import { toast } from 'sonner';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import { Server } from 'lucide-react';
import { KubernetesTypes } from '@/types';
import { Page } from '@/types';
import ContentWithMenu from '../../components/ContentWithMenu';
import { getClustersCached } from './clustersCache';

interface KubernetesWithContextSelectorProps {
  pages: Page[];
}

export default function KubernetesWithContextSelector({
  pages,
}: KubernetesWithContextSelectorProps): JSX.Element {
  const location = useLocation();
  const navigate = useNavigate();
  const [clusters, setClusters] = useState<KubernetesTypes.Cluster[]>([]);
  const [selectedContext, setSelectedContext] = useState<string>('');
  const [loading, setLoading] = useState(true);
  const [isLoadingClusters, setIsLoadingClusters] = useState(false);

  useEffect(() => {
    if (isLoadingClusters || clusters.length > 0) return;

    async function loadClusters() {
      setIsLoadingClusters(true);
      try {
        const data = await getClustersCached();
        setClusters(data);
        setLoading(false);
      } catch (error) {
        console.error('Failed to load clusters:', error);
        toast.error('Failed to load Kubernetes clusters');
        setLoading(false);
      } finally {
        setIsLoadingClusters(false);
      }
    }

    loadClusters();
  }, []);

  useEffect(() => {
    if (clusters.length === 0) return;

    const pathParts = location.pathname.split('/');
    const contextIndex = pathParts.findIndex((part) => part === 'k8s') + 1;
    const currentContext = pathParts[contextIndex];

    if (
      currentContext &&
      clusters.some((cluster) => cluster.context === currentContext)
    ) {
      setSelectedContext(currentContext);
    } else {
      setSelectedContext('');
    }
  }, [location.pathname, clusters]);

  const handleContextChange = (newContext: string) => {
    setSelectedContext(newContext);
    navigate(`/k8s/${newContext}`);
  };

  const pagesWithContext = pages;

  if (!loading && clusters.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 space-y-4">
        <Server className="h-12 w-12 text-muted-foreground" />
        <div className="text-center space-y-2">
          <h3 className="font-semibold text-lg">No Kubernetes Clusters</h3>
          <p className="text-muted-foreground max-w-md">
            No Kubernetes clusters are configured. Please check your
            configuration file.
          </p>
        </div>
      </div>
    );
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary" />
        <span className="ml-3">Loading clusters...</span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="border-b pb-4">
        <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
          <div className="space-y-1">
            <h1 className="text-2xl font-semibold tracking-tight">
              Kubernetes
            </h1>
            <p className="text-muted-foreground">
              Manage your Kubernetes clusters and workloads
            </p>
          </div>

          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Server className="h-4 w-4" />
              Cluster:
            </div>
            <Select value={selectedContext} onValueChange={handleContextChange}>
              <SelectTrigger className="w-[280px]">
                <SelectValue placeholder="Select a cluster" />
              </SelectTrigger>
              <SelectContent>
                {clusters.map((cluster) => (
                  <SelectItem key={cluster.context} value={cluster.context}>
                    <div className="flex items-center gap-2">
                      <span>{cluster.name}</span>
                      <Badge variant="secondary" className="text-xs">
                        {cluster.context}
                      </Badge>
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>
      </div>

      {!selectedContext ? (
        <div className="text-center py-8">
          <span className="text-muted-foreground">
            No cluster context provided
          </span>
        </div>
      ) : (
        <ContentWithMenu
          pages={pagesWithContext}
          paramName="context"
          contextValue={selectedContext}
        />
      )}
    </div>
  );
}
