import { Badge } from '@/components/ui/badge';

interface PodQoSBadgeProps {
  qosClass: string;
}

export default function PodQoSBadge({
  qosClass,
}: PodQoSBadgeProps): JSX.Element {
  const getQoSVariant = (qos: string) => {
    switch (qos.toLowerCase()) {
      case 'guaranteed':
        return 'default' as const;
      case 'burstable':
        return 'secondary' as const;
      case 'besteffort':
        return 'outline' as const;
      default:
        return 'outline' as const;
    }
  };

  return (
    <Badge variant={getQoSVariant(qosClass)} className="text-xs">
      {qosClass}
    </Badge>
  );
}
