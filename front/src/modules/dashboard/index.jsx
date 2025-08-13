import { DashboardOutlined } from "@ant-design/icons"
import DashboardPage from "./DashboardPage"

export default {
  menus: [
    {
      label: "Dashboard",
      icon: <DashboardOutlined />,
      key: "dashboard",
      link: "/",
    },
  ],
  routers: [
    {
      key: "dashboard",
      path: "/",
      element: <DashboardPage />,
    },
  ],
}
