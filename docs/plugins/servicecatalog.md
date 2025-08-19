# Service Catalog Plugin

The Service Catalog plugin provides a centralized registry for managing and monitoring all services in your infrastructure. It offers a comprehensive view of service metadata, configurations, and deployment information.

## Features

- **Service Registry**: Centralized catalog of all services
- **Tier Management**: Organize services by criticality (Tier 1, 2, 3)
- **Team Ownership**: Track service ownership and responsibility
- **Search & Filtering**: Find services by name, tags, team, or tier
- **Service Lifecycle**: Manage service status and deployment information
- **Tag System**: Categorize services with custom tags

## Configuration

Add the Service Catalog plugin to your `dash-ops.yaml`:

```yaml
plugins:
  - 'ServiceCatalog'

serviceCatalog:
  - name: 'Service Catalog'
    storage: 'file' # Storage type: file, database (future)
    catalogPath: './catalog/services' # Path for file storage
    permission:
      read: ['dash-ops*dev-team', 'dash-ops*platform-team']
      write: ['dash-ops*platform-team']
      admin: ['dash-ops*platform-team']
```

### Configuration Options

| Option             | Type   | Required | Description                                                     |
| ------------------ | ------ | -------- | --------------------------------------------------------------- |
| `name`             | string | Yes      | Display name for the catalog                                    |
| `storage`          | string | Yes      | Storage backend (`file` or `database`)                          |
| `catalogPath`      | string | No       | Directory path for file storage (default: `./catalog/services`) |
| `permission.read`  | array  | No       | Teams with read access                                          |
| `permission.write` | array  | No       | Teams with write access                                         |
| `permission.admin` | array  | No       | Teams with admin access                                         |

## Service Schema

Each service in the catalog follows this structure:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "name": "user-authentication-service",
  "displayName": "User Authentication Service",
  "description": "Service responsible for user authentication and authorization",
  "tier": "tier-1",
  "team": "Platform Team",
  "squad": "Auth Squad",
  "owner": "platform-team@company.com",
  "tags": ["authentication", "api", "microservice"],
  "regions": ["us-east-1", "eu-west-1"],
  "ingressType": "external",
  "status": "active",
  "createdAt": "2024-01-15T10:30:00Z",
  "updatedAt": "2024-01-20T14:45:00Z",
  "metadata": {}
}
```

### Field Descriptions

| Field         | Type   | Required | Description                                        |
| ------------- | ------ | -------- | -------------------------------------------------- |
| `id`          | string | Yes      | Unique service identifier (UUID)                   |
| `name`        | string | Yes      | Service name (kebab-case recommended)              |
| `displayName` | string | No       | Human-readable service name                        |
| `description` | string | Yes      | Service description                                |
| `tier`        | string | Yes      | Service tier: `tier-1`, `tier-2`, `tier-3`         |
| `team`        | string | Yes      | Owning team name                                   |
| `squad`       | string | Yes      | Responsible squad within the team                  |
| `owner`       | string | No       | Contact email for the service                      |
| `tags`        | array  | No       | Service tags for categorization                    |
| `regions`     | array  | No       | Deployment regions                                 |
| `ingressType` | string | No       | Access type: `internal` or `external`              |
| `status`      | string | No       | Service status: `active`, `inactive`, `deprecated` |
| `createdAt`   | string | Auto     | Service creation timestamp                         |
| `updatedAt`   | string | Auto     | Last update timestamp                              |
| `metadata`    | object | No       | Additional custom metadata                         |

## Service Tiers

Services are organized into three tiers based on criticality:

### Tier 1 - Critical

- **Color**: Red
- **Description**: Mission-critical services that directly impact business operations
- **Examples**: Payment processing, user authentication, core APIs
- **SLA**: 99.99% uptime, < 1 minute recovery time

### Tier 2 - Important

- **Color**: Orange
- **Description**: Important services that support business functions
- **Examples**: Notification services, reporting APIs, admin dashboards
- **SLA**: 99.9% uptime, < 5 minutes recovery time

### Tier 3 - Standard

- **Color**: Green
- **Description**: Supporting services and internal tools
- **Examples**: Development tools, internal dashboards, batch jobs
- **SLA**: 99% uptime, < 15 minutes recovery time

## API Endpoints

The Service Catalog plugin exposes the following REST API endpoints:

### List Services

```http
GET /api/v1/servicecatalog/services
```

Query parameters:

- `tier`: Filter by tier (`tier-1`, `tier-2`, `tier-3`)
- `team`: Filter by team name
- `status`: Filter by status (`active`, `inactive`, `deprecated`)
- `search`: Search in name, description, and tags

### Get Service

```http
GET /api/v1/servicecatalog/services/{id}
```

### Create Service

```http
POST /api/v1/servicecatalog/services
Content-Type: application/json

