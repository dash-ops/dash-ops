import { BaseEntity } from '@/types/api';

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  error?: string;
}

export type AuthProviderType = 'github' | 'google' | 'sso';

export interface AuthSetupConfig {
  provider: string;
  clientId?: string;
  clientSecret?: string;
  orgPermission?: string;
  redirectURL?: string;
  scopes?: string[];
}

export interface KubernetesSetupConfig {
  name: string;
  kubeconfig?: string;
  context?: string;
  connectionType?: 'kubeconfig' | 'remote';
  host?: string;
  token?: string;
  certificate?: string;
}

export interface AWSSetupConfig {
  name: string;
  region?: string;
  accessKeyId?: string;
  secretAccessKey?: string;
}

export interface ServiceCatalogStorageConfig {
  provider: 'filesystem' | 'github' | 's3';
  filesystem?: { directory: string };
  github?: { repository: string; branch?: string };
  s3?: { bucket: string };
}

export interface ServiceCatalogSetupConfig {
  storage: ServiceCatalogStorageConfig;
  versioning?: {
    enabled: boolean;
    provider?: string;
  };
}

export interface ProviderAuthConfig {
  type: string;
  username?: string;
  password?: string;
  token?: string;
  clearPassword?: boolean;
  clearToken?: boolean;
}

export interface ObservabilityProviderConfig {
  name: string;
  type: string;
  url: string;
  timeout?: string;
  retention?: string;
  enabled?: boolean;
  labels?: Record<string, string>;
  auth?: ProviderAuthConfig;
}

export interface ObservabilitySetupConfig {
  enabled: boolean;
  logs?: ObservabilityProviderConfig[];
  traces?: ObservabilityProviderConfig[];
  metrics?: ObservabilityProviderConfig[];
}

export interface SetupPluginsConfig {
  auth?: AuthSetupConfig | null;
  kubernetes?: KubernetesSetupConfig[] | null;
  aws?: AWSSetupConfig[] | null;
  service_catalog?: ServiceCatalogSetupConfig | null;
  observability?: ObservabilitySetupConfig | null;
}

export interface SetupConfigureRequest {
  config: {
    port: string;
    origin: string;
    headers: string[];
    front: string;
  };
  enabled_plugins: string[];
  plugins: SetupPluginsConfig;
}

export type SetupConfigureResponse = ApiResponse<{
  success: boolean;
  message: string;
  config_path?: string;
}>;

export interface AuthProviderSummary {
  provider: string;
  clientIdMasked?: string;
  orgPermission?: string;
  redirectURL?: string;
  scopes?: string[];
  hasClientSecret: boolean;
}

export interface KubernetesClusterSummary {
  name: string;
  connectionType?: 'kubeconfig' | 'remote';
  kubeconfig?: string;
  context?: string;
  host?: string;
  certificate?: string;
  hasToken: boolean;
}

export interface AwsAccountSummary {
  name: string;
  region?: string;
  accessKeyIdMasked?: string;
  hasSecretAccessKey: boolean;
}

export interface ObservabilityProviderAuthSummary {
  type: string;
  hasUsername: boolean;
  hasPassword: boolean;
  hasToken: boolean;
}

export interface ObservabilityProviderSummary {
  name: string;
  type: string;
  url: string;
  timeout?: string;
  retention?: string;
  enabled: boolean;
  labels?: Record<string, string>;
  auth?: ObservabilityProviderAuthSummary;
}

export interface ObservabilitySummary {
  enabled: boolean;
  logs?: ObservabilityProviderSummary[];
  traces?: ObservabilityProviderSummary[];
  metrics?: ObservabilityProviderSummary[];
}

export interface SettingsConfig {
  port: string;
  origin: string;
  headers: string[];
  front: string;
  plugins: string[];
  auth?: AuthProviderSummary[] | null;
  kubernetes?: KubernetesClusterSummary[] | null;
  aws?: AwsAccountSummary[] | null;
  service_catalog?: ServiceCatalogSetupConfig | null;
  observability?: ObservabilitySummary | null;
}

export interface Plugin extends BaseEntity {
  enabled: boolean;
}

export interface PluginsResponse {
  data: string[];
}

