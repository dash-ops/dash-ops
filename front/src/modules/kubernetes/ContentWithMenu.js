import { useState, useEffect } from "react"
import PropTypes from "prop-types"
import { Switch, Route, Link, useParams, useLocation } from "react-router-dom"
import { Row, Col, Menu } from "antd"

function ContentWithMenu({ pages }) {
  const { context } = useParams()
  const location = useLocation()
  const [current, setCurrent] = useState("/")

  useEffect(() => {
    setCurrent(location.pathname)
  }, [location.pathname])

  return (
    <Row gutter={16}>
      <Col xs={24} md={5} lg={4} xl={3}>
        <Menu
          onClick={(e) => setCurrent(e.key)}
          selectedKeys={[current]}
          mode="inline"
          theme="light"
        >
          {pages.map((route) => {
            return route.menu ? (
              <Menu.Item key={route.path.replace(/:context/, context)}>
                <Link to={route.path.replace(/:context/, context)}>{route.name}</Link>
              </Menu.Item>
            ) : null
          })}
        </Menu>
      </Col>
      <Col xs={24} md={19} lg={20} xl={21}>
        <Switch>
          {pages.map((route) => (
            <Route key={route.name} path={route.path} exact={route.exact}>
              <route.component />
            </Route>
          ))}
        </Switch>
      </Col>
    </Row>
  )
}

ContentWithMenu.propTypes = {
  pages: PropTypes.arrayOf({
    name: PropTypes.string,
    path: PropTypes.string,
    exact: PropTypes.bool,
  }).isRequired,
}

export default ContentWithMenu
