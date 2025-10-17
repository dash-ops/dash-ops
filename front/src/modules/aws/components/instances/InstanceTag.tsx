import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';
import { AWSTypes, InstanceState } from '@/types';

const instanceStyles: Record<InstanceState, string> = {
  pending: 'bg-yellow-100 text-yellow-800 border-yellow-300',
  running: 'bg-green-100 text-green-800 border-green-300',
  'shutting-down': 'bg-orange-100 text-orange-800 border-orange-300',
  terminated: 'bg-red-100 text-red-800 border-red-300',
  stopping: 'bg-purple-100 text-purple-800 border-purple-300',
  stopped: 'bg-red-100 text-red-800 border-red-300',
  loading: 'bg-blue-100 text-blue-800 border-blue-300',
};

function InstanceTag({ state }: AWSTypes.InstanceTagProps): JSX.Element {
  // Handle both string and object state formats
  const stateName = typeof state === 'string' ? state : state.name;
  const stateCode = typeof state === 'string' ? undefined : state.code;
  
  return (
    <Badge
      variant="outline"
      className={cn('capitalize', instanceStyles[stateName as InstanceState])}
    >
      {stateName}
    </Badge>
  );
}

export default InstanceTag;
