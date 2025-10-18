/**
 * Tests for useInstances hook
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useInstances } from '../hooks/useInstances';
import * as instanceResource from '../resources/instanceResource';
import * as instanceAdapter from '../adapters/instanceAdapter';

// Mock the dependencies
vi.mock('../resources/instanceResource');
vi.mock('../adapters/instanceAdapter');

describe('useInstances', () => {
  const mockFilter = { accountKey: 'prod' };
  const mockApiResponse = {
    data: [
      {
        instance_id: 'i-1',
        name: 'test-instance-1',
        state: { name: 'running', code: 16 },
        platform: 'linux',
        instance_type: 't3.micro',
        public_ip: '1.1.1.1',
        private_ip: '10.0.0.1',
        cpu: { vcpus: 2 },
        memory: { size_gb: 1 },
        tags: [],
        launch_time: '2023-01-01T00:00:00Z',
        account: 'prod',
        region: 'us-east-1',
        cost_estimate: 0,
      },
    ],
  } as any;

  const mockTransformedInstance = {
    id: 'i-1',
    name: 'test-instance-1',
    instance_id: 'i-1',
    state: { name: 'running', code: 16 },
    platform: 'linux',
    instance_type: 't3.micro',
    public_ip: '1.1.1.1',
    private_ip: '10.0.0.1',
    cpu: { vcpus: 2 },
    memory: { size_gb: 1 },
    tags: [],
    launch_time: '2023-01-01T00:00:00Z',
    account: 'prod',
    region: 'us-east-1',
    cost_estimate: 0,
  } as any;

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should fetch instances on mount', async () => {
    vi.mocked(instanceResource.getInstances).mockResolvedValue(mockApiResponse as any);
    vi.mocked(instanceAdapter.transformInstancesToDomain).mockReturnValue([mockTransformedInstance] as any);

    const { result } = renderHook(() => useInstances(mockFilter));

    expect(result.current.loading).toBe(true);
    expect(result.current.instances).toEqual([]);
    expect(result.current.error).toBeNull();

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(instanceResource.getInstances).toHaveBeenCalledWith(mockFilter);
    expect(instanceAdapter.transformInstancesToDomain).toHaveBeenCalledWith(mockApiResponse.data);
    expect(result.current.instances).toEqual([mockTransformedInstance]);
    expect(result.current.error).toBeNull();
  });

  it('should handle fetch error', async () => {
    const errorMessage = 'Failed to fetch instances';
    vi.mocked(instanceResource.getInstances).mockRejectedValue(new Error(errorMessage));

    const { result } = renderHook(() => useInstances(mockFilter));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    expect(result.current.instances).toEqual([]);
    expect(result.current.error).toBe(errorMessage);
  });

  it('should not fetch when accountKey is empty', () => {
    const emptyFilter = { accountKey: '' };
    
    renderHook(() => useInstances(emptyFilter));

    expect(instanceResource.getInstances).not.toHaveBeenCalled();
  });

  it('should refresh instances when refresh is called', async () => {
    vi.mocked(instanceResource.getInstances).mockResolvedValue(mockApiResponse as any);
    vi.mocked(instanceAdapter.transformInstancesToDomain).mockReturnValue([mockTransformedInstance] as any);

    const { result } = renderHook(() => useInstances(mockFilter));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    // Clear previous calls
    vi.clearAllMocks();
    vi.mocked(instanceResource.getInstances).mockResolvedValue(mockApiResponse as any);
    vi.mocked(instanceAdapter.transformInstancesToDomain).mockReturnValue([mockTransformedInstance] as any);

    // Call refresh
    await act(async () => {
      await result.current.refresh();
    });

    expect(instanceResource.getInstances).toHaveBeenCalledWith(mockFilter);
    expect(result.current.instances).toEqual([mockTransformedInstance]);
  });

  it('should start instance successfully', async () => {
    const mockStartResponse = { data: { current_state: 'running' } } as any;
    vi.mocked(instanceResource.getInstances).mockResolvedValue(mockApiResponse as any);
    vi.mocked(instanceAdapter.transformInstancesToDomain).mockReturnValue([mockTransformedInstance] as any);
    vi.mocked(instanceResource.startInstance).mockResolvedValue(mockStartResponse as any);

    const { result } = renderHook(() => useInstances(mockFilter));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    // Clear previous calls
    vi.clearAllMocks();
    vi.mocked(instanceResource.getInstances).mockResolvedValue(mockApiResponse as any);
    vi.mocked(instanceAdapter.transformInstancesToDomain).mockReturnValue([mockTransformedInstance] as any);

    await act(async () => {
      await result.current.startInstance('i-1');
    });

    expect(instanceResource.startInstance).toHaveBeenCalledWith('prod', 'i-1');
    expect(instanceResource.getInstances).toHaveBeenCalledWith(mockFilter);
  });

  it('should stop instance successfully', async () => {
    const mockStopResponse = { data: { current_state: 'stopped' } } as any;
    vi.mocked(instanceResource.getInstances).mockResolvedValue(mockApiResponse as any);
    vi.mocked(instanceAdapter.transformInstancesToDomain).mockReturnValue([mockTransformedInstance] as any);
    vi.mocked(instanceResource.stopInstance).mockResolvedValue(mockStopResponse as any);

    const { result } = renderHook(() => useInstances(mockFilter));

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    // Clear previous calls
    vi.clearAllMocks();
    vi.mocked(instanceResource.getInstances).mockResolvedValue(mockApiResponse as any);
    vi.mocked(instanceAdapter.transformInstancesToDomain).mockReturnValue([mockTransformedInstance] as any);

    await act(async () => {
      await result.current.stopInstance('i-1');
    });

    expect(instanceResource.stopInstance).toHaveBeenCalledWith('prod', 'i-1');
    expect(instanceResource.getInstances).toHaveBeenCalledWith(mockFilter);
  });

  it('should refetch instances when filter changes', async () => {
    const newFilter = { accountKey: 'staging' };
    vi.mocked(instanceResource.getInstances).mockResolvedValue(mockApiResponse as any);
    vi.mocked(instanceAdapter.transformInstancesToDomain).mockReturnValue([mockTransformedInstance] as any);

    const { result, rerender } = renderHook(
      ({ filter }) => useInstances(filter),
      { initialProps: { filter: mockFilter } }
    );

    await waitFor(() => {
      expect(result.current.loading).toBe(false);
    });

    // Change filter
    await act(async () => {
      rerender({ filter: newFilter });
    });

    expect(instanceResource.getInstances).toHaveBeenCalledWith(newFilter);
  });
});
