import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import { RefreshCw } from 'lucide-react';
import { setItem, getItem } from '../helpers/localStorage';
import { useInterval } from '@/hooks/use-interval';
import { RefreshProps } from '@/types';

function Refresh({ onReload }: RefreshProps): JSX.Element {
  const initState = getItem('auto-refresh') === 'true';
  const [isAutoRefresh, setIsAutoRefresh] = useState<boolean>(initState);

  useInterval(
    async () => {
      onReload();
    },
    isAutoRefresh ? 10000 : null
  );

  const onAutoRefresh = (checked: boolean) => {
    setIsAutoRefresh(checked);
    setItem('auto-refresh', checked.toString());
  };

  return (
    <div className="flex items-center gap-2">
      <div className="flex items-center space-x-2">
        <Checkbox
          id="auto-refresh"
          checked={isAutoRefresh}
          onCheckedChange={onAutoRefresh}
        />
        <label
          htmlFor="auto-refresh"
          className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
        >
          Auto Refresh 10s
        </label>
      </div>
      <Button variant="outline" size="sm" onClick={onReload}>
        <RefreshCw className="h-4 w-4" />
      </Button>
    </div>
  );
}

export default Refresh;
