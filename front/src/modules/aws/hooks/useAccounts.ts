/**
 * Custom hook for managing AWS accounts state
 */

import { useState, useEffect, useCallback } from 'react';
import { getAccounts } from '../resources/accountResource';

interface UseAccountsReturn {
  accounts: string[];
  loading: boolean;
  error: string | null;
  refresh: () => Promise<void>;
}

export function useAccounts(): UseAccountsReturn {
  const [accounts, setAccounts] = useState<string[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchAccounts = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await getAccounts();
      setAccounts(response.data || []);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to fetch accounts';
      setError(errorMessage);
      console.error('Error fetching accounts:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchAccounts();
  }, [fetchAccounts]);

  const refresh = useCallback(async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await getAccounts();
      setAccounts(response.data || []);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to refresh accounts';
      setError(errorMessage);
      console.error('Error refreshing accounts:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  return {
    accounts,
    loading,
    error,
    refresh,
  };
}
