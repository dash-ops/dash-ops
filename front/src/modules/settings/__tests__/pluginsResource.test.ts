import { describe, it, expect, vi } from 'vitest';
import * as pluginsResource from '../resources/pluginsResource';

vi.mock('@/helpers/http');

describe('pluginsResource', () => {
  it('fetches plugins from the API', async () => {
    const mockResponse = { data: ['aws', 'kubernetes', 'service-catalog'] };

    const http = await import('@/helpers/http');
    http.default.get = vi.fn().mockResolvedValue(mockResponse);

    const result = await pluginsResource.getPlugins();

    expect(http.default.get).toHaveBeenCalledWith('/config/plugins');
    expect(result.data).toEqual(['aws', 'kubernetes', 'service-catalog']);
  });

  it('handles empty plugin responses', async () => {
    const mockResponse = { data: [] };

    const http = await import('@/helpers/http');
    http.default.get = vi.fn().mockResolvedValue(mockResponse);

    const result = await pluginsResource.getPlugins();

    expect(result.data).toEqual([]);
  });

  it('normalizes null plugin responses to empty list', async () => {
    const mockResponse = { data: null };

    const http = await import('@/helpers/http');
    http.default.get = vi.fn().mockResolvedValue(mockResponse);

    const result = await pluginsResource.getPlugins();

    expect(result.data).toEqual([]);
  });
});

