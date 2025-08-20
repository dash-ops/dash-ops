import http from '../../helpers/http';
import { AxiosResponse } from 'axios';
import { OAuth2Types } from '@/types';

export function getUserData(): Promise<AxiosResponse<OAuth2Types.UserData>> {
  return http.get('/v1/me');
}

export function getUserPermissions(): Promise<
  AxiosResponse<OAuth2Types.UserPermission[]>
> {
  return http.get('/v1/me/permissions');
}
