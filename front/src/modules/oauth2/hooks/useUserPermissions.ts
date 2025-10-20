import { useState, useEffect, useCallback } from 'react';
import * as userResource from '../resources/userResource';
import * as userAdapter from '../adapters/userAdapter';
import type { UserPermission, PermissionData } from '../types';

interface UseUserPermissionsState {
  permissions: UserPermission[];
  permissionData: PermissionData[];
  teams: string[];
  organizations: string[];
  loading: boolean;
  error: string | null;
}

interface UseUserPermissionsReturn extends UseUserPermissionsState {
  fetchPermissions: () => Promise<void>;
  refresh: () => Promise<void>;
  hasPermission: (plugin: string, feature: string, action?: string) => boolean;
  getUserTeams: () => string[];
  getUserOrganizations: () => string[];
}

export const useUserPermissions = (): UseUserPermissionsReturn => {
  const [state, setState] = useState<UseUserPermissionsState>({
    permissions: [],
    permissionData: [],
    teams: [],
    organizations: [],
    loading: false,
    error: null,
  });

  const fetchPermissions = useCallback(async () => {
    try {
      setState(prev => ({ ...prev, loading: true, error: null }));
      
      const response = await userResource.getUserPermissions();
      const permissions = userAdapter.transformUserPermissionsToDomain(response.data);
      
      const allPermissions = permissions.reduce((acc, perm) => ({ ...acc, ...perm.permissions }), {});
      const permissionData = userAdapter.transformPermissionsToPermissionData(allPermissions);
      const teams = userAdapter.getUserTeams(permissions);
      const organizations = userAdapter.getUserOrganizations(permissions);
      
      setState(prev => ({ 
        ...prev, 
        permissions,
        permissionData,
        teams,
        organizations,
        loading: false 
      }));
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to fetch user permissions';
      setState(prev => ({ 
        ...prev, 
        error: errorMessage, 
        loading: false 
      }));
    }
  }, []);

  const refresh = useCallback(async () => {
    await fetchPermissions();
  }, [fetchPermissions]);

  const hasPermission = useCallback((plugin: string, feature: string, action?: string): boolean => {
    return userAdapter.hasPermission(state.permissions, plugin, feature, action);
  }, [state.permissions]);

  const getUserTeams = useCallback((): string[] => {
    return state.teams;
  }, [state.teams]);

  const getUserOrganizations = useCallback((): string[] => {
    return state.organizations;
  }, [state.organizations]);

  useEffect(() => {
    fetchPermissions();
  }, [fetchPermissions]);

  return {
    ...state,
    fetchPermissions,
    refresh,
    hasPermission,
    getUserTeams,
    getUserOrganizations,
  };
};
