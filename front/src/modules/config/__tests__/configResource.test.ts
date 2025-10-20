import { describe, it, expect, vi } from 'vitest';
import * as configResource from '../resources/configResource';

vi.mock('../../../helpers/http');

describe('configResource', () => {
  it('should fetch plugins', async () => {
    const mockResponse = { data: ['aws', 'kubernetes', 'service-catalog'] };
    
    // Mock the http.get method
    const http = await import('../../../helpers/http');
    http.default.get = vi.fn().mockResolvedValue(mockResponse);
    
    const result = await configResource.getPlugins();
    
    expect(http.default.get).toHaveBeenCalledWith('/config/plugins');
    expect(result.data).toEqual(['aws', 'kubernetes', 'service-catalog']);
  });

  it('should handle empty plugins response', async () => {
    const mockResponse = { data: [] };
    
    const http = await import('../../../helpers/http');
    http.default.get = vi.fn().mockResolvedValue(mockResponse);
    
    const result = await configResource.getPlugins();
    
    expect(result.data).toEqual([]);
  });

  it('should handle null plugins response', async () => {
    const mockResponse = { data: null };
    
    const http = await import('../../../helpers/http');
    http.default.get = vi.fn().mockResolvedValue(mockResponse);
    
    const result = await configResource.getPlugins();
    
    expect(result.data).toEqual([]);
  });
});
