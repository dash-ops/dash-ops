import { useState, useEffect, useCallback } from 'react';
import * as userResource from '../resources/userResource';
import * as userAdapter from '../adapters/userAdapter';
import type { UserData, UserPermission } from '../types';

interface UseUserState {
  user: UserData | null;
  loading: boolean;
  error: string | null;
}

interface UseUserReturn extends UseUserState {
  fetchUser: () => Promise<void>;
  refresh: () => Promise<void>;
}

export const useUser = (): UseUserReturn => {
  const [state, setState] = useState<UseUserState>({
    user: null,
    loading: false,
    error: null,
  });

  const fetchUser = useCallback(async () => {
    try {
      setState(prev => ({ ...prev, loading: true, error: null }));
      
      const response = await userResource.getUserData();
      const user = userAdapter.transformUserDataToDomain(response.data);
      
      setState(prev => ({ 
        ...prev, 
        user, 
        loading: false 
      }));
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to fetch user data';
      setState(prev => ({ 
        ...prev, 
        error: errorMessage, 
        loading: false 
      }));
    }
  }, []);

  const refresh = useCallback(async () => {
    await fetchUser();
  }, [fetchUser]);

  useEffect(() => {
    fetchUser();
  }, [fetchUser]);

  return {
    ...state,
    fetchUser,
    refresh,
  };
};
