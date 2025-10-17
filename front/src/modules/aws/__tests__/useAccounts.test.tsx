/**
 * Tests for useAccounts hook
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { useAccounts } from '../hooks/useAccounts';
import * as accountResource from '../resources/accountResource';

// Mock the dependencies
vi.mock('../resources/accountResource');

describe('useAccounts', () => {
  const mockApiResponse = {
    data: [
      {
        name: 'Production',
        key: 'prod',
        region: 'us-east-1',
        status: 'active',
      },
      {
        name: 'Staging',
        key: 'staging',
        region: 'us-west-2',
        status: 'active',
      },
    ],
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should fetch accounts on mount', async () => {
    vi.mocked(accountResource.getAccounts).mockResolvedValue(mockApiResponse);

    const { result } = renderHook(() => useAccounts());

    expect(result.current.loading).toBe(true);
    expect(result.current.accounts).toEqual([]);
    expect(result.current.error).toBeNull();

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(accountResource.getAccounts).toHaveBeenCalledWith();
    expect(result.current.accounts).toEqual(mockApiResponse.data);
    expect(result.current.error).toBeNull();
  });

  it('should handle fetch error', async () => {
    const errorMessage = 'Failed to fetch accounts';
    vi.mocked(accountResource.getAccounts).mockRejectedValue(new Error(errorMessage));

    const { result } = renderHook(() => useAccounts());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.accounts).toEqual([]);
    expect(result.current.error).toBe(errorMessage);
  });

  it('should refresh accounts when refresh is called', async () => {
    vi.mocked(accountResource.getAccounts).mockResolvedValue(mockApiResponse);

    const { result } = renderHook(() => useAccounts());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    // Clear previous calls
    vi.clearAllMocks();
    vi.mocked(accountResource.getAccounts).mockResolvedValue(mockApiResponse);

    // Call refresh
    await result.current.refresh();

    expect(accountResource.getAccounts).toHaveBeenCalledWith();
    expect(result.current.accounts).toEqual(mockApiResponse.data);
  });


  it('should handle empty accounts response', async () => {
    const emptyResponse = { data: [] };
    vi.mocked(accountResource.getAccounts).mockResolvedValue(emptyResponse);

    const { result } = renderHook(() => useAccounts());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.accounts).toEqual([]);
    expect(result.current.error).toBeNull();
  });

  it('should handle undefined accounts in response', async () => {
    const responseWithUndefinedAccounts = { data: undefined as any };
    vi.mocked(accountResource.getAccounts).mockResolvedValue(responseWithUndefinedAccounts);

    const { result } = renderHook(() => useAccounts());

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.accounts).toEqual([]);
    expect(result.current.error).toBeNull();
  });
});
