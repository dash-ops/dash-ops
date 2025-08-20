import LoginPage from './LoginPage';
import ProfilePage from './ProfilePage';
import { OAuth2Module } from '@/types';
import { createElement } from 'react';

const oauth2Module: OAuth2Module = {
  oAuth2: {
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

export default oauth2Module;
