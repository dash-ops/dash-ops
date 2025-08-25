import { KubernetesTypes } from '@/types';
import { getClusters } from './clusterResource';

let clustersCache: KubernetesTypes.Cluster[] | null = null;
let loadingPromise: Promise<KubernetesTypes.Cluster[]> | null = null;

export async function getClustersCached(): Promise<KubernetesTypes.Cluster[]> {
  if (clustersCache) {
    return clustersCache;
  }

  if (loadingPromise) {
    return loadingPromise;
  }

  loadingPromise = (async () => {
    try {
      const { data } = await getClusters();
      clustersCache = data;
      return data;
    } finally {
      loadingPromise = null;
    }
  })();

  return loadingPromise;
}

export function clearClustersCache(): void {
  clustersCache = null;
  loadingPromise = null;
}
