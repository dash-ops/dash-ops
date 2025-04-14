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
    label: name,
    icon: <DeploymentUnitOutlined />,
    key: `k8s-${context}`,
    link: `/k8s/${context}`,
  }))

  const pages = [
    {
      name: "Cluster",
      path: "/k8s/:context",
      menu: true,
      component: ClusterPage,
    },
    {
      name: "Deployments",
      path: "/k8s/:context/deployments",
      menu: true,
      component: DeploymentPage,
    },
    {
      name: "Pods",
      path: "/k8s/:context/pods",
      menu: true,
      component: PodPage,
    },
    {
      name: "PodLogs",
      path: "/k8s/:context/pod/logs",
      menu: false,
      component: PodLogPage,
    },
  ]

  return {
    menus,
    routers: [
      {
        key: "k8s",
        path: "/k8s/:context/*",
        element: <ContentWithMenu pages={pages} />,
      },
    ],
  }
}
