import { describe, it, expect } from 'vitest';
import * as userUtils from '../utils/userUtils';
import type { UserData, UserPermission } from '../types';

describe('userUtils', () => {
  const mockUser: UserData = {
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

  it('should get user display name', () => {
    expect(userUtils.getUserDisplayName(mockUser)).toBe('testuser');
    expect(userUtils.getUserDisplayName(null)).toBe('Unknown User');
    
    const userWithoutLogin = { ...mockUser, login: undefined };
    expect(userUtils.getUserDisplayName(userWithoutLogin)).toBe('test@example.com');
  });

  it('should get user initials', () => {
    expect(userUtils.getUserInitials(mockUser)).toBe('TE');
    expect(userUtils.getUserInitials(null)).toBe('UU');
    
    const userWithoutLogin = { ...mockUser, login: undefined };
    expect(userUtils.getUserInitials(userWithoutLogin)).toBe('TE');
  });

  it('should check if user is authenticated', () => {
    expect(userUtils.isUserAuthenticated(mockUser)).toBe(true);
    expect(userUtils.isUserAuthenticated(null)).toBe(false);
    
    const userWithoutId = { ...mockUser, id: undefined };
    expect(userUtils.isUserAuthenticated(userWithoutId)).toBe(false);
  });

  it('should format user location', () => {
    expect(userUtils.formatUserLocation(mockUser)).toBe('Test Location');
    expect(userUtils.formatUserLocation(null)).toBe('Unknown location');
  });

  it('should format user company', () => {
    expect(userUtils.formatUserCompany(mockUser)).toBe('Test Company');
    expect(userUtils.formatUserCompany(null)).toBe('No company');
  });

  it('should format user bio', () => {
    expect(userUtils.formatUserBio(mockUser)).toBe('Test bio');
    expect(userUtils.formatUserBio(null)).toBe('No bio available');
  });

  it('should check if user has blog', () => {
    expect(userUtils.hasUserBlog(mockUser)).toBe(true);
    expect(userUtils.hasUserBlog(null)).toBe(false);
    
    const userWithoutBlog = { ...mockUser, blog: undefined };
    expect(userUtils.hasUserBlog(userWithoutBlog)).toBe(false);
  });

  it('should check if user has profile URL', () => {
    expect(userUtils.hasUserProfileUrl(mockUser)).toBe(true);
    expect(userUtils.hasUserProfileUrl(null)).toBe(false);
    
    const userWithoutUrl = { ...mockUser, html_url: undefined };
    expect(userUtils.hasUserProfileUrl(userWithoutUrl)).toBe(false);
  });

  it('should format permission display name', () => {
    const permission = {
      name: 'aws.instances',
      resource: 'instances',
      actions: ['read', 'write'],
    };
    
    expect(userUtils.getPermissionDisplayName(permission)).toBe('aws › instances');
  });

  it('should format permission actions', () => {
    expect(userUtils.formatPermissionActions(['read'])).toBe('read');
    expect(userUtils.formatPermissionActions(['read', 'write'])).toBe('read and write');
    expect(userUtils.formatPermissionActions(['read', 'write', 'delete'])).toBe('read, write and delete');
    expect(userUtils.formatPermissionActions([])).toBe('None');
  });

  it('should get permission level', () => {
    const permissions: UserPermission[] = [
      {
        id: '1',
        permissions: {
          aws: {
            instances: ['read'],
          },
        },
      },
    ];
    
    expect(userUtils.getPermissionLevel(permissions, 'aws', 'instances')).toBe('read');
    expect(userUtils.getPermissionLevel(permissions, 'aws', 'nonexistent')).toBe('none');
  });

  it('should get permission color', () => {
    expect(userUtils.getPermissionColor('admin')).toContain('red');
    expect(userUtils.getPermissionColor('write')).toContain('yellow');
    expect(userUtils.getPermissionColor('read')).toContain('green');
    expect(userUtils.getPermissionColor('none')).toContain('gray');
  });

  it('should get permission label', () => {
    expect(userUtils.getPermissionLabel('admin')).toBe('Admin');
    expect(userUtils.getPermissionLabel('write')).toBe('Write');
    expect(userUtils.getPermissionLabel('read')).toBe('Read');
    expect(userUtils.getPermissionLabel('none')).toBe('None');
  });

  it('should validate user data', () => {
    expect(userUtils.validateUserData(mockUser)).toBe(true);
    expect(userUtils.validateUserData({ id: '123' })).toBe(false);
    expect(userUtils.validateUserData({ login: 'testuser' })).toBe(false);
    expect(userUtils.validateUserData(null)).toBe(false);
  });

  it('should sanitize user data', () => {
    const rawUser = {
      id: '123',
      login: 'testuser',
      extraField: 'should be removed',
    };
    
    const result = userUtils.sanitizeUserData(rawUser);
    
    expect(result.id).toBe('123');
    expect(result.login).toBe('testuser');
    expect(result.email).toBe('');
    expect('extraField' in result).toBe(false);
  });
});
