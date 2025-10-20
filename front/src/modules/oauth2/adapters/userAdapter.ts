import type { UserData, UserPermission, PermissionData, PluginPermissions } from '../types';

export const transformUserDataToDomain = (apiUser: any): UserData => {
  return {
    id: apiUser.id || '',
    login: apiUser.login,
    email: apiUser.email,
    avatar_url: apiUser.avatar_url,
    bio: apiUser.bio,
    location: apiUser.location,
    company: apiUser.company,
    blog: apiUser.blog,
    html_url: apiUser.html_url,
  };
};

export const transformUserPermissionsToDomain = (apiPermissions: any[]): UserPermission[] => {
  return apiPermissions.map(permission => ({
    id: permission.id || '',
    organization: permission.organization,
    teams: permission.teams?.map((team: any) => ({
      id: team.id || '',
      slug: team.slug,
    })) || [],
    permissions: permission.permissions || {},
  }));
};

export const transformPermissionsToPermissionData = (permissions: any): PermissionData[] => {
  const permissionData: PermissionData[] = [];
  
  Object.entries(permissions).forEach(([plugin, pluginPermissions]) => {
    Object.entries(pluginPermissions as any).forEach(([feature, actions]) => {
      permissionData.push({
        name: `${plugin}.${feature}`,
        resource: feature,
        actions: Array.isArray(actions) ? actions : [actions],
      });
    });
  });
  
  return permissionData;
};

export const transformUserToProfileData = (user: UserData, permissions: UserPermission[] = []) => {
  const allPermissions = permissions.reduce((acc, perm) => ({ ...acc, ...perm.permissions }), {});
  
  return {
    user,
    permissions: transformPermissionsToPermissionData(allPermissions),
    teams: permissions.flatMap(p => p.teams || []),
    organizations: permissions.map(p => p.organization).filter(Boolean),
  };
};

export const getUserDisplayName = (user: UserData): string => {
  return user.login || user.email || 'Unknown User';
};

export const getUserAvatarUrl = (user: UserData): string => {
  return user.avatar_url || '';
};

export const hasPermission = (permissions: UserPermission[], plugin: string, feature: string, action?: string): boolean => {
  return permissions.some(permission => {
    const pluginPerms = permission.permissions?.[plugin];
    if (!pluginPerms) return false;
    
    const featurePerms = pluginPerms[feature];
    if (!featurePerms) return false;
    
    if (action) {
      return Array.isArray(featurePerms) 
        ? featurePerms.includes(action)
        : featurePerms === action;
    }
    
    return true;
  });
};

export const getUserTeams = (permissions: UserPermission[]): string[] => {
  return permissions.flatMap(p => p.teams?.map(t => t.slug || '') || []).filter(Boolean);
};

export const getUserOrganizations = (permissions: UserPermission[]): string[] => {
  return permissions.map(p => p.organization).filter(Boolean) as string[];
};
