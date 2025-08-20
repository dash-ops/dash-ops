import { useState, useEffect, useCallback, useRef } from 'react';
import { useLocation, useNavigate, Link } from 'react-router';
import { toast } from 'sonner';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/ui/sidebar';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
  Building2,
  ChevronDown,
  LogOut,
  User,
  Settings,
  Users,
} from 'lucide-react';
import {
  getUserData,
  getUserPermissions,
} from '../modules/oauth2/userResource';
import { cleanToken } from '../helpers/oauth';
import { OAuth2Types, Menu } from '@/types';

interface AppSidebarProps {
  menus?: Menu[];
  oAuth2: boolean;
}

function AppSidebar({ menus = [], oAuth2 }: AppSidebarProps): JSX.Element {
  const location = useLocation();
  const navigate = useNavigate();
  const [user, setUser] = useState<OAuth2Types.UserData | null>(null);
  const [permissions, setPermissions] =
    useState<OAuth2Types.UserPermission | null>(null);
  const userDataFetched = useRef<boolean>(false);

  const logout = useCallback(() => {
    setUser(null);
    cleanToken();
    userDataFetched.current = false;
    navigate('/login');
  }, [navigate]);

  useEffect(() => {
    if (!oAuth2 || userDataFetched.current) {
      return;
    }

    async function fetchData() {
      try {
        const [userResult, permissionsResult] = await Promise.all([
          getUserData(),
          getUserPermissions(),
        ]);
        setUser(userResult.data);
        setPermissions(permissionsResult.data);
        userDataFetched.current = true;
      } catch {
        toast.error('Failed to fetch user data');
      }
    }

    fetchData();
  }, [oAuth2]);

  return (
    <Sidebar collapsible="icon">
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <button className="flex w-full items-center gap-2 rounded-md p-2 text-left text-sm outline-hidden ring-sidebar-ring transition-[width,height,padding] hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground h-12 group-data-[collapsible=icon]:size-8 group-data-[collapsible=icon]:p-0 group-data-[collapsible=icon]:justify-center">
                  <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                    <Building2 className="size-4" />
                  </div>
                  <div className="grid flex-1 text-left text-sm leading-tight group-data-[collapsible=icon]:hidden">
                    <span className="truncate font-semibold">
                      {permissions?.organization || 'DashOPS'}
                    </span>
                  </div>
                  <ChevronDown className="ml-auto size-4 group-data-[collapsible=icon]:hidden" />
                </button>
              </DropdownMenuTrigger>
              <DropdownMenuContent className="w-[240px]" align="start">
                <DropdownMenuItem className="flex-col items-start">
                  <div className="text-sm font-medium">
                    {permissions?.organization || 'DashOPS'}
                  </div>
                </DropdownMenuItem>
                <DropdownMenuSeparator />

                <div className="px-2 py-1">
                  <div className="flex items-center gap-2 text-sm font-medium mb-2">
                    <Users className="h-4 w-4" />
                    Teams
                  </div>
                  {permissions?.teams && permissions.teams.length > 0 ? (
                    permissions.teams.map((team, index) => (
                      <DropdownMenuItem
                        key={team.slug || index}
                        className="pl-6"
                      >
                        <div className="flex items-center justify-between w-full">
                          <span>{team.name || team.slug}</span>
                          <Badge variant="secondary" className="text-xs ml-2">
                            ⌘{index + 1}
                          </Badge>
                        </div>
                      </DropdownMenuItem>
                    ))
                  ) : (
                    <div className="pl-6 py-2 text-xs text-muted-foreground">
                      No teams found
                    </div>
                  )}
                </div>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Platform</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {menus.map((menu) => (
                <SidebarMenuItem key={menu.key}>
                  <SidebarMenuButton
                    asChild
                    isActive={location.pathname === menu.link}
                    tooltip={menu.label}
                  >
                    <Link to={menu.link}>
                      {menu.icon}
                      <span>{menu.label}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter>
        {user && (
          <SidebarMenu>
            <SidebarMenuItem>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <button className="flex w-full items-center gap-2 rounded-md p-2 text-left text-sm outline-hidden ring-sidebar-ring transition-[width,height,padding] hover:bg-sidebar-accent hover:text-sidebar-accent-foreground focus-visible:ring-2 active:bg-sidebar-accent active:text-sidebar-accent-foreground disabled:pointer-events-none disabled:opacity-50 aria-disabled:pointer-events-none aria-disabled:opacity-50 data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground h-12 group-data-[collapsible=icon]:size-8 group-data-[collapsible=icon]:p-0 group-data-[collapsible=icon]:justify-center">
                    <Avatar className="h-8 w-8 rounded-lg">
                      <AvatarImage
                        src={user.avatar_url}
                        alt={user.name || user.login}
                      />
                      <AvatarFallback className="rounded-lg">
                        {(user.name || user.login)?.charAt(0)}
                      </AvatarFallback>
                    </Avatar>
                    <div className="grid flex-1 text-left text-sm leading-tight group-data-[collapsible=icon]:hidden">
                      <span className="truncate font-semibold">
                        {user.name || user.login}
                      </span>
                      <span className="truncate text-xs">{user.email}</span>
                    </div>
                    <ChevronDown className="ml-auto size-4 group-data-[collapsible=icon]:hidden" />
                  </button>
                </DropdownMenuTrigger>
                <DropdownMenuContent className="w-[200px]" align="start">
                  <DropdownMenuItem asChild>
                    <Link to="/profile">
                      <User className="h-4 w-4 mr-2" />
                      Account
                    </Link>
                  </DropdownMenuItem>
                  <DropdownMenuItem disabled>
                    <Settings className="h-4 w-4 mr-2" />
                    Notifications
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem onClick={logout} className="text-red-600">
                    <LogOut className="h-4 w-4 mr-2" />
                    Log out
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </SidebarMenuItem>
          </SidebarMenu>
        )}
      </SidebarFooter>
    </Sidebar>
  );
}

export default AppSidebar;
