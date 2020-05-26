import React from "react"
import DeploymentPage from "./DeploymentPage"
import { DeploymentUnitOutlined } from "@ant-design/icons"

export default {
  routers: [
    {
      name: "K8S Deployments",
      icon: <DeploymentUnitOutlined />,
      path: "/k8s/deployments",
      component: DeploymentPage,
    },
  ],
}
