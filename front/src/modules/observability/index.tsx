import { Activity } from 'lucide-react';
import ObservabilityPageV2 from './components/ObservabilityPageV2';
import { ModuleConfig } from '@/types';

const ObservabilityModule = async (): Promise<ModuleConfig> => {
  const menus = [
    {
      label: 'Observability',
      icon: <Activity className="h-4 w-4" />,
      key: 'observability',
      link: '/observability',
    },
  ];

  return {
    menus,
    routers: [
      {
        key: 'observability',
        path: '/observability',
        element: <ObservabilityPageV2 />,
      },
    ],
  };
};

export default ObservabilityModule;
export * from './types';

