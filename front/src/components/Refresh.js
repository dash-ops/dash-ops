import { useState } from "react";
import PropTypes from "prop-types"
import { Button, Checkbox } from "antd"
import { ReloadOutlined } from "@ant-design/icons"
import { setItem, getItem } from "../helpers/localStorage"
import useInterval from "../helpers/useInterval"

function Refresh(props) {
  const { onReload } = props
  const initState = getItem("auto-refresh") === "true"
  const [isAutoRefresh, setIsAutoRefresh] = useState(initState)

  useInterval(
    async () => {
      onReload()
    },
    isAutoRefresh ? 10000 : null,
  )

  function onAutoRefresh(e) {
    setIsAutoRefresh(e.target.checked)
    setItem("auto-refresh", e.target.checked)
  }

  return (
    <>
      <Checkbox checked={isAutoRefresh} onChange={onAutoRefresh}>
        Auto Refresh 10s
      </Checkbox>
      <Button icon={<ReloadOutlined />} onClick={onReload} />
    </>
  )
}

Refresh.propTypes = {
  onReload: PropTypes.func.isRequired,
}

export default Refresh
