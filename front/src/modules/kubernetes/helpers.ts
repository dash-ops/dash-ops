/**
 * Shared helper functions for Kubernetes module
 */

/**
 * Parse resource quantity string (e.g., "1991880Ki", "1750m") to number
 * @param quantity - Resource quantity string
 * @returns Parsed number value
 */
export function parseResourceQuantity(quantity: string): number {
  if (!quantity || quantity === '0') return 0;
  
  // Handle different units
  const units: { [key: string]: number } = {
    'Ki': 1024,
    'Mi': 1024 * 1024,
    'Gi': 1024 * 1024 * 1024,
    'Ti': 1024 * 1024 * 1024 * 1024,
    'K': 1000,
    'M': 1000 * 1000,
    'G': 1000 * 1000 * 1000,
    'T': 1000 * 1000 * 1000 * 1000,
    'm': 0.001, // millicores
  };

  // Extract number and unit
  const match = quantity.match(/^(\d+(?:\.\d+)?)([a-zA-Z]*)$/);
  if (!match) return 0;

  const value = parseFloat(match[1]);
  const unit = match[2];

  if (unit === '') return value; // No unit, return as is
  if (units[unit]) return value * units[unit];
  
  return value; // Unknown unit, return as is
}

/**
 * Calculate usage percentage for resources
 * @param used - Used resource quantity string
 * @param capacity - Total capacity quantity string
 * @returns Usage percentage (0-100)
 */
export function calculateUsagePercentage(used: string, capacity: string): number {
  const usedNum = parseResourceQuantity(used);
  const capacityNum = parseResourceQuantity(capacity);
  
  if (capacityNum === 0) return 0;
  return Math.min((usedNum / capacityNum) * 100, 100);
}

/**
 * Format age string to show seconds if < 1m, otherwise show hours and minutes
 * @param age - Age string (e.g., "124h51m22.669937s")
 * @returns Formatted age string
 */
export function formatAge(age: string): string {
  if (!age) return '0s';
  
  // Parse the age string (e.g., "124h51m22.669937s")
  const hoursMatch = age.match(/(\d+)h/);
  const minutesMatch = age.match(/(\d+)m/);
  const secondsMatch = age.match(/(\d+(?:\.\d+)?)s/);
  
  const hours = hoursMatch && hoursMatch[1] ? parseInt(hoursMatch[1], 10) : 0;
  const minutes = minutesMatch && minutesMatch[1] ? parseInt(minutesMatch[1], 10) : 0;
  const seconds = secondsMatch && secondsMatch[1] ? parseFloat(secondsMatch[1]) : 0;
  
  if (hours > 0) {
    return `${hours}h${minutes > 0 ? `${minutes}m` : ''}`;
  } else if (minutes > 0) {
    return `${minutes}m`;
  } else {
    // Show seconds if less than 1 minute
    return `${Math.round(seconds)}s`;
  }
}
