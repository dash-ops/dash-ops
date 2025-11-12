import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Badge } from '@/components/ui/badge';
import { Eye, EyeOff, Github, Trash2 } from 'lucide-react';
import type { AuthProviderFormValue } from '../../types';
import { SECRET_PLACEHOLDER } from '../../utils/formUtils';

interface AuthProviderCardProps {
  provider: AuthProviderFormValue;
  onChange: (provider: AuthProviderFormValue) => void;
  onRemove?: (providerId: string) => void;
  secretsVisible?: boolean;
  onToggleSecrets?: () => void;
  disabled?: boolean;
  showStatusBadge?: boolean;
}

export function AuthProviderCard({
  provider,
  onChange,
  onRemove,
  secretsVisible = false,
  onToggleSecrets,
  disabled = false,
  showStatusBadge = true,
}: AuthProviderCardProps): JSX.Element {
  const handleFieldChange = <T extends keyof AuthProviderFormValue>(
    field: T,
    value: AuthProviderFormValue[T]
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
          <div className="flex items-center gap-2">
            <Github className="h-4 w-4 text-primary" />
            <CardTitle className="text-base">{provider.name}</CardTitle>
          </div>
          <div className="flex items-center gap-2">
            {showStatusBadge && (
              <Badge variant="secondary">Active</Badge>
            )}
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
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-4 md:grid-cols-2">
          <div className="space-y-2">
            <Label>Client ID</Label>
            <Input
              value={provider.clientId}
              placeholder={
                provider.clientId
                  ? undefined
                  : provider.clientIdMasked ?? 'Enter client ID'
              }
              onChange={(event) =>
                handleFieldChange('clientId', event.target.value)
              }
              disabled={disabled}
            />
          </div>
          <div className="space-y-2">
            <Label>Client secret</Label>
            <div className="flex items-center gap-2">
              <Input
                type={secretsVisible ? 'text' : 'password'}
                value={provider.clientSecretInput ?? ''}
                placeholder={
                  provider.hasClientSecret
                    ? `${SECRET_PLACEHOLDER} (stored)`
                    : 'Enter new secret'
                }
                onChange={(event) =>
                  handleFieldChange('clientSecretInput', event.target.value)
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
        <div className="space-y-2">
          <Label>Organization permission</Label>
          <Input
            value={provider.orgPermission ?? ''}
            onChange={(event) =>
              handleFieldChange('orgPermission', event.target.value)
            }
            disabled={disabled}
          />
        </div>
      </CardContent>
    </Card>
  );
}

