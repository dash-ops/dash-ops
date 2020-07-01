import React, { useState, useEffect } from "react"
import { Switch, Route, Link, useParams, useLocation } from "react-router-dom"
import { Row, Col, Menu } from "antd"

export default function ContentWithMenu({ pages }) {
  const { context } = useParams()
  const location = useLocation()
  const [current, setCurrent] = useState("/")

  useEffect(() => {
    setCurrent(location.pathname)
  }, [location.pathname])

  return (
    <Row gutter={16}>
      <Col xs={18} md={3}>
        <Menu
          onClick={(e) => setCurrent(e.key)}
          selectedKeys={[current]}
          mode="inline"
          theme="light"
        >
          {pages.map((menu) => (
            <Menu.Item key={menu.path.replace(/:context/, context)}>
              <Link to={menu.path.replace(/:context/, context)}>{menu.name}</Link>
            </Menu.Item>
          ))}
        </Menu>
      </Col>
      <Col xs={18} md={21}>
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
