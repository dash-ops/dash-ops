import { Progress } from '@/components/ui/progress';
import { ProgressDataProps } from '@/types';

function ProgressData({ percent }: ProgressDataProps): JSX.Element {
  return (
    <div className="w-full max-w-[170px]">
      <Progress value={percent} className="h-2" />
      <span className="text-xs text-muted-foreground">
        {percent?.toFixed(1)}%
      </span>
    </div>
  );
}

export default ProgressData;
