import http from '../../../helpers/http';
import { AxiosResponse } from 'axios';
import { AuthTypes } from '@/types';

export function getUserData(): Promise<AxiosResponse<AuthTypes.UserData>> {
  return http.get('/v1/me');
}

export function getUserPermissions(): Promise<
  AxiosResponse<AuthTypes.UserPermission[]>
> {
  return http.get('/v1/me/permissions');
}
