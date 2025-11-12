import { Settings } from 'lucide-react';
import SettingsPage from './components/SettingsPage';
import SetupPage from './components/SetupPage';
import { ModuleConfig } from '@/types';

interface SettingsModuleOptions {
  setupMode?: boolean;
}

const SetupPageWrapper = () => (
  <SetupPage onComplete={() => window.location.reload()} />
);

const loadSettingsModule = async (
  options: SettingsModuleOptions = {}
): Promise<ModuleConfig> => {
  if (options.setupMode) {
    return {
      menus: [],
      routers: [
        {
          key: 'setup-root',
          path: '/',
          element: <SetupPageWrapper />,
        },
        {
          key: 'setup',
          path: '*',
          element: <SetupPageWrapper />,
        },
      ],
    };
  }

  return {
    menus: [
      {
        label: 'Settings',
        icon: <Settings className="h-4 w-4" />,
        key: 'settings',
        link: '/settings',
      },
    ],
    routers: [
      {
        key: 'settings',
        path: '/settings',
        element: <SettingsPage />,
      },
    ],
  };
};

export default loadSettingsModule;
