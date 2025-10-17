/**
 * Global UI component mocks for testing
 * These mocks can be reused across all modules
 */

import { vi } from 'vitest';

// Mock React Router
export const mockReactRouter = {
  useParams: vi.fn(() => ({})),
  useLocation: vi.fn(() => ({ pathname: '/' })),
  useNavigate: vi.fn(() => vi.fn()),
};

// Mock Sonner toast
export const mockToast = {
  error: vi.fn(),
  success: vi.fn(),
  info: vi.fn(),
  warning: vi.fn(),
};

// Mock UI components
export const mockUIComponents = {
  Button: ({ children, ...props }) => {
    const React = require('react');
    return React.createElement('button', props, children);
  },
  
  Input: (props) => {
    const React = require('react');
    return React.createElement('input', props);
  },
  
  Table: ({ children }) => {
    const React = require('react');
    return React.createElement('table', null, children);
  },
  
  TableBody: ({ children }) => {
    const React = require('react');
    return React.createElement('tbody', null, children);
  },
  
  TableCell: ({ children }) => {
    const React = require('react');
    return React.createElement('td', null, children);
  },
  
  TableHead: ({ children }) => {
    const React = require('react');
    return React.createElement('th', null, children);
  },
  
  TableHeader: ({ children }) => {
    const React = require('react');
    return React.createElement('thead', null, children);
  },
  
  TableRow: ({ children }) => {
    const React = require('react');
    return React.createElement('tr', null, children);
  },
  
  Select: ({ children }) => {
    const React = require('react');
    return React.createElement('select', null, children);
  },
  
  SelectContent: ({ children }) => {
    const React = require('react');
    return React.createElement('div', null, children);
  },
  
  SelectItem: ({ children, value }) => {
    const React = require('react');
    return React.createElement('option', { value }, children);
  },
  
  SelectTrigger: ({ children }) => {
    const React = require('react');
    return React.createElement('div', null, children);
  },
  
  SelectValue: ({ placeholder }) => {
    const React = require('react');
    return React.createElement('span', null, placeholder);
  },
  
  Badge: ({ children, ...props }) => {
    const React = require('react');
    return React.createElement('span', props, children);
  },
};

// Mock Lucide React icons
export const mockIcons = {
  Cloud: () => {
    const React = require('react');
    return React.createElement('div', { 'data-testid': 'cloud-icon' });
  },
  RefreshCw: () => {
    const React = require('react');
    return React.createElement('div', { 'data-testid': 'refresh-icon' });
  },
  Play: () => {
    const React = require('react');
    return React.createElement('div', { 'data-testid': 'play-icon' });
  },
  Square: () => {
    const React = require('react');
    return React.createElement('div', { 'data-testid': 'square-icon' });
  },
  RotateCcw: () => {
    const React = require('react');
    return React.createElement('div', { 'data-testid': 'restart-icon' });
  },
  Trash2: () => {
    const React = require('react');
    return React.createElement('div', { 'data-testid': 'trash-icon' });
  },
};

// Mock HTTP helper
export const mockHttp = {
  get: vi.fn(),
  post: vi.fn(),
  put: vi.fn(),
  delete: vi.fn(),
};
