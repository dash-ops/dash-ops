import { useState } from 'react';
import { 
  Modal, 
  Steps, 
  Form, 
  Input, 
  Select, 
  Button, 
  Space, 
  Typography,
  Tag,
  message
} from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import PropTypes from 'prop-types';

const { TextArea } = Input;
const { Text } = Typography;

const steps = [
  {
    title: 'Basic Info',
    description: 'Service information',
  },
  {
    title: 'Infrastructure',
    description: 'Infrastructure settings',
  },
  {
    title: 'Networking',
    description: 'Network configuration',
  },
  {
    title: 'Features',
    description: 'Additional features',
  },
  {
    title: 'Review',
    description: 'Review and create',
  },
];

function CreateServiceModal({ visible, onCancel, onSubmit, loading }) {
  const [currentStep, setCurrentStep] = useState(0);
  const [form] = Form.useForm();
  const [tags, setTags] = useState([]);
  const [inputTag, setInputTag] = useState('');

  const handleNext = async () => {
    try {
      // Validate current step fields
      const stepFields = getStepFields(currentStep);
      await form.validateFields(stepFields);
      setCurrentStep(currentStep + 1);
    } catch (error) {
      console.log('Validation failed:', error);
    }
  };

  const handlePrev = () => {
    setCurrentStep(currentStep - 1);
  };

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields();
      const serviceData = {
        ...values,
        tags: tags,
      };
      await onSubmit(serviceData);
      handleReset();
    } catch (error) {
      console.log('Submit failed:', error);
    }
  };

  const handleReset = () => {
    setCurrentStep(0);
    form.resetFields();
    setTags([]);
    setInputTag('');
  };

  const handleCancel = () => {
    handleReset();
    onCancel();
  };

  const addTag = () => {
    if (inputTag && !tags.includes(inputTag)) {
      setTags([...tags, inputTag]);
      setInputTag('');
    }
  };

  const removeTag = (tagToRemove) => {
    setTags(tags.filter(tag => tag !== tagToRemove));
  };

  const getStepFields = (step) => {
    switch (step) {
      case 0: return ['name', 'description', 'tier'];
      case 1: return ['team', 'squad'];
      case 2: return []; // No required fields for networking in MVP
      case 3: return []; // No required fields for features in MVP
      case 4: return []; // Review step
      default: return [];
    }
  };

  const renderStepContent = () => {
    switch (currentStep) {
      case 0:
        return (
          <div>
            <Text type="secondary" style={{ display: 'block', marginBottom: 16 }}>
              Basic information about your service
            </Text>
            
            <Form.Item
              name="name"
              label="Service Name"
              rules={[{ required: true, message: 'Service name is required' }]}
            >
              <Input placeholder="e.g., user-authentication-service" />
            </Form.Item>

            <Form.Item
              name="tier"
              label="Service Tier"
              rules={[{ required: true, message: 'Service tier is required' }]}
            >
              <Select placeholder="Tier 2 - Important">
                <Select.Option value="tier-1">Tier 1 - Critical</Select.Option>
                <Select.Option value="tier-2">Tier 2 - Important</Select.Option>
                <Select.Option value="tier-3">Tier 3 - Standard</Select.Option>
              </Select>
            </Form.Item>

            <Form.Item
              name="description"
              label="Description"
              rules={[{ required: true, message: 'Description is required' }]}
            >
              <TextArea 
                rows={3} 
                placeholder="Describe what this service does..."
                showCount
                maxLength={200}
              />
            </Form.Item>
          </div>
        );

      case 1:
        return (
          <div>
            <Text type="secondary" style={{ display: 'block', marginBottom: 16 }}>
              Team and ownership information
            </Text>

            <Form.Item
              name="team"
              label="Team"
              rules={[{ required: true, message: 'Team is required' }]}
            >
              <Input placeholder="e.g., Platform Team" />
            </Form.Item>

            <Form.Item
              name="squad"
              label="Squad"
              rules={[{ required: true, message: 'Squad is required' }]}
            >
              <Input placeholder="e.g., Auth Squad" />
            </Form.Item>

            <div style={{ marginBottom: 16 }}>
              <Text strong>Tags</Text>
              <Text type="secondary" style={{ display: 'block', marginBottom: 8 }}>
                Add tags to help categorize and find your service
              </Text>
              
              <Space.Compact style={{ width: '100%', marginBottom: 8 }}>
                <Input
                  placeholder="Add a tag..."
                  value={inputTag}
                  onChange={(e) => setInputTag(e.target.value)}
                  onPressEnter={addTag}
                />
                <Button 
                  type="primary" 
                  icon={<PlusOutlined />} 
                  onClick={addTag}
                  disabled={!inputTag}
                />
              </Space.Compact>

              <div>
                {tags.map((tag, index) => (
                  <Tag
                    key={index}
                    closable
                    onClose={() => removeTag(tag)}
                    style={{ marginBottom: 4 }}
                  >
                    {tag}
                  </Tag>
                ))}
              </div>
            </div>
          </div>
        );

      case 2:
        return (
          <div>
            <Text type="secondary" style={{ display: 'block', marginBottom: 16 }}>
              Network and infrastructure settings (Coming soon)
            </Text>
            <div style={{ textAlign: 'center', padding: 40 }}>
              <Text type="secondary">
                Network configuration will be available in the next version.
                <br />
                For now, services will use default internal networking.
              </Text>
            </div>
          </div>
        );

      case 3:
        return (
          <div>
            <Text type="secondary" style={{ display: 'block', marginBottom: 16 }}>
              Additional features and integrations (Coming soon)
            </Text>
            <div style={{ textAlign: 'center', padding: 40 }}>
              <Text type="secondary">
                Features like logging, tracing, and monitoring configuration
                <br />
                will be available in future versions.
              </Text>
            </div>
          </div>
        );

      case 4:
        return (
          <div>
            <Text type="secondary" style={{ display: 'block', marginBottom: 16 }}>
              Review your service configuration before creating
            </Text>
            
            <div style={{ background: '#fafafa', padding: 16, borderRadius: 6 }}>
              <div style={{ marginBottom: 12 }}>
                <Text strong>Service Name:</Text> {form.getFieldValue('name')}
              </div>
              <div style={{ marginBottom: 12 }}>
                <Text strong>Tier:</Text> {form.getFieldValue('tier')}
              </div>
              <div style={{ marginBottom: 12 }}>
                <Text strong>Description:</Text> {form.getFieldValue('description')}
              </div>
              <div style={{ marginBottom: 12 }}>
                <Text strong>Team:</Text> {form.getFieldValue('team')}
              </div>
              <div style={{ marginBottom: 12 }}>
                <Text strong>Squad:</Text> {form.getFieldValue('squad')}
              </div>
              {tags.length > 0 && (
                <div>
                  <Text strong>Tags:</Text>{' '}
                  {tags.map((tag, index) => (
                    <Tag key={index} size="small" style={{ margin: '0 2px' }}>
                      {tag}
                    </Tag>
                  ))}
                </div>
              )}
            </div>
          </div>
        );

      default:
        return null;
    }
  };

  const renderFooter = () => {
    return (
      <div style={{ display: 'flex', justifyContent: 'space-between' }}>
        <Button onClick={handleCancel}>
          Cancel
        </Button>
        
        <Space>
          {currentStep > 0 && (
            <Button onClick={handlePrev}>
              Previous
            </Button>
          )}
          
          {currentStep < steps.length - 1 ? (
            <Button type="primary" onClick={handleNext}>
              Next
            </Button>
          ) : (
            <Button type="primary" onClick={handleSubmit} loading={loading}>
              Create Service
            </Button>
          )}
        </Space>
      </div>
    );
  };

  return (
    <Modal
      title="Create New Service"
      open={visible}
      onCancel={handleCancel}
      footer={renderFooter()}
      width={600}
      destroyOnClose
    >
      <div style={{ marginBottom: 24 }}>
        <Text type="secondary">
          Configure your new service with all necessary settings and deployments.
        </Text>
      </div>

      <Steps
        current={currentStep}
        items={steps}
        size="small"
        style={{ marginBottom: 24 }}
      />

      <Form
        form={form}
        layout="vertical"
        preserve={false}
      >
        {renderStepContent()}
      </Form>
    </Modal>
  );
}

CreateServiceModal.propTypes = {
  visible: PropTypes.bool.isRequired,
  onCancel: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
  loading: PropTypes.bool,
};

CreateServiceModal.defaultProps = {
  loading: false,
};

export default CreateServiceModal;
