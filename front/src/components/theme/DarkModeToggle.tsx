import { Moon, Sun } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useTheme } from '@/contexts/ThemeContext';

export function DarkModeToggle() {
  const { theme, toggleTheme } = useTheme();

  return (
    <Button
      variant="ghost"
      size="icon"
      onClick={toggleTheme}
      className="h-9 w-9"
      title={theme === 'light' ? 'Switch to dark mode' : 'Switch to light mode'}
    >
      {theme === 'light' ? (
        <Sun className="h-4 w-4" />
      ) : (
        <Moon className="h-4 w-4" />
      )}
      <span className="sr-only">Toggle theme</span>
    </Button>
  );
}
