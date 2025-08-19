# 🔌 Plugin Development Guide

## Overview

This guide provides comprehensive instructions for developing plugins for the dash-ops platform, including the universal adapter pattern and AI integration capabilities.

## Plugin Architecture

### Core Plugin Interface

```go
package plugin

import (
    "context"
    "time"
)

// Plugin defines the core interface that all plugins must implement
type Plugin interface {
    // Metadata
    Name() string
    Version() string
    Description() string
    Author() string

    // Lifecycle
    Install(config Config) error
    Uninstall() error
    Start(ctx context.Context) error
    Stop() error

    // Runtime
    GetRoutes() []Route
    GetPermissions() []Permission
    GetHealthCheck() HealthCheck

    // AI Integration
    GetAIContext(serviceID string) (map[string]interface{}, error)
    GetQuerySchema() QuerySchema
}

// Provider defines the universal adapter interface
type Provider interface {
    Name() string
    Type() string // "logs", "metrics", "traces", "cloud", etc.

    // Core functionality
    Query(ctx context.Context, query interface{}) (interface{}, error)
    Stream(ctx context.Context, query interface{}) (<-chan interface{}, error)

    // Metadata
    GetSchema() QuerySchema
    GetCapabilities() []Capability

    // Health and monitoring
    HealthCheck() error
    GetMetrics() map[string]interface{}
}
```

### Plugin Types

#### 1. Observability Plugins

Handle logs, metrics, traces, and APM data from various sources.

```go
type ObservabilityPlugin struct {
    BasePlugin
    providers map[string]Provider
    defaultProvider string
}

// Example: Logs Plugin
type LogsPlugin struct {
    ObservabilityPlugin
}

func (lp *LogsPlugin) AddProvider(name string, provider LogsProvider) {
    lp.providers[name] = provider
}

func (lp *LogsPlugin) Query(providerName string, query LogQuery) (*LogResult, error) {
    provider, exists := lp.providers[providerName]
    if !exists {
        provider = lp.providers[lp.defaultProvider]
    }

    return provider.Query(context.Background(), query)
}
```

#### 2. Infrastructure Plugins

Manage cloud resources, Kubernetes clusters, and infrastructure components.

```go
type InfrastructurePlugin struct {
    BasePlugin
    clients map[string]InfraClient
}

// Example: AWS Plugin
type AWSPlugin struct {
    InfrastructurePlugin
    accounts map[string]*AWSAccount
}

func (ap *AWSPlugin) ListInstances(accountID, region string) ([]EC2Instance, error) {
    account := ap.accounts[accountID]
    return account.EC2.ListInstances(region)
}
```

#### 3. Integration Plugins

Connect with external services like GitHub, GitLab, Jira, etc.

```go
type IntegrationPlugin struct {
    BasePlugin
    apiClients map[string]APIClient
}

// Example: GitHub Plugin
type GitHubPlugin struct {
    IntegrationPlugin
    orgs map[string]*GitHubOrg
}
```

## Universal Adapter Pattern

### Implementing Adapters

#### Logs Adapter Example

