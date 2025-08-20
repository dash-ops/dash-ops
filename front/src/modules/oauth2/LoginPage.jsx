import { useLocation } from "react-router"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Github } from "lucide-react"

function LoginPage() {
  const location = useLocation()
  const from = location.state || "/"
  const urlLoginGithub = `${import.meta.env.VITE_API_URL || 'http://localhost:8080/api'}/oauth?redirect_url=${from}`

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4">
      <div className="w-full max-w-md">
        <Card>
          <CardHeader>
            <CardTitle className="text-center text-2xl font-bold">
              DashOPS - Beta
            </CardTitle>
          </CardHeader>
          <CardContent>
            <Button 
              className="w-full gap-2" 
              size="lg" 
              asChild
            >
              <a href={urlLoginGithub}>
                <Github className="h-5 w-5" />
                Login with GitHub
              </a>
            </Button>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

export default LoginPage
