import React from "react"
import PropTypes from "prop-types"
import { Button, Tooltip } from "antd"
import { PlayCircleOutlined, PauseCircleOutlined } from "@ant-design/icons"

function DeploymentActions({ deployment, toUp, toDown }) {
  const showUpButton = deployment.pod_count === 0

  return (
    <Button.Group>
      {showUpButton && (
        <Tooltip title="Up deployment" placement="topRight">
          <Button
            type="primary"
            ghost
            size="small"
            icon={<PlayCircleOutlined />}
            disabled={deployment.pod_count > 0}
            onClick={toUp}
          >
            Up
          </Button>
        </Tooltip>
      )}
      {!showUpButton && (
        <Tooltip title="Down deployment" placement="topRight">
          <Button
            type="danger"
            ghost
            size="small"
            icon={<PauseCircleOutlined />}
            disabled={showUpButton}
            onClick={toDown}
          >
            Down
          </Button>
        </Tooltip>
      )}
    </Button.Group>
  )
}

DeploymentActions.propTypes = {
  deployment: PropTypes.shape({
    name: PropTypes.string,
    namespace: PropTypes.string,
    pod_count: PropTypes.number,
  }).isRequired,
  toUp: PropTypes.func.isRequired,
  toDown: PropTypes.func.isRequired,
}

export default DeploymentActions
