/**
 * Module-specific test setup for AWS module
 * This file contains mocks specific to the AWS module
 * and can be imported by individual test files when needed
 */

import { vi } from 'vitest';

// Mock AWS-specific dependencies
vi.mock('../resources/instanceResource', () => ({
  getInstances: vi.fn(),
  startInstance: vi.fn(),
  stopInstance: vi.fn(),
  restartInstance: vi.fn(),
  terminateInstance: vi.fn(),
  getInstanceDetails: vi.fn(),
}));

vi.mock('../resources/accountResource', () => ({
  getAccounts: vi.fn(),
}));

vi.mock('../resources/permissionResource', () => ({
  getPermissions: vi.fn(),
}));

vi.mock('../adapters/instanceAdapter', () => ({
  transformTags: vi.fn(),
  transformInstanceToDomain: vi.fn(),
  transformInstanceToApiRequest: vi.fn(),
  transformInstancesToDomain: vi.fn(),
  filterInstancesByState: vi.fn(),
  sortInstancesByLaunchTime: vi.fn(),
}));

vi.mock('../utils/accountsCache', () => ({
  getAccountsCached: vi.fn(),
  clearAccountsCache: vi.fn(),
}));

// Mock React Router with AWS-specific params
vi.mock('react-router', () => ({
  useParams: vi.fn(() => ({ key: 'test-account' })),
  useLocation: vi.fn(() => ({ pathname: '/aws/test-account' })),
  useNavigate: vi.fn(() => vi.fn()),
}));

// Mock Sonner toast
vi.mock('sonner', () => ({
  toast: {
    error: vi.fn(),
    success: vi.fn(),
    info: vi.fn(),
    warning: vi.fn(),
  },
}));

// Mock HTTP helper
vi.mock('../../../helpers/http', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    delete: vi.fn(),
  },
}));