```go
// Loki Adapter
type LokiAdapter struct {
    client   *loki.Client
    config   LokiConfig
    endpoint string
}

func NewLokiAdapter(config LokiConfig) *LokiAdapter {
    return &LokiAdapter{
        client:   loki.NewClient(config.Endpoint),
        config:   config,
        endpoint: config.Endpoint,
    }
}

func (la *LokiAdapter) Name() string { return "loki" }
func (la *LokiAdapter) Type() string { return "logs" }

func (la *LokiAdapter) Query(ctx context.Context, query interface{}) (interface{}, error) {
    logQuery, ok := query.(LogQuery)
    if !ok {
        return nil, errors.New("invalid query type for Loki adapter")
    }

    // Convert universal query to Loki-specific query
    lokiQuery := la.convertQuery(logQuery)

    // Execute query
    result, err := la.client.Query(ctx, lokiQuery)
    if err != nil {
        return nil, fmt.Errorf("loki query failed: %w", err)
    }

    // Convert result to universal format
    return la.convertResult(result), nil
}

func (la *LokiAdapter) GetSchema() QuerySchema {
    return QuerySchema{
        Fields: []Field{
            {Name: "query", Type: "string", Required: true, Description: "LogQL query string"},
            {Name: "start", Type: "timestamp", Required: false, Description: "Start time"},
            {Name: "end", Type: "timestamp", Required: false, Description: "End time"},
            {Name: "limit", Type: "integer", Required: false, Description: "Maximum number of entries"},
        },
        Examples: []QueryExample{
            {
                Name: "Error logs",
                Query: map[string]interface{}{
                    "query": `{service="user-api"} |= "ERROR"`,
                    "limit": 100,
                },
            },
        },
    }
}

// Elasticsearch Adapter
type ElasticsearchAdapter struct {
    client *elasticsearch.Client
    config ESConfig
}

func (ea *ElasticsearchAdapter) Query(ctx context.Context, query interface{}) (interface{}, error) {
    logQuery := query.(LogQuery)

    // Convert to Elasticsearch DSL
    esQuery := map[string]interface{}{
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "must": []map[string]interface{}{
                    {"match": map[string]interface{}{"service": logQuery.Service}},
                    {"range": map[string]interface{}{
                        "@timestamp": map[string]interface{}{
                            "gte": logQuery.StartTime,
                            "lte": logQuery.EndTime,
                        },
                    }},
                },
            },
        },
        "size": logQuery.Limit,
        "sort": []map[string]interface{}{
            {"@timestamp": map[string]interface{}{"order": "desc"}},
        },
    }

    // Execute query
    result, err := ea.client.Search(
        ea.client.Search.WithContext(ctx),
        ea.client.Search.WithIndex(ea.config.IndexPattern),
        ea.client.Search.WithBody(strings.NewReader(jsonEncode(esQuery))),
    )

    if err != nil {
        return nil, fmt.Errorf("elasticsearch query failed: %w", err)
    }

    return ea.convertResult(result), nil
}
```

### Universal Query Format

```go
// Universal query structures
type LogQuery struct {
    Service   string            `json:"service"`
    Level     string            `json:"level,omitempty"`
    Message   string            `json:"message,omitempty"`
    StartTime time.Time         `json:"start_time"`
    EndTime   time.Time         `json:"end_time"`
    Limit     int               `json:"limit,omitempty"`
    Labels    map[string]string `json:"labels,omitempty"`
}

type MetricQuery struct {
    Metric    string            `json:"metric"`
    Labels    map[string]string `json:"labels,omitempty"`
    StartTime time.Time         `json:"start_time"`
    EndTime   time.Time         `json:"end_time"`
    Step      time.Duration     `json:"step,omitempty"`
}

type TraceQuery struct {
    Service   string            `json:"service"`
    Operation string            `json:"operation,omitempty"`
    TraceID   string            `json:"trace_id,omitempty"`
    StartTime time.Time         `json:"start_time"`
    EndTime   time.Time         `json:"end_time"`
    Limit     int               `json:"limit,omitempty"`
}
```

## Plugin Configuration

### Configuration Schema

```yaml
# plugin-config.yaml
plugin:
  name: 'logs-plugin'
  version: '1.0.0'
  type: 'observability'

  providers:
    loki:
      type: 'grafana-loki'
      endpoint: 'http://loki:3100'
      auth:
        type: 'basic'
        username: '${LOKI_USERNAME}'
        password: '${LOKI_PASSWORD}'

    elasticsearch:
      type: 'elasticsearch'
      endpoint: 'https://es.company.com:9200'
      index_pattern: 'logs-*'
      auth:
        type: 'api_key'
        api_key: '${ES_API_KEY}'

    cloudwatch:
      type: 'aws-cloudwatch'
      region: 'us-east-1'
      log_groups: ['app-logs', 'infra-logs']
      auth:
        type: 'iam_role'
        role_arn: 'arn:aws:iam::123456789012:role/LogsAccess'

  settings:
    default_provider: 'loki'
    cache_ttl: '5m'
    max_query_size: 10000

  permissions:
    - 'logs:read'
    - 'logs:query'
```

### Environment-Specific Configuration

```yaml
# environments/production/plugins/logs.yaml
plugin_config:
  logs-plugin:
    providers:
      loki:
        endpoint: 'https://loki.prod.company.com'
        auth:
          type: 'oauth2'
          client_id: '${LOKI_CLIENT_ID}'
          client_secret: '${LOKI_CLIENT_SECRET}'

      splunk:
        type: 'splunk'
        endpoint: 'https://splunk.company.com:8089'
        auth:
          type: 'token'
          token: '${SPLUNK_TOKEN}'

    settings:
      default_provider: 'splunk'
      retention_days: 90
```

## AI Integration

### Providing Context for AI

