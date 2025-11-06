import { useState } from 'react';
import { format, subMinutes, subHours, subDays } from 'date-fns';
import { Calendar as CalendarIcon } from 'lucide-react';
import { type DateRange } from 'react-day-picker';
import { Button } from '@/components/ui/button';
import { Calendar } from '@/components/ui/calendar';
import { Input } from '@/components/ui/input';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { cn } from '@/lib/utils';

export type TimeRange = {
  value: string; // '5m' | '15m' | '1h' | '6h' | '24h' | '7d' | 'custom'
  label: string;
  from?: Date;
  to?: Date;
};

interface TimeRangePickerProps {
  value: TimeRange;
  onChange: (range: TimeRange) => void;
  className?: string;
}

const PRESET_RANGES = [
  { value: '5m', label: 'Last 5 minutes' },
  { value: '15m', label: 'Last 15 minutes' },
  { value: '1h', label: 'Last hour' },
  { value: '6h', label: 'Last 6 hours' },
  { value: '24h', label: 'Last 24 hours' },
  { value: '7d', label: 'Last 7 days' },
] as const;

export default function TimeRangePicker({ value, onChange, className }: TimeRangePickerProps) {
  const [isOpen, setIsOpen] = useState(false);
  const [customFrom, setCustomFrom] = useState<Date>();
  const [customTo, setCustomTo] = useState<Date>();
  const [fromTime, setFromTime] = useState('00:00');
  const [toTime, setToTime] = useState('00:00');

  const handlePresetSelect = (presetValue: string) => {
    const now = new Date();
    let from: Date;

    switch (presetValue) {
      case '5m':
        from = subMinutes(now, 5);
        break;
      case '15m':
        from = subMinutes(now, 15);
        break;
      case '1h':
        from = subHours(now, 1);
        break;
      case '6h':
        from = subHours(now, 6);
        break;
      case '24h':
        from = subHours(now, 24);
        break;
      case '7d':
        from = subDays(now, 7);
        break;
      default:
        return;
    }

    onChange({
      value: presetValue,
      label: getPresetLabel(presetValue),
      from,
      to: now,
    });

    setIsOpen(false);
  };

  const handleCustomApply = () => {
    if (customFrom && customTo) {
      // Combine date and time
      const fromParts = fromTime.split(':');
      const toParts = toTime.split(':');
      const fromHours = Number(fromParts[0]) || 0;
      const fromMinutes = Number(fromParts[1]) || 0;
      const toHours = Number(toParts[0]) || 0;
      const toMinutes = Number(toParts[1]) || 0;
      
      const fromDate = new Date(customFrom);
      fromDate.setHours(fromHours, fromMinutes, 0, 0);
      
      const toDate = new Date(customTo);
      toDate.setHours(toHours, toMinutes, 0, 0);
      
      onChange({
        value: 'custom',
        label: `${format(fromDate, 'MMM d, HH:mm')} - ${format(toDate, 'MMM d, HH:mm')}`,
        from: fromDate,
        to: toDate,
      });
      setIsOpen(false);
    }
  };

  const getPresetLabel = (presetValue: string): string => {
    const labels: Record<string, string> = {
      '5m': 'Last 5 minutes',
      '15m': 'Last 15 minutes',
      '1h': 'Last hour',
      '6h': 'Last 6 hours',
      '24h': 'Last 24 hours',
      '7d': 'Last 7 days',
      custom: value.from && value.to
        ? `${format(value.from, 'MMM d, HH:mm')} - ${format(value.to, 'MMM d, HH:mm')}`
        : 'Custom Range',
    };
    return labels[presetValue] || 'Custom Range';
  };

  return (
    <Popover open={isOpen} onOpenChange={setIsOpen}>
      <PopoverTrigger>
        <Button
          variant="outline"
          size="sm"
          className={cn(
            "gap-2",
            className
          )}
        >
          <CalendarIcon className="h-4 w-4" />
          {getPresetLabel(value.value)}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-auto p-0" align="start">
        <div className="flex">
          {/* Preset Ranges */}
          <div className="p-2 border-r w-40">
            <div className="text-xs font-medium mb-1.5">Presets</div>
            <div className="space-y-0.5">
              {PRESET_RANGES.map((preset) => (
                <Button
                  key={preset.value}
                  variant={value.value === preset.value ? "default" : "ghost"}
                  size="sm"
                  className="w-full justify-start h-7 px-2 py-0.5 text-xs"
                  onClick={() => handlePresetSelect(preset.value)}
                >
                  {preset.label}
                </Button>
              ))}
            </div>
          </div>

          {/* Custom Range */}
          <div className="p-3">
            <div className="text-sm font-medium mb-2">Custom Range</div>
            <Calendar
              mode="range"
              selected={{ from: customFrom, to: customTo }}
              onSelect={(range: DateRange | undefined) => {
                if (range?.from) {
                  setCustomFrom(range.from);
                  setFromTime(format(range.from, 'HH:mm'));
                }
                if (range?.to) {
                  setCustomTo(range.to);
                  setToTime(format(range.to, 'HH:mm'));
                }
              }}
              numberOfMonths={1}
              className="rounded-md border-0"
            />
            <div className="mt-3 flex items-center gap-2">
              <div className="flex-1">
                <label className="text-xs text-muted-foreground mb-1 block">From</label>
                <Input
                  type="time"
                  value={fromTime}
                  onChange={(e) => setFromTime(e.target.value)}
                  className="h-8 text-xs"
                />
              </div>
              <div className="flex-1">
                <label className="text-xs text-muted-foreground mb-1 block">To</label>
                <Input
                  type="time"
                  value={toTime}
                  onChange={(e) => setToTime(e.target.value)}
                  className="h-8 text-xs"
                />
              </div>
            </div>
            <div className="mt-3 flex items-center justify-end gap-2 border-t pt-3">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setIsOpen(false)}
              >
                Cancel
              </Button>
              <Button
                size="sm"
                onClick={handleCustomApply}
                disabled={!customFrom || !customTo}
              >
                Apply
              </Button>
            </div>
          </div>
        </div>
      </PopoverContent>
    </Popover>
  );
}

