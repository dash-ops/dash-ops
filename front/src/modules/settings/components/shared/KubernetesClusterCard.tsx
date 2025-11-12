import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { cn } from '@/lib/utils';
import { Globe, Server, Trash2 } from 'lucide-react';
import type { KubernetesClusterFormValue } from '../../types';
import { SECRET_PLACEHOLDER } from '../../utils/formUtils';

interface KubernetesClusterCardProps {
  cluster: KubernetesClusterFormValue;
  onChange: (cluster: KubernetesClusterFormValue) => void;
  onRemove?: (clusterId: string) => void;
  disabled?: boolean;
}

export function KubernetesClusterCard({
  cluster,
  onChange,
  onRemove,
  disabled = false,
}: KubernetesClusterCardProps): JSX.Element {
  const handleFieldChange = <T extends keyof KubernetesClusterFormValue>(
    field: T,
    value: KubernetesClusterFormValue[T]
  ) => {
    onChange({
      ...cluster,
      [field]: value,
    });
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="text-base font-medium">
            {cluster.name || 'Kubernetes cluster'}
          </CardTitle>
          {onRemove && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => onRemove(cluster.id)}
              disabled={disabled}
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          )}
        </div>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="grid gap-4 md:grid-cols-2">
          <div className="space-y-2">
            <Label>Cluster name</Label>
            <Input
              value={cluster.name}
              onChange={(event) =>
                handleFieldChange('name', event.target.value)
              }
              disabled={disabled}
            />
          </div>
          <div className="space-y-2">
            <Label>Connection type</Label>
            <div className="grid grid-cols-2 gap-3">
              <button
                type="button"
                onClick={() =>
                  handleFieldChange('connectionType', 'kubeconfig')
                }
                className={cn(
                  'rounded-lg border p-4 text-left text-sm transition-all',
                  cluster.connectionType === 'kubeconfig'
                    ? 'border-primary bg-primary/5 text-primary'
                    : 'border-border hover:border-primary/40'
                )}
                disabled={disabled}
              >
                <Server className="mb-2 h-4 w-4" />
                Kubeconfig file
              </button>
              <button
                type="button"
                onClick={() => handleFieldChange('connectionType', 'remote')}
                className={cn(
                  'rounded-lg border p-4 text-left text-sm transition-all',
                  cluster.connectionType === 'remote'
                    ? 'border-primary bg-primary/5 text-primary'
                    : 'border-border hover:border-primary/40'
                )}
                disabled={disabled}
              >
                <Globe className="mb-2 h-4 w-4" />
                Remote API
              </button>
            </div>
          </div>
        </div>

        {cluster.connectionType === 'kubeconfig' && (
          <div className="grid gap-4 md:grid-cols-2">
            <div className="space-y-2">
              <Label>Kubeconfig path</Label>
              <Input
                value={cluster.kubeconfig ?? ''}
                onChange={(event) =>
                  handleFieldChange('kubeconfig', event.target.value)
                }
                disabled={disabled}
              />
            </div>
            <div className="space-y-2">
              <Label>Context</Label>
              <Input
                value={cluster.context ?? ''}
                onChange={(event) =>
                  handleFieldChange('context', event.target.value)
                }
                disabled={disabled}
              />
            </div>
          </div>
        )}

        {cluster.connectionType === 'remote' && (
          <div className="space-y-4">
            <div className="space-y-2">
              <Label>API server host</Label>
              <Input
                value={cluster.host ?? ''}
                onChange={(event) =>
                  handleFieldChange('host', event.target.value)
                }
                disabled={disabled}
                placeholder="https://k8s.example.com:6443"
              />
            </div>
            <div className="space-y-2">
              <Label>Bearer token</Label>
              <Textarea
                rows={4}
                value={cluster.token ?? ''}
                placeholder={
                  cluster.hasToken ? `${SECRET_PLACEHOLDER} (stored)` : undefined
                }
                onChange={(event) =>
                  handleFieldChange('token', event.target.value)
                }
                disabled={disabled}
                className="font-mono text-xs"
              />
            </div>
            <div className="space-y-2">
              <Label>CA certificate (optional)</Label>
              <Textarea
                rows={4}
                value={cluster.certificate ?? ''}
                onChange={(event) =>
                  handleFieldChange('certificate', event.target.value)
                }
                disabled={disabled}
                className="font-mono text-xs"
              />
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}

