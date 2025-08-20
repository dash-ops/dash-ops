import PropTypes from "prop-types"
import { Progress } from "@/components/ui/progress"

function ProgressData({ percent }) {
  return (
    <div className="w-full max-w-[170px]">
      <Progress 
        value={percent} 
        className="h-2"
      />
      <span className="text-xs text-muted-foreground">{percent?.toFixed(1)}%</span>
    </div>
  )
}

ProgressData.propTypes = {
  percent: PropTypes.number.isRequired,
}

export default ProgressData
