import { describe, it, expect } from 'vitest';
import * as userAdapter from '../adapters/userAdapter';
import type { UserData, UserPermission } from '../types';

describe('userAdapter', () => {
  const mockApiUser = {
    id: '123',
    login: 'testuser',
    email: 'test@example.com',
    avatar_url: 'https://example.com/avatar.png',
    bio: 'Test bio',
    location: 'Test Location',
    company: 'Test Company',
    blog: 'https://testblog.com',
    html_url: 'https://github.com/testuser',
  };

  it('should transform API user to domain model', () => {
    const result = userAdapter.transformUserDataToDomain(mockApiUser);
    
    expect(result.id).toBe('123');
    expect(result.login).toBe('testuser');
    expect(result.email).toBe('test@example.com');
    expect(result.avatar_url).toBe('https://example.com/avatar.png');
    expect(result.bio).toBe('Test bio');
  });

  it('should transform user permissions to domain models', () => {
    const mockApiPermissions = [
      {
        id: '1',
        organization: 'test-org',
        teams: [{ id: '1', slug: 'test-team' }],
        permissions: {
          aws: {
            instances: ['read', 'write'],
          },
        },
      },
    ];
    
    const result = userAdapter.transformUserPermissionsToDomain(mockApiPermissions);
    
    expect(result).toHaveLength(1);
    expect(result[0].id).toBe('1');
    expect(result[0].organization).toBe('test-org');
    expect(result[0].teams).toHaveLength(1);
    expect(result[0].teams?.[0].slug).toBe('test-team');
  });

  it('should get user display name', () => {
    const user: UserData = {
      id: '123',
      login: 'testuser',
    };
    
    const result = userAdapter.getUserDisplayName(user);
    expect(result).toBe('testuser');
  });

  it('should get user avatar URL', () => {
    const user: UserData = {
      id: '123',
      avatar_url: 'https://example.com/avatar.png',
    };
    
    const result = userAdapter.getUserAvatarUrl(user);
    expect(result).toBe('https://example.com/avatar.png');
  });

  it('should check user permissions', () => {
    const permissions: UserPermission[] = [
      {
        id: '1',
        permissions: {
          aws: {
            instances: ['read', 'write'],
          },
        },
      },
    ];
    
    expect(userAdapter.hasPermission(permissions, 'aws', 'instances')).toBe(true);
    expect(userAdapter.hasPermission(permissions, 'aws', 'instances', 'read')).toBe(true);
    expect(userAdapter.hasPermission(permissions, 'aws', 'instances', 'delete')).toBe(false);
    expect(userAdapter.hasPermission(permissions, 'kubernetes', 'pods')).toBe(false);
  });

  it('should get user teams', () => {
    const permissions: UserPermission[] = [
      {
        id: '1',
        teams: [{ id: '1', slug: 'team1' }, { id: '2', slug: 'team2' }],
      },
    ];
    
    const result = userAdapter.getUserTeams(permissions);
    expect(result).toEqual(['team1', 'team2']);
  });

  it('should get user organizations', () => {
    const permissions: UserPermission[] = [
      {
        id: '1',
        organization: 'org1',
      },
      {
        id: '2',
        organization: 'org2',
      },
    ];
    
    const result = userAdapter.getUserOrganizations(permissions);
    expect(result).toEqual(['org1', 'org2']);
  });
});
