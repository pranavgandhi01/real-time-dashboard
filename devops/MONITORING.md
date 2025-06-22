# Monitoring and Observability Guide

## Overview

The Real-Time Flight Dashboard implements comprehensive observability using the three pillars: **Metrics**, **Logs**, and **Traces**. This guide covers monitoring setup, dashboard configuration, and operational procedures.

## Observability Stack

### Metrics (Prometheus + Grafana)
- **Collection**: Prometheus scrapes metrics from application services
- **Storage**: Time-series database with configurable retention
- **Visualization**: Grafana dashboards with alerting capabilities
- **Endpoints**: `/metrics` on each service

### Logs (ELK Stack)
- **Collection**: Filebeat collects Docker container logs
- **Processing**: Elasticsearch indexes and stores logs
- **Analysis**: Kibana provides search and visualization
- **Format**: JSON structured logging recommended

### Traces (Jaeger)
- **Collection**: OpenTelemetry compatible tracing
- **Storage**: In-memory (development) / Persistent (production)
- **Analysis**: Distributed request tracing and performance analysis

## Service Monitoring

### Application Services

#### API Gateway (Port 8080)
```yaml
# Prometheus scrape config
- job_name: 'api-gateway'
  static_configs:
    - targets: ['api-gateway:8080']
  metrics_path: '/metrics'
```

**Key Metrics**:
- `http_requests_total` - Request count by status code
- `http_request_duration_seconds` - Request latency
- `gateway_active_connections` - Active WebSocket connections

#### Flight Data Service (Port 8081)
```yaml
- job_name: 'flight-data-service'
  static_configs:
    - targets: ['flight-data-service:8081']
```

**Key Metrics**:
- `flight_data_requests_total` - API requests
- `flight_cache_hits_total` - Cache performance
- `external_api_calls_total` - Third-party API usage

#### WebSocket Service (Port 8082)
```yaml
- job_name: 'websocket-service'
  static_configs:
    - targets: ['websocket-service:8082']
```

**Key Metrics**:
- `websocket_connections_active` - Active connections
- `websocket_messages_sent_total` - Message throughput
- `websocket_connection_duration_seconds` - Connection lifetime

### Infrastructure Monitoring

#### Redis Monitoring
```bash
# Redis metrics via redis_exporter (optional)
docker run -d --name redis_exporter \
  -p 9121:9121 \
  oliver006/redis_exporter \
  --redis.addr=redis://redis:6379
```

**Key Metrics**:
- `redis_connected_clients` - Active connections
- `redis_memory_used_bytes` - Memory usage
- `redis_commands_processed_total` - Command throughput

#### Kafka Monitoring
```yaml
# Built-in JMX metrics via Strimzi
spec:
  kafka:
    metricsConfig:
      type: jmxPrometheusExporter
```

**Key Metrics**:
- `kafka_server_brokertopicmetrics_messagesin_total` - Message ingestion
- `kafka_server_brokertopicmetrics_bytesin_total` - Throughput
- `kafka_controller_kafkacontroller_activecontrollercount` - Controller status

## Dashboard Configuration

### Grafana Dashboard Setup

#### 1. Application Performance Dashboard
```json
{
  "dashboard": {
    "title": "Flight Dashboard - Application Performance",
    "panels": [
      {
        "title": "Request Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(http_requests_total[5m])",
            "legendFormat": "{{service}} - {{status}}"
          }
        ]
      },
      {
        "title": "Response Time",
        "type": "graph",
        "targets": [
          {
            "expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
            "legendFormat": "95th percentile"
          }
        ]
      }
    ]
  }
}
```

#### 2. Infrastructure Dashboard
```json
{
  "dashboard": {
    "title": "Flight Dashboard - Infrastructure",
    "panels": [
      {
        "title": "Redis Memory Usage",
        "type": "singlestat",
        "targets": [
          {
            "expr": "redis_memory_used_bytes / 1024 / 1024",
            "legendFormat": "Memory (MB)"
          }
        ]
      },
      {
        "title": "Kafka Message Rate",
        "type": "graph",
        "targets": [
          {
            "expr": "rate(kafka_server_brokertopicmetrics_messagesin_total[5m])",
            "legendFormat": "Messages/sec"
          }
        ]
      }
    ]
  }
}
```

### Kibana Index Patterns

#### 1. Application Logs
```bash
# Index pattern: flight-tracker-logs-*
# Time field: @timestamp
```

**Common Queries**:
```json
# Error logs
{
  "query": {
    "bool": {
      "must": [
        {"match": {"level": "error"}},
        {"range": {"@timestamp": {"gte": "now-1h"}}}
      ]
    }
  }
}

# Service-specific logs
{
  "query": {
    "bool": {
      "must": [
        {"match": {"container.name": "api-gateway"}}
      ]
    }
  }
}
```

#### 2. Performance Analysis
```json
# Slow requests (>1s response time)
{
  "query": {
    "bool": {
      "must": [
        {"range": {"response_time": {"gte": 1000}}}
      ]
    }
  }
}
```

## Alerting Configuration

### Grafana Alerts

#### 1. High Error Rate Alert
```yaml
alert:
  name: "High Error Rate"
  message: "Error rate is above 5% for 5 minutes"
  frequency: "10s"
  conditions:
    - query:
        queryType: ""
        refId: "A"
        model:
          expr: "rate(http_requests_total{status=~'5..'}[5m]) / rate(http_requests_total[5m]) * 100"
      reducer:
        type: "last"
        params: []
      evaluator:
        params: [5]
        type: "gt"
```

