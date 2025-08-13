import { useState, useEffect } from "react"
import PropTypes from "prop-types"
import { Routes, Route, useParams, useLocation, useNavigate } from "react-router"
import { Row, Col, Menu } from "antd"

function ContentWithMenu({ pages }) {
  const { context } = useParams()
  const location = useLocation()
  const navigate = useNavigate()
  const [current, setCurrent] = useState(location.pathname)

  useEffect(() => {
    setCurrent(location.pathname)
  }, [location.pathname])

  const onClick = (e) => {
    navigate(e.key)
  }

  const menuItems = pages
    .filter((page) => page.menu)
    .map((page) => ({
      key: page.path.replace(":context", context),
      label: page.name,
    }))

  return (
    <Row gutter={16}>
      <Col xs={24} md={5} lg={4} xl={3}>
        <Menu
          onClick={onClick}
          selectedKeys={[current]}
          mode="inline"
          theme="light"
          items={menuItems}
        />
      </Col>
      <Col xs={24} md={19} lg={20} xl={21}>
        <Routes>
          {pages.map((page) => {
            const path = page.path.replace(":context", context)
            const route = page.path.split(":context").pop()
            return (
              <Route
                key={path}
                path={route}
                element={page.element}
              />
            )
          })}
        </Routes>
      </Col>
    </Row>
  )
}

ContentWithMenu.propTypes = {
  pages: PropTypes.arrayOf(
    PropTypes.shape({
      name: PropTypes.string.isRequired,
      path: PropTypes.string.isRequired,
      menu: PropTypes.bool,
      element: PropTypes.object.isRequired,
    })
  ).isRequired,
}

export default ContentWithMenu
