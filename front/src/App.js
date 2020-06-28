import React, { useState, useEffect } from "react"
import { BrowserRouter as Router, Switch, Route } from "react-router-dom"
import { Layout, notification } from "antd"
import { loadModulesRouter } from "./helpers/loadModules"
import PrivateRoute from "./components/PrivateRoute"
import Sidebar from "./components/Sidebar"
import Topbar from "./components/Topbar"
import Footer from "./components/Footer"
import DashboardModule from "./modules/dashboard"
import Login from "./pages/Login"
import "./App.css"

export default function App() {
  const [routers, setRouters] = useState([...DashboardModule.routers])
  const [collapsed, setCollapsed] = useState(true)

  useEffect(() => {
    loadModulesRouter()
      .then((modules) => {
        setRouters([...routers, ...modules])
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
        <Route path="/login">
          <Login />
        </Route>
        <PrivateRoute path="/">
          <Layout className="dash-layout">
            <Layout.Sider collapsible collapsed={collapsed} onCollapse={onCollapse}>
              <Sidebar menus={routers} />
            </Layout.Sider>
            <Layout>
              <Layout.Header className="dash-header">
                <Topbar />
              </Layout.Header>
              <Layout.Content className="dash-content">
                <div className="dash-container" style={{ backgroundColor: "#fff" }}>
                  <Switch>
                    {routers.map((route) => (
                      <PrivateRoute key={route.name} path={route.path} exact={route.exact}>
                        <route.component {...route.props} />
                      </PrivateRoute>
                    ))}
                  </Switch>
                </div>
              </Layout.Content>
              <Layout.Footer className="dash-footer">
                <Footer />
              </Layout.Footer>
            </Layout>
          </Layout>
        </PrivateRoute>
      </Switch>
    </Router>
  )
}
