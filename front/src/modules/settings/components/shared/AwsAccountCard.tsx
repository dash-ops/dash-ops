import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Eye, EyeOff, Trash2 } from 'lucide-react';
import type { CloudProviderFormValue } from '../../types';
import { SECRET_PLACEHOLDER } from '../../utils/formUtils';

interface AwsAccountCardProps {
  account: CloudProviderFormValue;
  onChange: (account: CloudProviderFormValue) => void;
  onRemove?: (accountId: string) => void;
  secretsVisible?: boolean;
  onToggleSecrets?: () => void;
  disabled?: boolean;
}

export function AwsAccountCard({
  account,
  onChange,
  onRemove,
  secretsVisible = false,
  onToggleSecrets,
  disabled = false,
}: AwsAccountCardProps): JSX.Element {
  const handleFieldChange = <T extends keyof CloudProviderFormValue>(
    field: T,
    value: CloudProviderFormValue[T]
  ) => {
    onChange({
      ...account,
      [field]: value,
    });
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="text-base font-medium">
            {account.name || 'AWS account'}
          </CardTitle>
          {onRemove && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => onRemove(account.id)}
              disabled={disabled}
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <Label>Account name</Label>
          <Input
            value={account.name}
            onChange={(event) =>
              handleFieldChange('name', event.target.value)
            }
            disabled={disabled}
          />
        </div>
        <div className="space-y-2">
          <Label>Region</Label>
          <Input
            value={account.region ?? ''}
            onChange={(event) =>
              handleFieldChange('region', event.target.value)
            }
            disabled={disabled}
          />
        </div>
        <div className="grid gap-4 md:grid-cols-2">
          <div className="space-y-2">
            <Label>Access Key ID</Label>
            <Input
              value={account.accessKeyId ?? ''}
              placeholder={
                account.accessKeyId
                  ? undefined
                  : account.accessKeyIdMasked ?? 'AKIA...'
              }
              onChange={(event) =>
                handleFieldChange('accessKeyId', event.target.value)
              }
              disabled={disabled}
            />
          </div>
          <div className="space-y-2">
            <Label>Secret Access Key</Label>
            <div className="flex items-center gap-2">
              <Input
                type={secretsVisible ? 'text' : 'password'}
                value={account.secretAccessKeyInput ?? ''}
                placeholder={
                  account.hasSecretAccessKey
                    ? `${SECRET_PLACEHOLDER} (stored)`
                    : 'Enter new secret'
                }
                onChange={(event) =>
                  handleFieldChange('secretAccessKeyInput', event.target.value)
                }
                disabled={disabled}
              />
              {onToggleSecrets && (
                <Button
                  variant="outline"
                  size="icon"
                  onClick={onToggleSecrets}
                >
                  {secretsVisible ? (
                    <EyeOff className="h-4 w-4" />
                  ) : (
                    <Eye className="h-4 w-4" />
                  )}
                </Button>
              )}
            </div>
            <p className="text-xs text-muted-foreground">
              Leave blank to keep the existing secret.
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

