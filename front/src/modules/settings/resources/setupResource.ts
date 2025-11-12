import http from '@/helpers/http';
import {
  SetupConfigureRequest,
  SetupConfigureResponse,
} from '../types';

export async function configureSetup(
  payload: SetupConfigureRequest
): Promise<SetupConfigureResponse['data']> {
  const response = await http.post<SetupConfigureResponse>(
    '/settings/setup/configure',
    payload
  );

  if (!response.data.success) {
    throw new Error(response.data.error || 'Failed to configure setup');
  }

  return response.data.data;
}
