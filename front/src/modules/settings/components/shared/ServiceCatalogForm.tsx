import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import type { ServiceCatalogFormValue } from '../../types';
import { cn } from '@/lib/utils';
import { Database } from 'lucide-react';

interface ServiceCatalogFormProps {
  value: ServiceCatalogFormValue;
  onChange: (value: ServiceCatalogFormValue) => void;
  disabled?: boolean;
}

export function ServiceCatalogForm({
  value,
  onChange,
  disabled = false,
}: ServiceCatalogFormProps): JSX.Element {
  const handleChange = <T extends keyof ServiceCatalogFormValue>(
    field: T,
    fieldValue: ServiceCatalogFormValue[T]
  ) => {
    onChange({
      ...value,
      [field]: fieldValue,
    });
  };

  const handleProviderChange = (provider: ServiceCatalogFormValue['storageProvider']) => {
    onChange({
      ...value,
      storageProvider: provider,
    });
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base font-semibold">
          Service Catalog storage
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="space-y-2">
          <Label>Storage provider</Label>
          <div className="grid grid-cols-3 gap-3">
            {(['filesystem', 'github', 's3'] as const).map((provider) => (
              <button
                key={provider}
                type="button"
                onClick={() => handleProviderChange(provider)}
                className={cn(
                  'rounded-lg border p-4 text-center text-sm capitalize transition-all',
                  value.storageProvider === provider
                    ? 'border-primary bg-primary/5 text-primary'
                    : 'border-border hover:border-primary/40'
                )}
                disabled={disabled}
              >
                <Database className="mx-auto mb-2 h-5 w-5" />
                {provider}
              </button>
            ))}
          </div>
        </div>

        {value.storageProvider === 'filesystem' && (
          <div className="space-y-2">
            <Label>Directory path</Label>
            <Input
              value={value.directory}
              onChange={(event) =>
                handleChange('directory', event.target.value)
              }
              disabled={disabled}
              placeholder="../services"
            />
          </div>
        )}

        {value.storageProvider === 'github' && (
          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label>Repository</Label>
              <Input
                value={value.githubRepository}
                onChange={(event) =>
                  handleChange('githubRepository', event.target.value)
                }
                disabled={disabled}
                placeholder="your-org/service-definitions"
              />
            </div>
            <div className="space-y-2">
              <Label>Branch</Label>
              <Input
                value={value.githubBranch}
                onChange={(event) =>
                  handleChange('githubBranch', event.target.value)
                }
                disabled={disabled}
                placeholder="main"
              />
            </div>
          </div>
        )}

        {value.storageProvider === 's3' && (
          <div className="space-y-2">
            <Label>S3 bucket</Label>
            <Input
              value={value.s3Bucket}
              onChange={(event) =>
                handleChange('s3Bucket', event.target.value)
              }
              disabled={disabled}
              placeholder="company-service-definitions"
            />
          </div>
        )}

        <div className="flex items-center justify-between rounded-lg border p-4">
          <div>
            <p className="text-sm font-medium">Enable versioning</p>
            <p className="text-xs text-muted-foreground">
              Track changes to service definitions.
            </p>
          </div>
          <Switch
            checked={value.versioningEnabled}
            onCheckedChange={(checked) =>
              handleChange('versioningEnabled', checked)
            }
            disabled={disabled}
          />
        </div>
      </CardContent>
    </Card>
  );
}

