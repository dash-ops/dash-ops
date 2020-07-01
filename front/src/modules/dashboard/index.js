import React from "react"
import DashboardPage from "./DashboardPage"
import { DashboardOutlined } from "@ant-design/icons"

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
