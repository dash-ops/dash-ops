export const formatTimestamp = (timestamp: string | number): string => {
  const date = new Date(timestamp);
  if (Number.isNaN(date.getTime())) return 'Invalid Date';
  return date.toLocaleString();
};

export const formatDuration = (durationMs: number): string => {
  if (durationMs < 1000) return `${Math.round(durationMs)} ms`;
  return `${(durationMs / 1000).toFixed(2)} s`;
};

export const levelColor = (level: 'info' | 'warn' | 'error' | 'debug'): string => {
  switch (level) {
    case 'error':
      return 'text-red-600 bg-red-50 border-red-200';
    case 'warn':
      return 'text-yellow-600 bg-yellow-50 border-yellow-200';
    case 'info':
      return 'text-blue-600 bg-blue-50 border-blue-200';
    case 'debug':
    default:
      return 'text-gray-600 bg-gray-50 border-gray-200';
  }
};


