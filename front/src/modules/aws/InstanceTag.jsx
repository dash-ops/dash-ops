import PropTypes from "prop-types"
import { Badge } from "@/components/ui/badge"
import { cn } from "@/lib/utils"

const instanceStyles = {
  pending: "bg-yellow-100 text-yellow-800 border-yellow-300",
  running: "bg-green-100 text-green-800 border-green-300",
  "shutting-down": "bg-orange-100 text-orange-800 border-orange-300",
  terminated: "bg-red-100 text-red-800 border-red-300",
  stopping: "bg-purple-100 text-purple-800 border-purple-300",
  stopped: "bg-red-100 text-red-800 border-red-300",
  loading: "bg-blue-100 text-blue-800 border-blue-300",
}

function InstanceTag({ state }) {
  return (
    <Badge 
      variant="outline" 
      className={cn("capitalize", instanceStyles[state])}
    >
      {state}
    </Badge>
  )
}

InstanceTag.propTypes = {
  state: PropTypes.string.isRequired,
}

export default InstanceTag
