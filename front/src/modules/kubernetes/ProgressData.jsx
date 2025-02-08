import PropTypes from "prop-types"
import { Progress } from "antd"

function ProgressData({ percent }) {
  return (
    <div style={{ maxWidth: 170 }}>
      <Progress
        percent={percent.toFixed(1)}
        size="small"
        strokeColor={{
          "0%": "#ffad20",
          "100%": "#ff4d4f",
        }}
        status="active"
      />
    </div>
  )
}

ProgressData.propTypes = {
  percent: PropTypes.number.isRequired,
}

export default ProgressData
