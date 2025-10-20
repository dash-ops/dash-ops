import { describe, it, expect } from 'vitest';
import * as tracesAdapter from '../adapters/tracesAdapter';

describe('tracesAdapter', () => {
  it('transformTraceInfoToDomain maps fields and aggregates', () => {
    const api = {
      trace_id: 't1',
      spans: [
        { service: 'auth', status: 'ok', duration: 10 },
        { service: 'api', status: 'error', duration: 20 },
      ],
    };
    const info = tracesAdapter.transformTraceInfoToDomain(api);
    expect(info.traceId).toBe('t1');
    expect(info.serviceCount).toBe(2);
    expect(info.status).toBe('error');
  });

  it('transformTraceSpansToDomain maps spans', () => {
    const apiSpans = [
      { spanId: 's1', trace_id: 't1', operation: 'GET /x', service: 'api', duration: 30 },
    ];
    const spans = tracesAdapter.transformTraceSpansToDomain(apiSpans);
    expect(spans[0].id).toBe('s1');
    expect(spans[0].operationName).toContain('GET');
  });
});


