import { useState } from 'react';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
  FileText, 
  GitBranch 
} from 'lucide-react';
import TimeRangePicker, { type TimeRange } from '@/components/TimeRangePicker';
import LogsPage from './logs/LogsPage';
import TracesPage from './traces/TracesPage';

export default function ObservabilityPage(): JSX.Element {
  const [timeRange, setTimeRange] = useState<TimeRange>({ 
    value: '1h', 
    label: 'Last hour',
    from: new Date(Date.now() - 60 * 60 * 1000),
    to: new Date(),
  });
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
        <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-end">
          <TimeRangePicker value={timeRange} onChange={setTimeRange} />
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
