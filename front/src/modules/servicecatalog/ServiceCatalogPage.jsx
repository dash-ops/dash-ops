import { useState, useEffect } from 'react';
import { 
  Row, 
  Col, 
  Button, 
  Input, 
  Space, 
  Typography, 
  Spin, 
  Empty, 
  message,
  Modal,
  Tag
} from 'antd';
import { PlusOutlined, SearchOutlined, ExclamationCircleOutlined } from '@ant-design/icons';
import ServiceCard from './components/ServiceCard';
import CreateServiceModal from './components/CreateServiceModal';
import { getServices, createService, deleteService } from './resources/serviceCatalogResource';

const { Title, Text } = Typography;
const { Search } = Input;

function ServiceCatalogPage() {
  const [services, setServices] = useState([]);
  const [loading, setLoading] = useState(true);
  const [createModalVisible, setCreateModalVisible] = useState(false);
  const [createLoading, setCreateLoading] = useState(false);
  const [filters, setFilters] = useState({
    tier: 'all',
    search: '',
  });

  const tierOptions = [
    { key: 'all', label: 'All Tiers', color: 'default' },
    { key: 'tier-1', label: 'Tier 1', color: 'red' },
    { key: 'tier-2', label: 'Tier 2', color: 'orange' },
    { key: 'tier-3', label: 'Tier 3', color: 'green' },
  ];

  useEffect(() => {
    loadServices();
  }, [filters]);

  const loadServices = async () => {
    try {
      setLoading(true);
      const response = await getServices(filters);
      setServices(response.data.services || []);
    } catch (error) {
      message.error('Failed to load services');
      console.error('Error loading services:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateService = async (serviceData) => {
    try {
      setCreateLoading(true);
      await createService(serviceData);
      message.success('Service created successfully!');
      setCreateModalVisible(false);
      loadServices();
    } catch (error) {
      message.error('Failed to create service');
      console.error('Error creating service:', error);
    } finally {
      setCreateLoading(false);
    }
  };

  const handleDeleteService = (service) => {
    Modal.confirm({
      title: 'Delete Service',
      icon: <ExclamationCircleOutlined />,
      content: (
        <div>
          <p>Are you sure you want to delete <strong>{service.displayName || service.name}</strong>?</p>
          <p style={{ color: '#ff4d4f', fontSize: 12 }}>
            This action cannot be undone and will permanently remove the service from the catalog.
          </p>
        </div>
      ),
      okText: 'Delete',
      okType: 'danger',
      cancelText: 'Cancel',
      onOk: async () => {
        try {
          await deleteService(service.id);
          message.success('Service deleted successfully');
          loadServices();
        } catch (error) {
          message.error('Failed to delete service');
          console.error('Error deleting service:', error);
        }
      },
    });
  };

  const handleViewService = (service) => {
    // TODO: Navigate to service detail page
    message.info(`View details for ${service.displayName || service.name} (Coming soon)`);
  };

  const handleEditService = (service) => {
    // TODO: Open edit modal
    message.info(`Edit ${service.displayName || service.name} (Coming soon)`);
  };

  const handleTierFilter = (tier) => {
    setFilters(prev => ({ ...prev, tier }));
  };

  const handleSearch = (value) => {
    setFilters(prev => ({ ...prev, search: value }));
  };

  const getServiceStats = () => {
    const stats = {
      total: services.length,
      'tier-1': services.filter(s => s.tier === 'tier-1').length,
      'tier-2': services.filter(s => s.tier === 'tier-2').length,
      'tier-3': services.filter(s => s.tier === 'tier-3').length,
    };
    return stats;
  };

  const stats = getServiceStats();

  return (
    <div style={{ padding: 24 }}>
      {/* Header */}
      <div style={{ marginBottom: 24 }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
          <div>
            <Title level={2} style={{ margin: 0 }}>
              Services Catalog
            </Title>
            <Text type="secondary">
              Gerencie e monitore todos os seus serviços em um só lugar
            </Text>
          </div>
          <Button 
            type="primary" 
            icon={<PlusOutlined />}
            onClick={() => setCreateModalVisible(true)}
          >
            Add Service
          </Button>
        </div>
      </div>

      {/* Filters and Search */}
      <div style={{ marginBottom: 24 }}>
        <Row gutter={[16, 16]} align="middle">
          <Col xs={24} sm={12} md={8}>
            <Search
              placeholder="Buscar serviços, tags, equipes..."
              allowClear
              enterButton={<SearchOutlined />}
              onSearch={handleSearch}
              style={{ width: '100%' }}
            />
          </Col>
          <Col xs={24} sm={12} md={16}>
            <Space wrap>
              {tierOptions.map(option => (
                <Tag.CheckableTag
                  key={option.key}
                  checked={filters.tier === option.key}
                  onChange={() => handleTierFilter(option.key)}
                  style={{ 
                    padding: '4px 12px',
                    borderRadius: 16,
                    fontSize: 12
                  }}
                >
                  {option.label}
                  {option.key !== 'all' && stats[option.key] > 0 && (
                    <span style={{ marginLeft: 4, opacity: 0.7 }}>
                      ({stats[option.key]})
                    </span>
                  )}
                  {option.key === 'all' && (
                    <span style={{ marginLeft: 4, opacity: 0.7 }}>
                      ({stats.total})
                    </span>
                  )}
                </Tag.CheckableTag>
              ))}
            </Space>
          </Col>
        </Row>
      </div>

      {/* Services Grid */}
      <Spin spinning={loading}>
        {services.length === 0 && !loading ? (
          <Empty
            description={
              filters.search || filters.tier !== 'all' 
                ? "No services found matching your filters"
                : "No services in catalog yet"
            }
            style={{ padding: 60 }}
          >
            {(!filters.search && filters.tier === 'all') && (
              <Button 
                type="primary" 
                icon={<PlusOutlined />}
                onClick={() => setCreateModalVisible(true)}
              >
                Create Your First Service
              </Button>
            )}
          </Empty>
        ) : (
          <Row gutter={[16, 16]}>
            {services.map(service => (
              <Col xs={24} sm={12} lg={8} xl={6} key={service.id}>
                <ServiceCard
                  service={service}
                  onView={handleViewService}
                  onEdit={handleEditService}
                  onDelete={handleDeleteService}
                />
              </Col>
            ))}
          </Row>
        )}
      </Spin>

      {/* Create Service Modal */}
      <CreateServiceModal
        visible={createModalVisible}
        onCancel={() => setCreateModalVisible(false)}
        onSubmit={handleCreateService}
        loading={createLoading}
      />
    </div>
  );
}

export default ServiceCatalogPage;
