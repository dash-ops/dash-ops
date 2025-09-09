import LoginPage from './LoginPage';
import ProfilePage from './ProfilePage';
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
