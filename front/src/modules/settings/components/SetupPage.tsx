import { useMemo, useState, type ComponentType } from 'react';
import { toast } from 'sonner';
import {
  Activity,
  Check,
  CheckCircle2,
  Cloud,
  Container,
  FileText,
  LineChart,
  Loader2,
  Package,
  Plus,
  Settings as SettingsIcon,
  Shield,
  Workflow,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { Badge } from '@/components/ui/badge';
import { Switch } from '@/components/ui/switch';
import { cn } from '@/lib/utils';
import { configureSetup } from '../resources/setupResource';
import type { SetupConfigureRequest } from '../types';
import {
  AuthProviderCard,
  AwsAccountCard,
  GeneralSettingsForm,
  KubernetesClusterCard,
  ObservabilityProviderCard,
  PluginToggleGrid,
  ServiceCatalogForm,
} from './shared';
import { PLUGIN_OPTIONS } from '../constants';
import {
  createDefaultAuthProvider,
  createDefaultCloudProvider,
  createDefaultGeneralSettings,
  createDefaultKubernetesCluster,
  createDefaultObservability,
  createDefaultObservabilityProvider,
  createDefaultServiceCatalog,
  headersToArray,
} from '../utils/formUtils';
import type {
  AuthProviderFormValue,
  CloudProviderFormValue,
  GeneralSettingsFormValue,
  KubernetesClusterFormValue,
  ObservabilityFormValue,
  ServiceCatalogFormValue,
} from '../types';

interface SetupPageProps {
  onComplete?: () => Promise<void> | void;
}

const STEP_ICONS: Record<string, ComponentType<{ className?: string }>> = {
  general: SettingsIcon,
  plugins: Package,
  review: CheckCircle2,
  'config-Auth': Shield,
  'config-ServiceCatalog': Package,
  'config-Kubernetes': Container,
  'config-AWS': Cloud,
  'config-Observability': Activity,
};

type ObservabilityProviderGroup =
  | 'logsProviders'
  | 'tracesProviders'
  | 'metricsProviders';

export default function SetupPage({
  onComplete = async () => {},
}: SetupPageProps = {}): JSX.Element {
  const [currentStep, setCurrentStep] = useState(0);
  const [generalSettings, setGeneralSettings] =
    useState<GeneralSettingsFormValue>(createDefaultGeneralSettings());
  const [enabledPlugins, setEnabledPlugins] = useState<string[]>([]);
  const [authProviders, setAuthProviders] = useState<AuthProviderFormValue[]>([
    createDefaultAuthProvider(),
  ]);
  const [k8sClusters, setK8sClusters] = useState<KubernetesClusterFormValue[]>([
    createDefaultKubernetesCluster(),
  ]);
  const [cloudProviders, setCloudProviders] =
    useState<CloudProviderFormValue[]>([
      createDefaultCloudProvider('aws', {
    name: 'default-account',
        hasSecretAccessKey: false,
        secretAccessKeyInput: '',
      }),
    ]);
  const [serviceCatalog, setServiceCatalog] =
    useState<ServiceCatalogFormValue>(createDefaultServiceCatalog());
  const [observability, setObservability] =
    useState<ObservabilityFormValue>(createDefaultObservability());
  const [loading, setLoading] = useState(false);

  const steps = useMemo(() => {
    const pluginSteps = enabledPlugins.map((pluginId) => {
      const plugin = PLUGIN_OPTIONS.find((item) => item.id === pluginId);
      return {
        id: `config-${pluginId}`,
        title: plugin ? `Configure ${plugin.label}` : `Configure ${pluginId}`,
      };
    });

    return [
      { id: 'general', title: 'General Settings' },
      { id: 'plugins', title: 'Select Plugins' },
      ...pluginSteps,
      { id: 'review', title: 'Review & Finish' },
    ];
  }, [enabledPlugins]);

  const currentStepId = steps[currentStep]?.id ?? 'general';

  const togglePlugin = (pluginId: string) => {
    setEnabledPlugins((prev) =>
      prev.includes(pluginId)
        ? prev.filter((id) => id !== pluginId)
        : [...prev, pluginId]
    );
  };

  const handleAddObservabilityProvider = (
    group: ObservabilityProviderGroup,
    type: string
  ) => {
    const newProvider = createDefaultObservabilityProvider(type);
    setObservability((prev) => ({
      ...prev,
      [group]: [...prev[group], newProvider],
    }));
  };

  const handleRemoveObservabilityProvider = (
    group: ObservabilityProviderGroup,
    id: string
  ) => {
    setObservability((prev) => ({
      ...prev,
      [group]: prev[group].filter((provider) => provider.id !== id),
    }));
  };

  const handleSubmit = async () => {
    try {
      setLoading(true);

    const payload: SetupConfigureRequest = {
      config: {
        port: generalSettings.port,
        origin: generalSettings.origin,
        headers: headersToArray(generalSettings.headers),
        front: generalSettings.front,
      },
      enabled_plugins: enabledPlugins,
      plugins: {
        auth: enabledPlugins.includes('Auth')
          ? {
              provider: 'github',
              clientId: authProviders[0]?.clientId ?? '',
              clientSecret: authProviders[0]?.clientSecretInput ?? '',
              orgPermission: authProviders[0]?.orgPermission ?? '',
              redirectURL: `${generalSettings.origin.replace(/\/$/, '')}/api/oauth/redirect`,
              scopes: ['user', 'repo', 'read:org'],
            }
          : null,
        service_catalog: enabledPlugins.includes('ServiceCatalog')
          ? {
              storage: {
                provider: serviceCatalog.storageProvider,
                filesystem:
                  serviceCatalog.storageProvider === 'filesystem'
                    ? {
                        directory: serviceCatalog.directory,
                      }
                    : undefined,
                github:
                  serviceCatalog.storageProvider === 'github'
                    ? {
                        repository: serviceCatalog.githubRepository,
                        branch: serviceCatalog.githubBranch,
                      }
                    : undefined,
                s3:
                  serviceCatalog.storageProvider === 's3'
                    ? {
                        bucket: serviceCatalog.s3Bucket,
                      }
                    : undefined,
              },
              versioning: {
                enabled: serviceCatalog.versioningEnabled,
                provider: 'simple',
              },
            }
          : null,
        kubernetes: enabledPlugins.includes('Kubernetes')
          ? [
              {
                name: k8sClusters[0]?.name ?? '',
                connectionType: k8sClusters[0]?.connectionType ?? 'kubeconfig',
                kubeconfig:
                  k8sClusters[0]?.connectionType === 'kubeconfig'
                    ? k8sClusters[0]?.kubeconfig
                    : undefined,
                context:
                  k8sClusters[0]?.connectionType === 'kubeconfig'
                    ? k8sClusters[0]?.context
                    : undefined,
                host:
                  k8sClusters[0]?.connectionType === 'remote'
                    ? k8sClusters[0]?.host
                    : undefined,
                token:
                  k8sClusters[0]?.connectionType === 'remote'
                    ? k8sClusters[0]?.token
                    : undefined,
                certificate:
                  k8sClusters[0]?.connectionType === 'remote'
                    ? k8sClusters[0]?.certificate
                    : undefined,
              },
            ]
          : null,
        aws: enabledPlugins.includes('AWS')
          ? [
              {
                name: cloudProviders[0]?.name ?? '',
                region: cloudProviders[0]?.region ?? 'us-east-1',
                accessKeyId: cloudProviders[0]?.accessKeyId ?? '',
                secretAccessKey: cloudProviders[0]?.secretAccessKeyInput ?? '',
              },
            ]
          : null,
        observability: enabledPlugins.includes('Observability')
          ? {
              enabled: observability.enabled,
              logs:
                observability.logsProviders.length > 0
                  ? observability.logsProviders.map((provider) => ({
                      name: provider.name,
                      type: provider.type,
                      url: provider.url,
                      timeout: provider.timeout,
                      retention: provider.retention,
                      enabled: provider.enabled,
                    }))
                  : undefined,
              traces:
                observability.tracesProviders.length > 0
                  ? observability.tracesProviders.map((provider) => ({
                      name: provider.name,
                      type: provider.type,
                      url: provider.url,
                      timeout: provider.timeout,
                      retention: provider.retention,
                      enabled: provider.enabled,
                    }))
                  : undefined,
              metrics:
                observability.metricsProviders.length > 0
                  ? observability.metricsProviders.map((provider) => ({
                      name: provider.name,
                      type: provider.type,
                      url: provider.url,
                      timeout: provider.timeout,
                      retention: provider.retention,
                      enabled: provider.enabled,
                    }))
                  : undefined,
            }
          : null,
      },
    };

      await configureSetup(payload);
      toast.success('Setup completed! Restart the backend to load plugins.');
      await onComplete();
    } catch (error) {
      toast.error(
        error instanceof Error ? error.message : 'Failed to complete setup'
      );
    } finally {
      setLoading(false);
    }
  };

  const generalStep = (
    <GeneralSettingsForm
      value={generalSettings}
      onChange={setGeneralSettings}
    />
  );

  const pluginsStep = (
    <div className="space-y-6">
      <div>
        <h2 className="text-xl font-semibold tracking-tight">
          Select plugins to enable
        </h2>
        <p className="text-sm text-muted-foreground">
          Each module adds new capabilities to DashOps. You can configure them
          in the next steps.
        </p>
          </div>
      <PluginToggleGrid
        options={PLUGIN_OPTIONS}
        selected={enabledPlugins}
        onToggle={togglePlugin}
      />
    </div>
  );

  const authStep = (
    <div className="space-y-4">
      <AuthProviderCard
        provider={authProviders[0]}
        onChange={(provider) =>
          setAuthProviders([
            {
              ...provider,
              hasClientSecret:
                (provider.clientSecretInput?.trim().length ?? 0) > 0 ||
                provider.hasClientSecret,
            },
          ])
        }
        secretsVisible
        disabled={false}
      />
    </div>
  );

  const serviceCatalogStep = (
    <ServiceCatalogForm value={serviceCatalog} onChange={setServiceCatalog} />
  );

  const kubernetesStep = (
    <KubernetesClusterCard
      cluster={k8sClusters[0]}
      onChange={(cluster) => setK8sClusters([cluster])}
              />
  );

  const awsStep = (
    <AwsAccountCard
      account={cloudProviders[0]}
      onChange={(account) =>
        setCloudProviders([
          {
            ...account,
            hasSecretAccessKey:
              (account.secretAccessKeyInput?.trim().length ?? 0) > 0 ||
              account.hasSecretAccessKey,
          },
        ])
      }
      secretsVisible
      disabled={false}
    />
  );

  const observabilityStep = (
    <div className="space-y-6">
      <div className="flex items-center justify-between rounded-lg border p-4">
        <div>
          <p className="text-sm font-medium">Enable observability</p>
          <p className="text-xs text-muted-foreground">
            Toggle to activate logs, traces and metrics integrations.
          </p>
        </div>
        <Switch
          checked={observability.enabled}
          onCheckedChange={(enabled) =>
            setObservability((prev) => ({ ...prev, enabled }))
          }
        />
      </div>

      <ObservabilityGroup
        title="Logs providers"
        description="Configure the logging backend (e.g., Loki)."
        icon={<FileText className="h-4 w-4 text-muted-foreground" />}
        labelPrefix="Logs"
        providers={observability.logsProviders}
        onChange={(providers) =>
          setObservability((prev) => ({
            ...prev,
            logsProviders: providers,
          }))
        }
        onAdd={() => handleAddObservabilityProvider('logsProviders', 'loki')}
        onRemove={(id) =>
          handleRemoveObservabilityProvider('logsProviders', id)
        }
      />

      <ObservabilityGroup
        title="Traces providers"
        description="Configure the tracing backend (e.g., Tempo)."
        icon={<Workflow className="h-4 w-4 text-muted-foreground" />}
        labelPrefix="Traces"
        providers={observability.tracesProviders}
        onChange={(providers) =>
          setObservability((prev) => ({
            ...prev,
            tracesProviders: providers,
          }))
        }
        onAdd={() =>
          handleAddObservabilityProvider('tracesProviders', 'tempo')
        }
        onRemove={(id) =>
          handleRemoveObservabilityProvider('tracesProviders', id)
        }
      />

      <ObservabilityGroup
        title="Metrics providers"
        description="Configure the metrics backend (e.g., Prometheus)."
        icon={<LineChart className="h-4 w-4 text-muted-foreground" />}
        labelPrefix="Metrics"
        providers={observability.metricsProviders}
        onChange={(providers) =>
          setObservability((prev) => ({
            ...prev,
            metricsProviders: providers,
          }))
        }
        onAdd={() =>
          handleAddObservabilityProvider('metricsProviders', 'prometheus')
        }
        onRemove={(id) =>
          handleRemoveObservabilityProvider('metricsProviders', id)
        }
      />
    </div>
  );

  const reviewStep = (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>General</CardTitle>
          <CardDescription>
            Overview of the core settings that will be applied.
          </CardDescription>
        </CardHeader>
        <CardContent className="grid gap-4 md:grid-cols-2">
          <SummaryField label="Port" value={generalSettings.port} />
          <SummaryField label="Origin" value={generalSettings.origin} />
          <SummaryField
            label="Headers"
            value={generalSettings.headers}
          />
          <SummaryField label="Frontend path" value={generalSettings.front} />
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <CardTitle>Selected plugins</CardTitle>
          <CardDescription>
            Modules that will be enabled once the backend restarts.
          </CardDescription>
        </CardHeader>
        <CardContent className="flex flex-wrap gap-2">
          {enabledPlugins.length === 0 ? (
            <p className="text-sm text-muted-foreground">
              No plugins selected. You can enable modules later in Settings.
            </p>
          ) : (
            enabledPlugins.map((pluginId) => {
              const plugin = PLUGIN_OPTIONS.find(
                (option) => option.id === pluginId
              );
              return (
                <Badge key={pluginId} variant="secondary">
                  {plugin?.label ?? pluginId}
                </Badge>
              );
            })
          )}
        </CardContent>
      </Card>
    </div>
  );

  let stepContent: JSX.Element = generalStep;

  if (currentStepId === 'general') {
    stepContent = generalStep;
  } else if (currentStepId === 'plugins') {
    stepContent = pluginsStep;
  } else if (currentStepId === 'review') {
    stepContent = reviewStep;
  } else if (currentStepId === 'config-Auth') {
    stepContent = authStep;
  } else if (currentStepId === 'config-ServiceCatalog') {
    stepContent = serviceCatalogStep;
  } else if (currentStepId === 'config-Kubernetes') {
    stepContent = kubernetesStep;
  } else if (currentStepId === 'config-AWS') {
    stepContent = awsStep;
  } else if (currentStepId === 'config-Observability') {
    stepContent = observabilityStep;
  }

  const StepIcon =
    STEP_ICONS[currentStepId] ??
    STEP_ICONS[currentStepId.replace('config-', '')] ??
    Package;

  return (
    <div className="min-h-screen bg-background">
      <div className="mx-auto flex min-h-screen w-full max-w-5xl flex-col gap-8 px-6 py-12">
        <div className="text-center space-y-3">
          <Badge variant="secondary" className="self-center">
            DashOps Setup
          </Badge>
          <h1 className="text-3xl font-semibold tracking-tight">
            Configure your DashOps instance
          </h1>
          <p className="text-sm text-muted-foreground">
            Provide the initial configuration, select the modules you want to
            enable and review the summary before generating the configuration
            file. You can change everything later in Settings.
          </p>
        </div>

        <div className="flex flex-wrap items-center justify-center gap-4">
              {steps.map((step, index) => {
                const isActive = index === currentStep;
                const isCompleted = index < currentStep;
            const Icon = STEP_ICONS[step.id] ??
              STEP_ICONS[step.id.replace('config-', '')] ??
              Package;

                return (
              <div key={step.id} className="flex items-center gap-3">
                    <div
                  className={cn(
                    'flex h-12 w-12 items-center justify-center rounded-full border-2 transition-all',
                        isActive || isCompleted
                      ? 'border-primary bg-primary text-primary-foreground shadow-sm'
                          : 'border-border bg-background text-muted-foreground'
                  )}
                    >
                  {isCompleted ? (
                    <Check className="h-5 w-5" />
                  ) : (
                    <Icon className="h-5 w-5" />
                  )}
                    </div>
                    <span
                  className={cn(
                    'text-sm font-medium',
                        isActive ? 'text-primary' : 'text-muted-foreground'
                  )}
                    >
                      {step.title}
                    </span>
                    {index < steps.length - 1 && (
                  <Separator orientation="vertical" className="hidden h-6 md:block" />
                    )}
                  </div>
                );
              })}
            </div>

        <Card className="flex-1">
          <CardHeader className="flex flex-row items-center justify-between gap-4">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-full bg-primary/10 text-primary">
                <StepIcon className="h-5 w-5" />
              </div>
              <div>
                <CardTitle className="text-lg font-semibold">
                  {steps[currentStep]?.title ?? 'Setup'}
                </CardTitle>
                <CardDescription>
                  {currentStepId === 'general' &&
                    'Define the core settings for your DashOps deployment.'}
                  {currentStepId === 'plugins' &&
                    'Pick the modules that you want to activate.'}
                  {currentStepId.startsWith('config-') &&
                    'Provide the configuration required by this module.'}
                  {currentStepId === 'review' &&
                    'Validate all information before finishing the setup.'}
                </CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent className="space-y-8">
            {stepContent}
            <Separator />
            <div className="flex items-center justify-between">
              <Button
                variant="outline"
                onClick={() =>
                  setCurrentStep((step) => (step > 0 ? step - 1 : step))
                }
                disabled={currentStep === 0 || loading}
                className="gap-2"
              >
                Back
              </Button>
              {currentStep === steps.length - 1 ? (
                <Button
                  onClick={handleSubmit}
                  disabled={loading}
                  className="gap-2"
                >
                  {loading && <Loader2 className="h-4 w-4 animate-spin" />}
                  Finish setup
                </Button>
              ) : (
                <Button
                  onClick={() =>
                    setCurrentStep((step) =>
                      step < steps.length - 1 ? step + 1 : step
                    )
                  }
                  className="gap-2"
                  disabled={
                    loading ||
                    (currentStepId === 'plugins' && enabledPlugins.length === 0)
                  }
                >
                  Continue
                </Button>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

interface ObservabilityGroupProps {
  title: string;
  description: string;
  icon: JSX.Element;
  labelPrefix: string;
  providers: ObservabilityProviderFormValue[];
  onChange: (providers: ObservabilityProviderFormValue[]) => void;
  onAdd: () => void;
  onRemove: (id: string) => void;
}

function ObservabilityGroup({
  title,
  description,
  icon,
  labelPrefix,
  providers,
  onChange,
  onAdd,
  onRemove,
}: ObservabilityGroupProps): JSX.Element {
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          {icon}
          <div>
            <h4 className="text-sm font-semibold">{title}</h4>
            <p className="text-xs text-muted-foreground">{description}</p>
          </div>
        </div>
        <Button
          variant="outline"
          size="sm"
          className="gap-2"
          onClick={onAdd}
        >
          <Plus className="h-4 w-4" />
          Add provider
        </Button>
      </div>
      <div className="space-y-4">
        {providers.length === 0 && (
          <div className="rounded-lg border border-dashed p-6 text-center text-sm text-muted-foreground">
            No providers configured yet.
          </div>
        )}
        {providers.map((provider) => (
          <ObservabilityProviderCard
            key={provider.id}
            titlePrefix={labelPrefix}
            provider={provider}
            onChange={(updated) =>
              onChange(
                providers.map((current) =>
                  current.id === provider.id ? updated : current
                )
              )
            }
            onRemove={
              providers.length > 1
                ? (id) => onRemove(id)
                : undefined
            }
          />
        ))}
      </div>
    </div>
  );
}

interface SummaryFieldProps {
  label: string;
  value: string;
}

function SummaryField({ label, value }: SummaryFieldProps) {
  return (
    <div className="space-y-1 rounded-lg border p-4">
      <p className="text-xs font-medium text-muted-foreground">{label}</p>
      <p className="text-sm font-semibold">{value}</p>
    </div>
  );
}

