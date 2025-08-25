import { Badge } from '@/components/ui/badge';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { getConditionStatus } from './nodeHelpers';
import { KubernetesTypes } from '@/types';

interface NodeConditionsProps {
  conditions: KubernetesTypes.NodeCondition[];
}

const IMPORTANT_CONDITIONS = [
  { type: 'Ready', label: 'Ready' },
  { type: 'MemoryPressure', label: 'Memory' },
  { type: 'DiskPressure', label: 'Disk' },
  { type: 'PIDPressure', label: 'PID' },
  { type: 'NetworkUnavailable', label: 'Network' },
];

export default function NodeConditions({
  conditions,
}: NodeConditionsProps): JSX.Element {
  return (
    <TooltipProvider>
      <div className="flex items-center gap-1 flex-wrap">
        {IMPORTANT_CONDITIONS.map(({ type, label }) => {
          const status = getConditionStatus(type, conditions);
          const condition = conditions.find((c) => c.type === type);

          let variant: 'default' | 'secondary' | 'destructive' | 'outline' =
            'outline';

          if (type === 'Ready') {
            variant = status === 'True' ? 'default' : 'destructive';
          } else {
            // For pressure conditions, True means there IS pressure (bad)
            variant =
              status === 'True'
                ? 'destructive'
                : status === 'False'
                  ? 'secondary'
                  : 'outline';
          }

          return (
            <Tooltip key={type}>
              <TooltipTrigger asChild>
                <Badge variant={variant} className="text-xs cursor-help">
                  {label}: {status}
                </Badge>
              </TooltipTrigger>
              {condition?.message && (
                <TooltipContent>
                  <p className="max-w-xs">{condition.message}</p>
                </TooltipContent>
              )}
            </Tooltip>
          );
        })}
      </div>
    </TooltipProvider>
  );
}
