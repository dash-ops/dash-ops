import { useState, useEffect, useMemo, useCallback } from 'react';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '../../components/ui/card';
import { Button } from '../../components/ui/button';
import { Input } from '../../components/ui/input';
import { Badge } from '../../components/ui/badge';
import { Avatar, AvatarFallback } from '../../components/ui/avatar';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '../../components/ui/select';
import {
  Plus,
  Search,
  Globe,
  Server,
  Shield,
  Clock,
  MoreVertical,
  Edit3,
  AlertTriangle,
  Loader2,
  RefreshCw,
} from 'lucide-react';
import { cn } from '../../lib/utils';
import type { Service, ServiceHealth, ServiceStats } from './types';
import { getServices, getServiceHealthBatch } from './serviceCatalogResource';
import { ServiceFormModal } from './ServiceFormModal';

export function ServicesCatalogPage() {
  const [services, setServices] = useState<Service[]>([]);
  const [serviceHealths, setServiceHealths] = useState<
    Record<string, ServiceHealth>
  >({});
  const [loading, setLoading] = useState(true);
  const [healthLoading, setHealthLoading] = useState(false);
  const [error, setError] = useState<string>('');

  // Filters
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedTier, setSelectedTier] = useState<string>('');
  const [teamFilter, setTeamFilter] = useState<string>('my-team');
  const [statusFilter, setStatusFilter] = useState<string>('');
  const [sortBy, setSortBy] = useState<'name' | 'tier' | 'team' | 'updated_at'>(
    'name'
  );

  // Modal state
  const [createModalOpen, setCreateModalOpen] = useState(false);
  const [editingService, setEditingService] = useState<Service | undefined>(
    undefined
  );

  // Pagination
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 12;

  // Mock user context (in real app, this would come from auth)
  const currentUser = {
    teams: ['auth-squad'], // This would come from GitHub teams via OAuth2
    username: 'current-user',
  };

  const loadServices = useCallback(async () => {
    try {
      setLoading(true);
      setError('');
      const result = await getServices();
      setServices(result.services || []);
    } catch (err) {
      console.error('Failed to load services:', err);
      setError('Failed to load services. Please try again.');
    } finally {
      setLoading(false);
    }
  }, []);

  const loadHealthData = useCallback(async (serviceNames: string[]) => {
    try {
      setHealthLoading(true);
      const healthData = await getServiceHealthBatch(serviceNames);
      setServiceHealths(healthData);
    } catch (err) {
      console.error('Failed to load health data:', err);
      // Don't show error for health data - it's optional
    } finally {
      setHealthLoading(false);
    }
  }, []);

  // Load services on component mount
  useEffect(() => {
    loadServices();
  }, [loadServices]);

  // Load health data when services change
  useEffect(() => {
    if (services.length > 0) {
      const serviceNames = services.map((s) => s.metadata.name);
      loadHealthData(serviceNames);
    }
  }, [services, loadHealthData]);

  // Filter and sort services
  const filteredServices = useMemo(() => {
    let filtered = [...services];

    // Team filter
    if (teamFilter === 'my-team') {
      filtered = filtered.filter((service) =>
        currentUser.teams.includes(service.spec.team.github_team)
      );
    }

    // Search filter
    if (searchTerm) {
      const search = searchTerm.toLowerCase();
      filtered = filtered.filter(
        (service) =>
          service.metadata.name.toLowerCase().includes(search) ||
          service.spec.description.toLowerCase().includes(search) ||
          service.spec.team.github_team.toLowerCase().includes(search) ||
          service.spec.technology?.language?.toLowerCase().includes(search)
      );
    }

    // Tier filter
    if (selectedTier) {
      filtered = filtered.filter(
        (service) => service.metadata.tier === selectedTier
      );
    }

    // Status filter (based on health data)
    if (statusFilter) {
      filtered = filtered.filter((service) => {
        const health = serviceHealths[service.metadata.name];
        return health?.overall_status === statusFilter;
      });
    }

    // Sort services
    filtered.sort((a, b) => {
      // Prioritize services user can edit
      const aCanEdit = currentUser.teams.includes(a.spec.team.github_team);
      const bCanEdit = currentUser.teams.includes(b.spec.team.github_team);

      if (aCanEdit && !bCanEdit) return -1;
      if (!aCanEdit && bCanEdit) return 1;

      // Then sort by selected criteria
      switch (sortBy) {
        case 'name':
          return a.metadata.name.localeCompare(b.metadata.name);
        case 'tier':
          return a.metadata.tier.localeCompare(b.metadata.tier);
        case 'team':
          return a.spec.team.github_team.localeCompare(b.spec.team.github_team);
        case 'updated_at':
          return (
            new Date(b.metadata.updated_at || 0).getTime() -
            new Date(a.metadata.updated_at || 0).getTime()
          );
        default:
          return 0;
      }
    });

    return filtered;
  }, [
    services,
    searchTerm,
    selectedTier,
    teamFilter,
    statusFilter,
    sortBy,
    serviceHealths,
    currentUser.teams,
  ]);

  // Pagination
  const totalPages = Math.ceil(filteredServices.length / itemsPerPage);
  const paginatedServices = filteredServices.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );

  // Calculate statistics
  const stats = useMemo((): ServiceStats => {
    const myTeamServices = services.filter((s) =>
      currentUser.teams.includes(s.spec.team.github_team)
    );
    const tier1Services = filteredServices.filter(
      (s) => s.metadata.tier === 'TIER-1'
    );
    const tier2Services = filteredServices.filter(
      (s) => s.metadata.tier === 'TIER-2'
    );
    const tier3Services = filteredServices.filter(
      (s) => s.metadata.tier === 'TIER-3'
    );
    const criticalServices = filteredServices.filter((s) => {
      const health = serviceHealths[s.metadata.name];
      return health?.overall_status === 'critical';
    });

    return {
      total: services.length,
      myTeam: myTeamServices.length,
      tier1: tier1Services.length,
      tier2: tier2Services.length,
      tier3: tier3Services.length,
      critical: criticalServices.length,
      editable: myTeamServices.length,
    };
  }, [services, filteredServices, serviceHealths, currentUser.teams]);

  // Helper functions
  const getStatusColor = (status?: string) => {
    switch (status) {
      case 'healthy':
        return 'bg-green-500';
      case 'degraded':
        return 'bg-yellow-500';
      case 'critical':
        return 'bg-red-500';
      case 'down':
        return 'bg-red-600';
      default:
        return 'bg-gray-400';
    }
  };

  const getTierColor = (tier: string) => {
    switch (tier) {
      case 'TIER-1':
        return 'bg-red-100 text-red-800 border-red-200';
      case 'TIER-2':
        return 'bg-yellow-100 text-yellow-800 border-yellow-200';
      case 'TIER-3':
        return 'bg-green-100 text-green-800 border-green-200';
      default:
        return 'bg-gray-100 text-gray-800 border-gray-200';
    }
  };

  const formatLastUpdate = (dateString?: string) => {
    if (!dateString) return 'Unknown';
    const date = new Date(dateString);
    const now = new Date();
    const diffInHours = Math.floor(
      (now.getTime() - date.getTime()) / (1000 * 60 * 60)
    );

    if (diffInHours < 24) {
      return `${diffInHours}h ago`;
    } else {
      return `${Math.floor(diffInHours / 24)}d ago`;
    }
  };

  const getUniqueTeams = () => {
    const teams = services.map((s) => s.spec.team.github_team);
    return [...new Set(teams)].sort();
  };

  if (loading) {
    return (
      <div className="flex-1 flex items-center justify-center">
        <div className="text-center">
          <Loader2 className="h-8 w-8 animate-spin mx-auto mb-4" />
          <p>Loading services...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex-1 flex items-center justify-center">
        <div className="text-center">
          <AlertTriangle className="h-8 w-8 mx-auto mb-4 text-red-500" />
          <p className="text-red-600 mb-4">{error}</p>
          <Button onClick={loadServices}>Try Again</Button>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 flex flex-col overflow-hidden">
      {/* Header */}
      <div className="flex-none p-6 border-b bg-background">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h1 className="text-2xl font-bold">Services Catalog</h1>
            <p className="text-muted-foreground">
              Manage {stats.total.toLocaleString()} services •{' '}
              {filteredServices.length} shown • {stats.editable} editable
            </p>
          </div>

          <div className="flex gap-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => {
                const serviceNames = services.map((s) => s.metadata.name);
                loadHealthData(serviceNames);
              }}
              disabled={healthLoading}
            >
              <RefreshCw
                className={cn('h-4 w-4 mr-2', healthLoading && 'animate-spin')}
              />
              Refresh
            </Button>
            <Button
              className="gap-2"
              onClick={() => {
                setEditingService(undefined);
                setCreateModalOpen(true);
              }}
            >
              <Plus className="h-4 w-4" />
              Add Service
            </Button>
          </div>
        </div>

        {/* Filters */}
        <div className="grid grid-cols-1 md:grid-cols-6 gap-4 mt-6">
          <div className="relative md:col-span-2">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
            <Input
              placeholder="Search services, teams, technologies..."
              value={searchTerm}
              onChange={(e) => {
                setSearchTerm(e.target.value);
                setCurrentPage(1);
              }}
              className="pl-10"
            />
          </div>

          <Select
            value={teamFilter}
            onValueChange={(value) => {
              setTeamFilter(value);
              setCurrentPage(1);
            }}
          >
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="my-team">My Team ({stats.myTeam})</SelectItem>
              <SelectItem value="all">All Teams ({stats.total})</SelectItem>
              {getUniqueTeams().map((team) => (
                <SelectItem key={team} value={team}>
                  {team}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Select
            value={selectedTier || 'all'}
            onValueChange={(value) => {
              setSelectedTier(value === 'all' ? '' : value);
              setCurrentPage(1);
            }}
          >
            <SelectTrigger>
              <SelectValue placeholder="All Tiers" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Tiers</SelectItem>
              <SelectItem value="TIER-1">TIER-1 ({stats.tier1})</SelectItem>
              <SelectItem value="TIER-2">TIER-2 ({stats.tier2})</SelectItem>
              <SelectItem value="TIER-3">TIER-3 ({stats.tier3})</SelectItem>
            </SelectContent>
          </Select>

          <Select
            value={statusFilter || 'all'}
            onValueChange={(value) => {
              setStatusFilter(value === 'all' ? '' : value);
              setCurrentPage(1);
            }}
          >
            <SelectTrigger>
              <SelectValue placeholder="All Status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Status</SelectItem>
              <SelectItem value="healthy">Healthy</SelectItem>
              <SelectItem value="degraded">Degraded</SelectItem>
              <SelectItem value="critical">Critical</SelectItem>
              <SelectItem value="down">Down</SelectItem>
            </SelectContent>
          </Select>

          <Select
            value={sortBy}
            onValueChange={(value: any) => setSortBy(value)}
          >
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="name">Sort by Name</SelectItem>
              <SelectItem value="tier">Sort by Tier</SelectItem>
              <SelectItem value="team">Sort by Team</SelectItem>
              <SelectItem value="updated_at">Sort by Updated</SelectItem>
            </SelectContent>
          </Select>
        </div>

        {/* Critical Services Alert */}
        {stats.critical > 0 && (
          <div className="mt-4 p-3 bg-red-50 border border-red-200 rounded-lg">
            <div className="flex items-center gap-2">
              <AlertTriangle className="h-4 w-4 text-red-600" />
              <span className="text-sm font-medium text-red-800">
                {stats.critical} services need immediate attention
              </span>
              <Button
                variant="outline"
                size="sm"
                onClick={() => {
                  setStatusFilter('critical');
                  setCurrentPage(1);
                }}
                className="ml-auto"
              >
                View Critical
              </Button>
            </div>
          </div>
        )}
      </div>

      {/* Services Grid */}
      <div className="flex-1 overflow-auto p-6">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {paginatedServices.map((service) => {
            const health = serviceHealths[service.metadata.name];
            const canEdit = currentUser.teams.includes(
              service.spec.team.github_team
            );
            const environments = service.spec.kubernetes?.environments || [];

            return (
              <Card
                key={service.metadata.name}
                className="cursor-pointer hover:shadow-md transition-shadow"
              >
                <CardHeader className="pb-3">
                  <div className="flex items-start justify-between">
                    <div className="flex items-center gap-2">
                      <div
                        className={cn(
                          'w-3 h-3 rounded-full',
                          getStatusColor(health?.overall_status)
                        )}
                      />
                      <Badge
                        variant="outline"
                        className={getTierColor(service.metadata.tier)}
                      >
                        {service.metadata.tier}
                      </Badge>
                    </div>
                    <div className="flex items-center gap-1">
                      {canEdit ? (
                        <Button
                          variant="secondary"
                          size="sm"
                          className="text-xs gap-1 h-6 px-2"
                          onClick={() => {
                            setEditingService(service);
                            setCreateModalOpen(true);
                          }}
                        >
                          <Edit3 className="h-3 w-3" />
                        </Button>
                      ) : (
                        <Badge variant="outline" className="text-xs gap-1">
                          <Shield className="h-3 w-3" />
                        </Badge>
                      )}
                      <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
                        <MoreVertical className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>

                  <CardTitle className="text-base leading-tight">
                    {service.metadata.name}
                  </CardTitle>
                  <CardDescription className="line-clamp-2 text-xs">
                    {service.spec.description}
                  </CardDescription>
                </CardHeader>

                <CardContent className="space-y-3">
                  {/* Team */}
                  <div className="flex items-center gap-2">
                    <Avatar className="h-6 w-6">
                      <AvatarFallback className="text-xs">
                        {service.spec.team.github_team
                          .split('-')
                          .map((word) => word[0]?.toUpperCase())
                          .join('')
                          .slice(0, 2)}
                      </AvatarFallback>
                    </Avatar>
                    <div className="min-w-0 flex-1">
                      <p className="text-xs font-medium truncate">
                        {service.spec.team.github_team}
                      </p>
                      {service.spec.technology?.language && (
                        <p className="text-xs text-muted-foreground truncate">
                          {service.spec.technology.language}
                        </p>
                      )}
                    </div>
                  </div>

                  {/* Environment Info */}
                  <div className="flex items-center justify-between text-xs">
                    <div className="flex items-center gap-1">
                      <Globe className="h-3 w-3 text-muted-foreground" />
                      <span>
                        {environments.length} env
                        {environments.length !== 1 ? 's' : ''}
                      </span>
                    </div>
                    <div className="flex items-center gap-1">
                      <Server className="h-3 w-3 text-muted-foreground" />
                      <span>{service.spec.business?.impact || 'medium'}</span>
                    </div>
                  </div>

                  {/* Technology & Features */}
                  <div className="flex flex-wrap gap-1">
                    {service.spec.technology?.framework && (
                      <Badge
                        variant="secondary"
                        className="text-xs px-1.5 py-0.5"
                      >
                        {service.spec.technology.framework}
                      </Badge>
                    )}
                    {service.spec.kubernetes && (
                      <Badge
                        variant="secondary"
                        className="text-xs px-1.5 py-0.5"
                      >
                        k8s
                      </Badge>
                    )}
                    {service.spec.observability?.metrics && (
                      <Badge
                        variant="secondary"
                        className="text-xs px-1.5 py-0.5"
                      >
                        metrics
                      </Badge>
                    )}
                  </div>

                  {/* Last Update */}
                  <div className="flex items-center gap-1 text-xs text-muted-foreground">
                    <Clock className="h-3 w-3" />
                    <span>
                      Updated {formatLastUpdate(service.metadata.updated_at)}
                    </span>
                  </div>
                </CardContent>
              </Card>
            );
          })}
        </div>

        {/* Empty State */}
        {paginatedServices.length === 0 && (
          <div className="text-center py-12">
            <Server className="h-12 w-12 mx-auto mb-4 text-muted-foreground opacity-50" />
            <p className="text-muted-foreground mb-4">
              No services found with the current filters.
            </p>
            <div className="flex justify-center gap-2">
              <Button
                variant="outline"
                onClick={() => {
                  setSearchTerm('');
                  setSelectedTier('');
                  setStatusFilter('');
                  setTeamFilter('all');
                  setCurrentPage(1);
                }}
              >
                Clear Filters
              </Button>
              <Button onClick={() => setCreateModalOpen(true)}>
                Create New Service
              </Button>
            </div>
          </div>
        )}

        {/* Simple Pagination */}
        {totalPages > 1 && (
          <div className="flex justify-center items-center gap-2 mt-8">
            <Button
              variant="outline"
              size="sm"
              disabled={currentPage === 1}
              onClick={() => setCurrentPage(currentPage - 1)}
            >
              Previous
            </Button>
            <span className="text-sm text-muted-foreground">
              Page {currentPage} of {totalPages}
            </span>
            <Button
              variant="outline"
              size="sm"
              disabled={currentPage === totalPages}
              onClick={() => setCurrentPage(currentPage + 1)}
            >
              Next
            </Button>
          </div>
        )}
      </div>

      {/* Create/Edit Service Modal */}
      <ServiceFormModal
        open={createModalOpen}
        onOpenChange={(open) => {
          setCreateModalOpen(open);
          if (!open) {
            setEditingService(undefined);
          }
        }}
        onServiceCreated={() => {
          loadServices();
          setCreateModalOpen(false);
          setEditingService(undefined);
        }}
        editingService={editingService}
      />
    </div>
  );
}
