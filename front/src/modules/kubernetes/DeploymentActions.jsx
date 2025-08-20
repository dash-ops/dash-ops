import PropTypes from "prop-types"
import { NavLink } from "react-router"
import { Button } from "@/components/ui/button"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip"
import { Play, Pause } from "lucide-react"

function DeploymentActions({ context, deployment, toUp, toDown }) {
  const showUpButton = deployment.pod_info.current === 0

  return (
    <TooltipProvider>
      <div className="flex gap-1">
        <Tooltip>
          <TooltipTrigger asChild>
            <Button variant="outline" size="sm" disabled={showUpButton} asChild>
              <NavLink to={`/k8s/${context}/pods?name=${deployment.name}&namespace=${deployment.namespace}`}>
                Pods
              </NavLink>
            </Button>
          </TooltipTrigger>
          <TooltipContent>
            <p>{showUpButton ? "No Pods" : "Pods"}</p>
          </TooltipContent>
        </Tooltip>
        
        {showUpButton && (
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant="outline"
                size="sm"
                disabled={deployment.pod_count > 0}
                onClick={toUp}
                className="gap-2"
              >
                <Play className="h-4 w-4" />
                Up
              </Button>
            </TooltipTrigger>
            <TooltipContent>
              <p>Up deployment</p>
            </TooltipContent>
          </Tooltip>
        )}
        
        {!showUpButton && (
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant="destructive"
                size="sm"
                disabled={deployment.pod_count === 0}
                onClick={toDown}
                className="gap-2"
              >
                <Pause className="h-4 w-4" />
                Down
              </Button>
            </TooltipTrigger>
            <TooltipContent>
              <p>Down deployment</p>
            </TooltipContent>
          </Tooltip>
        )}
      </div>
    </TooltipProvider>
  )
}

DeploymentActions.propTypes = {
  context: PropTypes.string.isRequired,
  deployment: PropTypes.shape({
    name: PropTypes.string,
    namespace: PropTypes.string,
    pod_count: PropTypes.number,
    pod_info: PropTypes.shape({
      current: PropTypes.number,
    }),
  }).isRequired,
  toUp: PropTypes.func.isRequired,
  toDown: PropTypes.func.isRequired,
}

export default DeploymentActions
