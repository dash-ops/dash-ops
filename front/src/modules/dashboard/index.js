import { DashboardOutlined } from "@ant-design/icons"
import DashboardPage from "./DashboardPage"

export default {
  menus: [
    {
      name: "Dashboard",
      icon: <DashboardOutlined />,
      link: "/",
    },
  ],
  routers: [
    {
      key: "dashboard",
      path: "/",
      exact: true,
      component: DashboardPage,
    },
  ],
}
