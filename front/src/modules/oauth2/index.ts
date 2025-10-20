import LoginPage from './components/pages/LoginPage';
import ProfilePage from './components/pages/ProfilePage';
import { AuthModule } from '@/types';
import { createElement } from 'react';

const authModule: AuthModule = {
  auth: {
    active: true,
    LoginPage,
  },
  routers: [
    {
      key: 'profile',
      path: '/profile',
      element: createElement(ProfilePage),
    },
  ],
};

export default authModule;

// Export only public types for external use
export * from './types';
