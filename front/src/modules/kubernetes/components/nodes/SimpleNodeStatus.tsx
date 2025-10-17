import { Badge } from '@/components/ui/badge';
import { KubernetesTypes } from '@/types';

interface SimpleNodeStatusProps {
  conditions: KubernetesTypes.NodeCondition[];
}

export default function SimpleNodeStatus({
  conditions,
}: SimpleNodeStatusProps): JSX.Element {
  const readyCondition = conditions.find((c) => c.type === 'Ready');
  const isReady = readyCondition?.status === 'True';

  return (
    <Badge variant={isReady ? 'default' : 'destructive'} className="text-xs">
      {isReady ? 'Ready' : 'Not Ready'}
    </Badge>
  );
}
