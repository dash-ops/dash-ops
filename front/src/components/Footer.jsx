import { Github } from "lucide-react"

function Footer() {
  return (
    <div className="flex items-center justify-center gap-2 text-sm text-muted-foreground">
      <a 
        href="https://github.com/dash-ops/dash-ops" 
        alt="DashOps Repository"
        className="hover:text-foreground transition-colors flex items-center gap-1"
      >
        <Github className="h-4 w-4" />
      </a>
      <a 
        href="https://dash-ops.github.io/" 
        alt="DashOps WebSite"
        className="hover:text-foreground transition-colors"
      >
        DashOPS
      </a>
      <span>©{new Date().getFullYear()}</span>
    </div>
  )
}

export default Footer
