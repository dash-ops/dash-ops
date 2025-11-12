import type { LucideIcon } from 'lucide-react';
import {
  Activity,
  Cloud,
  Container,
  Package,
  Shield,
} from 'lucide-react';

export interface PluginOption {
  id: string;
  label: string;
  description: string;
  icon: LucideIcon;
}

export const PLUGIN_OPTIONS: PluginOption[] = [
  {
    id: 'Auth',
    label: 'Authentication',
    description: 'GitHub OAuth and login flow',
    icon: Shield,
  },
  {
    id: 'ServiceCatalog',
    label: 'Service Catalog',
    description: 'Service registry and metadata',
    icon: Package,
  },
  {
    id: 'Kubernetes',
    label: 'Kubernetes',
    description: 'Cluster inventory and workloads',
    icon: Container,
  },
  {
    id: 'AWS',
    label: 'AWS',
    description: 'Cloud accounts and instances',
    icon: Cloud,
  },
  {
    id: 'Observability',
    label: 'Observability',
    description: 'Logs, traces, and metrics providers',
    icon: Activity,
  },
];

