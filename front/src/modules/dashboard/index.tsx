import { LayoutDashboard } from 'lucide-react';
import DashboardPage from './DashboardPage';
import { ModuleConfig } from '@/types';

const dashboardModule: ModuleConfig = {
  menus: [
    {
      label: 'Dashboard',
      icon: <LayoutDashboard className="h-4 w-4" />,
      key: 'dashboard',
      link: '/',
    },
  ],
  routers: [
    {
      key: 'dashboard',
      path: '/',
      element: <DashboardPage />,
    },
  ],
};

export default dashboardModule;
