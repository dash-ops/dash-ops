import http from '../../helpers/http';

export function getPlugins() {
  return http
    .get('/config/plugins')
    .then((resp) => (resp.data ? resp : { ...resp, data: [] }));
}
