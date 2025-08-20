import { useState, useEffect } from "react"
import { toast } from "sonner"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Separator } from "@/components/ui/separator"
import { User, Settings, Github, Edit, Users, Key } from "lucide-react"
import { getUserData, getUserPermissions } from "./userResource"

function ProfilePage() {
  const [user, setUser] = useState(null)
  const [permissions, setPermissions] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function fetchData() {
      try {
        const [userResult, permissionsResult] = await Promise.all([
          getUserData(),
          getUserPermissions()
        ])
        setUser(userResult.data)
        setPermissions(permissionsResult.data)
      } catch (error) {
        console.error('Failed to fetch data:', error)
        toast.error("Failed to load profile data")
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [])

  if (loading) {
    return (
      <div className="p-6">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center space-x-4">
              <div className="h-12 w-12 bg-gray-200 rounded-full animate-pulse" />
              <div className="space-y-2">
                <div className="h-4 bg-gray-200 rounded w-32 animate-pulse" />
                <div className="h-3 bg-gray-200 rounded w-24 animate-pulse" />
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    )
  }

  if (!user) {
    return (
      <div className="p-6">
        <Card>
          <CardContent className="p-6">
            <div className="text-center py-10">
              <p className="text-muted-foreground">Failed to load user data</p>
            </div>
          </CardContent>
        </Card>
      </div>
    )
  }

  return (
    <div className="p-6">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="md:col-span-1">
          <Card>
            <CardContent className="p-6 text-center">
              <Avatar className="h-28 w-28 mx-auto mb-4">
                <AvatarImage src={user.avatar_url} alt={user.name || user.login} />
                <AvatarFallback className="text-lg">
                  <User className="h-8 w-8" />
                </AvatarFallback>
              </Avatar>
              <h3 className="text-lg font-semibold mb-1">{user.name || user.login}</h3>
              <p className="text-muted-foreground text-sm">@{user.login}</p>
              {user.bio && (
                <div className="mt-4">
                  <p className="text-sm">{user.bio}</p>
                </div>
              )}
              <Separator className="my-4" />
              <div className="space-y-2">
                <Button className="w-full gap-2" disabled>
                  <Edit className="h-4 w-4" />
                  Edit Profile
                </Button>
                <Button variant="outline" className="w-full gap-2" disabled>
                  <Settings className="h-4 w-4" />
                  OAuth2 Settings
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
        
        <div className="md:col-span-2 space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <User className="h-5 w-5" />
                User Information
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="grid grid-cols-3 gap-2 py-2 border-b">
                  <span className="font-medium text-sm">Username:</span>
                  <span className="col-span-2 text-sm font-mono bg-muted px-2 py-1 rounded">{user.login}</span>
                </div>
                <div className="grid grid-cols-3 gap-2 py-2 border-b">
                  <span className="font-medium text-sm">Full Name:</span>
                  <span className="col-span-2 text-sm">{user.name || 'Not provided'}</span>
                </div>
                <div className="grid grid-cols-3 gap-2 py-2 border-b">
                  <span className="font-medium text-sm">Email:</span>
                  <span className="col-span-2 text-sm">{user.email || 'Not provided'}</span>
                </div>
                <div className="grid grid-cols-3 gap-2 py-2 border-b">
                  <span className="font-medium text-sm">Location:</span>
                  <span className="col-span-2 text-sm">{user.location || 'Not provided'}</span>
                </div>
                <div className="grid grid-cols-3 gap-2 py-2 border-b">
                  <span className="font-medium text-sm">Company:</span>
                  <span className="col-span-2 text-sm">{user.company || 'Not provided'}</span>
                </div>
                <div className="grid grid-cols-3 gap-2 py-2 border-b">
                  <span className="font-medium text-sm">Blog:</span>
                  <span className="col-span-2 text-sm">
                    {user.blog ? (
                      <a href={user.blog} target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline">
                        {user.blog}
                      </a>
                    ) : 'Not provided'}
                  </span>
                </div>
                <div className="grid grid-cols-3 gap-2 py-2">
                  <span className="font-medium text-sm">GitHub Profile:</span>
                  <span className="col-span-2 text-sm">
                    <a href={user.html_url} target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline flex items-center gap-1">
                      <Github className="h-4 w-4" /> View on GitHub
                    </a>
                  </span>
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Key className="h-5 w-5" />
                User Permissions
              </CardTitle>
            </CardHeader>
            <CardContent>
            {permissions && (
              <div className="space-y-4">
                <div className="flex items-center gap-2">
                  <span className="font-medium text-sm">Organization:</span>
                  <Badge variant="default">{permissions.organization}</Badge>
                </div>
                
                <div className="flex items-start gap-2">
                  <span className="font-medium text-sm">Teams:</span>
                  <div className="flex flex-wrap gap-1">
                    {permissions.teams && permissions.teams.length > 0 ? (
                      permissions.teams.map((team, index) => (
                        <Badge key={index} variant="secondary">
                          {team.name || team.slug}
                        </Badge>
                      ))
                    ) : (
                      <span className="text-muted-foreground text-sm">No teams found</span>
                    )}
                  </div>
                </div>

                <Separator />

                <div>
                  <h4 className="font-medium text-sm mb-3">Plugin Permissions:</h4>
                  <div className="space-y-3">
                    {Object.entries(permissions.permissions || {}).map(([plugin, pluginPerms]) => (
                      <div key={plugin} className="border rounded-lg p-3">
                        <h5 className="font-medium text-sm uppercase mb-2">{plugin}:</h5>
                        <div className="space-y-2">
                          {Object.entries(pluginPerms).map(([feature, actions]) => (
                            <div key={feature} className="ml-4">
                              <span className="text-muted-foreground text-sm">{feature}:</span>
                              <div className="flex flex-wrap gap-1 mt-1">
                                {Array.isArray(actions) ? actions.map((action, idx) => (
                                  <Badge key={idx} variant="outline" className="text-xs">
                                    {action}
                                  </Badge>
                                )) : (
                                  <Badge variant="outline" className="text-xs">
                                    {actions}
                                  </Badge>
                                )}
                              </div>
                            </div>
                          ))}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            )}
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Settings className="h-5 w-5" />
                OAuth2 Configuration
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-center py-6">
                <p className="text-muted-foreground text-sm mb-2">
                  OAuth2 configuration settings will be available here in future updates.
                </p>
                <p className="text-muted-foreground text-sm">
                  This will include client ID management, scopes, and redirect URL configurations.
                </p>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}

export default ProfilePage
