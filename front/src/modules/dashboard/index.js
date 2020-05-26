import React from "react"
import DashboardPage from "./DashboardPage"
import { DashboardOutlined } from "@ant-design/icons"

export default {
  routers: [
    {
      name: "Dashboard",
      icon: <DashboardOutlined />,
      path: "/",
      exact: true,
      component: DashboardPage,
    },
  ],
}