```go
func (lp *LogsPlugin) GetAIContext(serviceID string) (map[string]interface{}, error) {
    // Gather recent logs for the service
    recentLogs, err := lp.Query("", LogQuery{
        Service:   serviceID,
        StartTime: time.Now().Add(-1 * time.Hour),
        EndTime:   time.Now(),
        Limit:     100,
    })
    if err != nil {
        return nil, err
    }

    // Analyze log patterns
    errorCount := lp.countErrorLogs(recentLogs)
    commonErrors := lp.extractCommonErrors(recentLogs)
    logVolume := lp.calculateLogVolume(recentLogs)

    return map[string]interface{}{
        "recent_logs": recentLogs,
        "error_count": errorCount,
        "common_errors": commonErrors,
        "log_volume": logVolume,
        "log_level_distribution": lp.getLogLevelDistribution(recentLogs),
    }, nil
}

func (lp *LogsPlugin) GetQuerySchema() QuerySchema {
    return QuerySchema{
        Description: "Query logs from various providers",
        Fields: []Field{
            {
                Name: "service",
                Type: "string",
                Required: true,
                Description: "Service name to query logs for",
            },
            {
                Name: "level",
                Type: "enum",
                Options: []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"},
                Description: "Log level filter",
            },
            {
                Name: "time_range",
                Type: "duration",
                Description: "Time range for log query (e.g., '1h', '30m')",
            },
        },
        AIPrompts: []AIPrompt{
            {
                Pattern: "show me errors for {service}",
                Query: map[string]interface{}{
                    "service": "{service}",
                    "level": "ERROR",
                    "time_range": "1h",
                },
            },
            {
                Pattern: "what happened to {service} in the last {duration}",
                Query: map[string]interface{}{
                    "service": "{service}",
                    "time_range": "{duration}",
                },
            },
        },
    }
}
```

## Plugin Development Workflow

### 1. Setup Development Environment

```bash
# Clone plugin template
git clone https://github.com/dash-ops/plugin-template my-plugin
cd my-plugin

# Install dependencies
go mod init github.com/company/dash-ops-my-plugin
go mod tidy

# Setup development tools
make setup-dev
```

### 2. Implement Plugin Interface

```go
package main

import (
    "context"
    "github.com/dash-ops/dash-ops/pkg/plugin"
)

type MyPlugin struct {
    plugin.BasePlugin
    config MyPluginConfig
}

func (mp *MyPlugin) Name() string { return "my-plugin" }
func (mp *MyPlugin) Version() string { return "1.0.0" }
func (mp *MyPlugin) Description() string { return "My custom plugin" }

func (mp *MyPlugin) Install(config plugin.Config) error {
    // Parse configuration
    mp.config = parseConfig(config)

    // Initialize resources
    return mp.initialize()
}

func (mp *MyPlugin) Start(ctx context.Context) error {
    // Start background processes
    go mp.backgroundWorker(ctx)
    return nil
}

func (mp *MyPlugin) GetRoutes() []plugin.Route {
    return []plugin.Route{
        {
            Method: "GET",
            Path: "/api/my-plugin/status",
            Handler: mp.handleStatus,
        },
        {
            Method: "POST",
            Path: "/api/my-plugin/action",
            Handler: mp.handleAction,
        },
    }
}

// Register plugin
func init() {
    plugin.Register(&MyPlugin{})
}
```

### 3. Testing

```go
package main

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/dash-ops/dash-ops/pkg/plugin/testing"
)

func TestMyPlugin(t *testing.T) {
    // Create test plugin instance
    p := &MyPlugin{}

    // Test configuration
    config := plugin.Config{
        "endpoint": "http://test:8080",
        "timeout": "30s",
    }

    err := p.Install(config)
    assert.NoError(t, err)

    // Test functionality
    result, err := p.Query(context.Background(), TestQuery{})
    assert.NoError(t, err)
    assert.NotNil(t, result)
}

func TestMyPluginIntegration(t *testing.T) {
    // Use plugin testing framework
    suite := testing.NewPluginTestSuite(&MyPlugin{})

    suite.TestInstallation()
    suite.TestRoutes()
    suite.TestPermissions()
    suite.TestAIIntegration()
}
```

### 4. Documentation

````markdown
# My Plugin

## Overview

Brief description of what the plugin does.

## Configuration

```yaml
plugin_config:
  my-plugin:
    endpoint: 'https://api.service.com'
    timeout: '30s'
    auth:
      type: 'api_key'
      key: '${API_KEY}'
```
````

## API Endpoints

- `GET /api/my-plugin/status` - Get plugin status
- `POST /api/my-plugin/action` - Perform action

