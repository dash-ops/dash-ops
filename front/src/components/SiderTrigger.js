import { MenuUnfoldOutlined, MenuFoldOutlined } from "@ant-design/icons"
import "./SiderTrigger.css"

function SiderTrigger({ collapsed, onCollapse }) {
  if (collapsed) {
    return <MenuUnfoldOutlined className="trigger" onClick={() => onCollapse(!collapsed)} />
  }
  return <MenuFoldOutlined className="trigger" onClick={() => onCollapse(!collapsed)} />
}

export default SiderTrigger
