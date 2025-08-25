import { Container } from 'lucide-react';
import NodesPage from './NodesPage';
import DeploymentPage from './DeploymentPage';
import PodPage from './PodPage';
import PodLogPage from './PodLogPage';
import KubernetesWithContextSelector from './KubernetesWithContextSelector';
import { Menu, Page, ModuleConfig } from '@/types';

const KubernetesModule = async (): Promise<ModuleConfig> => {
  const menus: Menu[] = [
    {
      label: 'Kubernetes',
      icon: <Container className="h-4 w-4" />,
      key: 'kubernetes',
      link: '/k8s',
    },
  ];

  const pages: Page[] = [
    {
      name: 'Nodes',
      path: '/k8s/:context',
      menu: true,
      element: <NodesPage />,
    },
    {
      name: 'Deployments',
      path: '/k8s/:context/deployments',
      menu: true,
      element: <DeploymentPage />,
    },
    {
      name: 'Pods',
      path: '/k8s/:context/pods',
      menu: true,
      element: <PodPage />,
    },
    {
      name: 'PodLogs',
      path: '/k8s/:context/pod/logs',
      menu: false,
      element: <PodLogPage />,
    },
  ];

  return {
    menus,
    routers: [
      {
        key: 'k8s-root',
        path: '/k8s',
        element: <KubernetesWithContextSelector pages={pages} />,
      },
      {
        key: 'k8s',
        path: '/k8s/:context/*',
        element: <KubernetesWithContextSelector pages={pages} />,
      },
    ],
  };
};

export default KubernetesModule;