## AI Integration

The plugin provides context for AI analysis including:

- Service health metrics
- Recent events
- Configuration status

````

### 5. Publishing

```bash
# Build plugin
make build

# Run tests
make test

# Package plugin
make package

# Publish to registry
dash-ops plugin publish my-plugin-1.0.0.tar.gz
````

## Best Practices

### Performance

1. **Efficient Queries**: Implement query optimization and caching
2. **Connection Pooling**: Reuse connections to external services
3. **Async Operations**: Use goroutines for non-blocking operations
4. **Resource Management**: Properly cleanup resources

```go
type PluginWithCache struct {
    BasePlugin
    cache *cache.Cache
    pool  *connectionPool
}

func (p *PluginWithCache) Query(ctx context.Context, query interface{}) (interface{}, error) {
    // Check cache first
    if cached := p.cache.Get(query); cached != nil {
        return cached, nil
    }

    // Execute query with connection pooling
    conn := p.pool.Get()
    defer p.pool.Put(conn)

    result, err := conn.Query(ctx, query)
    if err != nil {
        return nil, err
    }

    // Cache result
    p.cache.Set(query, result, 5*time.Minute)

    return result, nil
}
```

### Security

1. **Input Validation**: Validate all inputs
2. **Authentication**: Secure API endpoints
3. **Secrets Management**: Use secure secret storage
4. **Audit Logging**: Log all significant actions

```go
func (p *MyPlugin) validateQuery(query interface{}) error {
    // Implement validation logic
    if query == nil {
        return errors.New("query cannot be nil")
    }

    // Validate query structure
    return p.validator.Validate(query)
}

func (p *MyPlugin) auditLog(action string, user string, details map[string]interface{}) {
    p.logger.Info("plugin action",
        "plugin", p.Name(),
        "action", action,
        "user", user,
        "details", details,
    )
}
```

### Error Handling

1. **Graceful Degradation**: Handle failures gracefully
2. **Retry Logic**: Implement exponential backoff
3. **Circuit Breakers**: Prevent cascade failures
4. **Meaningful Errors**: Provide actionable error messages

```go
func (p *MyPlugin) QueryWithRetry(ctx context.Context, query interface{}) (interface{}, error) {
    var lastErr error

    for attempt := 0; attempt < 3; attempt++ {
        result, err := p.Query(ctx, query)
        if err == nil {
            return result, nil
        }

        lastErr = err

        // Exponential backoff
        backoff := time.Duration(attempt*attempt) * time.Second
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        case <-time.After(backoff):
            continue
        }
    }

    return nil, fmt.Errorf("query failed after 3 attempts: %w", lastErr)
}
```

## Plugin Registry

### Submitting Plugins

1. **Code Review**: All plugins undergo code review
2. **Security Scan**: Automated security scanning
3. **Testing**: Comprehensive test suite required
4. **Documentation**: Complete documentation required

### Plugin Metadata

```yaml
# plugin.yaml
metadata:
  name: 'my-plugin'
  version: '1.0.0'
  description: 'Plugin description'
  author: 'Your Name <email@company.com>'
  license: 'MIT'
  homepage: 'https://github.com/company/my-plugin'

  tags: ['observability', 'logs', 'monitoring']
  categories: ['observability']

  compatibility:
    dash_ops_version: '>=2.0.0'
    go_version: '>=1.21'

  dependencies:
    - name: 'prometheus-client'
      version: '>=1.0.0'

security:
  permissions:
    - 'logs:read'
    - 'metrics:read'

  secrets:
    - name: 'API_KEY'
      description: 'API key for external service'
      required: true

support:
  documentation: 'https://docs.company.com/my-plugin'
  issues: 'https://github.com/company/my-plugin/issues'
  community: 'https://discord.gg/dash-ops'
```

## Community Guidelines

### Contributing

1. **Follow Standards**: Adhere to Go and React coding standards
2. **Write Tests**: Maintain >90% test coverage
3. **Document Everything**: Comprehensive documentation required
4. **Be Responsive**: Respond to issues and PRs promptly

### Support

1. **Community Forum**: Use GitHub Discussions for questions
2. **Bug Reports**: Use GitHub Issues for bugs
3. **Feature Requests**: Use GitHub Issues with enhancement label
4. **Security Issues**: Report privately to security@dash-ops.io

This guide provides the foundation for developing robust, scalable plugins for the dash-ops platform. Follow these patterns and best practices to create plugins that integrate seamlessly with the universal adapter architecture and AI capabilities.
