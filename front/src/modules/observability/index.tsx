import { Activity } from 'lucide-react';
import LogsPage from './components/logs/LogsPage';
import TracesPage from './components/traces/TracesPage';
import { ModuleConfig } from '@/types';

const ObservabilityModule = async (): Promise<ModuleConfig> => {
  const menus = [
    {
      label: 'Observability',
      icon: <Activity className="h-4 w-4" />,
      key: 'observability',
      link: '/observability/logs',
    },
  ];

  return {
    menus,
    routers: [
      {
        key: 'observability-logs',
        path: '/observability/logs',
        element: <LogsPage />,
      },
      {
        key: 'observability-traces',
        path: '/observability/traces',
        element: <TracesPage />,
      },
    ],
  };
};

export default ObservabilityModule;
export * from './types';

