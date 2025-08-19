import { AppstoreOutlined } from '@ant-design/icons';
import ServiceCatalogPage from './ServiceCatalogPage';

export default {
  menus: [
    {
      label: 'Services Catalog',
      icon: <AppstoreOutlined />,
      key: 'servicecatalog',
      link: '/servicecatalog',
    },
  ],
  routers: [
    {
      key: 'servicecatalog',
      path: '/servicecatalog',
      element: <ServiceCatalogPage />,
    },
  ],
};
