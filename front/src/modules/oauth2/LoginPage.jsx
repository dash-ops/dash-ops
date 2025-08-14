import { useLocation } from "react-router"
import { Row, Col, Card, Button } from "antd"
import { GithubFilled } from "@ant-design/icons"

function LoginPage() {
  const location = useLocation()
  const from = location.state || "/"
  const urlLoginGithub = `${import.meta.env.VITE_API_URL || 'http://localhost:8080/api'}/oauth?redirect_url=${from}`

  return (
    <Row gutter={16}>
      <Col span={8} />
      <Col span={8}>
        <Card title="Dash-OPS - Beta" bordered={false} style={{ top: 40 }}>
          <Button type="primary" block icon={<GithubFilled />} size="large" href={urlLoginGithub}>
            Login
          </Button>
        </Card>
      </Col>
      <Col span={8} />
    </Row>
  )
}

export default LoginPage
