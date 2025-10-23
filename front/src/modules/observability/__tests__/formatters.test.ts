import { describe, it, expect } from 'vitest';
import { formatTimestamp, formatDuration, levelColor } from '../utils/formatters';

describe('formatters', () => {
  it('formatTimestamp returns string', () => {
    const out = formatTimestamp('2025-01-01T00:00:00Z');
    expect(typeof out).toBe('string');
  });

  it('formatDuration formats ms and s', () => {
    expect(formatDuration(500)).toContain('ms');
    expect(formatDuration(1500)).toContain('s');
  });

  it('levelColor returns css class', () => {
    expect(levelColor('error')).toContain('red');
    expect(levelColor('warn')).toContain('yellow');
  });
});


