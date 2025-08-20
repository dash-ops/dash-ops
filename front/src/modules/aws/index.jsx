import { Cloud } from "lucide-react"
import InstancePage from "./InstancePage"
import { getAccounts } from "./accountResource"
import ContentWithMenu from "../../components/ContentWithMenu"

const AwsModule = async () => {
  const { data } = await getAccounts()
  const menus = data.map(({ name, key }) => ({
    label: name,
    icon: <Cloud className="h-4 w-4" />,
    key: `aws-${key}`,
    link: `/aws/${key}`,
  }))

  const pages = [
    {
      name: "EC2 Instances",
      path: "/aws/:key/ec2",
      menu: true,
      element: <InstancePage />,
    },
  ]

  return {
    menus,
    routers: [
      {
        key: "aws",
        path: "/aws/:key/*",
        element: <ContentWithMenu pages={pages} paramName="key" />,
      },
    ],
  }
}

export default AwsModule;
