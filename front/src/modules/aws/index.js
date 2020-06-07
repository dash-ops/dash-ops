import React from "react"
import { CloudOutlined } from "@ant-design/icons"
import InstancePage from "./InstancePage"

export default {
  routers: [
    {
      name: "EC2 Instance",
      icon: <CloudOutlined />,
      path: "/ec2/instance",
      component: InstancePage,
    },
  ],
}
