import React, { useState } from "react"
import PropTypes from "prop-types"
import { Link, useLocation } from "react-router-dom"
import { Menu } from "antd"
import "./Sidebar.css"
// import logo from "../logo.png"

function Sidebar({ menus }) {
  const location = useLocation()
  const [current, setCurrent] = useState(location.pathname)

  return (
    <>
      <div className="dash-logo">
        {/* <img src={logo} alt="DashOPS - Beta" /> */}
        DashOPS
      </div>
      <Menu onClick={(e) => setCurrent(e.key)} selectedKeys={[current]} mode="inline" theme="dark">
        {menus.map((menu) => (
          <Menu.Item key={menu.path}>
            <Link to={menu.path}>
              {menu.icon ?? <></>}
              <span>{menu.name}</span>
            </Link>
          </Menu.Item>
        ))}
      </Menu>
    </>
  )
}

Sidebar.propTypes = {
  menus: PropTypes.arrayOf(
    PropTypes.shape({
      name: PropTypes.string,
      icon: PropTypes.element,
      path: PropTypes.string,
      component: PropTypes.func,
    }),
  ),
}

Sidebar.defaultProps = {
  menus: [],
}

export default Sidebar
