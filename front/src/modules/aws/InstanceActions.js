import PropTypes from "prop-types"
import { Button, Tooltip, notification } from "antd"
import { DesktopOutlined, PlayCircleOutlined, PauseCircleOutlined } from "@ant-design/icons"

function ssh(instance) {
  if (instance.platform === "windows") {
    notification.error({
      message: `Sorry... I'm afraid I can't do that...`,
      description: `
        Windows does not provides a method to connect a Remote Desktop via URL.
        You can try to connect via command line using on Windows: mstsc /v:${instance.name}
      `,
    })
    return
  }
  window.location = `ssh://${instance.name}`
}

function InstanceActions({ instance, toStart, toStop }) {
  const showPlayButton = instance.state !== "running"

  return (
    <Button.Group>
      <Tooltip title="SSH access">
        <Button type="primary" ghost size="small" icon={<DesktopOutlined />} onClick={() => ssh(instance)} />
      </Tooltip>
      {showPlayButton && (
        <Tooltip title="Start instance" placement="topRight">
          <Button
            type="primary"
            ghost
            size="small"
            icon={<PlayCircleOutlined />}
            disabled={instance.state !== "stopped"}
            onClick={toStart}
          >
            Start
          </Button>
        </Tooltip>
      )}
      {!showPlayButton && (
        <Tooltip title="Stop instance" placement="topRight">
          <Button
            type="danger"
            ghost
            size="small"
            icon={<PauseCircleOutlined />}
            disabled={showPlayButton}
            onClick={toStop}
          >
            Stop
          </Button>
        </Tooltip>
      )}
    </Button.Group>
  )
}

InstanceActions.propTypes = {
  instance: PropTypes.objectOf(PropTypes.string).isRequired,
  toStart: PropTypes.func.isRequired,
  toStop: PropTypes.func.isRequired,
}

export default InstanceActions