export type SettingsConfigResponse = ApiResponse<{
  config: SettingsConfig;
  plugins: string[];
  can_edit: boolean;
}>;

export interface SettingsConfigPayload {
  config: SettingsConfig;
  plugins: string[];
  canEdit: boolean;
}

export interface UpdateGeneralConfigPayload {
  port?: string;
  origin?: string;
  headers?: string[];
  front?: string;
}

export interface UpdateAuthProviderPayload {
  provider: string;
  clientId?: string;
  clientSecret?: string;
  orgPermission?: string;
  redirectURL?: string;
  scopes?: string[];
  clearClientSecret?: boolean;
}

export interface UpdateKubernetesClusterPayload {
  name: string;
  connectionType?: 'kubeconfig' | 'remote';
  kubeconfig?: string;
  context?: string;
  host?: string;
  token?: string;
  certificate?: string;
  clearToken?: boolean;
}

export interface UpdateAwsAccountPayload {
  name: string;
  region?: string;
  accessKeyId?: string;
  secretAccessKey?: string;
  clearSecretAccessKey?: boolean;
}

export type UpdateServiceCatalogPayload = ServiceCatalogSetupConfig;

export interface UpdateObservabilityPayload {
  enabled?: boolean;
  logs?: ObservabilityProviderConfig[];
  traces?: ObservabilityProviderConfig[];
  metrics?: ObservabilityProviderConfig[];
}

export interface UpdatePluginsPayload {
  auth?: UpdateAuthProviderPayload[] | null;
  kubernetes?: UpdateKubernetesClusterPayload[] | null;
  aws?: UpdateAwsAccountPayload[] | null;
  service_catalog?: UpdateServiceCatalogPayload | null;
  observability?: UpdateObservabilityPayload | null;
}

export interface UpdateSettingsRequest {
  enabled_plugins: string[];
  config?: UpdateGeneralConfigPayload;
  plugins?: UpdatePluginsPayload;
}

export type UpdateSettingsResponse = ApiResponse<{
  success: boolean;
  message: string;
  requires_restart?: boolean;
}>;

export interface AuthProviderFormValue {
  id: string;
  provider: AuthProviderType;
  name: string;
  enabled: boolean;
  clientId: string;
  clientIdMasked?: string;
  orgPermission?: string;
  hasClientSecret?: boolean;
  clientSecretInput?: string;
}

export interface KubernetesClusterFormValue {
  id: string;
  name: string;
  connectionType: 'kubeconfig' | 'remote';
  kubeconfig?: string;
  context?: string;
  host?: string;
  token?: string;
  certificate?: string;
  hasToken?: boolean;
}

export type CloudProviderType = 'aws' | 'gcp' | 'azure';

export interface CloudProviderFormValue {
  id: string;
  type: CloudProviderType;
  name: string;
  region?: string;
  accessKeyId?: string;
  accessKeyIdMasked?: string;
  hasSecretAccessKey?: boolean;
  secretAccessKeyInput?: string;
  projectId?: string;
  serviceAccountKey?: string;
  subscriptionId?: string;
  tenantId?: string;
  clientId?: string;
  clientSecretInput?: string;
}

export interface ServiceCatalogFormValue {
  storageProvider: 'filesystem' | 'github' | 's3';
  directory: string;
  githubRepository: string;
  githubBranch: string;
  s3Bucket: string;
  versioningEnabled: boolean;
}

export interface ObservabilityProviderFormValue {
  id: string;
  name: string;
  type: string;
  url: string;
  timeout?: string;
  retention?: string;
  enabled?: boolean;
  labels?: Record<string, string>;
  authType?: string;
  hasAuthPassword?: boolean;
  hasAuthToken?: boolean;
  username?: string;
  passwordInput?: string;
  tokenInput?: string;
}

export interface ObservabilityFormValue {
  enabled: boolean;
  logsProviders: ObservabilityProviderFormValue[];
  tracesProviders: ObservabilityProviderFormValue[];
  metricsProviders: ObservabilityProviderFormValue[];
}

export interface GeneralSettingsFormValue {
  port: string;
  origin: string;
  headers: string;
  front: string;
}

