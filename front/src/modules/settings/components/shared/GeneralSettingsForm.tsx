import { Card } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import type { GeneralSettingsFormValue } from '../../types';

interface GeneralSettingsFormProps {
  value: GeneralSettingsFormValue;
  onChange: (value: GeneralSettingsFormValue) => void;
  disabled?: boolean;
  description?: string;
  title?: string;
}

export function GeneralSettingsForm({
  value,
  onChange,
  disabled = false,
  description = 'Core settings for your DashOps instance.',
  title = 'General configuration',
}: GeneralSettingsFormProps): JSX.Element {
  const handleChange = <T extends keyof GeneralSettingsFormValue>(
    field: T,
    fieldValue: GeneralSettingsFormValue[T]
  ) => {
    onChange({
      ...value,
      [field]: fieldValue,
    });
  };

  return (
    <Card>
      <div className="space-y-6 p-6">
        <div>
          <h2 className="text-lg font-semibold">{title}</h2>
          <p className="text-sm text-muted-foreground">{description}</p>
        </div>
        <div className="grid gap-4 md:grid-cols-2">
          <div className="space-y-2">
            <Label htmlFor="settings-port">Backend port</Label>
            <Input
              id="settings-port"
              value={value.port}
              onChange={(event) => handleChange('port', event.target.value)}
              disabled={disabled}
              placeholder="8080"
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="settings-origin">Allowed origin</Label>
            <Input
              id="settings-origin"
              value={value.origin}
              onChange={(event) => handleChange('origin', event.target.value)}
              disabled={disabled}
              placeholder="http://localhost:5173"
            />
          </div>
          <div className="space-y-2 md:col-span-2">
            <Label htmlFor="settings-headers">CORS headers</Label>
            <Input
              id="settings-headers"
              value={value.headers}
              onChange={(event) => handleChange('headers', event.target.value)}
              disabled={disabled}
              placeholder="Content-Type, Authorization"
            />
            <p className="text-xs text-muted-foreground">
              Comma separated list of allowed headers.
            </p>
          </div>
          <div className="space-y-2 md:col-span-2">
            <Label htmlFor="settings-front">Frontend build path</Label>
            <Input
              id="settings-front"
              value={value.front}
              onChange={(event) => handleChange('front', event.target.value)}
              disabled={disabled}
              placeholder="front/dist"
            />
          </div>
        </div>
      </div>
    </Card>
  );
}


