import { useState } from 'react';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { RotateCcw, Scale } from 'lucide-react';
import { KubernetesTypes } from '@/types';

interface ModernDeploymentActionsProps {
  context: string;
  deployment: KubernetesTypes.Deployment;
  onRestart: () => Promise<void>;
  onScale: (replicas: number) => Promise<void>;
}

export default function ModernDeploymentActions({
  deployment,
  onRestart,
  onScale,
}: ModernDeploymentActionsProps): JSX.Element {
  const [scaleDialogOpen, setScaleDialogOpen] = useState(false);
  const [newReplicas, setNewReplicas] = useState(deployment.replicas.desired);
  const [loading, setLoading] = useState(false);

  const handleRestart = async () => {
    setLoading(true);
    try {
      await onRestart();
    } finally {
      setLoading(false);
    }
  };

  const handleScale = async () => {
    if (newReplicas < 0 || newReplicas > 50) {
      return;
    }

    setLoading(true);
    try {
      await onScale(newReplicas);
      setScaleDialogOpen(false);
    } finally {
      setLoading(false);
    }
  };

  return (
    <TooltipProvider>
      <div className="flex gap-2">
        {/* Restart Action with AlertDialog */}
        <AlertDialog>
          <AlertDialogTrigger asChild>
            <div>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="outline"
                    size="sm"
                    disabled={loading}
                    className="h-8 w-8 p-0"
                  >
                    <RotateCcw className="h-3 w-3" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>
                  <p>Restart deployment</p>
                </TooltipContent>
              </Tooltip>
            </div>
          </AlertDialogTrigger>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Restart Deployment</AlertDialogTitle>
              <AlertDialogDescription>
                Are you sure you want to restart the deployment{' '}
                <span className="font-semibold">{deployment.name}</span> in
                namespace{' '}
                <span className="font-semibold">{deployment.namespace}</span>?
                <br />
                <br />
                This will trigger a rolling restart of all pods in this
                deployment.
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel disabled={loading}>Cancel</AlertDialogCancel>
              <AlertDialogAction
                onClick={handleRestart}
                disabled={loading}
                className="bg-orange-600 hover:bg-orange-700"
              >
                {loading ? 'Restarting...' : 'Restart'}
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>

        {/* Scale Action with Dialog */}
        <Dialog open={scaleDialogOpen} onOpenChange={setScaleDialogOpen}>
          <DialogTrigger asChild>
            <div>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="outline"
                    size="sm"
                    className="h-8 w-8 p-0"
                    disabled={loading}
                  >
                    <Scale className="h-3 w-3" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>
                  <p>Scale deployment</p>
                </TooltipContent>
              </Tooltip>
            </div>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[425px]">
            <DialogHeader>
              <DialogTitle>Scale Deployment</DialogTitle>
              <DialogDescription>
                Adjust the number of replicas for{' '}
                <span className="font-semibold">{deployment.name}</span> in
                namespace{' '}
                <span className="font-semibold">{deployment.namespace}</span>.
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
              <div className="space-y-2">
                <Label htmlFor="replicas">Number of Replicas</Label>
                <Input
                  id="replicas"
                  type="number"
                  min="0"
                  max="50"
                  value={newReplicas}
                  onChange={(e) => setNewReplicas(Number(e.target.value))}
                  placeholder="Enter number of replicas (0-50)"
                />
                <div className="flex justify-between text-xs text-muted-foreground">
                  <span>
                    Current: {deployment.replicas.ready}/
                    {deployment.replicas.desired}
                  </span>
                  <span>Min: 0 | Max: 50</span>
                </div>
              </div>
            </div>
            <DialogFooter className="gap-2">
              <Button
                variant="outline"
                onClick={() => setScaleDialogOpen(false)}
                disabled={loading}
              >
                Cancel
              </Button>
              <Button
                onClick={handleScale}
                disabled={loading || newReplicas < 0 || newReplicas > 50}
                className="bg-blue-600 hover:bg-blue-700"
              >
                {loading ? 'Scaling...' : `Scale to ${newReplicas}`}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </div>
    </TooltipProvider>
  );
}
