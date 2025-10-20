import http from '../../../helpers/http';
import { getUserData } from '../resources/userResource';

vi.mock('../../../helpers/http');

it('should return logged in user', async () => {
  const mockResponse = {
    login: 'usergithub',
    id: 666,
    node_id: '666=',
    avatar_url: 'https://avatars1.githubusercontent.com/u/666?v=4',
    name: 'User Github',
    email: 'user.github@gmail.com',
    created_at: '2019-11-11T11:19:04Z',
    updated_at: '2019-12-20T19:08:00Z',
  };
  http.get.mockResolvedValue({
    data: mockResponse,
  });

  const resp = await getUserData();

  expect(http.get).toBeCalledWith(`/v1/me`);
  expect(resp.data).toEqual(mockResponse);
});
