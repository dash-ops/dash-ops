import React, { useState } from "react"
import { BrowserRouter as Router, Switch, Route } from "react-router-dom"
import { Layout } from "antd"
import "./App.css"
import PrivateRoute from "./components/PrivateRoute"
import Sidebar from "./components/Sidebar"
import Topbar from "./components/Topbar"
import Footer from "./components/Footer"
import DashboardModule from "./modules/dashboard"
import KubernetesModule from "./modules/kubernetes"
import Login from "./pages/Login"

export default function App() {
  const routers = [...DashboardModule.routers, ...KubernetesModule.routers]
  const [collapsed, setCollapsed] = useState(true)
  
  const onCollapse = data => {
    setCollapsed(data)
  };

  return (
    <Router>
      <Switch>
        <Route path="/login">
          <Login />
        </Route>
        <PrivateRoute path="/">
          <Layout className="layout">
            <Layout.Sider collapsible collapsed={collapsed} onCollapse={onCollapse}>
              <Sidebar menus={routers} />
            </Layout.Sider>
            <Layout>
              <Layout.Header className="header">
                <Topbar />
              </Layout.Header>
              <Layout.Content className="content">
                <div className="container">
                  <Switch>
                    {routers.map(route => (
                      <PrivateRoute key={route.name} path={route.path} exact={route.exact}>
                        <route.component />
                      </PrivateRoute>
                    ))}
                  </Switch>
                </div>
              </Layout.Content>
              <Layout.Footer className="footer">
                <Footer />
              </Layout.Footer>
            </Layout>
          </Layout>
        </PrivateRoute>
      </Switch>
    </Router>
  )
}
