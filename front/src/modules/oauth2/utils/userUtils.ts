import type { UserData, UserPermission, PermissionData } from '../types';

export const getUserDisplayName = (user: UserData | null): string => {
  if (!user) return 'Unknown User';
  return user.login || user.email || 'Unknown User';
};

export const getUserInitials = (user: UserData | null): string => {
  if (!user) return 'UU';
  
  const name = user.login || user.email || '';
  const parts = name.split(/[@._-]/);
  
  if (parts.length >= 2) {
    return parts[0].charAt(0).toUpperCase() + parts[1].charAt(0).toUpperCase();
  }
  
  return name.charAt(0).toUpperCase() + name.charAt(1).toUpperCase();
};

export const getUserAvatarUrl = (user: UserData | null): string => {
  return user?.avatar_url || '';
};

export const isUserAuthenticated = (user: UserData | null): boolean => {
  return !!(user?.id);
};

export const formatUserLocation = (user: UserData | null): string => {
  return user?.location || 'Unknown location';
};

export const formatUserCompany = (user: UserData | null): string => {
  return user?.company || 'No company';
};

export const formatUserBio = (user: UserData | null): string => {
  return user?.bio || 'No bio available';
};

export const formatUserBlog = (user: UserData | null): string => {
  return user?.blog || '';
};

export const hasUserBlog = (user: UserData | null): boolean => {
  return !!(user?.blog);
};

export const getUserProfileUrl = (user: UserData | null): string => {
  return user?.html_url || '';
};

export const hasUserProfileUrl = (user: UserData | null): boolean => {
  return !!(user?.html_url);
};

export const getPermissionDisplayName = (permission: PermissionData): string => {
  return permission.name.replace(/\./g, ' › ');
};

export const getPermissionResource = (permission: PermissionData): string => {
  return permission.resource;
};

export const getPermissionActions = (permission: PermissionData): string[] => {
  return permission.actions;
};

export const formatPermissionActions = (actions: string[]): string => {
  if (actions.length === 0) return 'None';
  if (actions.length === 1) return actions[0];
  if (actions.length === 2) return actions.join(' and ');
  return `${actions.slice(0, -1).join(', ')} and ${actions[actions.length - 1]}`;
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

export const hasTeamMembership = (permissions: UserPermission[], teamSlug: string): boolean => {
  return getUserTeams(permissions).includes(teamSlug);
};

export const hasOrganizationMembership = (permissions: UserPermission[], organization: string): boolean => {
  return getUserOrganizations(permissions).includes(organization);
};

export const getPermissionLevel = (permissions: UserPermission[], plugin: string, feature: string): 'none' | 'read' | 'write' | 'admin' => {
  const permission = permissions.find(p => p.permissions?.[plugin]?.[feature]);
  if (!permission) return 'none';
  
  const actions = permission.permissions[plugin][feature];
  const actionArray = Array.isArray(actions) ? actions : [actions];
  
  if (actionArray.includes('admin')) return 'admin';
  if (actionArray.includes('write') || actionArray.includes('create') || actionArray.includes('update') || actionArray.includes('delete')) return 'write';
  if (actionArray.includes('read') || actionArray.includes('view')) return 'read';
  
  return 'none';
};

export const getPermissionColor = (level: 'none' | 'read' | 'write' | 'admin'): string => {
  switch (level) {
    case 'admin':
      return 'bg-red-100 text-red-800 border-red-200';
    case 'write':
      return 'bg-yellow-100 text-yellow-800 border-yellow-200';
    case 'read':
      return 'bg-green-100 text-green-800 border-green-200';
    case 'none':
    default:
      return 'bg-gray-100 text-gray-800 border-gray-200';
  }
};

export const getPermissionLabel = (level: 'none' | 'read' | 'write' | 'admin'): string => {
  switch (level) {
    case 'admin':
      return 'Admin';
    case 'write':
      return 'Write';
    case 'read':
      return 'Read';
    case 'none':
    default:
      return 'None';
  }
};

export const validateUserData = (user: any): boolean => {
  return !!(user?.id && (user.login || user.email));
};

export const sanitizeUserData = (user: any): UserData => {
  return {
    id: user.id || '',
    login: user.login || '',
    email: user.email || '',
    avatar_url: user.avatar_url || '',
    bio: user.bio || '',
    location: user.location || '',
    company: user.company || '',
    blog: user.blog || '',
    html_url: user.html_url || '',
  };
};
