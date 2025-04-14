import PropTypes from "prop-types"
import { NavLink } from "react-router"
import { Button, Tooltip } from "antd"
import { PlayCircleOutlined, PauseCircleOutlined } from "@ant-design/icons"

function DeploymentActions({ context, deployment, toUp, toDown }) {
  const showUpButton = deployment.pod_info.current === 0

  return (
    <Button.Group>
      <Tooltip title={showUpButton ? "No Pods" : "Pods"}>
        <NavLink to={`/k8s/${context}/pods?name=${deployment.name}&namespace=${deployment.namespace}`}>
          <Button type="primary" ghost size="small" disabled={showUpButton}>
            Pods
          </Button>
        </NavLink>
      </Tooltip>
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
            disabled={deployment.pod_count === 0}
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
