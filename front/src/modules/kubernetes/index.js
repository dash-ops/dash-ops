import React from "react"
import { DeploymentUnitOutlined } from "@ant-design/icons"
import DeploymentPage from "./DeploymentPage"

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
