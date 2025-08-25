interface ResourceProgressBarProps {
  label: string;
  percentage: number;
  color: 'cpu' | 'memory' | 'disk';
  className?: string;
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
}: ResourceProgressBarProps): JSX.Element {
  const safePercentage = Math.max(0, Math.min(100, percentage));

  return (
    <div className={`space-y-1 ${className}`}>
      <div className="text-xs text-muted-foreground font-medium">{label}</div>
      <div className="w-full bg-muted rounded-full h-2">
        <div
          className={`h-2 rounded-full transition-all duration-300 ${colorClasses[color]}`}
          style={{ width: `${safePercentage}%` }}
        />
      </div>
      <div className="text-xs text-muted-foreground">
        {safePercentage.toFixed(1)}%
      </div>
    </div>
  );
}
