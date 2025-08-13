
import { useState, useEffect, useRef } from "react"
import { Routes, Route } from "react-router"
import { Layout, notification } from "antd"
import { loadModulesConfig } from "./helpers/loadModules"
import { verifyToken } from "./helpers/oauth"
import Sidebar from "./components/Sidebar"
import Toolbar from "./components/Toolbar"
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
  const initialized = useRef(false)

  useEffect(() => {
    if (initialized.current) return;
    initialized.current = true;
    
    verifyToken()
    loadModulesConfig()
      .then((modules) => {
        setOAuth2(modules.oAuth2)
        setMenus([...DashboardModule.menus, ...modules.menus])
        setRouters([...DashboardModule.routers, ...modules.routers])
      })
      .catch(() => {
        notification.error({ message: "Failed to load plugins" })
      })

  }, [])

  const onCollapse = (data) => {
    setCollapsed(data)
  }

  return (
    <Routes>
      {oAuth2.active && (
        <Route path="/login" element={<oAuth2.LoginPage />} />
      )}
      <Route
        path="*"
        element={
          <Layout className="dash-layout">
            <Layout.Header className="dash-header">
              <SiderTrigger collapsed={collapsed} onCollapse={onCollapse} />
              <Logo />
              <Toolbar oAuth2={oAuth2.active} />
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
                    <Routes>
                      {routers.map((route) => (
                        <Route
                          key={route.key}
                          path={route.path}
                          element={route.element}
                        />
                      ))}
                    </Routes>
                  </div>
                </Layout.Content>
                <Layout.Footer className="dash-footer">
                  <Footer />
                </Layout.Footer>
              </Layout>
            </Layout>
          </Layout>
        }
      />
    </Routes>
  )
}
