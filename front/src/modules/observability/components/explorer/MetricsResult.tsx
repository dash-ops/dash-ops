import { BarChart3 } from 'lucide-react';

interface MetricsResultProps {
  metrics: any[];
}

export default function MetricsResult({ metrics }: MetricsResultProps): JSX.Element {
  return (
    <div className="flex flex-col items-center justify-center h-full text-center">
      <BarChart3 className="h-12 w-12 mx-auto mb-4 opacity-50 text-muted-foreground" />
      <p className="text-lg font-medium mb-2">Metrics visualization</p>
      <p className="text-sm text-muted-foreground">Coming soon</p>
      <p className="text-xs text-muted-foreground mt-2">{metrics.length} metrics found</p>
    </div>
  );
}

