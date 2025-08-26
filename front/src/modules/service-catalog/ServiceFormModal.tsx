import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '../../components/ui/dialog';
import { Button } from '../../components/ui/button';
import { Input } from '../../components/ui/input';
import { Label } from '../../components/ui/label';
import { Badge } from '../../components/ui/badge';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '../../components/ui/select';
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '../../components/ui/card';
import { toast } from 'sonner';
import { Plus, X, Info, Server } from 'lucide-react';
import { createService, updateService } from './serviceCatalogResource';
import type { Service, ServiceFormModalProps, ServiceFormData } from './types';

export function ServiceFormModal({
  open,
  onOpenChange,
  onServiceCreated,
  editingService,
}: ServiceFormModalProps) {
  const [isLoading, setIsLoading] = useState(false);
  const [customDependencies, setCustomDependencies] = useState<string[]>([]);
  const [newDependency, setNewDependency] = useState('');

  // Check if we're in edit mode
  const isEditing = !!editingService;

  // Helper function to convert Service to FormData for editing
  const serviceToFormData = (service: Service): Partial<ServiceFormData> => {
    const env = service.spec.kubernetes?.environments?.[0];
    const deployment = env?.resources?.deployments?.[0];

    return {
      name: service.metadata.name,
      description: service.spec.description,
      tier: service.metadata.tier,
      github_team: service.spec.team.github_team,
      impact: service.spec.business.impact || 'medium',
      sla_target: service.spec.business.sla_target || '',
      language: service.spec.technology?.language || '',
      framework: service.spec.technology?.framework || '',
      env_name: env?.name || 'local',
      env_context: env?.context || 'docker-desktop',
      env_namespace: env?.namespace || '',
      deployment_name: deployment?.name || '',
      deployment_replicas: deployment?.replicas || 3,
      cpu_request: deployment?.resources?.requests?.cpu || '100m',
      memory_request: deployment?.resources?.requests?.memory || '128Mi',
      cpu_limit: deployment?.resources?.limits?.cpu || '500m',
      memory_limit: deployment?.resources?.limits?.memory || '256Mi',
      metrics_url: service.spec.observability?.metrics || '',
      logs_url: service.spec.observability?.logs || '',
    };
  };

  const {
    register,
    handleSubmit,
    watch,
    setValue,
    reset,
    formState: { errors },
  } = useForm<ServiceFormData>({
    defaultValues:
      isEditing && editingService
        ? ({
            ...serviceToFormData(editingService),
          } as ServiceFormData)
        : {
            tier: 'TIER-2',
            impact: 'medium',
            sla_target: '99.9%',
            env_name: 'local',
            env_context: 'docker-desktop',
            env_namespace: '',
            deployment_name: '',
            deployment_replicas: 3,
            cpu_request: '100m',
            memory_request: '128Mi',
            cpu_limit: '500m',
            memory_limit: '256Mi',
          },
  });

  const watchedTier = watch('tier');
  const watchedName = watch('name');

  // Initialize dependencies when editing
  useEffect(() => {
    if (isEditing && editingService) {
      const deps = editingService.spec.business.dependencies || [];
      setCustomDependencies(deps);
    }
  }, [isEditing, editingService]);

  const onSubmit = async (data: ServiceFormData) => {
    if (
      !data.name ||
      !data.description ||
      !data.github_team ||
      !data.env_namespace
    ) {
      toast.error('Preencha todos os campos obrigatórios');
      return;
    }

    if (!data.deployment_name) {
      toast.error('Nome do deployment é obrigatório');
      return;
    }

    setIsLoading(true);

    try {
      // Transform form data to Service structure
      const service: Service = {
        apiVersion: 'v1',
        kind: 'Service',
        metadata: {
          name: data.name,
          tier: data.tier,
        },
        spec: {
          description: data.description,
          team: {
            github_team: data.github_team,
          },
          business: {
            impact: data.impact,
            ...(data.sla_target ? { sla_target: data.sla_target } : {}),
            ...(customDependencies.length > 0
              ? { dependencies: customDependencies }
              : {}),
          },
          ...(data.language || data.framework
            ? {
                technology: {
                  ...(data.language ? { language: data.language } : {}),
                  ...(data.framework ? { framework: data.framework } : {}),
                },
              }
            : {}),
          kubernetes: {
            environments: [
              {
                name: data.env_name,
                context: data.env_context,
                namespace: data.env_namespace,
                resources: {
                  deployments: [
                    {
                      name: data.deployment_name,
                      replicas: data.deployment_replicas,
                      resources: {
                        requests: {
                          cpu: data.cpu_request,
                          memory: data.memory_request,
                        },
                        limits: {
                          cpu: data.cpu_limit,
                          memory: data.memory_limit,
                        },
                      },
                    },
                  ],
                  services: [`${data.name}-svc`],
                },
              },
            ],
          },
          ...(data.metrics_url || data.logs_url
            ? {
                observability: {
                  ...(data.metrics_url ? { metrics: data.metrics_url } : {}),
                  ...(data.logs_url ? { logs: data.logs_url } : {}),
                },
              }
            : {}),
        },
      };

      console.log(
        isEditing ? 'Updating service:' : 'Creating service:',
        service
      );

      if (isEditing) {
        await updateService(data.name, service);
        toast.success(`Serviço '${data.name}' atualizado com sucesso!`);
      } else {
        await createService(service);
        toast.success(`Serviço '${data.name}' criado com sucesso!`);
      }

      handleClose();
      onServiceCreated?.();
    } catch (error: any) {
      console.error('Error creating service:', error);
      toast.error(error.response?.data?.error || 'Erro ao criar serviço');
    } finally {
      setIsLoading(false);
    }
  };

  const handleClose = () => {
    reset();
    setCustomDependencies([]);
    setNewDependency('');
    onOpenChange(false);
  };

  const addDependency = () => {
    if (
      newDependency.trim() &&
      !customDependencies.includes(newDependency.trim())
    ) {
      setCustomDependencies((prev) => [...prev, newDependency.trim()]);
      setNewDependency('');
    }
  };

  const removeDependency = (dependency: string) => {
    setCustomDependencies((prev) => prev.filter((d) => d !== dependency));
  };

  // Auto-fill namespace and deployment name based on service name
  const handleNameChange = (name: string) => {
    setValue('name', name);
    if (name) {
      if (!watch('env_namespace')) {
        setValue('env_namespace', name.toLowerCase());
      }
      if (!watch('deployment_name')) {
        setValue('deployment_name', `${name.toLowerCase()}-api`);
      }
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-3xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>
            {isEditing ? 'Editar Serviço' : 'Criar Novo Serviço'}
          </DialogTitle>
          <DialogDescription>
            {isEditing
              ? 'Atualize as configurações do seu serviço. Os campos obrigatórios são marcados com *.'
              : 'Configure seu novo serviço com integração ao Kubernetes. Os campos obrigatórios são marcados com *.'}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          {/* Basic Information */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Info className="h-5 w-5" />
                Informações Básicas
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="name">Nome do Serviço *</Label>
                  <Input
                    id="name"
                    placeholder="ex: user-authentication"
                    disabled={isEditing}
                    {...register('name', { required: 'Nome é obrigatório' })}
                    onChange={(e) => handleNameChange(e.target.value)}
                  />
                  {errors.name && (
                    <p className="text-sm text-destructive mt-1">
                      {errors.name.message}
                    </p>
                  )}
                </div>

                <div>
                  <Label htmlFor="tier">Tier do Serviço *</Label>
                  <Select
                    value={watchedTier}
                    onValueChange={(value) => setValue('tier', value as any)}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="TIER-1">TIER-1 - Crítico</SelectItem>
                      <SelectItem value="TIER-2">
                        TIER-2 - Importante
                      </SelectItem>
                      <SelectItem value="TIER-3">TIER-3 - Padrão</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>

              <div>
                <Label htmlFor="description">Descrição *</Label>
                <Input
                  id="description"
                  placeholder="Descreva o que esse serviço faz..."
                  {...register('description', {
                    required: 'Descrição é obrigatória',
                  })}
                />
                {errors.description && (
                  <p className="text-sm text-destructive mt-1">
                    {errors.description.message}
                  </p>
                )}
              </div>

              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <Label htmlFor="github_team">GitHub Team *</Label>
                  <Input
                    id="github_team"
                    placeholder="ex: auth-squad"
                    {...register('github_team', {
                      required: 'GitHub team é obrigatório',
                    })}
                  />
                  {errors.github_team && (
                    <p className="text-sm text-destructive mt-1">
                      {errors.github_team.message}
                    </p>
                  )}
                </div>

                <div>
                  <Label htmlFor="language">Linguagem</Label>
                  <Input
                    id="language"
                    placeholder="ex: Go"
                    {...register('language')}
                  />
                </div>

                <div>
                  <Label htmlFor="framework">Framework</Label>
                  <Input
                    id="framework"
                    placeholder="ex: gin-gonic"
                    {...register('framework')}
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="impact">Impacto de Negócio</Label>
                  <Select
                    value={watch('impact')}
                    onValueChange={(value) => setValue('impact', value as any)}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="high">Alto</SelectItem>
                      <SelectItem value="medium">Médio</SelectItem>
                      <SelectItem value="low">Baixo</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div>
                  <Label htmlFor="sla_target">SLA Target</Label>
                  <Input
                    id="sla_target"
                    placeholder="ex: 99.9%"
                    {...register('sla_target')}
                  />
                </div>
              </div>

              {/* Dependencies */}
              <div>
                <Label>Dependências</Label>
                <div className="flex gap-2 mb-3">
                  <Input
                    placeholder="ex: user-database"
                    value={newDependency}
                    onChange={(e) => setNewDependency(e.target.value)}
                    onKeyPress={(e) =>
                      e.key === 'Enter' && (e.preventDefault(), addDependency())
                    }
                  />
                  <Button type="button" onClick={addDependency} size="sm">
                    <Plus className="h-4 w-4" />
                  </Button>
                </div>

                <div className="flex flex-wrap gap-2">
                  {customDependencies.map((dep) => (
                    <Badge key={dep} variant="secondary" className="gap-1">
                      {dep}
                      <button
                        type="button"
                        onClick={() => removeDependency(dep)}
                        className="ml-1 hover:text-destructive"
                      >
                        <X className="h-3 w-3" />
                      </button>
                    </Badge>
                  ))}
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Kubernetes Configuration */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Server className="h-5 w-5" />
                Configuração Kubernetes
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <Label htmlFor="env_name">Nome do Ambiente</Label>
                  <Input
                    id="env_name"
                    placeholder="ex: local"
                    {...register('env_name')}
                  />
                </div>

                <div>
                  <Label htmlFor="env_context">Context K8s *</Label>
                  <Input
                    id="env_context"
                    placeholder="ex: docker-desktop"
                    {...register('env_context', {
                      required: 'Context é obrigatório',
                    })}
                  />
                  {errors.env_context && (
                    <p className="text-sm text-destructive mt-1">
                      {errors.env_context.message}
                    </p>
                  )}
                </div>

                <div>
                  <Label htmlFor="env_namespace">Namespace *</Label>
                  <Input
                    id="env_namespace"
                    placeholder="Preenchido automaticamente"
                    {...register('env_namespace', {
                      required: 'Namespace é obrigatório',
                    })}
                  />
                  {errors.env_namespace && (
                    <p className="text-sm text-destructive mt-1">
                      {errors.env_namespace.message}
                    </p>
                  )}
                </div>
              </div>

              <div className="space-y-4">
                <h4 className="font-medium">Deployment Principal</h4>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <Label htmlFor="deployment_name">
                      Nome do Deployment *
                    </Label>
                    <Input
                      id="deployment_name"
                      placeholder={`ex: ${watchedName || 'app'}-api`}
                      {...register('deployment_name', {
                        required: 'Nome do deployment é obrigatório',
                      })}
                    />
                    {errors.deployment_name && (
                      <p className="text-sm text-destructive mt-1">
                        {errors.deployment_name.message}
                      </p>
                    )}
                  </div>

                  <div>
                    <Label htmlFor="deployment_replicas">Réplicas</Label>
                    <Input
                      id="deployment_replicas"
                      type="number"
                      min="1"
                      max="10"
                      {...register('deployment_replicas', {
                        valueAsNumber: true,
                        min: 1,
                        max: 10,
                      })}
                    />
                  </div>
                </div>

                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                  <div>
                    <Label htmlFor="cpu_request">CPU Request</Label>
                    <Input
                      id="cpu_request"
                      placeholder="ex: 100m"
                      {...register('cpu_request')}
                    />
                  </div>

                  <div>
                    <Label htmlFor="memory_request">Memory Request</Label>
                    <Input
                      id="memory_request"
                      placeholder="ex: 128Mi"
                      {...register('memory_request')}
                    />
                  </div>

                  <div>
                    <Label htmlFor="cpu_limit">CPU Limit</Label>
                    <Input
                      id="cpu_limit"
                      placeholder="ex: 500m"
                      {...register('cpu_limit')}
                    />
                  </div>

                  <div>
                    <Label htmlFor="memory_limit">Memory Limit</Label>
                    <Input
                      id="memory_limit"
                      placeholder="ex: 256Mi"
                      {...register('memory_limit')}
                    />
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Observability */}
          <Card>
            <CardHeader>
              <CardTitle>Observabilidade (Opcional)</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="metrics_url">URL de Métricas</Label>
                  <Input
                    id="metrics_url"
                    placeholder="ex: https://grafana.company.com/d/service"
                    {...register('metrics_url')}
                  />
                </div>

                <div>
                  <Label htmlFor="logs_url">URL de Logs</Label>
                  <Input
                    id="logs_url"
                    placeholder="ex: https://kibana.company.com/app"
                    {...register('logs_url')}
                  />
                </div>
              </div>
            </CardContent>
          </Card>

          <div className="flex justify-end gap-2 pt-6 border-t">
            <Button
              type="button"
              variant="outline"
              onClick={handleClose}
              disabled={isLoading}
            >
              Cancelar
            </Button>
            <Button type="submit" disabled={isLoading}>
              {isLoading
                ? isEditing
                  ? 'Salvando...'
                  : 'Criando...'
                : isEditing
                  ? 'Salvar Alterações'
                  : 'Criar Serviço'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
