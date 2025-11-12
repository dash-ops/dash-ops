import { Switch } from '@/components/ui/switch';
import { cn } from '@/lib/utils';
import type { PluginOption } from '../../constants';

interface PluginToggleGridProps {
  options: PluginOption[];
  selected: string[];
  onToggle: (pluginId: string) => void;
  disabled?: boolean;
}

export function PluginToggleGrid({
  options,
  selected,
  onToggle,
  disabled = false,
}: PluginToggleGridProps): JSX.Element {
  return (
    <div className="grid gap-4 md:grid-cols-2">
      {options.map((plugin) => {
        const Icon = plugin.icon;
        const enabled = selected.includes(plugin.id);

        return (
          <div
            key={plugin.id}
            className={cn(
              'flex items-start justify-between rounded-lg border p-4 transition-all',
              enabled ? 'border-primary bg-primary/5' : 'hover:border-primary/40'
            )}
          >
            <div className="space-y-1">
              <div className="flex items-center gap-2">
                <Icon className="h-4 w-4 text-primary" />
                <p className="text-sm font-medium">{plugin.label}</p>
              </div>
              <p className="text-xs text-muted-foreground">
                {plugin.description}
              </p>
            </div>
            <Switch
              checked={enabled}
              onCheckedChange={() => onToggle(plugin.id)}
              disabled={disabled}
            />
          </div>
        );
      })}
    </div>
  );
}

