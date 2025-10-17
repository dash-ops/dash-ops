import { AWSTypes } from '@/types';
import { getAccounts } from '../resources/accountResource';

let accountsCache: AWSTypes.Account[] | null = null;
let loadingPromise: Promise<AWSTypes.Account[]> | null = null;

export async function getAccountsCached(): Promise<AWSTypes.Account[]> {
  if (accountsCache) {
    return accountsCache;
  }

  if (loadingPromise) {
    return loadingPromise;
  }

  loadingPromise = (async () => {
    try {
      const { data } = await getAccounts();
      accountsCache = data;
      return data;
    } finally {
      loadingPromise = null;
    }
  })();

  return loadingPromise;
}

export function clearAccountsCache(): void {
  accountsCache = null;
  loadingPromise = null;
}
