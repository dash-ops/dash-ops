import { Badge } from '@/components/ui/badge';
import { KubernetesTypes } from '@/types';

interface SimpleDeploymentStatusProps {
  conditions: KubernetesTypes.DeploymentCondition[];
  replicas: KubernetesTypes.DeploymentReplicas;
}

export default function SimpleDeploymentStatus({
  conditions,
  replicas,
}: SimpleDeploymentStatusProps): JSX.Element {
  const safeConditions = conditions || [];
  const availableCondition = safeConditions.find((c) => c.type === 'Available');
  const progressingCondition = safeConditions.find((c) => c.type === 'Progressing');

  const isAvailable = availableCondition?.status === 'True';
  const isProgressing = progressingCondition?.status === 'True';
  const isReady = replicas.ready === replicas.desired && replicas.desired > 0;

  if (isReady && isAvailable) {
    return (
      <div className="flex gap-1">
        <Badge variant="default" className="text-xs">
          Available
        </Badge>
      </div>
    );
  }

  if (isProgressing) {
    return (
      <div className="flex gap-1">
        <Badge variant="secondary" className="text-xs">
          Progressing
        </Badge>
      </div>
    );
  }

  return (
    <div className="flex gap-1">
      <Badge variant="destructive" className="text-xs">
        Failed
      </Badge>
    </div>
  );
}
