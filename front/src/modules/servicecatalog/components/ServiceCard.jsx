import { Card, Tag, Avatar, Space, Typography, Tooltip, Dropdown } from 'antd';
import { 
  MoreOutlined, 
  GlobalOutlined, 
  TeamOutlined,
  ClockCircleOutlined,
  EditOutlined,
  DeleteOutlined,
  EyeOutlined
} from '@ant-design/icons';
import PropTypes from 'prop-types';

const { Text, Title } = Typography;

function ServiceCard({ service, onEdit, onDelete, onView }) {
  const getTierColor = (tier) => {
    switch (tier) {
      case 'tier-1': return 'red';
      case 'tier-2': return 'orange';
      case 'tier-3': return 'green';
      default: return 'default';
    }
  };

  const getIngressIcon = (type) => {
    return type === 'external' ? <GlobalOutlined /> : <TeamOutlined />;
  };

  const getTeamInitials = (team) => {
    return team
      .split(' ')
      .map(word => word.charAt(0).toUpperCase())
      .join('')
      .substring(0, 2);
  };

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    const now = new Date();
    const diffTime = Math.abs(now - date);
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    
    if (diffDays === 1) return 'Deploy 1 day ago';
    if (diffDays < 30) return `Deploy ${diffDays} days ago`;
    if (diffDays < 365) return `Deploy ${Math.ceil(diffDays / 30)} months ago`;
    return `Deploy ${Math.ceil(diffDays / 365)} years ago`;
  };

  const menuItems = [
    {
      key: 'view',
      label: (
        <Space>
          <EyeOutlined />
          View Details
        </Space>
      ),
      onClick: () => onView?.(service),
    },
    {
      key: 'edit',
      label: (
        <Space>
          <EditOutlined />
          Edit Service
        </Space>
      ),
      onClick: () => onEdit?.(service),
    },
    {
      type: 'divider',
    },
    {
      key: 'delete',
      label: (
        <Space>
          <DeleteOutlined />
          Delete Service
        </Space>
      ),
      danger: true,
      onClick: () => onDelete?.(service),
    },
  ];

  return (
    <Card
      hoverable
      style={{ marginBottom: 16 }}
      actions={[
        <Tooltip title={`${service.regions?.length || 0} regions`}>
          <Space>
            <GlobalOutlined />
            <Text type="secondary">{service.regions?.length || 0} regions</Text>
          </Space>
        </Tooltip>,
        <Tooltip title={service.ingressType === 'external' ? 'External access' : 'Internal access'}>
          <Space>
            {getIngressIcon(service.ingressType)}
            <Text type="secondary">{service.ingressType}</Text>
          </Space>
        </Tooltip>,
        <Tooltip title="Last deployment">
          <Space>
            <ClockCircleOutlined />
            <Text type="secondary">{formatDate(service.updatedAt)}</Text>
          </Space>
        </Tooltip>,
      ]}
    >
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
        <div style={{ flex: 1 }}>
          {/* Header with Tier and Menu */}
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
            <Tag color={getTierColor(service.tier)} style={{ margin: 0 }}>
              {service.tier?.toUpperCase()}
            </Tag>
            <Dropdown menu={{ items: menuItems }} trigger={['click']} placement="bottomRight">
              <MoreOutlined 
                style={{ 
                  fontSize: 16, 
                  padding: 4, 
                  cursor: 'pointer',
                  color: '#8c8c8c'
                }} 
              />
            </Dropdown>
          </div>

          {/* Service Name */}
          <Title level={4} style={{ margin: '8px 0', fontSize: 16 }}>
            {service.displayName || service.name}
          </Title>

          {/* Description */}
          <Text type="secondary" style={{ display: 'block', marginBottom: 12, fontSize: 13 }}>
            {service.description}
          </Text>

          {/* Team Info */}
          <div style={{ display: 'flex', alignItems: 'center', marginBottom: 12 }}>
            <Avatar 
              size="small" 
              style={{ 
                backgroundColor: '#1890ff', 
                marginRight: 8,
                fontSize: 10
              }}
            >
              {getTeamInitials(service.team)}
            </Avatar>
            <div>
              <Text strong style={{ fontSize: 12 }}>{service.team}</Text>
              <br />
              <Text type="secondary" style={{ fontSize: 11 }}>{service.squad}</Text>
            </div>
          </div>

          {/* Tags */}
          <div style={{ marginTop: 12 }}>
            <Text type="secondary" style={{ fontSize: 11, marginRight: 8 }}>tags</Text>
            {service.tags?.slice(0, 3).map((tag, index) => (
              <Tag key={index} size="small" style={{ fontSize: 10, margin: '0 2px' }}>
                {tag}
              </Tag>
            ))}
            {service.tags?.length > 3 && (
              <Tag size="small" style={{ fontSize: 10 }}>
                +{service.tags.length - 3}
              </Tag>
            )}
          </div>
        </div>
      </div>
    </Card>
  );
}

ServiceCard.propTypes = {
  service: PropTypes.shape({
    id: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired,
    displayName: PropTypes.string,
    description: PropTypes.string,
    tier: PropTypes.string,
    team: PropTypes.string,
    squad: PropTypes.string,
    tags: PropTypes.arrayOf(PropTypes.string),
    regions: PropTypes.arrayOf(PropTypes.string),
    ingressType: PropTypes.string,
    status: PropTypes.string,
    updatedAt: PropTypes.string,
  }).isRequired,
  onEdit: PropTypes.func,
  onDelete: PropTypes.func,
  onView: PropTypes.func,
};

export default ServiceCard;
