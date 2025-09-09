# AWS Plugin

> **⚠️ Alpha Plugin** - Enhanced EC2 operations with modern UI. Improved performance but still not production-ready.

The AWS plugin provides a simplified interface for common AWS operations, reducing cognitive load for developers who need basic cloud resource management without deep AWS expertise.

## 🎯 Features

### **Current Capabilities (Alpha)**

- **EC2 Instance Management** - List, start, and stop EC2 instances with modern interface
- **Multi-account Support** - Connect multiple AWS accounts with unified sidebar selector
- **Optimized Performance** - Intelligent caching prevents redundant API calls
- **Single Account Selector** - Unified account selection similar to kubectl context switching
- **Permission Controls** - Team-based instance operation permissions

### **Recent Updates (v0.2.0-alpha)**

**Sidebar & Navigation Improvements:**

- ✅ **Single Menu Item**: Unified "AWS" menu with account selector dropdown
- ✅ **Account Switching**: Modern UI for switching between multiple AWS accounts
- ✅ **Optimized API Calls**: Intelligent caching prevents redundant account list requests
- ✅ **Performance**: Significant reduction in initial load and navigation times

**Infrastructure Optimization:**

- ✅ **Shared Account Cache**: Similar pattern to Kubernetes cluster caching
- ✅ **Hook Dependencies**: Fixed React Hook dependency arrays for better performance
- ✅ **Loading States**: Improved loading indicators during account switches

### **Planned Features**

- **Auto Scaling Groups** - Scale group management
- **Load Balancers** - ELB/ALB monitoring and configuration
- **RDS Management** - Database instance operations
- **Cost Analytics** - Real-time spending and optimization
- **CloudWatch Integration** - Metrics and monitoring dashboards

## 🔧 Configuration

### **1. AWS Credentials Setup**

#### **Method 1: Environment Variables** (Recommended)

```bash
export AWS_ACCESS_KEY_ID="your-aws-access-key"
export AWS_SECRET_ACCESS_KEY="your-aws-secret-key"
```

#### **Method 2: AWS CLI Configuration**

```bash
aws configure
# Follow prompts to set credentials
```

### **2. DashOPS Configuration**

Add AWS configuration to your `dash-ops.yaml`:

```yaml
# Enable AWS plugin
plugins:
  - 'AWS'

# AWS configuration
aws:
  - name: 'Production Account'
    region: us-east-1
    accessKeyId: ${AWS_ACCESS_KEY_ID}
    secretAccessKey: ${AWS_SECRET_ACCESS_KEY}
    ec2Config:
      skipList:
        - 'EKSWorkerAutoScalingGroupSpot' # Hide specific instances
        - 'DatabaseReplica'

  - name: 'Development Account'
    region: us-west-2
    accessKeyId: ${AWS_DEV_ACCESS_KEY_ID}
    secretAccessKey: ${AWS_DEV_SECRET_ACCESS_KEY}
```

### **3. Permission Configuration**

Control which teams can perform operations:

```yaml
aws:
  - name: 'Development Account'
    region: us-east-1
    accessKeyId: ${AWS_ACCESS_KEY_ID}
    secretAccessKey: ${AWS_SECRET_ACCESS_KEY}
    permission:
      ec2:
        start: ['your-org*developers', 'your-org*sre'] # Teams that can start instances
        stop: ['your-org*developers', 'your-org*sre'] # Teams that can stop instances
```

## 🖥️ EC2 Management

### **Instance Operations**

#### **Viewing Instances**

- **Instance List** - All EC2 instances with current status
- **Real-time Status** - Live instance state updates
- **Instance Details** - Type, region, launch time, tags
- **Filtering** - Filter by state, type, or tags

#### **Instance Actions** (Development Only)

> **⚠️ Warning**: Instance start/stop operations should only be used in development environments.

- **Start Instance** - Power on stopped instances
- **Stop Instance** - Gracefully shut down running instances
- **Instance Logs** - CloudWatch log integration (planned)

### **Instance Filtering**

Hide specific instances from the interface:

```yaml
aws:
  - name: 'Production Account'
    ec2Config:
      skipList:
        - 'EKSWorkerAutoScalingGroupSpot' # Auto-scaling worker nodes
        - 'DatabasePrimary' # Critical database instances
        - 'LoadBalancer' # Infrastructure components
```

### **Multi-Account Support**

Manage multiple AWS accounts simultaneously:

```yaml
aws:
  - name: 'Production US-East'
    region: us-east-1
    accessKeyId: ${AWS_PROD_EAST_KEY_ID}
    secretAccessKey: ${AWS_PROD_EAST_SECRET}

  - name: 'Production EU-West'
    region: eu-west-1
    accessKeyId: ${AWS_PROD_EU_KEY_ID}
    secretAccessKey: ${AWS_PROD_EU_SECRET}

  - name: 'Development'
    region: us-west-2
    accessKeyId: ${AWS_DEV_KEY_ID}
    secretAccessKey: ${AWS_DEV_SECRET}
```

## 🔐 Security & Permissions

### **IAM Requirements**

#### **Minimum Required Permissions**

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "ec2:DescribeInstances",
        "ec2:DescribeInstanceStatus",
        "ec2:DescribeInstanceTypes"
      ],
      "Resource": "*"
    }
  ]
}
```

#### **Additional Permissions for Instance Management**

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["ec2:StartInstances", "ec2:StopInstances"],
      "Resource": "arn:aws:ec2:*:*:instance/*",
      "Condition": {
        "StringEquals": {
          "ec2:InstanceTag/Environment": ["development", "staging"]
        }
      }
    }
  ]
}
```

