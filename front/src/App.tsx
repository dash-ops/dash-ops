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
import { Menu, Router, OAuth2Config } from '@/types';
import './App.css';

export default function App(): JSX.Element {
  const [oAuth2, setOAuth2] = useState<OAuth2Config>({ active: false });
  const [menus, setMenus] = useState<Menu[]>([
    ...(DashboardModule.menus || []),
  ]);
  const [routers, setRouters] = useState<Router[]>([
    ...(DashboardModule.routers || []),
  ]);
  const initialized = useRef<boolean>(false);

  useEffect(() => {
    if (initialized.current) return;
    initialized.current = true;

    verifyToken();
    loadModulesConfig()
      .then((modules) => {
        setOAuth2(modules.oAuth2);
        setMenus([...(DashboardModule.menus || []), ...modules.menus]);
        setRouters([...(DashboardModule.routers || []), ...modules.routers]);
      })
      .catch((error) => {
        console.error('Failed to load modules:', error);
        toast.error('Failed to load plugins');
      });
  }, []);

  return (
    <ThemeProvider>
      <Routes>
        {oAuth2.active && oAuth2.LoginPage && (
          <Route path="/login" element={<oAuth2.LoginPage />} />
        )}
        <Route
          path="*"
          element={
            <SidebarProvider>
              <AppSidebar menus={menus} oAuth2={oAuth2.active} />
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
                    {routers.map((route) => (
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
