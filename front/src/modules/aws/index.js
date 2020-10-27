import { CloudOutlined } from "@ant-design/icons"
import InstancePage from "./InstancePage"
import { getAccounts } from "./accountResource"

export default async () => {
  const { data } = await getAccounts()
  const menus = data.map(({ name, key }) => ({
    name,
    icon: <CloudOutlined />,
    link: `/aws/${key}/ec2`,
  }))

  return {
    menus,
    routers: [
      {
        key: "awsEc2",
        path: "/aws/:key/ec2",
        component: InstancePage,
      },
    ],
  }
}
