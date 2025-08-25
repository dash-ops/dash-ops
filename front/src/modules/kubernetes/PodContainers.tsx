import { Badge } from '@/components/ui/badge';
import { KubernetesTypes } from '@/types';

interface PodContainersProps {
  containers: KubernetesTypes.PodContainer[];
}

export default function PodContainers({
  containers,
}: PodContainersProps): JSX.Element {
  const readyCount = containers.filter((c) => c.ready).length;
  const totalCount = containers.length;

  if (containers.length === 0) {
    return (
      <Badge variant="outline" className="text-xs">
        0/0
      </Badge>
    );
  }

  return (
    <Badge
      variant={readyCount === totalCount && totalCount > 0 ? 'default' : 'destructive'}
      className="text-xs"
    >
      {readyCount}/{totalCount}
    </Badge>
  );
}
