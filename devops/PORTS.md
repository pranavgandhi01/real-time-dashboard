# Port Allocation

## DevOps Infrastructure Ports

| Service | Port | Purpose |
|---------|------|---------|
| Redis | 6379 | Cache & Sessions |
| Prometheus | 9090 | Metrics Collection |
| Grafana | 3000 | Metrics Dashboards |
| Jaeger | 16686 | Distributed Tracing UI |
| Elasticsearch | 9200 | Log Storage API |
| Kibana | 5601 | Log Analysis UI |
| Kafka | 9094 | Event Streaming (K8s NodePort: 32092) |

## Application Ports

| Service | Port | Purpose |
|---------|------|---------|
| Frontend | 3000 | Next.js Application |
| API Gateway | 8080 | Request Routing |
| Flight Data Service | 8081 | Flight Data REST API |
| WebSocket Service | 8082 | Real-time Updates |
| Mock Data Service | 8083 | Testing Data |

## Port Conflicts

- **Grafana vs Frontend**: Both use 3000 - Run separately or change one
- **Reserved Ranges**: 
  - 3000-3099: Frontend applications
  - 8080-8089: Backend services
  - 9000-9099: Monitoring tools
  - 5000-5999: Databases & storage