import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { LayoutDashboard, Search, Plus } from 'lucide-react';

interface DashboardsTabProps {
  onOpenExplorer: () => void;
}

export default function DashboardsTab({ onOpenExplorer }: DashboardsTabProps): JSX.Element {
  return (
    <div className="flex-1 flex flex-col p-6 gap-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-semibold flex items-center gap-2">
            <LayoutDashboard className="h-5 w-5" />
            Dashboards
          </h2>
          <p className="text-muted-foreground text-sm">
            Manage and view your custom observability dashboards
          </p>
        </div>
        <div className="flex gap-2">
          <Button onClick={onOpenExplorer} variant="outline" className="gap-2">
            <Search className="h-4 w-4" />
            Open Explorer
          </Button>
          <Button className="gap-2">
            <Plus className="h-4 w-4" />
            New Dashboard
          </Button>
        </div>
      </div>

      {/* Empty State */}
      <Card className="border-dashed">
        <CardContent className="flex flex-col items-center justify-center py-16 text-center">
          <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center mb-4">
            <LayoutDashboard className="h-8 w-8 text-primary" />
          </div>
          <h3 className="text-lg font-medium mb-2">No Dashboards Yet</h3>
          <p className="text-muted-foreground mb-6 max-w-md">
            Create your first dashboard or start exploring your data to build custom visualizations
          </p>
          <div className="flex gap-2">
            <Button onClick={onOpenExplorer} variant="outline" className="gap-2">
              <Search className="h-4 w-4" />
              Open Explorer
            </Button>
            <Button className="gap-2">
              <Plus className="h-4 w-4" />
              Create Dashboard
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
