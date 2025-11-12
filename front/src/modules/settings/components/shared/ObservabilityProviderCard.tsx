import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { Trash2 } from 'lucide-react';
import type { ObservabilityProviderFormValue } from '../../types';

interface ObservabilityProviderCardProps {
  provider: ObservabilityProviderFormValue;
  onChange: (provider: ObservabilityProviderFormValue) => void;
  onRemove?: (providerId: string) => void;
  disabled?: boolean;
  titlePrefix: string;
}

export function ObservabilityProviderCard({
  provider,
  onChange,
  onRemove,
  disabled = false,
  titlePrefix,
}: ObservabilityProviderCardProps): JSX.Element {
  const handleFieldChange = <T extends keyof ObservabilityProviderFormValue>(
    field: T,
    value: ObservabilityProviderFormValue[T]
  ) => {
    onChange({
      ...provider,
      [field]: value,
    });
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="text-base font-medium">
            {titlePrefix} - {provider.name || provider.id}
          </CardTitle>
          {onRemove && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => onRemove(provider.id)}
              disabled={disabled}
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <Label>Name</Label>
          <Input
            value={provider.name}
            onChange={(event) =>
              handleFieldChange('name', event.target.value)
            }
            disabled={disabled}
          />
        </div>
        <div className="grid gap-4 md:grid-cols-2">
          <div className="space-y-2">
            <Label>Type</Label>
            <Input
              value={provider.type}
              onChange={(event) =>
                handleFieldChange('type', event.target.value)
              }
              disabled={disabled}
            />
          </div>
          <div className="space-y-2">
            <Label>URL</Label>
            <Input
              value={provider.url}
              onChange={(event) =>
                handleFieldChange('url', event.target.value)
              }
              disabled={disabled}
            />
          </div>
        </div>
        <div className="grid gap-4 md:grid-cols-2">
          <div className="space-y-2">
            <Label>Timeout</Label>
            <Input
              value={provider.timeout ?? ''}
              onChange={(event) =>
                handleFieldChange('timeout', event.target.value)
              }
              disabled={disabled}
            />
          </div>
          <div className="space-y-2">
            <Label>Retention</Label>
            <Input
              value={provider.retention ?? ''}
              onChange={(event) =>
                handleFieldChange('retention', event.target.value)
              }
              disabled={disabled}
            />
          </div>
        </div>
        <div className="flex items-center justify-between rounded-lg border p-3">
          <div>
            <p className="text-sm font-medium">Enabled</p>
            <p className="text-xs text-muted-foreground">
              Toggle to activate this provider.
            </p>
          </div>
          <Switch
            checked={provider.enabled ?? true}
            onCheckedChange={(checked) =>
              handleFieldChange('enabled', checked)
            }
            disabled={disabled}
          />
        </div>
      </CardContent>
    </Card>
  );
}

