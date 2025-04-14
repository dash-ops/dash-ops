import { useState, useEffect } from "react"
import PropTypes from "prop-types"
import { useLocation, useNavigate } from "react-router"
import { Menu } from "antd"

function Sidebar({ menus = [] }) {
  const location = useLocation()
  const navigate = useNavigate();
  const [current, setCurrent] = useState("")

  useEffect(() => {
    setCurrent(location.pathname)
  }, [location.pathname])

  const onClick = (e) => {
    const menuItem = menus.find(menu => menu.key === e.key)
    if (menuItem) {
      navigate(menuItem.link)
    }
  };

  return (
    <Menu 
      onClick={onClick} 
      selectedKeys={[current]} 
      mode="inline" 
      theme="dark" 
      items={menus.map(menu => ({
        key: menu.key,
        icon: menu.icon,
        label: menu.label
      }))} 
    />
  )
}

Sidebar.propTypes = {
  menus: PropTypes.arrayOf(
    PropTypes.shape({
      label: PropTypes.string,
      icon: PropTypes.element,
      key: PropTypes.string,
      link: PropTypes.string,
    }),
  ),
}

export default Sidebar
