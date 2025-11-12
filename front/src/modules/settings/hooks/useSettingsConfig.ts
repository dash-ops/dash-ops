import { useCallback, useEffect, useState } from 'react';
import {
  getSettingsConfig,
  updateSettingsConfig,
} from '../resources/settingsResource';
import { SettingsConfig, UpdateSettingsRequest } from '../types';

interface UseSettingsConfigReturn {
  config: SettingsConfig | null;
  loading: boolean;
  saving: boolean;
  error: string | null;
  canEdit: boolean;
  refresh: () => Promise<void>;
  save: (payload: UpdateSettingsRequest) => Promise<void>;
}

export function useSettingsConfig(): UseSettingsConfigReturn {
  const [config, setConfig] = useState<SettingsConfig | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [saving, setSaving] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [canEdit, setCanEdit] = useState<boolean>(true);

  const refresh = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await getSettingsConfig();
      setConfig(data.config);
      setCanEdit(data.canEdit);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load settings');
    } finally {
      setLoading(false);
    }
  }, []);

  const save = useCallback(async (payload: UpdateSettingsRequest) => {
    try {
      setSaving(true);
      setError(null);
      await updateSettingsConfig(payload);
      await refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update settings');
    } finally {
      setSaving(false);
    }
  }, [refresh]);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  return { config, loading, saving, error, canEdit, refresh, save };
}