### **Team-based Access Control**

```yaml
aws:
  - name: 'Development Account'
    permission:
      ec2:
        start: ['dash-ops*developers'] # Only developers can start
        stop: ['dash-ops*developers', 'dash-ops*sre'] # Developers and SRE can stop
```

## 🚨 Beta Limitations

### **Current Restrictions**

❌ **Not Production Ready**

- **Limited AWS services** - Only EC2 supported
- **Basic permissions** - Simple team-based access only
- **No cost controls** - No spend monitoring or limits
- **Missing audit trail** - Limited operation logging
- **Credential exposure** - Keys stored in configuration files

### **Security Concerns**

- **Plain text credentials** - AWS keys in YAML configuration
- **No credential rotation** - Manual key management required
- **Limited RBAC** - Basic team permission model only
- **No MFA enforcement** - Single-factor authentication
- **Missing rate limiting** - API calls not rate-limited

## 🛣️ Roadmap

### **Short-term (Q1 2025)**

- **Enhanced EC2 management** - Tagging, security groups
- **Auto Scaling Groups** - ASG monitoring and scaling
- **CloudWatch integration** - Metrics and alarms
- **Improved permissions** - Resource-level access control

### **Medium-term (Q2-Q3 2025)**

- **Additional AWS services** - RDS, ECS, Lambda support
- **Cost management** - Spending alerts and budget controls
- **Resource discovery** - Automatic resource inventory
- **Compliance controls** - Governance and policy enforcement

### **Long-term (Q4 2025+)**

- **FinOps integration** - Complete cost optimization platform
- **Advanced automation** - Infrastructure lifecycle management
- **Multi-cloud support** - Unified interface across cloud providers

## 📊 API Endpoints

### **Account Operations**

```
GET /api/aws/accounts
```

**Response:**

```json
{
  "data": [
    {
      "name": "Production Account",
      "region": "us-east-1",
      "accountId": "123456789012"
    }
  ]
}
```

### **EC2 Operations**

```
GET /api/aws/instances?account={name}
```

**Response:**

```json
{
  "data": [
    {
      "instanceId": "i-1234567890abcdef0",
      "name": "web-server-1",
      "state": "running",
      "instanceType": "t3.micro",
      "launchTime": "2024-01-15T10:30:00Z",
      "privateIpAddress": "10.0.1.100",
      "publicIpAddress": "54.123.45.67"
    }
  ]
}
```

```
POST /api/aws/instances/{instanceId}/start
POST /api/aws/instances/{instanceId}/stop
```

**Response:**

```json
{
  "data": {
    "current_state": "pending"
  },
  "success": true
}
```

## 🐛 Troubleshooting

### **Common Issues**

#### **"Access Denied" Errors**

- ✅ Verify AWS IAM permissions
- ✅ Check environment variables are set correctly
- ✅ Ensure AWS credentials are valid and not expired

#### **Instances Not Appearing**

- ✅ Check if instances are in the `skipList`
- ✅ Verify correct AWS region configuration
- ✅ Confirm instances exist in the specified account

#### **Permission Errors**

- ✅ Verify GitHub team membership
- ✅ Check team naming format: `org*team`
- ✅ Ensure user is authenticated via Auth

### **Debug Mode**

Enable AWS plugin debugging:

```bash
# Backend debug logs
AWS_DEBUG=true go run main.go

# Frontend debug
localStorage.setItem('dash-ops:aws-debug', 'true')
```

## 🧪 Testing

### **Safe Testing Practices**

> **⚠️ Important**: Always test in non-production environments first.

1. **Use development AWS accounts** - Never test on production resources
2. **Limited instance types** - Test only with t2.micro/t3.micro instances
3. **Tag-based restrictions** - Use conditional IAM policies
4. **Monitoring** - Watch CloudTrail for all operations

### **Test Configuration Example**

```yaml
aws:
  - name: 'Development Testing'
    region: us-west-2
    accessKeyId: ${AWS_DEV_ACCESS_KEY_ID}
    secretAccessKey: ${AWS_DEV_SECRET_ACCESS_KEY}
    ec2Config:
      skipList:
        - 'production' # Hide any instance with 'production' in name
        - 'database' # Hide database instances
    permission:
      ec2:
        start: ['dash-ops*developers']
        stop: ['dash-ops*developers']
```

## 🤝 Contributing

### **Priority Areas**

1. **🔒 Security** - Implement secure credential management
2. **🧪 Testing** - Add comprehensive AWS integration tests
3. **📊 Monitoring** - CloudWatch and cost tracking integration
4. **🔌 Services** - Support for additional AWS services
5. **📖 Documentation** - Usage guides and best practices

### **Development Setup**

```bash
# 1. Set up AWS credentials for testing
export AWS_ACCESS_KEY_ID="test-key"
export AWS_SECRET_ACCESS_KEY="test-secret"

# 2. Run backend with debug logging
AWS_DEBUG=true go run main.go

# 3. Test API endpoints
curl http://localhost:8080/api/aws/accounts
```

## 📚 Resources

- **[AWS SDK for Go Documentation](https://docs.aws.amazon.com/sdk-for-go/)**
- **[AWS EC2 API Reference](https://docs.aws.amazon.com/ec2/latest/api/)**
- **[AWS IAM Best Practices](https://docs.aws.amazon.com/iam/latest/userguide/best-practices.html)**

---

**⚠️ Beta Notice**: This plugin is experimental and intended for development use only. Do not use in production environments.