{
  "name": "my-service",
  "description": "Service description",
  "tier": "tier-2",
  "team": "Platform Team",
  "squad": "Backend Squad",
  "tags": ["api", "microservice"]
}
```

### Update Service

```http
PUT /api/v1/servicecatalog/services/{id}
Content-Type: application/json

{
  "description": "Updated description",
  "tier": "tier-1"
}
```

### Delete Service

```http
DELETE /api/v1/servicecatalog/services/{id}
```

### Get Statistics

```http
GET /api/v1/servicecatalog/stats
```

Returns:

```json
{
  "total": 15,
  "by_tier": {
    "tier-1": 3,
    "tier-2": 7,
    "tier-3": 5
  },
  "by_status": {
    "active": 12,
    "inactive": 2,
    "deprecated": 1
  }
}
```

## Usage Examples

### Creating a Service via API

```bash
curl -X POST http://localhost:8080/api/v1/servicecatalog/services \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "payment-processor",
    "displayName": "Payment Processor",
    "description": "Handles payment transactions and billing",
    "tier": "tier-1",
    "team": "Payments Team",
    "squad": "Core Payments",
    "tags": ["payments", "critical", "pci-compliant"]
  }'
```

### Searching Services

```bash
# Search by name or description
curl "http://localhost:8080/api/v1/servicecatalog/services?search=payment"

# Filter by tier
curl "http://localhost:8080/api/v1/servicecatalog/services?tier=tier-1"

# Filter by team
curl "http://localhost:8080/api/v1/servicecatalog/services?team=Platform%20Team"

# Combined filters
curl "http://localhost:8080/api/v1/servicecatalog/services?tier=tier-1&search=api"
```

## Best Practices

### Service Naming

- Use kebab-case for service names: `user-authentication-service`
- Keep names descriptive but concise
- Avoid abbreviations that might be unclear

### Tagging Strategy

- Use consistent tag naming conventions
- Common tags: `api`, `microservice`, `batch-job`, `frontend`, `backend`
- Technology tags: `nodejs`, `python`, `react`, `postgres`
- Feature tags: `authentication`, `payments`, `notifications`

### Team Organization

- Align team names with your organizational structure
- Use squad names to identify specific responsible groups
- Keep ownership information up to date

### Tier Assignment

- **Tier 1**: Services that cause immediate business impact if down
- **Tier 2**: Services that cause significant user experience degradation
- **Tier 3**: Services that have minimal immediate impact

## Integration with Other Plugins

### Kubernetes Integration

The Service Catalog can integrate with the Kubernetes plugin to:

- Sync deployment status
- Validate resource configurations
- Track service health metrics

### AWS Integration

Integration with AWS plugin provides:

- Cost tracking per service
- Resource utilization metrics
- Load balancer synchronization

## Roadmap

### Phase 1 (Current - MVP)

- ✅ Basic service CRUD operations
- ✅ Tier-based organization
- ✅ Search and filtering
- ✅ File-based storage

### Phase 2 (Next)

- 🔄 Service detail pages
- 🔄 Environment-specific configurations
- 🔄 Deployment history tracking
- 🔄 Integration with K8s/AWS plugins

### Phase 3 (Future)

- 🔄 Database storage backend
- 🔄 Service dependency mapping
- 🔄 Automated service discovery
- 🔄 Compliance and governance checks

### Phase 4 (Advanced)

- 🔄 CI/CD pipeline integration
- 🔄 Monitoring and alerting configuration
- 🔄 Cost optimization recommendations
- 🔄 Service mesh integration

## Troubleshooting

### Common Issues

**Services not appearing in the catalog**

- Check if the ServiceCatalog plugin is enabled in your configuration
- Verify the `catalogPath` directory exists and is writable
- Check backend logs for any storage errors

**Permission denied errors**

- Ensure your user is part of a team with appropriate permissions
- Check the `permission` configuration in your YAML file
- Verify OAuth2 integration is working correctly

**Search not working**

- Search is case-insensitive and searches across name, description, and tags
- Try using partial matches instead of exact terms
- Check if there are any special characters that might interfere

### Debug Mode

Enable debug logging by setting the log level in your configuration:

```yaml
logging:
  level: debug
```

This will provide detailed information about service operations and API calls.

## Contributing

To contribute to the Service Catalog plugin:

1. Follow the existing code structure in `pkg/servicecatalog/`
2. Add tests for new functionality
3. Update documentation for any new features
4. Ensure backward compatibility with existing service definitions

For more information, see the main [CONTRIBUTING.md](../CONTRIBUTING.md) guide.