#### 2. Service Down Alert
```yaml
alert:
  name: "Service Down"
  message: "Service {{$labels.job}} is down"
  frequency: "10s"
  conditions:
    - query:
        expr: "up == 0"
      evaluator:
        params: [0]
        type: "eq"
```

#### 3. High Memory Usage Alert
```yaml
alert:
  name: "High Memory Usage"
  message: "Redis memory usage is above 80%"
  conditions:
    - query:
        expr: "(redis_memory_used_bytes / redis_memory_max_bytes) * 100"
      evaluator:
        params: [80]
        type: "gt"
```

### Notification Channels

#### Slack Integration
```json
{
  "name": "slack-alerts",
  "type": "slack",
  "settings": {
    "url": "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK",
    "channel": "#alerts",
    "title": "Flight Dashboard Alert",
    "text": "{{ range .Alerts }}{{ .Annotations.summary }}{{ end }}"
  }
}
```

#### Email Notifications
```json
{
  "name": "email-alerts",
  "type": "email",
  "settings": {
    "addresses": "ops-team@company.com",
    "subject": "Flight Dashboard Alert: {{ .GroupLabels.alertname }}"
  }
}
```

## Log Management

### Structured Logging Format

#### Application Logs
```json
{
  "@timestamp": "2024-01-15T10:30:00.000Z",
  "level": "info",
  "service": "api-gateway",
  "message": "Request processed",
  "request_id": "req-123456",
  "user_id": "user-789",
  "duration_ms": 150,
  "status_code": 200,
  "endpoint": "/api/flights"
}
```

#### Error Logs
```json
{
  "@timestamp": "2024-01-15T10:30:00.000Z",
  "level": "error",
  "service": "flight-data-service",
  "message": "External API timeout",
  "error": {
    "type": "TimeoutError",
    "message": "Request timeout after 5000ms",
    "stack": "..."
  },
  "request_id": "req-123456",
  "external_service": "flight-api"
}
```

### Log Retention Policies

#### Elasticsearch Index Management
```json
{
  "policy": {
    "phases": {
      "hot": {
        "actions": {
          "rollover": {
            "max_size": "1GB",
            "max_age": "1d"
          }
        }
      },
      "warm": {
        "min_age": "7d",
        "actions": {
          "allocate": {
            "number_of_replicas": 0
          }
        }
      },
      "delete": {
        "min_age": "30d"
      }
    }
  }
}
```

## Performance Monitoring

### Key Performance Indicators (KPIs)

#### Application KPIs
- **Availability**: 99.9% uptime target
- **Response Time**: 95th percentile < 500ms
- **Error Rate**: < 1% of total requests
- **Throughput**: Requests per second capacity

#### Infrastructure KPIs
- **Redis Hit Rate**: > 95%
- **Kafka Lag**: < 1000 messages
- **Memory Usage**: < 80% of allocated
- **CPU Usage**: < 70% average

### SLA Monitoring

#### Service Level Objectives (SLOs)
```yaml
slos:
  - name: "API Response Time"
    target: 0.95  # 95% of requests
    threshold: "500ms"
    
  - name: "Service Availability"
    target: 0.999  # 99.9% uptime
    measurement_window: "30d"
    
  - name: "Error Budget"
    target: 0.01  # 1% error rate
    measurement_window: "7d"
```

## Troubleshooting Runbook

### Common Issues

#### 1. High Response Times
**Symptoms**: 95th percentile > 1s
**Investigation**:
```bash
# Check service metrics
curl http://localhost:9090/api/v1/query?query=histogram_quantile(0.95,rate(http_request_duration_seconds_bucket[5m]))

# Check logs for slow queries
# Kibana: response_time:>1000 AND @timestamp:[now-1h TO now]
```

**Resolution**:
- Scale service instances
- Optimize database queries
- Check external API performance

#### 2. Service Unavailable
**Symptoms**: Service returning 503 errors
**Investigation**:
```bash
# Check service health
curl -f http://service:port/health

# Check container status
docker ps | grep service-name

# Check logs
docker logs service-container
```

**Resolution**:
- Restart service container
- Check resource constraints
- Verify dependencies

#### 3. Memory Leaks
**Symptoms**: Increasing memory usage over time
**Investigation**:
```bash
# Monitor memory trends
# Grafana: container_memory_usage_bytes{name="service-name"}

# Check for memory leaks in logs
# Kibana: level:warn AND message:"memory"
```

**Resolution**:
- Restart affected service
- Review code for memory leaks
- Adjust memory limits

## Maintenance Procedures

### Daily Checks
```bash
# Run status check
./scripts/status.sh

# Check disk usage
df -h
docker system df

# Review error logs
# Kibana: level:error AND @timestamp:[now-24h TO now]
```

### Weekly Maintenance
```bash
# Update monitoring stack
cd devops/observability
docker-compose pull
docker-compose up -d

# Clean up old logs
docker system prune -f

# Review performance trends
# Grafana: Weekly performance review dashboard
```

### Monthly Reviews
- Review and update alerting thresholds
- Analyze performance trends
- Update monitoring dashboards
- Review log retention policies
- Capacity planning based on growth trends