import http from '@/helpers/http';
import {
  SettingsConfigPayload,
  SettingsConfigResponse,
  UpdateSettingsRequest,
  UpdateSettingsResponse,
} from '../types';

export async function getSettingsConfig(): Promise<SettingsConfigPayload> {
  const response = await http.get<SettingsConfigResponse>('/settings/config');
  if (!response.data.success) {
    throw new Error(response.data.error || 'Failed to load settings');
  }
  const { config, plugins, can_edit } = response.data.data;
  return {
    config,
    plugins,
    canEdit: can_edit,
  };
}

export async function updateSettingsConfig(
  payload: UpdateSettingsRequest
): Promise<UpdateSettingsResponse['data']> {
  const response = await http.put<UpdateSettingsResponse>(
    '/settings/config',
    payload
  );

  if (!response.data.success) {
    throw new Error(response.data.error || 'Failed to update settings');
  }

  return response.data.data;
}
