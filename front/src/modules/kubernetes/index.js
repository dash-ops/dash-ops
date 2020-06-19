import React from "react"
import { DeploymentUnitOutlined } from "@ant-design/icons"
import DeploymentPage from "./DeploymentPage"
import ContentWithMenu from "../../components/ContentWIthMenu"

export default {
  routers: [
    {
      name: "Kubernetes",
      icon: <DeploymentUnitOutlined />,
      path: "/k8s",
      component: ContentWithMenu,
      props: {
        routers: [
          // {
          //   name: "Nodes",
          //   path: "/k8s",
          //   exact: true,
          //   component: DeploymentPage,
          // },
          {
            name: "Deployments",
            path: "/k8s/deployments",
            exact: true,
            component: DeploymentPage,
          },
        ],
      },
    },
  ],
}
