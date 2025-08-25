export function formatDiskInfo(capacity: {
  storage?: string;
  ephemeral_storage?: string;
  disk_pressure: boolean;
}): string {
  const parts = [];
  if (capacity.storage) parts.push(capacity.storage);
  if (capacity.disk_pressure) parts.push('Pressure');

  return parts.length > 0 ? parts.join(' | ') : 'N/A';
}

export function getConditionStatus(
  type: string,
  conditions?: Array<{ type: string; status: string }>
): 'True' | 'False' | 'Unknown' {
  const condition = conditions?.find((c) => c.type === type);
  return (condition?.status as 'True' | 'False' | 'Unknown') || 'Unknown';
}
