import { useState, useEffect } from "react"
import PropTypes from "prop-types"
import { useLocation, Link } from "react-router"
import { Menu } from "antd"

function Sidebar({ menus = [] }) {
  const location = useLocation()
  const [current, setCurrent] = useState([])

  useEffect(() => {
    const parent = location.pathname.substring(0, location.pathname.lastIndexOf("/"))
    const grandparent = parent.substring(0, parent.lastIndexOf("/"))
    setCurrent([grandparent, parent, location.pathname])
  }, [location.pathname])
  console.log("menus", menus)

  return (
    <Menu onClick={(e) => setCurrent(e.key)} selectedKeys={[current]} mode="inline" theme="dark" items={menus}>
      {menus.map((menu, index) => (
        <Menu.Item key={index.toString()}>
          <Link to={menu.link}>
            {menu.icon ?? <></>}
            <span>{menu.label}</span>
          </Link>
        </Menu.Item>
      ))}
    </Menu>
  )
}

Sidebar.propTypes = {
  menus: PropTypes.arrayOf(
    PropTypes.shape({
      label: PropTypes.string,
      icon: PropTypes.element,
      key: PropTypes.string,
    }),
  ),
}

export default Sidebar
