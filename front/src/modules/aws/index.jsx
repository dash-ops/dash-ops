import { CloudOutlined } from "@ant-design/icons"
import InstancePage from "./InstancePage"
import { getAccounts } from "./accountResource"
import ContentWithMenu from "../kubernetes/ContentWithMenu"

const AwsModule = async () => {
  const { data } = await getAccounts()
  const menus = data.map(({ name }) => ({
    label: name,
    icon: <CloudOutlined />,
    key: `aws-${name}`,
    link: `/aws/${name}`,
  }))

  const pages = [
    {
      name: "EC2 Instances",
      path: "/aws/:key/ec2",
      menu: true,
      component: InstancePage,
    },
  ]

  return {
    menus,
    routers: [
      {
        key: "aws",
        path: "/aws/:key/*",
        element: <ContentWithMenu pages={pages} />,
      },
    ],
  }
}

export default AwsModule;
