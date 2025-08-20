import { render, screen, cleanup, act } from '@testing-library/react';
import { MemoryRouter } from 'react-router';
import * as userResource from '../modules/oauth2/userResource';
import * as oauth from '../helpers/oauth';
import App from '../App';

vi.mock('../modules/oauth2/userResource');
vi.mock('../helpers/oauth');

afterEach(cleanup);

it('should render app footer when logged in user', async () => {
  oauth.verifyToken.mockReturnValue(true);
  userResource.getUserData.mockResolvedValue({ name: 'Bla' });

  await act(async () => {
    render(
      <MemoryRouter>
        <App />
      </MemoryRouter>
    );
  });

  const footer = screen.getByRole('contentinfo');
  expect(footer).toBeInTheDocument();
  expect(
    screen.getByRole('link', { name: 'DashOps Repository' })
  ).toBeInTheDocument();
});
