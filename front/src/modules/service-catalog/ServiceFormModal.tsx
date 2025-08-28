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
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from '../../components/ui/tabs';
import { ScrollArea } from '../../components/ui/scroll-area';
import { Separator } from '../../components/ui/separator';
import { toast } from 'sonner';
import { Plus, X, Info, Server, Activity, Settings } from 'lucide-react';
import {
  createService,
  updateService,
  getService,
} from './serviceCatalogResource';
import type { Service, ServiceFormModalProps, ServiceFormData } from './types';

export function ServiceFormModal({
  open,
  onOpenChange,
  onServiceCreated,
  editingServiceName,
}: ServiceFormModalProps) {
  const [isLoading, setIsLoading] = useState(false);
  const [isLoadingService, setIsLoadingService] = useState(false);
  const [serviceData, setServiceData] = useState<Service | null>(null);
  const [customDependencies, setCustomDependencies] = useState<string[]>([]);
  const [newDependency, setNewDependency] = useState('');
  const [activeTab, setActiveTab] = useState('basic');

  // Check if we're in edit mode
  const isEditing = !!editingServiceName;

  // Helper function to convert Service to FormData for editing
  const serviceToFormData = (service: Service): ServiceFormData => {
    const env = service.spec.kubernetes?.environments?.[0];
    const deployment = env?.resources?.deployments?.[0];

    const formData: ServiceFormData = {
      name: service.metadata.name || '',
      description: service.spec.description || '',
      tier: service.metadata.tier || 'TIER-2',
      github_team: service.spec.team?.github_team || '',
      impact: service.spec.business?.impact || 'medium',
      sla_target: service.spec.business?.sla_target || '',
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

    return formData;
  };

  // Default form values - always start with defaults, then use setValue to populate
  const defaultFormValues: ServiceFormData = {
    name: '',
    description: '',
    tier: 'TIER-2',
    github_team: '',
    impact: 'medium',
    sla_target: '99.9%',
    language: '',
    framework: '',
    env_name: 'local',
    env_context: 'docker-desktop',
    env_namespace: '',
    deployment_name: '',
    deployment_replicas: 3,
    cpu_request: '100m',
    memory_request: '128Mi',
    cpu_limit: '500m',
    memory_limit: '256Mi',
    metrics_url: '',
    logs_url: '',
  };

  const {
    register,
    handleSubmit,
    watch,
    setValue,
    reset,
    formState: { errors },
  } = useForm<ServiceFormData>({
    // defaultValues: defaultFormValues,
  });

  const watchedTier = watch('tier');
  const watchedName = watch('name');

  // Load service data when editing
  useEffect(() => {
    const loadServiceData = async () => {
      if (open && isEditing && editingServiceName) {
        setIsLoadingService(true);
        try {
          const loadedServiceData = await getService(editingServiceName);

          // Store service data in state
          setServiceData(loadedServiceData);
          reset(serviceToFormData(loadedServiceData));

          // Set dependencies
          const deps = loadedServiceData.spec.business?.dependencies || [];
          setCustomDependencies(deps);
        } catch (error) {
          console.error('Failed to load service data:', error);
          toast.error(`Failed to load service data: ${editingServiceName}`);
          onOpenChange(false);
        } finally {
          setIsLoadingService(false);
        }
      } else if (open && !isEditing) {
        // Clear service data and dependencies for new service
        setServiceData(null);
        setCustomDependencies([]);
      }

      // Always reset active tab to basic when modal opens
      if (open) {
        setActiveTab('basic');
      } else {
        // Reset state when modal closes
        setServiceData(null);
        setCustomDependencies([]);
      }
    };

    loadServiceData();
  }, [open, isEditing, editingServiceName, onOpenChange, reset]);

  // Set form values when service data is loaded
  // useEffect(() => {
  //   if (serviceData && isEditing) {
  //     const formData = serviceToFormData(serviceData);
  //     console.log('🔄 Setting form values:', formData);

  //     // Set each field individually using setValue
  //     setValue('name', formData.name);
  //     setValue('description', formData.description);
  //     setValue('tier', formData.tier);
  //     setValue('github_team', formData.github_team);
  //     setValue('impact', formData.impact);
  //     setValue('sla_target', formData.sla_target);
  //     setValue('language', formData.language);
  //     setValue('framework', formData.framework);
  //     setValue('env_name', formData.env_name);
  //     setValue('env_context', formData.env_context);
  //     setValue('env_namespace', formData.env_namespace);
  //     setValue('deployment_name', formData.deployment_name);
  //     setValue('deployment_replicas', formData.deployment_replicas);
  //     setValue('cpu_request', formData.cpu_request);
  //     setValue('memory_request', formData.memory_request);
  //     setValue('cpu_limit', formData.cpu_limit);
  //     setValue('memory_limit', formData.memory_limit);
  //     setValue('metrics_url', formData.metrics_url);
  //     setValue('logs_url', formData.logs_url);

  //     console.log('✅ Form values set, current form state:', watch());
  //   }
  // }, [serviceData, isEditing, setValue, watch]);

  // Auto-fill namespace and deployment name based on service name
  useEffect(() => {
    const currentName = watchedName;
    if (currentName && !isEditing) {
      // Auto-fill namespace if empty
      if (!watch('env_namespace')) {
        setValue('env_namespace', currentName.toLowerCase());
      }
      // Auto-fill deployment name if empty
      if (!watch('deployment_name')) {
        setValue('deployment_name', `${currentName.toLowerCase()}-api`);
      }
    }
  }, [watchedName, isEditing, watch, setValue]);

  const onSubmit = async (data: ServiceFormData) => {
    if (
      !data.name ||
      !data.description ||
      !data.github_team ||
      !data.env_namespace
    ) {
      toast.error('Fill in all required fields');
      return;
    }

    if (!data.deployment_name) {
      toast.error('Deployment name is required');
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

      if (isEditing && editingServiceName) {
        await updateService(editingServiceName, service);
        toast.success(`Service '${editingServiceName}' updated successfully!`);
      } else {
        await createService(service);
        toast.success(`Service '${data.name}' created successfully!`);
      }

      handleClose();
      onServiceCreated?.();
    } catch (error: unknown) {
      console.error('Error creating service:', error);
      toast.error(
        (error as { response?: { data?: { error?: string } } })?.response?.data
          ?.error || 'Failed to create service'
      );
    } finally {
      setIsLoading(false);
    }
  };

  const handleClose = () => {
    // Reset form to default values
    reset(defaultFormValues);
    setCustomDependencies([]);
    setNewDependency('');
    setActiveTab('basic');
    setServiceData(null);
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

  // Don't render modal at all when closed to prevent unnecessary computations
  if (!open) {
    return null;
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-4xl max-h-[90vh] overflow-hidden flex flex-col">
        <DialogHeader className="flex-none">
          <DialogTitle>
            {isEditing ? 'Edit Service' : 'Create New Service'}
          </DialogTitle>
          <DialogDescription>
            {isEditing
              ? 'Update your service configuration. Required fields are marked with *.'
              : 'Configure your new service with Kubernetes integration. Required fields are marked with *.'}
          </DialogDescription>
        </DialogHeader>

        {isLoadingService || (isEditing && !serviceData) ? (
          <div className="flex-1 flex items-center justify-center">
            <div className="text-center">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4" />
              <p className="text-muted-foreground">
                {isEditing ? 'Loading service data...' : 'Preparing form...'}
              </p>
            </div>
          </div>
        ) : (
          <form
            onSubmit={handleSubmit(onSubmit)}
            className="flex-1 flex flex-col overflow-hidden"
          >
            <Tabs
              value={activeTab}
              onValueChange={setActiveTab}
              className="flex-1 flex flex-col overflow-hidden"
            >
              <TabsList className="grid w-full grid-cols-4 flex-none">
                <TabsTrigger value="basic">Basic Info</TabsTrigger>
                <TabsTrigger value="infrastructure">Kubernetes</TabsTrigger>
                <TabsTrigger value="observability">Observability</TabsTrigger>
                <TabsTrigger value="review">Review</TabsTrigger>
              </TabsList>

              <ScrollArea className="flex-1 mt-4">
                <div className="space-y-6 px-1">
                  {/* Basic Information Tab */}
                  <TabsContent value="basic" className="space-y-6 m-0">
                    <Card>
                      <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                          <Info className="h-5 w-5" />
                          Basic Information
                        </CardTitle>
                      </CardHeader>
                      <CardContent className="space-y-4">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                          <div>
                            <Label htmlFor="name">Service Name *</Label>
                            <Input
                              id="name"
                              placeholder="ex: user-authentication"
                              readOnly={isEditing}
                              className={isEditing ? 'bg-muted' : ''}
                              {...register('name', {
                                required: 'Service name is required',
                              })}
                            />
                            {errors.name && (
                              <p className="text-sm text-destructive mt-1">
                                {errors.name.message}
                              </p>
                            )}
                          </div>

                          <div>
                            <Label htmlFor="tier">Service Tier *</Label>
                            <Select
                              value={watchedTier}
                              onValueChange={(value) =>
                                setValue(
                                  'tier',
                                  value as 'TIER-1' | 'TIER-2' | 'TIER-3'
                                )
                              }
                            >
                              <SelectTrigger>
                                <SelectValue />
                              </SelectTrigger>
                              <SelectContent>
                                <SelectItem value="TIER-1">
                                  TIER-1 - Critical
                                </SelectItem>
                                <SelectItem value="TIER-2">
                                  TIER-2 - Important
                                </SelectItem>
                                <SelectItem value="TIER-3">
                                  TIER-3 - Standard
                                </SelectItem>
                              </SelectContent>
                            </Select>
                          </div>
                        </div>

                        <div>
                          <Label htmlFor="description">Description *</Label>
                          <Input
                            id="description"
                            placeholder="Describe what this service does..."
                            {...register('description', {
                              required: 'Description is required',
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
                                required: 'GitHub team is required',
                              })}
                            />
                            {errors.github_team && (
                              <p className="text-sm text-destructive mt-1">
                                {errors.github_team.message}
                              </p>
                            )}
                          </div>

                          <div>
                            <Label htmlFor="language">Language</Label>
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
                            <Label htmlFor="impact">Business Impact</Label>
                            <Select
                              value={watch('impact')}
                              onValueChange={(value) =>
                                setValue(
                                  'impact',
                                  value as 'high' | 'medium' | 'low'
                                )
                              }
                            >
                              <SelectTrigger>
                                <SelectValue />
                              </SelectTrigger>
                              <SelectContent>
                                <SelectItem value="high">High</SelectItem>
                                <SelectItem value="medium">Medium</SelectItem>
                                <SelectItem value="low">Low</SelectItem>
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
                          <Label>Dependencies</Label>
                          <div className="flex gap-2 mb-3">
                            <Input
                              placeholder="ex: user-database"
                              value={newDependency}
                              onChange={(e) => setNewDependency(e.target.value)}
                              onKeyPress={(e) =>
                                e.key === 'Enter' &&
                                (e.preventDefault(), addDependency())
                              }
                            />
                            <Button
                              type="button"
                              onClick={addDependency}
                              size="sm"
                            >
                              <Plus className="h-4 w-4" />
                            </Button>
                          </div>

                          <div className="flex flex-wrap gap-2">
                            {customDependencies.map((dep) => (
                              <Badge
                                key={dep}
                                variant="secondary"
                                className="gap-1"
                              >
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
                  </TabsContent>

                  {/* Infrastructure Tab */}
                  <TabsContent value="infrastructure" className="space-y-6 m-0">
                    <Card>
                      <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                          <Server className="h-5 w-5" />
                          Kubernetes Configuration
                        </CardTitle>
                      </CardHeader>
                      <CardContent className="space-y-4">
                        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                          <div>
                            <Label htmlFor="env_name">Environment Name</Label>
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
                                required: 'Context is required',
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
                              placeholder="Automatically filled"
                              {...register('env_namespace', {
                                required: 'Namespace is required',
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
                          <h4 className="font-medium">Main Deployment</h4>

                          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div>
                              <Label htmlFor="deployment_name">
                                Deployment Name *
                              </Label>
                              <Input
                                id="deployment_name"
                                placeholder={`ex: ${watchedName || 'app'}-api`}
                                {...register('deployment_name', {
                                  required: 'Deployment name is required',
                                })}
                              />
                              {errors.deployment_name && (
                                <p className="text-sm text-destructive mt-1">
                                  {errors.deployment_name.message}
                                </p>
                              )}
                            </div>

                            <div>
                              <Label htmlFor="deployment_replicas">
                                Replicas
                              </Label>
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
                              <Label htmlFor="memory_request">
                                Memory Request
                              </Label>
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
                  </TabsContent>

                  {/* Observability Tab */}
                  <TabsContent value="observability" className="space-y-6 m-0">
                    <Card>
                      <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                          <Activity className="h-5 w-5" />
                          Observability (Optional)
                        </CardTitle>
                      </CardHeader>
                      <CardContent className="space-y-4">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                          <div>
                            <Label htmlFor="metrics_url">Metrics URL</Label>
                            <Input
                              id="metrics_url"
                              placeholder="ex: https://grafana.company.com/d/service"
                              {...register('metrics_url')}
                            />
                          </div>

                          <div>
                            <Label htmlFor="logs_url">Logs URL</Label>
                            <Input
                              id="logs_url"
                              placeholder="ex: https://kibana.company.com/app"
                              {...register('logs_url')}
                            />
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  </TabsContent>

                  {/* Review Tab */}
                  <TabsContent value="review" className="space-y-6 m-0">
                    <Card>
                      <CardHeader>
                        <CardTitle className="flex items-center gap-2">
                          <Settings className="h-5 w-5" />
                          Review Configuration
                        </CardTitle>
                      </CardHeader>
                      <CardContent className="space-y-6">
                        <div>
                          <h4 className="font-medium mb-3">
                            Basic Information
                          </h4>
                          <div className="grid grid-cols-2 gap-4 text-sm">
                            <div>
                              <span className="text-muted-foreground">
                                Name:
                              </span>
                              <p className="font-medium">
                                {watch('name') ||
                                  serviceData?.metadata.name ||
                                  'Not specified'}
                              </p>
                            </div>
                            <div>
                              <span className="text-muted-foreground">
                                Tier:
                              </span>
                              <p className="font-medium">
                                {watchedTier || serviceData?.metadata.tier}
                              </p>
                            </div>
                            <div>
                              <span className="text-muted-foreground">
                                GitHub Team:
                              </span>
                              <p className="font-medium">
                                {watch('github_team') ||
                                  serviceData?.spec.team.github_team ||
                                  'Not specified'}
                              </p>
                            </div>
                            <div>
                              <span className="text-muted-foreground">
                                Impact:
                              </span>
                              <p className="font-medium">
                                {watch('impact') ||
                                  serviceData?.spec.business?.impact ||
                                  'medium'}
                              </p>
                            </div>
                          </div>
                        </div>

                        <Separator />

                        <div>
                          <h4 className="font-medium mb-3">
                            Kubernetes Configuration
                          </h4>
                          <div className="grid grid-cols-2 gap-4 text-sm">
                            <div>
                              <span className="text-muted-foreground">
                                Context:
                              </span>
                              <p className="font-medium">
                                {watch('env_context') ||
                                  serviceData?.spec.kubernetes
                                    ?.environments?.[0]?.context ||
                                  'Not specified'}
                              </p>
                            </div>
                            <div>
                              <span className="text-muted-foreground">
                                Namespace:
                              </span>
                              <p className="font-medium">
                                {watch('env_namespace') ||
                                  serviceData?.spec.kubernetes
                                    ?.environments?.[0]?.namespace ||
                                  'Not specified'}
                              </p>
                            </div>
                            <div>
                              <span className="text-muted-foreground">
                                Deployment:
                              </span>
                              <p className="font-medium">
                                {watch('deployment_name') ||
                                  serviceData?.spec.kubernetes
                                    ?.environments?.[0]?.resources
                                    ?.deployments?.[0]?.name ||
                                  'Not specified'}
                              </p>
                            </div>
                            <div>
                              <span className="text-muted-foreground">
                                Replicas:
                              </span>
                              <p className="font-medium">
                                {watch('deployment_replicas')}
                              </p>
                            </div>
                          </div>
                        </div>

                        <Separator />

                        <div>
                          <h4 className="font-medium mb-3">Dependencies</h4>
                          <div className="flex flex-wrap gap-1">
                            {customDependencies.length > 0 ? (
                              customDependencies.map((dep) => (
                                <Badge key={dep} variant="secondary">
                                  {dep}
                                </Badge>
                              ))
                            ) : (
                              <p className="text-sm text-muted-foreground">
                                No dependencies added
                              </p>
                            )}
                          </div>
                        </div>

                        {(watch('metrics_url') || watch('logs_url')) && (
                          <>
                            <Separator />
                            <div>
                              <h4 className="font-medium mb-3">
                                Observability
                              </h4>
                              <div className="grid grid-cols-2 gap-4 text-sm">
                                {watch('metrics_url') && (
                                  <div>
                                    <span className="text-muted-foreground">
                                      Metrics:
                                    </span>
                                    <p className="font-medium break-all">
                                      {watch('metrics_url')}
                                    </p>
                                  </div>
                                )}
                                {watch('logs_url') && (
                                  <div>
                                    <span className="text-muted-foreground">
                                      Logs:
                                    </span>
                                    <p className="font-medium break-all">
                                      {watch('logs_url')}
                                    </p>
                                  </div>
                                )}
                              </div>
                            </div>
                          </>
                        )}
                      </CardContent>
                    </Card>
                  </TabsContent>
                </div>
              </ScrollArea>
            </Tabs>

            <div className="flex-none flex justify-between pt-6 border-t">
              <Button
                type="button"
                variant="outline"
                onClick={handleClose}
                disabled={isLoading}
              >
                Cancel
              </Button>
              <div className="flex gap-2">
                {activeTab !== 'basic' && (
                  <Button
                    type="button"
                    variant="outline"
                    disabled={isLoading}
                    onClick={() => {
                      const tabs = [
                        'basic',
                        'infrastructure',
                        'observability',
                        'review',
                      ];
                      const currentIndex = tabs.indexOf(activeTab);
                      if (currentIndex > 0) {
                        setActiveTab(tabs[currentIndex - 1] as string);
                      }
                    }}
                  >
                    Previous
                  </Button>
                )}
                {activeTab !== 'review' ? (
                  <Button
                    type="button"
                    disabled={isLoading}
                    onClick={() => {
                      const tabs = [
                        'basic',
                        'infrastructure',
                        'observability',
                        'review',
                      ];
                      const currentIndex = tabs.indexOf(activeTab);
                      if (currentIndex < tabs.length - 1) {
                        setActiveTab(tabs[currentIndex + 1] as string);
                      }
                    }}
                  >
                    Next
                  </Button>
                ) : (
                  <Button type="submit" disabled={isLoading}>
                    {isLoading
                      ? isEditing
                        ? 'Saving...'
                        : 'Creating...'
                      : isEditing
                        ? 'Save Changes'
                        : 'Create Service'}
                  </Button>
                )}
              </div>
            </div>
          </form>
        )}
      </DialogContent>
    </Dialog>
  );
}
