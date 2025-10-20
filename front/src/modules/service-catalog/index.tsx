import { Layers3 } from 'lucide-react';
import { ServicesCatalogPage } from './components/ServicesCatalogPage';
import { Menu, ModuleConfig } from '@/types';

const ServiceCatalogModule = async (): Promise<ModuleConfig> => {
  const menus: Menu[] = [
    {
      label: 'Services Catalog',
      icon: <Layers3 className="h-4 w-4" />,
      key: 'service-catalog',
      link: '/services',
    },
  ];

  return {
    menus,
    routers: [
      {
        key: 'service-catalog',
        path: '/services',
        element: <ServicesCatalogPage />,
      },
      {
        key: 'service-catalog-wildcard',
        path: '/services/*',
        element: <ServicesCatalogPage />,
      },
    ],
  };
};

export default ServiceCatalogModule;

// Export only public types for external use
export * from './types';
