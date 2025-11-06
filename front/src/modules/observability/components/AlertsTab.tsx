import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Bell, Plus } from 'lucide-react';

export default function AlertsTab(): JSX.Element {
  return (
    <div className="flex-1 flex flex-col p-6 gap-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-semibold flex items-center gap-2">
            <Bell className="h-5 w-5" />
            Alerts
          </h2>
          <p className="text-muted-foreground text-sm">
            Configure and manage alerts for your observability data
          </p>
        </div>
        <Button className="gap-2">
          <Plus className="h-4 w-4" />
          Create Alert
        </Button>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Total Alerts</CardDescription>
            <CardTitle>0</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Active</CardDescription>
            <CardTitle className="text-green-600">0</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Triggered</CardDescription>
            <CardTitle className="text-red-600">0</CardTitle>
          </CardHeader>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Paused</CardDescription>
            <CardTitle className="text-gray-600">0</CardTitle>
          </CardHeader>
        </Card>
      </div>

      {/* Empty State */}
      <Card className="border-dashed">
        <CardContent className="flex flex-col items-center justify-center py-16 text-center">
          <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center mb-4">
            <Bell className="h-8 w-8 text-primary" />
          </div>
          <h3 className="text-lg font-medium mb-2">No Alerts Configured</h3>
          <p className="text-muted-foreground mb-6 max-w-md">
            Create your first alert to get notified when specific conditions are met in your observability data
          </p>
          <Button className="gap-2">
            <Plus className="h-4 w-4" />
            Create Alert
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
