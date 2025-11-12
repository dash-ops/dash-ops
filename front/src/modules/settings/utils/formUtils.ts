import type {
  AuthProviderFormValue,
  AuthProviderType,
  CloudProviderFormValue,
  CloudProviderType,
  GeneralSettingsFormValue,
  KubernetesClusterFormValue,
  ObservabilityFormValue,
  ObservabilityProviderFormValue,
  ServiceCatalogFormValue,
} from '../types';

export const SECRET_PLACEHOLDER = '********';

export const createId = (): string =>
  typeof crypto !== 'undefined' && 'randomUUID' in crypto
    ? crypto.randomUUID()
    : `${Date.now().toString(36)}-${Math.random().toString(36).slice(2, 8)}`;

export const headersToString = (headers?: string[]): string =>
  headers && headers.length > 0 ? headers.join(', ') : '';

export const headersToArray = (headers: string): string[] =>
  headers
    .split(',')
    .map((value) => value.trim())
    .filter(Boolean);

export const createDefaultAuthProvider = (
  provider: AuthProviderType = 'github',
  overrides: Partial<AuthProviderFormValue> = {}
): AuthProviderFormValue => {
  const names: Record<AuthProviderType, string> = {
    github: 'GitHub OAuth',
    google: 'Google OAuth',
    sso: 'SSO',
  };

  return {
    id: overrides.id ?? createId(),
    provider,
    name: overrides.name ?? names[provider],
    enabled: overrides.enabled ?? true,
    clientId: overrides.clientId ?? '',
    clientIdMasked: overrides.clientIdMasked,
    orgPermission: overrides.orgPermission ?? '',
    hasClientSecret: overrides.hasClientSecret ?? false,
    clientSecretInput: overrides.clientSecretInput,
  };
};

export const createDefaultKubernetesCluster = (
  overrides: Partial<KubernetesClusterFormValue> = {}
): KubernetesClusterFormValue => ({
  id: overrides.id ?? createId(),
  name: overrides.name ?? 'production-cluster',
  connectionType: overrides.connectionType ?? 'kubeconfig',
  kubeconfig: overrides.kubeconfig ?? '$HOME/.kube/config',
  context: overrides.context ?? 'kind-dashops-dev',
  host: overrides.host ?? '',
  token: overrides.token,
  certificate: overrides.certificate ?? '',
  hasToken: overrides.hasToken ?? false,
});

export const createDefaultCloudProvider = (
  type: CloudProviderType = 'aws',
  overrides: Partial<CloudProviderFormValue> = {}
): CloudProviderFormValue => ({
  id: overrides.id ?? createId(),
  type,
  name:
    overrides.name ??
    (type === 'aws'
      ? 'AWS Account'
      : type === 'gcp'
      ? 'GCP Project'
      : 'Azure Subscription'),
  region: overrides.region ?? (type === 'aws' ? 'us-east-1' : undefined),
  accessKeyId: overrides.accessKeyId ?? '',
  accessKeyIdMasked: overrides.accessKeyIdMasked,
  hasSecretAccessKey: overrides.hasSecretAccessKey ?? false,
  secretAccessKeyInput: overrides.secretAccessKeyInput,
  projectId: overrides.projectId,
  serviceAccountKey: overrides.serviceAccountKey,
  subscriptionId: overrides.subscriptionId,
  tenantId: overrides.tenantId,
  clientId: overrides.clientId,
  clientSecretInput: overrides.clientSecretInput,
});

export const createDefaultServiceCatalog = (
  overrides: Partial<ServiceCatalogFormValue> = {}
): ServiceCatalogFormValue => ({
  storageProvider: overrides.storageProvider ?? 'filesystem',
  directory: overrides.directory ?? '../services',
  githubRepository: overrides.githubRepository ?? '',
  githubBranch: overrides.githubBranch ?? 'main',
  s3Bucket: overrides.s3Bucket ?? '',
  versioningEnabled: overrides.versioningEnabled ?? false,
});

export const createDefaultObservabilityProvider = (
  type: string,
  overrides: Partial<ObservabilityProviderFormValue> = {}
): ObservabilityProviderFormValue => ({
  id: overrides.id ?? createId(),
  name: overrides.name ?? `${type}-provider`,
  type: overrides.type ?? type,
  url: overrides.url ?? '',
  timeout: overrides.timeout ?? '30s',
  retention: overrides.retention ?? (type === 'prometheus' ? '90d' : '30d'),
  enabled: overrides.enabled ?? true,
  labels: overrides.labels ?? {},
  authType: overrides.authType,
  hasAuthPassword: overrides.hasAuthPassword ?? false,
  hasAuthToken: overrides.hasAuthToken ?? false,
  username: overrides.username ?? '',
  passwordInput: overrides.passwordInput ?? '',
  tokenInput: overrides.tokenInput ?? '',
});

export const createDefaultObservability = (
  overrides: Partial<ObservabilityFormValue> = {}
): ObservabilityFormValue => ({
  enabled: overrides.enabled ?? true,
  logsProviders:
    overrides.logsProviders ??
    [
      createDefaultObservabilityProvider('loki', {
        name: 'loki-local',
        url: 'http://localhost:30100',
      }),
    ],
  tracesProviders:
    overrides.tracesProviders ??
    [
      createDefaultObservabilityProvider('tempo', {
        name: 'tempo-local',
        url: 'http://localhost:30200',
        retention: '7d',
      }),
    ],
  metricsProviders:
    overrides.metricsProviders ??
    [
      createDefaultObservabilityProvider('prometheus', {
        name: 'prometheus-local',
        url: 'http://localhost:30090',
        retention: '90d',
      }),
    ],
});

export const createDefaultGeneralSettings = (
  overrides: Partial<GeneralSettingsFormValue> = {}
): GeneralSettingsFormValue => ({
  port: overrides.port ?? '8080',
  origin: overrides.origin ?? 'http://localhost:5173',
  headers: overrides.headers ?? 'Content-Type, Authorization',
  front: overrides.front ?? 'front/dist',
});

