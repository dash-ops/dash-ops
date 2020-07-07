import React, { useState, useEffect } from "react"
import { BrowserRouter as Router, Switch, Route } from "react-router-dom"
import { Layout, notification } from "antd"
import { loadModulesConfig } from "./helpers/loadModules"
import InternalRoute from "./components/InternalRoute"
import Sidebar from "./components/Sidebar"
import Topbar from "./components/Topbar"
import Footer from "./components/Footer"
import SiderTrigger from "./components/SiderTrigger"
import Logo from "./components/Logo"
import DashboardModule from "./modules/dashboard"
import "./App.css"

export default function App() {
  const [oAuth2, setOAuth2] = useState({ active: false })
  const [menus, setMenus] = useState([...DashboardModule.menus])
  const [routers, setRouters] = useState([...DashboardModule.routers])
  const [collapsed, setCollapsed] = useState(false)

  useEffect(() => {
    loadModulesConfig()
      .then((modules) => {
        setOAuth2(modules.oAuth2)
        setMenus([...menus, ...modules.menus])
        setRouters([...routers, ...modules.routers])
      })
      .catch(() => {
        notification.error({ message: "Failed to load plugins" })
      })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const onCollapse = (data) => {
    setCollapsed(data)
  }

  return (
    <Router>
      <Switch>
        {oAuth2.active && (
          <Route path="/login">
            <oAuth2.LoginPage />
          </Route>
        )}
        <InternalRoute oAuth2={oAuth2.active} path="/">
          <Layout className="dash-layout">
            <Layout.Header className="dash-header">
              <SiderTrigger collapsed={collapsed} onCollapse={onCollapse} />
              <Logo />
              <Topbar oAuth2={oAuth2.active} />
            </Layout.Header>
            <Layout>
              <Layout.Sider
                trigger={null}
                breakpoint="lg"
                collapsedWidth="0"
                collapsible
                collapsed={collapsed}
                onCollapse={onCollapse}
              >
                <Sidebar menus={menus} />
              </Layout.Sider>
              <Layout>
                <Layout.Content className="dash-content">
                  <div className="dash-container">
                    <Switch>
                      {routers.map((route) => (
                        <InternalRoute
                          oAuth2={oAuth2.active}
                          key={route.key}
                          path={route.path}
                          exact={route.exact}
                        >
                          <route.component {...route.props} />
                        </InternalRoute>
                      ))}
                    </Switch>
                  </div>
                </Layout.Content>
                <Layout.Footer className="dash-footer">
                  <Footer />
                </Layout.Footer>
              </Layout>
            </Layout>
          </Layout>
        </InternalRoute>
      </Switch>
    </Router>
  )
}
