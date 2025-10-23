import { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { 
  Bot, 
  Download, 
  Calendar, 
  RefreshCw, 
  Play, 
  Pause, 
  FileText, 
  GitBranch 
} from 'lucide-react';
import LogsPage from './logs/LogsPage';
import TracesPage from './traces/TracesPage';

export default function ObservabilityPage(): JSX.Element {
  const [selectedService, setSelectedService] = useState<string>('all');
  const [timeRange, setTimeRange] = useState<string>('1h');
  const [isAutoRefresh, setIsAutoRefresh] = useState(true);
  const [activeTab, setActiveTab] = useState<string>('logs');

  return (
    <div className="flex-1 flex flex-col overflow-hidden">
      {/* Global Header & Filters */}
      <div className="flex-none p-6 border-b bg-background">
        <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4 mb-6">
          <div>
            <h1 className="text-2xl font-bold">Observability & Troubleshooting</h1>
            <p className="text-muted-foreground">
              Comprehensive monitoring with logs, metrics, and distributed tracing
            </p>
          </div>
          
          <div className="flex items-center gap-2">
            {/* <Button variant="outline" size="sm" className="gap-2">
              <Bot className="h-4 w-4" />
              AI Analysis
            </Button>
            <Button variant="outline" size="sm" className="gap-2">
              <Download className="h-4 w-4" />
              Export
            </Button> */}
          </div>
        </div>

        {/* Global Filters */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Select value={selectedService} onValueChange={setSelectedService}>
            <SelectTrigger>
              <SelectValue placeholder="Select Service" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Services</SelectItem>
              {/* TODO: Add dynamic services from API */}
            </SelectContent>
          </Select>

          <Select value={timeRange} onValueChange={setTimeRange}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="5m">Last 5 minutes</SelectItem>
              <SelectItem value="15m">Last 15 minutes</SelectItem>
              <SelectItem value="1h">Last hour</SelectItem>
              <SelectItem value="6h">Last 6 hours</SelectItem>
              <SelectItem value="24h">Last 24 hours</SelectItem>
              <SelectItem value="7d">Last 7 days</SelectItem>
            </SelectContent>
          </Select>

          <div className="flex items-center gap-2">
            <Button
              variant={isAutoRefresh ? "default" : "outline"}
              size="sm"
              onClick={() => setIsAutoRefresh(!isAutoRefresh)}
              className="gap-2"
            >
              {isAutoRefresh ? <Pause className="h-4 w-4" /> : <Play className="h-4 w-4" />}
              {isAutoRefresh ? 'Live' : 'Paused'}
            </Button>
            <Button variant="outline" size="sm" className="gap-2">
              <RefreshCw className="h-4 w-4" />
            </Button>
          </div>

          <div className="flex items-center gap-2">
            <Button variant="outline" size="sm" className="gap-2">
              <Calendar className="h-4 w-4" />
              Custom Range
            </Button>
          </div>
        </div>
      </div>

      {/* Tabs Content */}
      <div className="flex-1 overflow-hidden">
        <Tabs value={activeTab} onValueChange={setActiveTab} className="h-full">
          <div className="flex-none px-6 pt-4">
            <TabsList className="grid w-full max-w-md grid-cols-2">
              <TabsTrigger value="logs" className="gap-2">
                <FileText className="h-4 w-4" />
                Logs
              </TabsTrigger>
              <TabsTrigger value="traces" className="gap-2">
                <GitBranch className="h-4 w-4" />
                Traces
              </TabsTrigger>
            </TabsList>
          </div>

          <TabsContent value="logs" className="h-full mt-0">
            <LogsPage />
          </TabsContent>

          <TabsContent value="traces" className="h-full mt-0">
            <TracesPage />
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}
