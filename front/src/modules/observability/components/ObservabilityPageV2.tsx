import { useState, useEffect } from 'react';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Badge } from '@/components/ui/badge';
import { 
  LayoutDashboard,
  Bell,
  Database,
  ChevronUp
} from 'lucide-react';
import DashboardsTab from './DashboardsTab';
import AlertsTab from './AlertsTab';
import ExplorerDrawer from './ExplorerDrawer';

export default function ObservabilityPageV2(): JSX.Element {
  const [activeTab, setActiveTab] = useState<string>('dashboards');
  const [isExplorerOpen, setIsExplorerOpen] = useState(false);

  // Keyboard shortcut for opening Explorer (Ctrl/Cmd + Shift + O)
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.shiftKey && e.key === 'o') {
        e.preventDefault();
        setIsExplorerOpen(true);
      }
      // Escape to close
      if (e.key === 'Escape' && isExplorerOpen) {
        setIsExplorerOpen(false);
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isExplorerOpen]);

  return (
    <div className="flex-1 flex flex-col overflow-hidden relative">
      {/* Global Header & Filters */}
      <div className="flex-none p-6 border-b bg-background">
        <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4 mb-6">
          <div>
            <h1 className="text-2xl font-bold">Observability Platform</h1>
            <p className="text-muted-foreground">
              360° view of your services with AI-powered insights
            </p>
          </div>
        </div>
      </div>

      {/* Main Content with Tabs */}
      <div className="flex-1 overflow-hidden">
        <Tabs value={activeTab} onValueChange={setActiveTab} className="h-full flex flex-col">
          <div className="flex-none px-6 pt-4">
            <TabsList className="grid w-full max-w-md grid-cols-2">
              <TabsTrigger value="dashboards" className="gap-2">
                <LayoutDashboard className="h-4 w-4" />
                Dashboards
              </TabsTrigger>
              <TabsTrigger value="alerts" className="gap-2">
                <Bell className="h-4 w-4" />
                Alerts
              </TabsTrigger>
            </TabsList>
          </div>

          <TabsContent value="dashboards" className="flex-1 overflow-auto m-0">
            <DashboardsTab onOpenExplorer={() => setIsExplorerOpen(true)} />
          </TabsContent>

          <TabsContent value="alerts" className="flex-1 overflow-auto m-0">
            <AlertsTab />
          </TabsContent>
        </Tabs>
      </div>

      {/* Explorer Footer Bar - Always visible when drawer is closed */}
      {!isExplorerOpen && (
        <div 
          className="fixed bottom-0 left-0 right-0 z-40 bg-card border-t border-border shadow-2xl cursor-pointer hover:bg-accent transition-all duration-200"
          onClick={() => setIsExplorerOpen(true)}
        >
          <div className="container mx-auto px-6 py-3">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-4">
                <div className="flex items-center gap-3">
                  <div className="w-8 h-8 rounded-lg bg-primary flex items-center justify-center">
                    <Database className="h-4 w-4 text-primary-foreground" />
                  </div>
                  <div>
                    <p className="text-sm font-medium">Query your data</p>
                    <p className="text-xs text-muted-foreground">Click or press Ctrl+Shift+O to open Explorer</p>
                  </div>
                </div>
                
                <div className="hidden md:flex items-center gap-2">
                  <Badge variant="outline" className="bg-muted/50">
                    <kbd className="text-xs">Ctrl</kbd>
                  </Badge>
                  <span className="text-muted-foreground">+</span>
                  <Badge variant="outline" className="bg-muted/50">
                    <kbd className="text-xs">Shift</kbd>
                  </Badge>
                  <span className="text-muted-foreground">+</span>
                  <Badge variant="outline" className="bg-muted/50">
                    <kbd className="text-xs">O</kbd>
                  </Badge>
                </div>
              </div>

              <div className="flex items-center gap-3">
                <Badge variant="secondary">
                  Explorer
                </Badge>
                <ChevronUp className="h-5 w-5 text-muted-foreground animate-bounce" />
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Explorer Drawer */}
      <ExplorerDrawer 
        isOpen={isExplorerOpen} 
        onClose={() => setIsExplorerOpen(false)} 
      />
    </div>
  );
}

