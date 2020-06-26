import React, { useState } from "react"
import { Switch, Route, Link, useLocation } from "react-router-dom"
import { Row, Col, Menu } from "antd"

export default function ContentWithMenu({ routers }) {
  const location = useLocation()
  const [current, setCurrent] = useState(location.pathname)

  return (
    <Row gutter={16}>
      <Col xs={18} md={3}>
        <Menu
          onClick={(e) => setCurrent(e.key)}
          selectedKeys={[current]}
          mode="inline"
          theme="light"
        >
          {routers.map((route) => (
            <Menu.Item key={route.path}>
              <Link to={route.path}>{route.name}</Link>
            </Menu.Item>
          ))}
        </Menu>
      </Col>
      <Col xs={18} md={21}>
        <Switch>
          {routers.map((route) => (
            <Route key={route.name} path={route.path} exact={route.exact}>
              <route.component />
            </Route>
          ))}
        </Switch>
      </Col>
    </Row>
  )
}
