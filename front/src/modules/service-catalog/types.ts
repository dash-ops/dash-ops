// Service Catalog Types - Frontend
// Based on backend types from pkg/service-catalog/types.go

export interface Service {
  apiVersion: string;
  kind: string;
  metadata: ServiceMetadata;
  spec: ServiceSpec;
}

export interface ServiceMetadata {
  name: string;
  tier: 'TIER-1' | 'TIER-2' | 'TIER-3';
  created_at?: string;
  created_by?: string;
  updated_at?: string;
  updated_by?: string;
  version?: number;
}

export interface ServiceSpec {
  description: string;
  team: ServiceTeam;
  business: ServiceBusiness;
  technology?: ServiceTechnology;
  kubernetes?: ServiceKubernetes;
  observability?: ServiceObservability;
  runbooks?: ServiceRunbook[];
}

export interface ServiceTeam {
  github_team: string;
}

export interface ServiceBusiness {
  sla_target?: string;
  dependencies?: string[];
  impact?: 'high' | 'medium' | 'low';
}

export interface ServiceTechnology {
  language?: string;
  framework?: string;
}

export interface ServiceKubernetes {
  environments: KubernetesEnvironment[];
}

export interface KubernetesEnvironment {
  name: string;
  context: string;
  namespace: string;
  resources: KubernetesEnvironmentResources;
}

export interface KubernetesEnvironmentResources {
  deployments: KubernetesDeployment[];
  services?: string[];
  configmaps?: string[];
}

export interface KubernetesDeployment {
  name: string;
  replicas: number;
  resources: KubernetesResourceRequests;
}

export interface KubernetesResourceRequests {
  requests: KubernetesResourceSpec;
  limits: KubernetesResourceSpec;
}

export interface KubernetesResourceSpec {
  cpu: string;
  memory: string;
}

export interface ServiceObservability {
  metrics?: string;
  logs?: string;
  traces?: string;
}

export interface ServiceRunbook {
  name: string;
  url: string;
}

export interface ServiceList {
  services: Service[];
  total: number;
}

export interface ServiceHealth {
  service_name: string;
  overall_status:
    | 'healthy'
    | 'degraded'
    | 'down'
    | 'critical'
    | 'drift'
    | 'unknown';
  environments: EnvironmentHealth[];
  last_updated: string;
}

export interface EnvironmentHealth {
  name: string;
  context: string;
  status: 'healthy' | 'degraded' | 'down' | 'drift';
  deployments: DeploymentHealth[];
}

export interface DeploymentHealth {
  name: string;
  ready_replicas: number;
  desired_replicas: number;
  status: string;
  last_updated: string;
}

// UI-specific types
export interface ServiceCardData {
  service: Service;
  health?: ServiceHealth;
  hasWriteAccess: boolean;
  isMyTeam: boolean;
}

export interface ServiceFilters {
  search: string;
  tier: string;
  team: string;
  status: string;
  sortBy: 'name' | 'tier' | 'team' | 'updated_at';
}

export interface ServiceStats {
  total: number;
  myTeam: number;
  tier1: number;
  tier2: number;
  tier3: number;
  critical: number;
  editable: number;
}

// Modal-specific types
export interface ServiceFormModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onServiceCreated?: () => void;
  editingService?: Service | undefined; // When provided, modal is in edit mode
}

export interface ServiceFormData {
  // Basic info
  name: string;
  description: string;
  tier: 'TIER-1' | 'TIER-2' | 'TIER-3';

  // Team
  github_team: string;

  // Business (optional)
  impact: 'high' | 'medium' | 'low';
  sla_target: string;

  // Technology (optional)
  language: string;
  framework: string;

  // Kubernetes environment
  env_name: string;
  env_context: string;
  env_namespace: string;

  // Single deployment for simplicity
  deployment_name: string;
  deployment_replicas: number;
  cpu_request: string;
  memory_request: string;
  cpu_limit: string;
  memory_limit: string;

  // Observability (optional)
  metrics_url: string;
  logs_url: string;
}
