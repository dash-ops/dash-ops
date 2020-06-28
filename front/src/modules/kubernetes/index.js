import React from "react"
import { DeploymentUnitOutlined } from "@ant-design/icons"
import ClusterPage from "./ClusterPage"
import DeploymentPage from "./DeploymentPage"
import PodPage from "./PodPage"
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
          {
            name: "Cluster",
            path: "/k8s",
            exact: true,
            component: ClusterPage,
          },
          {
            name: "Deployments",
            path: "/k8s/deployments",
            exact: true,
            component: DeploymentPage,
          },
          {
            name: "Pods",
            path: "/k8s/pods",
            exact: true,
            component: PodPage,
          },
        ],
      },
    },
  ],
}
