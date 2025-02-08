import { CloudOutlined } from "@ant-design/icons"
import InstancePage from "./InstancePage"
import { getAccounts } from "./accountResource"

export default async () => {
  const { data } = await getAccounts()
  const menus = data.map(({ name }) => ({
    name,
    icon: <CloudOutlined />,
    key: 'aws',
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
