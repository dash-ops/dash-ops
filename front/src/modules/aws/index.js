import React from "react"
import { CloudOutlined } from "@ant-design/icons"
import InstancePage from "./InstancePage"

export default {
  menus: [
    {
      name: "EC2 Instance",
      icon: <CloudOutlined />,
      link: "/ec2/instance",
    },
  ],
  routers: [
    {
      key: "ec2_instance",
      path: "/ec2/instance",
      component: InstancePage,
    },
  ],
}
