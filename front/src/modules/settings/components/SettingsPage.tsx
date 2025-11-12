import { useEffect, useMemo, useState } from 'react';
import { toast } from 'sonner';
import {
  Activity,
  Building2,
  Chrome,
  Cloud,
  Container,
  Github,
  LayoutGrid,
  Loader2,
  Package,
  Plus,
  Save,
  Shield,
  Trash2,
  Workflow,
  LineChart,
  FileText,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Switch } from '@/components/ui/switch';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { cn } from '@/lib/utils';
import { useSettingsConfig } from '../hooks/useSettingsConfig';
import {
  AuthProviderCard,
  AwsAccountCard,
  DisabledTabNotice,
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
  headersToString,
} from '../utils/formUtils';
import type {
  AuthProviderFormValue,
  AuthProviderType,
  CloudProviderFormValue,
  GeneralSettingsFormValue,
  KubernetesClusterFormValue,
  ObservabilityFormValue,
  ObservabilityProviderFormValue,
  ServiceCatalogFormValue,
} from '../types';
import type { UpdatePluginsPayload, UpdateSettingsRequest } from '../types';

const TAB_ORDER = ['auth', 'kubernetes', 'cloud', 'catalog', 'observability'] as const;

export default function SettingsPage(): JSX.Element {
  const { config, loading, saving, error, canEdit, save } = useSettingsConfig();
  const [generalSettings, setGeneralSettings] =
    useState<GeneralSettingsFormValue>(createDefaultGeneralSettings());
  const [enabledPlugins, setEnabledPlugins] = useState<string[]>([]);
  const [authProviders, setAuthProviders] = useState<AuthProviderFormValue[]>([]);
  const [k8sClusters, setK8sClusters] = useState<KubernetesClusterFormValue[]>([]);
  const [cloudProviders, setCloudProviders] = useState<CloudProviderFormValue[]>([]);
  const [serviceCatalog, setServiceCatalog] =
    useState<ServiceCatalogFormValue>(createDefaultServiceCatalog());
  const [observability, setObservability] =
    useState<ObservabilityFormValue>(createDefaultObservability());
  const [activeTab, setActiveTab] = useState<string>('auth');

  useEffect(() => {
    if (!config) return;

    setGeneralSettings(
      createDefaultGeneralSettings({
        port: config.port ?? '8080',
        origin: config.origin ?? 'http://localhost:5173',
        headers: headersToString(config.headers),
        front: config.front ?? 'front/dist',
      })
    );

    const plugins = config.plugins ?? [];
    setEnabledPlugins(plugins);

    if (config.auth && config.auth.length > 0) {
      setAuthProviders(
        config.auth.map((provider, index) => {
          const providerId = provider.provider?.toLowerCase() ?? 'github';
          const type =
            providerId === 'google'
              ? 'google'
              : providerId === 'sso'
              ? 'sso'
              : 'github';
          const providerName =
            type === 'google'
              ? 'Google OAuth'
              : type === 'sso'
              ? 'SSO'
              : 'GitHub OAuth';

          return createDefaultAuthProvider(type, {
            id: `auth-${index}-${provider.provider ?? 'github'}`,
            name: providerName,
            clientId: '',
            clientIdMasked: provider.clientIdMasked,
            orgPermission: provider.orgPermission ?? '',
            hasClientSecret: provider.hasClientSecret,
          });
        })
      );
    } else if (plugins.includes('Auth')) {
      setAuthProviders([createDefaultAuthProvider()]);
    } else {
      setAuthProviders([]);
    }

    if (config.kubernetes && config.kubernetes.length > 0) {
      setK8sClusters(
        config.kubernetes.map((cluster, index) =>
          createDefaultKubernetesCluster({
            id: `k8s-${index}-${cluster.name ?? 'cluster'}`,
            name: cluster.name ?? '',
            connectionType: cluster.connectionType ?? 'kubeconfig',
            kubeconfig: cluster.kubeconfig ?? '$HOME/.kube/config',
            context: cluster.context ?? '',
            host: cluster.host ?? '',
            certificate: cluster.certificate ?? '',
            hasToken: cluster.hasToken,
          })
        )
      );
    } else if (plugins.includes('Kubernetes')) {
      setK8sClusters([createDefaultKubernetesCluster()]);
    } else {
      setK8sClusters([]);
    }

    if (config.aws && config.aws.length > 0) {
      setCloudProviders(
        config.aws.map((account, index) =>
          createDefaultCloudProvider('aws', {
            id: `aws-${index}-${account.name ?? 'aws'}`,
            name: account.name ?? '',
            region: account.region ?? 'us-east-1',
            accessKeyId: '',
            accessKeyIdMasked: account.accessKeyIdMasked,
            hasSecretAccessKey: account.hasSecretAccessKey,
          })
        )
      );
    } else if (plugins.includes('AWS')) {
      setCloudProviders([createDefaultCloudProvider('aws')]);
    } else {
      setCloudProviders([]);
    }

    if (config.service_catalog) {
      setServiceCatalog(
        createDefaultServiceCatalog({
          storageProvider:
            config.service_catalog.storage?.provider ?? 'filesystem',
          directory:
            config.service_catalog.storage?.filesystem?.directory ??
            '../services',
          githubRepository:
            config.service_catalog.storage?.github?.repository ?? '',
          githubBranch: config.service_catalog.storage?.github?.branch ?? 'main',
          s3Bucket: config.service_catalog.storage?.s3?.bucket ?? '',
          versioningEnabled: config.service_catalog.versioning?.enabled ?? false,
        })
      );
    } else if (plugins.includes('ServiceCatalog')) {
      setServiceCatalog(createDefaultServiceCatalog());
    } else {
      setServiceCatalog(createDefaultServiceCatalog());
    }

    if (config.observability) {
      setObservability({
        enabled: config.observability.enabled ?? true,
        logsProviders:
          config.observability.logs?.map((provider, index) =>
            createDefaultObservabilityProvider(provider.type ?? 'loki', {
              id: `logs-${index}-${provider.name ?? 'provider'}`,
              name: provider.name ?? '',
              url: provider.url ?? '',
              timeout: provider.timeout ?? '30s',
              retention: provider.retention ?? '30d',
              enabled: provider.enabled ?? true,
              labels: provider.labels ?? {},
              authType: provider.auth?.type,
              hasAuthPassword: provider.auth?.hasPassword,
              hasAuthToken: provider.auth?.hasToken,
            })
          ) ?? [createDefaultObservabilityProvider('loki')],
        tracesProviders:
          config.observability.traces?.map((provider, index) =>
            createDefaultObservabilityProvider(provider.type ?? 'tempo', {
              id: `traces-${index}-${provider.name ?? 'provider'}`,
              name: provider.name ?? '',
              url: provider.url ?? '',
              timeout: provider.timeout ?? '30s',
              retention: provider.retention ?? '7d',
              enabled: provider.enabled ?? true,
              labels: provider.labels ?? {},
              authType: provider.auth?.type,
              hasAuthPassword: provider.auth?.hasPassword,
              hasAuthToken: provider.auth?.hasToken,
            })
          ) ?? [createDefaultObservabilityProvider('tempo', { retention: '7d' })],
        metricsProviders:
          config.observability.metrics?.map((provider, index) =>
            createDefaultObservabilityProvider(provider.type ?? 'prometheus', {
              id: `metrics-${index}-${provider.name ?? 'provider'}`,
              name: provider.name ?? '',
              url: provider.url ?? '',
              timeout: provider.timeout ?? '30s',
              retention: provider.retention ?? '90d',
              enabled: provider.enabled ?? true,
              labels: provider.labels ?? {},
              authType: provider.auth?.type,
              hasAuthPassword: provider.auth?.hasPassword,
              hasAuthToken: provider.auth?.hasToken,
            })
          ) ??
          [
            createDefaultObservabilityProvider('prometheus', {
              retention: '90d',
            }),
          ],
      });
    } else if (plugins.includes('Observability')) {
      setObservability(createDefaultObservability());
    } else {
      setObservability(createDefaultObservability({ enabled: false }));
    }

    const firstEnabledTab =
      TAB_ORDER.find((tab) => {
        if (tab === 'auth') return plugins.includes('Auth');
        if (tab === 'kubernetes') return plugins.includes('Kubernetes');
        if (tab === 'cloud') return plugins.includes('AWS');
        if (tab === 'catalog') return plugins.includes('ServiceCatalog');
        if (tab === 'observability') return plugins.includes('Observability');
        return false;
      }) ?? 'auth';

    setActiveTab(firstEnabledTab);
  }, [config]);

  useEffect(() => {
    if (error) {
      toast.error(error);
    }
  }, [error]);

  const isPluginEnabled = (pluginId: string) =>
    enabledPlugins.includes(pluginId);

  const handleTogglePlugin = (pluginId: string) => {
    if (!canEdit) return;

    setEnabledPlugins((prev) => {
      const enabled = prev.includes(pluginId);
      if (enabled) {
        return prev.filter((id) => id !== pluginId);
      }

      if (pluginId === 'Auth' && authProviders.length === 0) {
        setAuthProviders([createDefaultAuthProvider()]);
      }
      if (pluginId === 'Kubernetes' && k8sClusters.length === 0) {
        setK8sClusters([createDefaultKubernetesCluster()]);
      }
      if (pluginId === 'AWS' && cloudProviders.length === 0) {
        setCloudProviders([createDefaultCloudProvider('aws')]);
      }
      if (pluginId === 'ServiceCatalog') {
        setServiceCatalog(createDefaultServiceCatalog());
      }
      if (pluginId === 'Observability') {
        setObservability(createDefaultObservability());
      }

      return [...prev, pluginId];
    });
  };

  const handleAddAuthProvider = (provider: AuthProviderType) => {
    setAuthProviders((prev) => [...prev, createDefaultAuthProvider(provider)]);
  };

  const handleRemoveAuthProvider = (id: string) => {
    setAuthProviders((prev) => prev.filter((provider) => provider.id !== id));
  };

  const handleAddK8sCluster = () => {
    setK8sClusters((prev) => [...prev, createDefaultKubernetesCluster()]);
  };

  const handleRemoveK8sCluster = (id: string) => {
    setK8sClusters((prev) => prev.filter((cluster) => cluster.id !== id));
  };

  const handleAddCloudProvider = () => {
    setCloudProviders((prev) => [...prev, createDefaultCloudProvider('aws')]);
  };

  const handleRemoveCloudProvider = (id: string) => {
    setCloudProviders((prev) => prev.filter((provider) => provider.id !== id));
  };

  const handleAddObservabilityProvider = (
    group: keyof ObservabilityFormValue,
    type: string
  ) => {
    const newProvider = createDefaultObservabilityProvider(type);
    setObservability((prev) => ({
      ...prev,
      [group]: [...prev[group], newProvider],
    }));
  };

  const handleRemoveObservabilityProvider = (
    group: keyof ObservabilityFormValue,
    id: string
  ) => {
    setObservability((prev) => ({
      ...prev,
      [group]: prev[group].filter((provider) => provider.id !== id),
    }));
  };

  const activeTabDisabled = useMemo(() => {
    if (activeTab === 'auth') return !isPluginEnabled('Auth');
    if (activeTab === 'kubernetes') return !isPluginEnabled('Kubernetes');
    if (activeTab === 'cloud') return !isPluginEnabled('AWS');
    if (activeTab === 'catalog') return !isPluginEnabled('ServiceCatalog');
    if (activeTab === 'observability')
      return !isPluginEnabled('Observability');
    return false;
  }, [activeTab, enabledPlugins]);

  const handleSave = async () => {
    if (!config) return;

    const pluginsPayload: UpdatePluginsPayload = {};

    if (isPluginEnabled('Auth')) {
      pluginsPayload.auth = authProviders.map((provider) => {
        const clientIdInput = provider.clientId ?? '';
        const trimmedClientId = clientIdInput.trim();
        const hasSecretInput = provider.clientSecretInput !== undefined;
        const trimmedSecret = provider.clientSecretInput?.trim() ?? '';

        return {
          provider: provider.provider,
          clientId: trimmedClientId.length > 0 ? trimmedClientId : undefined,
          clientSecret: trimmedSecret.length > 0 ? trimmedSecret : undefined,
          clearClientSecret:
            hasSecretInput && trimmedSecret === '' && provider.hasClientSecret
              ? true
              : undefined,
          orgPermission: provider.orgPermission,
          redirectURL:
            provider.provider === 'github'
              ? `${generalSettings.origin.replace(/\/$/, '')}/api/oauth/redirect`
              : undefined,
          scopes:
            provider.provider === 'github'
              ? ['user', 'repo', 'read:org']
              : undefined,
        };
      });
    } else if (config.auth && config.auth.length > 0) {
      pluginsPayload.auth = [];
    }

    if (isPluginEnabled('Kubernetes')) {
      pluginsPayload.kubernetes = k8sClusters.map((cluster) => {
        const isRemote = cluster.connectionType === 'remote';
        const hasTokenInput = cluster.token !== undefined;
        const trimmedToken = cluster.token?.trim() ?? '';

        return {
          name: cluster.name,
          connectionType: cluster.connectionType,
          kubeconfig:
            cluster.connectionType === 'kubeconfig'
              ? cluster.kubeconfig
              : undefined,
          context:
            cluster.connectionType === 'kubeconfig'
              ? cluster.context
              : undefined,
          host: isRemote ? cluster.host : undefined,
          token: isRemote && trimmedToken.length > 0 ? trimmedToken : undefined,
          certificate: isRemote ? cluster.certificate : undefined,
          clearToken:
            isRemote && hasTokenInput && trimmedToken === '' && cluster.hasToken
              ? true
              : undefined,
        };
      });
    } else if (config.kubernetes && config.kubernetes.length > 0) {
      pluginsPayload.kubernetes = [];
    }

    if (isPluginEnabled('AWS')) {
      pluginsPayload.aws = cloudProviders.map((provider) => {
        const trimmedAccessKeyId = provider.accessKeyId?.trim() ?? '';
        const hasSecretInput = provider.secretAccessKeyInput !== undefined;
        const trimmedSecret = provider.secretAccessKeyInput?.trim() ?? '';

        return {
          name: provider.name,
          region: provider.region,
          accessKeyId:
            trimmedAccessKeyId.length > 0 ? trimmedAccessKeyId : undefined,
          secretAccessKey:
            trimmedSecret.length > 0 ? trimmedSecret : undefined,
          clearSecretAccessKey:
            hasSecretInput &&
            trimmedSecret === '' &&
            provider.hasSecretAccessKey
              ? true
              : undefined,
        };
      });
    } else if (config.aws && config.aws.length > 0) {
      pluginsPayload.aws = [];
    }

    if (isPluginEnabled('ServiceCatalog')) {
      pluginsPayload.service_catalog = {
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
      };
    } else if (config.service_catalog) {
      pluginsPayload.service_catalog = null;
    }

    if (isPluginEnabled('Observability')) {
      const mapProviders = (providers: ObservabilityProviderFormValue[]) =>
        providers.map((provider) => {
          const hasPasswordInput = provider.passwordInput !== undefined;
          const hasTokenInput = provider.tokenInput !== undefined;
          const trimmedPassword = provider.passwordInput?.trim() ?? '';
          const trimmedToken = provider.tokenInput?.trim() ?? '';

          const auth =
            provider.authType != null
              ? {
                  type: provider.authType,
                  username: provider.username,
                  password:
                    trimmedPassword.length > 0 ? trimmedPassword : undefined,
                  token: trimmedToken.length > 0 ? trimmedToken : undefined,
                  clearPassword:
                    hasPasswordInput &&
                    trimmedPassword === '' &&
                    provider.hasAuthPassword
                      ? true
                      : undefined,
                  clearToken:
                    hasTokenInput &&
                    trimmedToken === '' &&
                    provider.hasAuthToken
                      ? true
                      : undefined,
                }
              : undefined;

          return {
            name: provider.name,
            type: provider.type,
            url: provider.url,
            timeout: provider.timeout,
            retention: provider.retention,
            enabled: provider.enabled ?? true,
            labels:
              provider.labels && Object.keys(provider.labels).length > 0
                ? provider.labels
                : undefined,
            auth,
          };
        });

      pluginsPayload.observability = {
        enabled: observability.enabled,
        logs:
          observability.logsProviders.length > 0
            ? mapProviders(observability.logsProviders)
            : undefined,
        traces:
          observability.tracesProviders.length > 0
            ? mapProviders(observability.tracesProviders)
            : undefined,
        metrics:
          observability.metricsProviders.length > 0
            ? mapProviders(observability.metricsProviders)
            : undefined,
      };
    } else if (config.observability) {
      pluginsPayload.observability = null;
    }

    const updatePayload: UpdateSettingsRequest = {
      enabled_plugins: enabledPlugins,
      config: {
        port: generalSettings.port,
        origin: generalSettings.origin,
        headers: headersToArray(generalSettings.headers),
        front: generalSettings.front,
      },
    };

    if (Object.keys(pluginsPayload).length > 0) {
      updatePayload.plugins = pluginsPayload;
    }

    try {
      await save(updatePayload);
      toast.success('Settings saved successfully');
    } catch (saveError) {
      toast.error(
        saveError instanceof Error
          ? saveError.message
          : 'Failed to update settings'
      );
    }
  };

  if (loading || !config) {
    return (
      <div className="flex h-full items-center justify-center">
        <span className="flex items-center gap-2 text-sm text-muted-foreground">
          <Loader2 className="h-4 w-4 animate-spin" />
          Loading configuration...
        </span>
      </div>
    );
  }

  return (
    <div className="flex h-full flex-col overflow-hidden">
      <div className="flex items-center justify-between border-b bg-background px-6 py-5">
        <div className="space-y-1">
          <h1 className="text-2xl font-semibold tracking-tight flex items-center gap-2">
            <LayoutGrid className="h-5 w-5 text-primary" />
            Settings
          </h1>
          <p className="text-sm text-muted-foreground">
            Review and adjust your DashOps configuration. Changes require a
            backend restart to take effect.
          </p>
        </div>
        <Button
          onClick={handleSave}
          disabled={saving || !canEdit}
          className="gap-2"
        >
          {saving && <Loader2 className="h-4 w-4 animate-spin" />}
          <Save className="h-4 w-4" />
          Save changes
        </Button>
      </div>

      <div className="flex-1 overflow-auto p-6 space-y-6">
        <GeneralSettingsForm
          value={generalSettings}
          onChange={setGeneralSettings}
          disabled={!canEdit}
        />

      <Card>
        <CardHeader>
            <CardTitle>Enabled plugins</CardTitle>
          <CardDescription>
              Toggle which modules are available. Configuration tabs will unlock
              once enabled.
          </CardDescription>
        </CardHeader>
          <CardContent>
            <PluginToggleGrid
              options={PLUGIN_OPTIONS}
              selected={enabledPlugins}
              onToggle={handleTogglePlugin}
              disabled={!canEdit}
            />
          </CardContent>
        </Card>

        <Tabs
          value={activeTab}
          onValueChange={setActiveTab}
          className="space-y-6"
        >
          <TabsList className="flex w-full justify-start gap-2 overflow-x-auto">
            <TabsTrigger value="auth" disabled={!isPluginEnabled('Auth')}>
              <Shield className="mr-2 h-4 w-4" />
              Auth
            </TabsTrigger>
            <TabsTrigger
              value="kubernetes"
              disabled={!isPluginEnabled('Kubernetes')}
            >
              <Container className="mr-2 h-4 w-4" />
              Kubernetes
            </TabsTrigger>
            <TabsTrigger value="cloud" disabled={!isPluginEnabled('AWS')}>
              <Cloud className="mr-2 h-4 w-4" />
              Cloud
            </TabsTrigger>
            <TabsTrigger
              value="catalog"
              disabled={!isPluginEnabled('ServiceCatalog')}
            >
              <Package className="mr-2 h-4 w-4" />
              Catalog
            </TabsTrigger>
            <TabsTrigger
              value="observability"
              disabled={!isPluginEnabled('Observability')}
            >
              <Activity className="mr-2 h-4 w-4" />
              Observability
            </TabsTrigger>
          </TabsList>

          <TabsContent value="auth">
            {!isPluginEnabled('Auth') ? (
              <DisabledTabNotice />
            ) : (
              <div className="space-y-6">
              <div className="flex items-center justify-between">
                <div>
                    <h3 className="text-lg font-semibold">
                      Authentication providers
                    </h3>
                    <p className="text-sm text-muted-foreground">
                      Configure OAuth providers used for login.
                    </p>
                  </div>
                  <div className="flex items-center gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      className="gap-2"
                      onClick={() => handleAddAuthProvider('github')}
                      disabled={!canEdit}
                    >
                      <Github className="h-4 w-4" />
                      Add GitHub
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      className="gap-2"
                      disabled
                      title="Coming soon"
                    >
                      <Chrome className="h-4 w-4" />
                      Add Google
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      className="gap-2"
                      disabled
                      title="Coming soon"
                    >
                      <Building2 className="h-4 w-4" />
                      Add SSO
                    </Button>
                  </div>
                </div>
                <div className="space-y-4">
                  {authProviders.length === 0 && (
                    <div className="rounded-lg border border-dashed p-8 text-center text-sm text-muted-foreground">
                      No authentication providers configured yet.
                    </div>
                  )}
                  {authProviders.map((provider) => (
                    <AuthProviderCard
                      key={provider.id}
                      provider={provider}
                      onChange={(updated) =>
                        setAuthProviders((prev) =>
                          prev.map((current) =>
                            current.id === provider.id ? updated : current
                          )
                        )
                      }
                      onRemove={
                        authProviders.length === 1
                          ? undefined
                          : () =>
                              setAuthProviders((prev) =>
                                prev.filter((item) => item.id !== provider.id)
                              )
                      }
                      secretsVisible={false}
                      disabled={!canEdit}
                    />
                  ))}
                </div>
              </div>
            )}
          </TabsContent>

          <TabsContent value="kubernetes">
            {!isPluginEnabled('Kubernetes') ? (
              <DisabledTabNotice />
            ) : (
              <div className="space-y-6">
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="text-lg font-semibold">Clusters</h3>
                    <p className="text-sm text-muted-foreground">
                      Manage connected Kubernetes clusters.
                    </p>
                  </div>
                  <Button
                    variant="outline"
                    size="sm"
                    className="gap-2"
                    onClick={handleAddK8sCluster}
                    disabled={!canEdit}
                  >
                    <Plus className="h-4 w-4" />
                    Add cluster
                  </Button>
            </div>
                <div className="space-y-4">
                  {k8sClusters.length === 0 && (
                    <div className="rounded-lg border border-dashed p-8 text-center text-sm text-muted-foreground">
                      No clusters configured yet.
            </div>
                  )}
                  {k8sClusters.map((cluster) => (
                    <KubernetesClusterCard
                      key={cluster.id}
                      cluster={cluster}
                      onChange={(updated) =>
                        setK8sClusters((prev) =>
                          prev.map((current) =>
                            current.id === cluster.id ? updated : current
                          )
                        )
                      }
                      onRemove={
                        k8sClusters.length > 1 && canEdit
                          ? () => handleRemoveK8sCluster(cluster.id)
                          : undefined
                      }
                      disabled={!canEdit}
                    />
                  ))}
                </div>
              </div>
            )}
          </TabsContent>

          <TabsContent value="cloud">
            {!isPluginEnabled('AWS') ? (
              <DisabledTabNotice />
            ) : (
              <div className="space-y-6">
                <div className="flex items-center justify-between">
                  <div>
                    <h3 className="text-lg font-semibold">AWS accounts</h3>
                    <p className="text-sm text-muted-foreground">
                      Manage AWS credentials used by DashOps.
                    </p>
                  </div>
                  <Button
                    variant="outline"
                    size="sm"
                    className="gap-2"
                    onClick={handleAddCloudProvider}
                    disabled={!canEdit}
                  >
                    <Plus className="h-4 w-4" />
                    Add account
                  </Button>
                </div>
                <div className="space-y-4">
                  {cloudProviders.length === 0 && (
                    <div className="rounded-lg border border-dashed p-8 text-center text-sm text-muted-foreground">
                      No AWS accounts configured yet.
            </div>
          )}
                  {cloudProviders.map((account) => (
                    <AwsAccountCard
                      key={account.id}
                      account={account}
                      onChange={(updated) =>
                        setCloudProviders((prev) =>
                          prev.map((current) =>
                            current.id === account.id ? updated : current
                          )
                        )
                      }
                      onRemove={
                        cloudProviders.length === 1
                          ? undefined
                          : () =>
                              setCloudProviders((prev) =>
                                prev.filter((item) => item.id !== account.id)
                              )
                      }
                      secretsVisible={false}
                      disabled={!canEdit}
                    />
                  ))}
                </div>
              </div>
            )}
          </TabsContent>

          <TabsContent value="catalog">
            {!isPluginEnabled('ServiceCatalog') ? (
              <DisabledTabNotice />
            ) : (
              <ServiceCatalogForm
                value={serviceCatalog}
                onChange={setServiceCatalog}
                disabled={!canEdit}
              />
            )}
          </TabsContent>

          <TabsContent value="observability">
            {!isPluginEnabled('Observability') ? (
              <DisabledTabNotice />
            ) : (
              <div className="space-y-6">
                <div className="flex items-center justify-between rounded-lg border p-4">
                  <div>
                    <h3 className="text-sm font-medium">Enable observability</h3>
          <p className="text-xs text-muted-foreground">
                      Toggle to activate logs, traces and metrics integrations.
                    </p>
                  </div>
                  <Switch
                    checked={observability.enabled}
                    onCheckedChange={(enabled) =>
                      setObservability((prev) => ({ ...prev, enabled }))
                    }
                    disabled={!canEdit}
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
                  onAdd={() =>
                    handleAddObservabilityProvider('logsProviders', 'loki')
                  }
                  onRemove={(id) =>
                    handleRemoveObservabilityProvider('logsProviders', id)
                  }
                  disabled={!canEdit}
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
                  disabled={!canEdit}
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
                  disabled={!canEdit}
                />
              </div>
            )}
          </TabsContent>
        </Tabs>
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
  disabled?: boolean;
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
  disabled = false,
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
          disabled={disabled}
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
              providers.length > 1 && !disabled
                ? (id) => onRemove(id)
                : undefined
            }
            disabled={disabled}
          />
        ))}
      </div>
    </div>
  );
}

