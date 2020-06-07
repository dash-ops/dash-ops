import React from "react"
import PropTypes from "prop-types"
import { Tag } from "antd"

const instanceColors = {
  pending: "gold",
  running: "green",
  "shutting-down": "orange",
  terminated: "volcano",
  stopping: "purple",
  stopped: "red",
  loading: "blue",
}

function InstanceTag({ state }) {
  return <Tag color={instanceColors[state]}>{state}</Tag>
}

InstanceTag.propTypes = {
  state: PropTypes.string.isRequired,
}

export default InstanceTag
