import { useState, useEffect, useRef } from 'react';
import { Routes, Route } from 'react-router';
import { toast } from 'sonner';
import { loadModulesConfig } from './helpers/loadModules';
import { verifyToken } from './helpers/oauth';
import {
  SidebarProvider,
  SidebarInset,
  SidebarTrigger,
} from '@/components/ui/sidebar';
import AppSidebar from './components/AppSidebar';
import Footer from './components/Footer';
import Logo from './components/Logo';
import { DarkModeToggle } from './components/theme/DarkModeToggle';
import { ThemeSelector } from './components/theme/ThemeSelector';
import { ThemeProvider } from './contexts/ThemeContext';
import DashboardModule from './modules/dashboard';
import { Toaster } from '@/components/ui/sonner';
import { Menu, Router, AuthConfig } from '@/types';
import './App.css';

export default function App(): JSX.Element {
  const [auth, setAuth] = useState<AuthConfig>({ active: false });
  const [menus, setMenus] = useState<Menu[]>([]);
  const [routers, setRouters] = useState<Router[]>([]);
  const [setupMode, setSetupMode] = useState<boolean | null>(null);
  const [modulesReady, setModulesReady] = useState<boolean>(false);
  const initialized = useRef<boolean>(false);

  useEffect(() => {
    if (initialized.current) return;
    initialized.current = true;

    verifyToken();
    loadModulesConfig()
      .then((modules) => {
        setAuth(modules.auth);
        setSetupMode(modules.setupMode);

        if (modules.setupMode) {
          setMenus(modules.menus || []);
          setRouters(modules.routers || []);
        } else {
          setMenus([
            ...(DashboardModule.menus || []),
            ...(modules.menus || []),
          ]);
          setRouters([
            ...(DashboardModule.routers || []),
            ...(modules.routers || []),
          ]);
        }
      })
      .catch((error) => {
        console.error('Failed to load modules:', error);
        // ToDo: create new page for error
        toast.error('Failed to load plugins');
        setSetupMode(false);
        setMenus([...(DashboardModule.menus || [])]);
        setRouters([...(DashboardModule.routers || [])]);
      })
      .finally(() => {
        setModulesReady(true);
      });
  }, []);

  if (!modulesReady) {
    return (
      <ThemeProvider>
        <div className="flex h-screen items-center justify-center">
          <span className="text-sm text-muted-foreground">
            Loading DashOps...
          </span>
        </div>
        <Toaster />
      </ThemeProvider>
    );
  }

  if (setupMode) {
    return (
      <ThemeProvider>
        <Routes>
          {routers.map((route) => (
            <Route
              key={route.key}
              path={route.path}
              element={route.element}
            />
          ))}
        </Routes>
        <Toaster />
      </ThemeProvider>
    );
  }

  const sidebarMenus = menus.length > 0 ? menus : (DashboardModule.menus || []);
  const appRouters =
    routers.length > 0 ? routers : (DashboardModule.routers || []);

  return (
    <ThemeProvider>
      <Routes>
        {auth.active && auth.LoginPage && (
          <Route path="/login" element={<auth.LoginPage />} />
        )}
        <Route
          path="*"
          element={
            <SidebarProvider>
              <AppSidebar menus={sidebarMenus} oAuth2={auth.active} />
              <SidebarInset>
                <header className="flex h-16 shrink-0 items-center justify-between gap-2 border-b px-4">
                  <div className="flex items-center gap-2">
                    <SidebarTrigger className="-ml-1" />
                    <Logo />
                  </div>
                  <div className="flex items-center gap-1">
                    <ThemeSelector />
                    <DarkModeToggle />
                  </div>
                </header>
                <main className="flex flex-1 flex-col gap-4 p-4 pt-0 overflow-auto">
                  <Routes>
                    {appRouters.map((route) => (
                      <Route
                        key={route.key}
                        path={route.path}
                        element={route.element}
                      />
                    ))}
                  </Routes>
                </main>
                <footer className="border-t bg-background p-4">
                  <Footer />
                </footer>
              </SidebarInset>
            </SidebarProvider>
          }
        />
      </Routes>
      <Toaster />
    </ThemeProvider>
  );
}
