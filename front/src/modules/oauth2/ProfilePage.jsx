import { useState, useEffect } from "react"
import { Card, Row, Col, Avatar, Descriptions, Button, Space, Divider, Typography, Tag, List } from "antd"
import { UserOutlined, SettingOutlined, GithubOutlined, EditOutlined, TeamOutlined, KeyOutlined } from "@ant-design/icons"
import { getUserData, getUserPermissions } from "./userResource"

const { Title, Text } = Typography

function ProfilePage() {
  const [user, setUser] = useState(null)
  const [permissions, setPermissions] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function fetchData() {
      try {
        const [userResult, permissionsResult] = await Promise.all([
          getUserData(),
          getUserPermissions()
        ])
        setUser(userResult.data)
        setPermissions(permissionsResult.data)
      } catch (error) {
        console.error('Failed to fetch data:', error)
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [])

  if (loading) {
    return (
      <Card loading={true} style={{ margin: 24 }}>
        <div style={{ height: 200 }} />
      </Card>
    )
  }

  if (!user) {
    return (
      <Card style={{ margin: 24 }}>
        <div style={{ textAlign: 'center', padding: 40 }}>
          <Text type="secondary">Failed to load user data</Text>
        </div>
      </Card>
    )
  }

  return (
    <div style={{ padding: 24 }}>
      <Row gutter={[24, 24]}>
        <Col xs={24} md={8}>
          <Card>
            <div style={{ textAlign: 'center', padding: 20 }}>
              <Avatar 
                size={120} 
                src={user.avatar_url} 
                icon={<UserOutlined />}
                style={{ marginBottom: 16 }}
              />
              <Title level={3} style={{ margin: 0 }}>{user.name || user.login}</Title>
              <Text type="secondary">@{user.login}</Text>
              {user.bio && (
                <div style={{ marginTop: 16 }}>
                  <Text>{user.bio}</Text>
                </div>
              )}
              <Divider />
              <Space direction="vertical" style={{ width: '100%' }}>
                <Button 
                  type="primary" 
                  icon={<EditOutlined />}
                  block
                  disabled
                >
                  Edit Profile
                </Button>
                <Button 
                  icon={<SettingOutlined />}
                  block
                  disabled
                >
                  OAuth2 Settings
                </Button>
              </Space>
            </div>
          </Card>
        </Col>
        
        <Col xs={24} md={16}>
          <Card title="User Information" icon={<UserOutlined />}>
            <Descriptions column={1} bordered>
              <Descriptions.Item label="Username">
                <Text code>{user.login}</Text>
              </Descriptions.Item>
              <Descriptions.Item label="Full Name">
                {user.name || 'Not provided'}
              </Descriptions.Item>
              <Descriptions.Item label="Email">
                {user.email || 'Not provided'}
              </Descriptions.Item>
              <Descriptions.Item label="Location">
                {user.location || 'Not provided'}
              </Descriptions.Item>
              <Descriptions.Item label="Company">
                {user.company || 'Not provided'}
              </Descriptions.Item>
              <Descriptions.Item label="Blog">
                {user.blog ? (
                  <a href={user.blog} target="_blank" rel="noopener noreferrer">
                    {user.blog}
                  </a>
                ) : 'Not provided'}
              </Descriptions.Item>
              <Descriptions.Item label="GitHub Profile">
                <a href={user.html_url} target="_blank" rel="noopener noreferrer">
                  <GithubOutlined /> View on GitHub
                </a>
              </Descriptions.Item>
            </Descriptions>
          </Card>

          <Card title="User Permissions" style={{ marginTop: 24 }} icon={<KeyOutlined />}>
            {permissions && (
              <>
                <div style={{ marginBottom: 16 }}>
                  <Text strong>Organization: </Text>
                  <Tag color="blue">{permissions.organization}</Tag>
                </div>
                
                <div style={{ marginBottom: 16 }}>
                  <Text strong>Teams: </Text>
                  {permissions.teams && permissions.teams.length > 0 ? (
                    permissions.teams.map((team, index) => (
                      <Tag key={index} color="green" style={{ marginLeft: 8 }}>
                        {team.name || team.slug}
                      </Tag>
                    ))
                  ) : (
                    <Text type="secondary">No teams found</Text>
                  )}
                </div>

                <Divider />

                <div>
                  <Text strong>Plugin Permissions:</Text>
                  <List
                    size="small"
                    dataSource={Object.entries(permissions.permissions || {})}
                    renderItem={([plugin, pluginPerms]) => (
                      <List.Item>
                        <div style={{ width: '100%' }}>
                          <Text strong style={{ textTransform: 'uppercase' }}>{plugin}:</Text>
                          {Object.entries(pluginPerms).map(([feature, actions]) => (
                            <div key={feature} style={{ marginLeft: 16, marginTop: 8 }}>
                              <Text type="secondary">{feature}:</Text>
                              {Array.isArray(actions) ? actions.map((action, idx) => (
                                <Tag key={idx} color="orange" style={{ marginLeft: 8 }}>
                                  {action}
                                </Tag>
                              )) : (
                                <Tag color="orange" style={{ marginLeft: 8 }}>
                                  {actions}
                                </Tag>
                              )}
                            </div>
                          ))}
                        </div>
                      </List.Item>
                    )}
                  />
                </div>
              </>
            )}
          </Card>

          <Card title="OAuth2 Configuration" style={{ marginTop: 24 }} icon={<SettingOutlined />}>
            <div style={{ padding: 20, textAlign: 'center' }}>
              <Text type="secondary">
                OAuth2 configuration settings will be available here in future updates.
              </Text>
              <br />
              <Text type="secondary">
                This will include client ID management, scopes, and redirect URL configurations.
              </Text>
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default ProfilePage
