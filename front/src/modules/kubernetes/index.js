import React from "react"
import { DeploymentUnitOutlined } from "@ant-design/icons"
import ClusterPage from "./ClusterPage"
import DeploymentPage from "./DeploymentPage"
import PodPage from "./PodPage"
import PodLogPage from "./PodLogPage"
import ContentWithMenu from "./ContentWithMenu"
import { getClusters } from "./clusterResource"

export default async () => {
  const { data } = await getClusters()
  const menus = data.map(({ name, context }) => ({
    name,
    icon: <DeploymentUnitOutlined />,
    link: `/k8s/${context}`,
  }))

  return {
    menus,
    routers: [
      {
        key: "k8s",
        path: "/k8s/:context",
        component: ContentWithMenu,
        props: {
          pages: [
            {
              name: "Cluster",
              path: "/k8s/:context",
              exact: true,
              menu: true,
              component: ClusterPage,
            },
            {
              name: "Deployments",
              path: "/k8s/:context/deployments",
              exact: true,
              menu: true,
              component: DeploymentPage,
            },
            {
              name: "Pods",
              path: "/k8s/:context/pods",
              exact: true,
              menu: true,
              component: PodPage,
            },
            {
              name: "PodLogs",
              path: "/k8s/:context/pod/logs",
              exact: true,
              menu: false,
              component: PodLogPage,
            },
          ],
        },
      },
    ],
  }
}
