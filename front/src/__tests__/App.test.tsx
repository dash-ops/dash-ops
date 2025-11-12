import { render, screen, cleanup, act, waitFor } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import * as userResource from '../modules/oauth2/resources/userResource';
import * as oauth from '../helpers/oauth';
import * as modulesHelper from '../helpers/loadModules';
import App from '../App';

vi.mock('../modules/oauth2/resources/userResource');
vi.mock('../helpers/oauth');
vi.mock('../helpers/loadModules');

afterEach(cleanup);

it('should render app footer when logged in user', async () => {
  oauth.verifyToken.mockReturnValue(true);
  userResource.getUserData.mockResolvedValue({ name: 'Bla' });
  modulesHelper.loadModulesConfig.mockResolvedValue({
    auth: { active: false },
    menus: [],
    routers: [],
    setupMode: false,
  });

  await act(async () => {
    render(
      <MemoryRouter>
        <App />
      </MemoryRouter>
    );
  });

  const footer = await screen.findByRole('contentinfo');
  expect(footer).toBeInTheDocument();
  await waitFor(() =>
    expect(
      screen.getByRole('link', { name: 'DashOps Repository' })
    ).toBeInTheDocument()
  );
});
