import { describe, it, expect } from 'vitest';
import * as logsAdapter from '../adapters/logsAdapter';

describe('logsAdapter', () => {
  it('transformLogToDomain maps fields correctly', () => {
    const api = {
      id: '1',
      timestamp: '2025-01-01T00:00:00Z',
      level: 'ERROR',
      service: 'auth',
      message: 'failed login',
      trace_id: 't1',
      span_id: 's1',
    };
    const log = logsAdapter.transformLogToDomain(api);
    expect(log.id).toBe('1');
    expect(log.level).toBe('error');
    expect(log.service).toBe('auth');
    expect(log.message).toBe('failed login');
    expect(log.traceId).toBe('t1');
    expect(log.spanId).toBe('s1');
  });

  it('transformLogsResponseToDomain handles arrays', () => {
    const apiResp = { items: [{ id: '1' }, { id: '2' }] };
    const res = logsAdapter.transformLogsResponseToDomain(apiResp);
    expect(res.items.length).toBe(2);
    expect(res.total).toBe(2);
  });
});


