import { toast } from 'sonner';
import { Button } from '@/components/ui/button';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { Monitor, Play, Square } from 'lucide-react';
import { AWSTypes } from '@/types';

function ssh(instance: AWSTypes.Instance): void {
  if (instance.platform === 'windows') {
    toast.error(`Sorry... I'm afraid I can't do that...`, {
      description: `Windows does not provides a method to connect a Remote Desktop via URL. You can try to connect via command line using on Windows: mstsc /v:${instance.name}`,
    });
    return;
  }
  window.location.href = `ssh://${instance.name}`;
}

function InstanceActions({
  instance,
  toStart,
  toStop,
}: AWSTypes.InstanceActionsProps): JSX.Element {
  const showPlayButton = instance.state !== 'running';

  return (
    <TooltipProvider>
      <div className="flex gap-1">
        <Tooltip>
          <TooltipTrigger asChild>
            <Button variant="outline" size="sm" onClick={() => ssh(instance)}>
              <Monitor className="h-4 w-4" />
            </Button>
          </TooltipTrigger>
          <TooltipContent>
            <p>SSH access</p>
          </TooltipContent>
        </Tooltip>

        {showPlayButton && (
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant="outline"
                size="sm"
                disabled={instance.state !== 'stopped'}
                onClick={toStart}
                className="gap-1"
              >
                <Play className="h-4 w-4" />
                Start
              </Button>
            </TooltipTrigger>
            <TooltipContent>
              <p>Start instance</p>
            </TooltipContent>
          </Tooltip>
        )}

        {!showPlayButton && (
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant="destructive"
                size="sm"
                disabled={showPlayButton}
                onClick={toStop}
                className="gap-1"
              >
                <Square className="h-4 w-4" />
                Stop
              </Button>
            </TooltipTrigger>
            <TooltipContent>
              <p>Stop instance</p>
            </TooltipContent>
          </Tooltip>
        )}
      </div>
    </TooltipProvider>
  );
}

export default InstanceActions;
