import { getNamespaces } from './namespaceResource';
import { KubernetesTypes } from '@/types';

let namespacesCache: Array<{
  context: string;
  namespaces: KubernetesTypes.Namespace[];
}> = [];
let loadingPromises: Map<
  string,
  Promise<KubernetesTypes.Namespace[]>
> = new Map();

export async function getNamespacesCached(
  context: string
): Promise<KubernetesTypes.Namespace[]> {
  const cached = namespacesCache.find((item) => item.context === context);
  if (cached) {
    return cached.namespaces;
  }

  const existingPromise = loadingPromises.get(context);
  if (existingPromise) {
    return existingPromise;
  }

  const promise = getNamespaces({ context })
    .then((result) => {
      const namespaces = result.data;
      namespacesCache.push({ context, namespaces });
      loadingPromises.delete(context);
      return namespaces;
    })
    .catch((error) => {
      loadingPromises.delete(context);
      throw error;
    });

  loadingPromises.set(context, promise);
  return promise;
}

export function clearNamespacesCache(context?: string): void {
  if (context) {
    namespacesCache = namespacesCache.filter(
      (item) => item.context !== context
    );
    loadingPromises.delete(context);
  } else {
    namespacesCache = [];
    loadingPromises.clear();
  }
}
