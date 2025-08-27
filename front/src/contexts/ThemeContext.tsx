import {
  createContext,
  useContext,
  useEffect,
  useState,
  ReactNode,
} from 'react';

export type Theme = 'light' | 'dark';
export type ThemeColor =
  | 'neutral'
  | 'red'
  | 'rose'
  | 'orange'
  | 'green'
  | 'blue'
  | 'yellow'
  | 'violet'
  | 'slate';

interface ThemeContextType {
  theme: Theme;
  themeColor: ThemeColor;
  setTheme: (theme: Theme) => void;
  setThemeColor: (color: ThemeColor) => void;
  toggleTheme: () => void;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

interface ThemeProviderProps {
  children: ReactNode;
}

const STORAGE_KEYS = {
  theme: 'dashops-theme',
  color: 'dashops-theme-color',
} as const;

export function ThemeProvider({ children }: ThemeProviderProps) {
  const [theme, setTheme] = useState<Theme>(() => {
    // Check localStorage first
    const stored = localStorage.getItem(STORAGE_KEYS.theme) as Theme;
    if (stored && ['light', 'dark'].includes(stored)) {
      return stored;
    }
    // Fallback to system preference
    return window.matchMedia('(prefers-color-scheme: dark)').matches
      ? 'dark'
      : 'light';
  });

  const [themeColor, setThemeColor] = useState<ThemeColor>(() => {
    const stored = localStorage.getItem(STORAGE_KEYS.color) as ThemeColor;
    if (
      stored &&
      [
        'neutral',
        'red',
        'rose',
        'orange',
        'green',
        'blue',
        'yellow',
        'violet',
        'slate',
      ].includes(stored)
    ) {
      return stored;
    }
    return 'neutral';
  });

  // Apply theme class to document
  useEffect(() => {
    const root = document.documentElement;
    root.classList.remove('light', 'dark');
    root.classList.add(theme);

    // Save to localStorage
    localStorage.setItem(STORAGE_KEYS.theme, theme);
  }, [theme]);

  // Apply theme color data attribute to document
  useEffect(() => {
    const root = document.documentElement;
    root.setAttribute('data-theme-color', themeColor);

    // Save to localStorage
    localStorage.setItem(STORAGE_KEYS.color, themeColor);
  }, [themeColor]);

  const toggleTheme = () => {
    setTheme((prev) => (prev === 'light' ? 'dark' : 'light'));
  };

  return (
    <ThemeContext.Provider
      value={{
        theme,
        themeColor,
        setTheme,
        setThemeColor,
        toggleTheme,
      }}
    >
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  const context = useContext(ThemeContext);
  if (context === undefined) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }
  return context;
}
