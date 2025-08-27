import { Palette, Check } from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { useTheme, type ThemeColor } from '@/contexts/ThemeContext';

interface ThemeOption {
  value: ThemeColor;
  label: string;
  preview: string; // CSS color for preview dot
}

const themeOptions: ThemeOption[] = [
  { value: 'neutral', label: 'Neutral', preview: 'oklch(0.205 0 0)' },
  { value: 'red', label: 'Red', preview: 'oklch(0.627 0.202 29.234)' },
  { value: 'rose', label: 'Rose', preview: 'oklch(0.646 0.188 12.178)' },
  { value: 'orange', label: 'Orange', preview: 'oklch(0.646 0.222 41.116)' },
  { value: 'green', label: 'Green', preview: 'oklch(0.543 0.137 145.224)' },
  { value: 'blue', label: 'Blue', preview: 'oklch(0.554 0.186 263.389)' },
  { value: 'yellow', label: 'Yellow', preview: 'oklch(0.769 0.188 70.08)' },
  { value: 'violet', label: 'Violet', preview: 'oklch(0.571 0.19 303.9)' },
  { value: 'slate', label: 'Slate', preview: 'oklch(0.208 0.042 265.755)' },
];

export function ThemeSelector() {
  const { themeColor, setThemeColor } = useTheme();

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="ghost"
          size="icon"
          className="h-9 w-9"
          title="Change theme color"
        >
          <Palette className="h-4 w-4" />
          <span className="sr-only">Change theme color</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-48">
        {themeOptions.map((option) => (
          <DropdownMenuItem
            key={option.value}
            onClick={() => setThemeColor(option.value)}
            className="flex items-center justify-between cursor-pointer"
          >
            <div className="flex items-center gap-3">
              <div
                className="w-4 h-4 rounded-full border border-border/50"
                style={{ backgroundColor: option.preview }}
                title={`${option.label} theme`}
              />
              <span>{option.label}</span>
            </div>
            {themeColor === option.value && <Check className="h-4 w-4" />}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
