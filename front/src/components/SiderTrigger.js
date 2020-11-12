import PropTypes from "prop-types"
import { MenuUnfoldOutlined, MenuFoldOutlined } from "@ant-design/icons"
import "./SiderTrigger.css"

function SiderTrigger({ collapsed, onCollapse }) {
  if (collapsed) {
    return <MenuUnfoldOutlined className="trigger" onClick={() => onCollapse(!collapsed)} />
  }
  return <MenuFoldOutlined className="trigger" onClick={() => onCollapse(!collapsed)} />
}

SiderTrigger.propTypes = {
  collapsed: PropTypes.number.isRequired,
  onCollapse: PropTypes.func.isRequired,
}

export default SiderTrigger
