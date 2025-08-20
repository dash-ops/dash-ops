/**
 * Common types shared across multiple modules
 */

import { ReactElement } from 'react';

// Module configuration types
export interface Menu {
  label: string;
  icon: ReactElement;
  key: string;
  link: string;
}

export interface Router {
  key: string;
  path: string;
  element: ReactElement;
}

export interface Page {
  name: string;
  path: string;
  menu: boolean;
  element: ReactElement;
}

// Module configuration interfaces
export interface ModuleConfig {
  menus?: Menu[];
  routers?: Router[];
}

export interface OAuth2Config {
  active: boolean;
  LoginPage?: () => JSX.Element;
}

export interface OAuth2Module {
  oAuth2: OAuth2Config;
  routers?: Router[];
}

// Generic filter interface
export interface Filter {
  [key: string]: unknown;
}

// Context-based filters
export interface ContextFilter extends Filter {
  context: string;
}

export interface NamespaceFilter extends ContextFilter {
  namespace: string;
}

export interface AccountFilter extends Filter {
  accountKey: string;
}
