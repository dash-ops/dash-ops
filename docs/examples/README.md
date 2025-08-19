# Service Catalog Examples

This directory contains example configurations and service definitions for the Service Catalog plugin.

## Files Overview

### Configuration Examples

- **`service-catalog-config.yaml`** - Complete dash-ops.yaml configuration example with Service Catalog plugin enabled
- **`service-definitions/`** - Directory containing example service definitions for different scenarios

### Service Definition Examples

The `service-definitions/` directory contains examples of different types of services:

#### 1. Tier 1 - Critical Services

**File**: `tier1-critical-api.json`

- **Service**: Core User API
- **Characteristics**:
  - 99.99% uptime SLA
  - Multi-region deployment
  - External access
  - High security requirements
  - Comprehensive monitoring

#### 2. Tier 2 - Business Services

**File**: `tier2-business-service.json`

- **Service**: Order Processing Service
- **Characteristics**:
  - 99.9% uptime SLA
  - Internal service with business logic
  - Integration with multiple services
  - Business metrics tracking

#### 3. Tier 3 - Support Services

**File**: `tier3-support-service.json`

- **Service**: Admin Dashboard Backend
- **Characteristics**:
  - 99% uptime SLA
  - Internal administrative tool
  - Single region deployment
  - Admin-only access

#### 4. Batch Jobs

**File**: `batch-job-service.json`

- **Service**: Daily Report Generator
- **Characteristics**:
  - Scheduled execution
  - Resource-intensive processing
  - File output generation
  - Email notifications

## Using These Examples

### 1. Copy Service Definitions

To use these examples in your Service Catalog:

```bash
# Copy example services to your catalog directory
cp docs/examples/service-definitions/*.json catalog/services/

# Or copy individual services
cp docs/examples/service-definitions/tier1-critical-api.json catalog/services/
```

### 2. Customize for Your Environment

Edit the copied files to match your environment:

```json
{
  "name": "your-service-name",
  "team": "Your Team Name",
  "squad": "Your Squad Name",
  "owner": "your-team@yourcompany.com",
  "regions": ["your-region-1", "your-region-2"]
  // ... other customizations
}
```

### 3. Validate Service Definitions

Use the API to validate your service definitions:

```bash
# Test creating a service
curl -X POST http://localhost:8080/api/v1/servicecatalog/services \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d @docs/examples/service-definitions/tier1-critical-api.json
```

## Service Definition Best Practices

Based on these examples, follow these best practices:

### Naming Conventions

- **Service names**: Use kebab-case (`user-authentication-service`)
- **Display names**: Use Title Case (`User Authentication Service`)
- **Tags**: Use lowercase with hyphens (`batch-job`, `user-management`)

### Tier Assignment Guidelines

**Tier 1 - Critical**

- Services that cause immediate business impact if unavailable
- Customer-facing APIs and core business functions
- SLA: 99.99% uptime, <100ms response time
- Multi-region deployment recommended

**Tier 2 - Important**

- Services that support business operations
- Internal APIs and business logic services
- SLA: 99.9% uptime, <500ms response time
- At least 2 regions for redundancy

**Tier 3 - Standard**

- Supporting services and internal tools
- Batch jobs and administrative services
- SLA: 99% uptime, <2s response time
- Single region deployment acceptable

### Metadata Usage

Use the `metadata` field to store additional service information:

```json
{
  "metadata": {
    "sla_uptime": "99.99%",
    "max_response_time": "100ms",
    "security_level": "high",
    "compliance": ["SOC2", "GDPR"],
    "monitoring": {
      "health_check": "/health",
      "metrics_endpoint": "/metrics"
    },
    "scaling": {
      "min_replicas": 3,
      "max_replicas": 20
    }
  }
}
```

### Tag Strategy

Use consistent tags across your services:

**Technology Tags**

- `nodejs`, `python`, `java`, `go`
- `react`, `vue`, `angular`
- `postgres`, `mongodb`, `redis`

**Function Tags**

- `api`, `microservice`, `batch-job`
- `authentication`, `payments`, `notifications`
- `frontend`, `backend`, `database`

**Operational Tags**

- `critical`, `monitoring`, `logging`
- `high-availability`, `multi-region`
- `pci-compliant`, `gdpr-compliant`

## Testing Your Configuration

### 1. Start the Service

```bash
DASH_CONFIG=docs/examples/service-catalog-config.yaml go run main.go
```

### 2. Load Example Services

```bash
# Copy examples to your catalog directory
mkdir -p catalog/services
cp docs/examples/service-definitions/*.json catalog/services/
```

### 3. Access the UI

Open your browser and navigate to:

- Main catalog: http://localhost:5173/servicecatalog
- API endpoint: http://localhost:8080/api/v1/servicecatalog/services

### 4. Test API Endpoints

```bash
# List all services
curl http://localhost:8080/api/v1/servicecatalog/services

# Filter by tier
curl "http://localhost:8080/api/v1/servicecatalog/services?tier=tier-1"

# Search services
curl "http://localhost:8080/api/v1/servicecatalog/services?search=api"

# Get statistics
curl http://localhost:8080/api/v1/servicecatalog/stats
```

## Next Steps

1. **Customize Examples**: Modify the example services to match your actual services
2. **Add More Services**: Create additional service definitions for your infrastructure
3. **Integrate with CI/CD**: Automate service registration as part of your deployment pipeline
4. **Set Up Monitoring**: Configure alerts and monitoring for your service catalog
5. **Team Training**: Train your teams on using the service catalog for service discovery and management

For more detailed information, see the [Service Catalog Plugin Documentation](../plugins/servicecatalog.md).
