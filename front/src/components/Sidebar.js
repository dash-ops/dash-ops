import React, { useState, useEffect } from "react"
import PropTypes from "prop-types"
import { Link, useLocation } from "react-router-dom"
import { Menu } from "antd"
import "./Sidebar.css"
// import logo from "../logo.png"

function Sidebar({ menus }) {
  const location = useLocation()
  const [current, setCurrent] = useState([])

  useEffect(() => {
    const parent = location.pathname.substring(0, location.pathname.lastIndexOf("/"))
    setCurrent([parent, location.pathname])
  }, [location.pathname])

  return (
    <>
      <div className="dash-logo">
        {/* <img src={logo} alt="DashOPS - Beta" /> */}
        DashOPS
      </div>
      <Menu onClick={(e) => setCurrent(e.key)} selectedKeys={current} mode="inline" theme="dark">
        {menus.map((menu) => (
          <Menu.Item key={menu.link}>
            <Link to={menu.link}>
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
      link: PropTypes.string,
    }),
  ),
}

Sidebar.defaultProps = {
  menus: [],
}

export default Sidebar
