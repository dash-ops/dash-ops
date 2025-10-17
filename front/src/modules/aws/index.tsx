import { Cloud } from 'lucide-react';
import InstancePage from './components/instances/InstancePage';
import AWSWithAccountSelector from './components/AWSWithAccountSelector';
import { Menu, Page, ModuleConfig } from '@/types';

const AwsModule = async (): Promise<ModuleConfig> => {
  const menus: Menu[] = [
    {
      label: 'AWS',
      icon: <Cloud className="h-4 w-4" />,
      key: 'aws',
      link: '/aws',
    },
  ];

  const pages: Page[] = [
    {
      name: 'EC2 Instances',
      path: '/aws/:key',
      menu: true,
      element: <InstancePage />,
    },
  ];

  return {
    menus,
    routers: [
      {
        key: 'aws-root',
        path: '/aws',
        element: <AWSWithAccountSelector pages={pages} />,
      },
      {
        key: 'aws',
        path: '/aws/:key/*',
        element: <AWSWithAccountSelector pages={pages} />,
      },
    ],
  };
};

export default AwsModule;

export * from './types';
