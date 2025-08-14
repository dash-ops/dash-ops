import LoginPage from './LoginPage';
import ProfilePage from './ProfilePage';

export default {
  oAuth2: {
    active: true,
    LoginPage,
  },
  routers: [
    {
      key: 'profile',
      path: '/profile',
      element: <ProfilePage />,
    },
  ],
};
