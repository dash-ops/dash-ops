interface ResourceProgressBarProps {
  label: string;
  percentage: number;
  color: 'cpu' | 'memory' | 'disk';
  className?: string;
  tooltip?: string;
}

const colorClasses = {
  cpu: 'bg-blue-500',
  memory: 'bg-pink-500',
  disk: 'bg-yellow-500',
};

export default function ResourceProgressBar({
  label,
  percentage,
  color,
  className = '',
  tooltip,
}: ResourceProgressBarProps): JSX.Element {
  const safePercentage = Math.max(0, Math.min(100, percentage));

  return (
    <div className={`space-y-1 ${className}`}>
      <div className="text-xs text-muted-foreground font-medium">{label}</div>
      <div
        className="w-full bg-muted rounded-full h-2 relative group"
        title={tooltip}
      >
        <div
          className={`h-2 rounded-full transition-all duration-300 ${colorClasses[color]}`}
          style={{ width: `${safePercentage}%` }}
        />
        {tooltip && (
          <div className="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-2 py-1 bg-gray-900 text-white text-xs rounded opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none whitespace-nowrap z-10">
            {tooltip}
          </div>
        )}
      </div>
      <div className="text-xs text-muted-foreground">
        {safePercentage.toFixed(1)}%
      </div>
    </div>
  );
}
