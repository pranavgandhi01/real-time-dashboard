# Security Guide

## Overview

This document outlines security best practices, configurations, and procedures for the Real-Time Flight Dashboard DevOps infrastructure. It covers container security, network isolation, secrets management, and monitoring security.

## Security Architecture

### Defense in Depth Strategy
```
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                        │
│  • Input validation  • Authentication  • Authorization     │
├─────────────────────────────────────────────────────────────┤
│                    Container Layer                          │
│  • Image scanning  • Runtime security  • Resource limits   │
├─────────────────────────────────────────────────────────────┤
│                    Network Layer                            │
│  • Network isolation  • TLS encryption  • Firewall rules   │
├─────────────────────────────────────────────────────────────┤
│                    Infrastructure Layer                     │
│  • Host hardening  • Access controls  • Audit logging      │
└─────────────────────────────────────────────────────────────┘
```

## Container Security

### Image Security

#### Base Image Selection
```dockerfile
# Use official, minimal base images
FROM node:18-alpine  # Preferred over full images
FROM redis:7-alpine   # Alpine variants for smaller attack surface
```

#### Image Scanning
```bash
# Scan images for vulnerabilities
docker scout cves grafana/grafana:latest
docker scout cves prom/prometheus:latest

# Use specific versions, not 'latest'
image: grafana/grafana:10.2.0  # Instead of :latest
```

#### Multi-stage Builds
```dockerfile
# Example secure multi-stage build
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

FROM node:18-alpine AS runtime
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nextjs -u 1001
USER nextjs
COPY --from=builder --chown=nextjs:nodejs /app .
```

### Runtime Security

#### Resource Limits
```yaml
# docker-compose.yml security configurations
services:
  redis:
    image: redis:7-alpine
    deploy:
      resources:
        limits:
          memory: 512M
          cpus: '0.5'
    security_opt:
      - no-new-privileges:true
    read_only: true
    tmpfs:
      - /tmp
```

#### User Privileges
```yaml
# Run containers as non-root
services:
  app:
    user: "1001:1001"  # Non-root user
    cap_drop:
      - ALL
    cap_add:
      - NET_BIND_SERVICE  # Only if needed
```

#### Security Options
```yaml
services:
  secure-service:
    security_opt:
      - no-new-privileges:true
      - apparmor:docker-default
    read_only: true
    tmpfs:
      - /tmp:noexec,nosuid,size=100m
```

## Network Security

### Network Isolation

#### Docker Networks
```yaml
# Separate networks for different tiers
networks:
  frontend:
    driver: bridge
    internal: false  # Internet access
  backend:
    driver: bridge
    internal: true   # No internet access
  monitoring:
    driver: bridge
    internal: true
```

#### Service Communication
```yaml
services:
  api-gateway:
    networks:
      - frontend
      - backend
  
  database:
    networks:
      - backend  # Only backend access
  
  monitoring:
    networks:
      - monitoring  # Isolated monitoring
```

### TLS/SSL Configuration

#### Grafana TLS
```yaml
# grafana.ini
[server]
protocol = https
cert_file = /etc/ssl/certs/grafana.crt
cert_key = /etc/ssl/private/grafana.key
```

#### Elasticsearch TLS
```yaml
services:
  elasticsearch:
    environment:
      - xpack.security.enabled=true
      - xpack.security.http.ssl.enabled=true
      - xpack.security.transport.ssl.enabled=true
    volumes:
      - ./certs:/usr/share/elasticsearch/config/certs:ro
```

#### Kafka TLS
```yaml
# kafka-cluster.yaml
spec:
  kafka:
    listeners:
      - name: tls
        port: 9093
        type: internal
        tls: true
        authentication:
          type: tls
```

### Firewall Rules

#### Host-level Firewall
```bash
# UFW configuration example
sudo ufw default deny incoming
sudo ufw default allow outgoing

# Allow specific ports
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw allow 3000/tcp  # Grafana (if external)

# Deny direct access to internal services
sudo ufw deny 9200/tcp   # Elasticsearch
sudo ufw deny 6379/tcp   # Redis
```

## Secrets Management

### Environment Variables

#### Secure Environment Files
```bash
# .env.production (never commit to git)
GRAFANA_ADMIN_PASSWORD=<strong-password>
REDIS_PASSWORD=<redis-password>
ELASTICSEARCH_PASSWORD=<es-password>
JWT_SECRET=<jwt-secret>
```

#### Docker Secrets
```yaml
# docker-compose.yml
services:
  grafana:
    environment:
      - GF_SECURITY_ADMIN_PASSWORD_FILE=/run/secrets/grafana_password
    secrets:
      - grafana_password

secrets:
  grafana_password:
    file: ./secrets/grafana_password.txt
```

### Kubernetes Secrets

#### Secret Creation
```bash
# Create secrets for Kafka
kubectl create secret generic kafka-credentials \
  --from-literal=username=admin \
  --from-literal=password=<strong-password> \
  -n kafka
```

#### Secret Usage
```yaml
# kafka-cluster.yaml
spec:
  kafka:
    authorization:
      type: simple
    authentication:
      type: scram-sha-512
      passwordSecret:
        secretName: kafka-credentials
        password: password
```

## Authentication and Authorization

### Service Authentication

#### Grafana Security
```yaml
# grafana.ini
[security]
admin_user = admin
admin_password = <strong-password>
secret_key = <random-secret-key>

[auth]
disable_login_form = false
disable_signout_menu = false

[auth.anonymous]
enabled = false
```

#### Elasticsearch Security
```yaml
services:
  elasticsearch:
    environment:
      - ELASTIC_PASSWORD=<strong-password>
      - xpack.security.enabled=true
      - xpack.security.authc.api_key.enabled=true
```

#### Redis Authentication
```yaml
services:
  redis:
    command: redis-server --requirepass <redis-password>
    environment:
      - REDIS_PASSWORD=<redis-password>
```

### API Security

#### Rate Limiting
```yaml
# nginx.conf for API Gateway
http {
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    
    server {
        location /api/ {
            limit_req zone=api burst=20 nodelay;
            proxy_pass http://backend;
        }
    }
}
```

#### CORS Configuration
```javascript
// Express.js CORS setup
const cors = require('cors');
app.use(cors({
  origin: ['https://yourdomain.com'],
  credentials: true,
  optionsSuccessStatus: 200
}));
```

## Monitoring Security

### Security Monitoring

#### Failed Authentication Attempts
```json
{
  "query": {
    "bool": {
      "must": [
        {"match": {"event.type": "authentication"}},
        {"match": {"event.outcome": "failure"}},
        {"range": {"@timestamp": {"gte": "now-1h"}}}
      ]
    }
  }
}
```

#### Suspicious Network Activity
```json
{
  "query": {
    "bool": {
      "should": [
        {"match": {"source.ip": "suspicious-ip"}},
        {"range": {"network.bytes": {"gte": 1000000}}},
        {"match": {"http.response.status_code": "403"}}
      ]
    }
  }
}
```

### Security Alerts

#### Grafana Security Alerts
```yaml
alert:
  name: "Multiple Failed Logins"
  message: "Multiple failed login attempts detected"
  conditions:
    - query:
        expr: "increase(grafana_api_login_post_total{status='error'}[5m])"
      evaluator:
        params: [5]
        type: "gt"
```

#### Intrusion Detection
```yaml
alert:
  name: "Unusual Network Traffic"
  message: "Unusual network traffic pattern detected"
  conditions:
    - query:
        expr: "rate(container_network_receive_bytes_total[5m])"
      evaluator:
        params: [10000000]  # 10MB/s
        type: "gt"
```

## Compliance and Auditing

### Audit Logging

#### Docker Audit
```bash
# Enable Docker daemon audit logging
# /etc/docker/daemon.json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  },
  "audit-log-format": "json",
  "audit-log-maxage": 30
}
```

#### Kubernetes Audit
```yaml
# audit-policy.yaml
apiVersion: audit.k8s.io/v1
kind: Policy
rules:
- level: Metadata
  namespaces: ["kafka"]
  resources:
  - group: "kafka.strimzi.io"
    resources: ["kafkas", "kafkatopics"]
```

### Compliance Checks

#### CIS Benchmarks
```bash
# Docker CIS benchmark check
docker run --rm -it \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /usr/bin/docker:/usr/bin/docker \
  docker/docker-bench-security
```

#### Security Scanning
```bash
# Container vulnerability scanning
trivy image grafana/grafana:latest
trivy image prom/prometheus:latest
trivy image redis:7-alpine
```

## Incident Response

### Security Incident Playbook

#### 1. Detection and Analysis
```bash
# Check for security events
# Kibana: level:error AND (authentication OR authorization OR security)

# Review access logs
docker logs nginx-container | grep -E "(403|401|404)"

# Check system integrity
docker diff container-name
```

#### 2. Containment
```bash
# Isolate compromised container
docker network disconnect network-name container-name

# Stop suspicious processes
docker exec container-name pkill -f suspicious-process

# Create forensic snapshot
docker commit container-name forensic-snapshot:latest
```

#### 3. Recovery
```bash
# Rotate compromised secrets
kubectl delete secret compromised-secret
kubectl create secret generic new-secret --from-literal=key=new-value

# Restart affected services
docker-compose restart affected-service

# Update security configurations
```

### Security Hardening Checklist

#### Container Hardening
- [ ] Use minimal base images (Alpine)
- [ ] Run containers as non-root users
- [ ] Enable read-only root filesystem
- [ ] Set resource limits
- [ ] Remove unnecessary packages
- [ ] Scan images for vulnerabilities
- [ ] Use specific image tags, not 'latest'

#### Network Hardening
- [ ] Implement network segmentation
- [ ] Enable TLS for all communications
- [ ] Configure firewall rules
- [ ] Disable unnecessary ports
- [ ] Use private networks for internal communication

#### Access Control
- [ ] Change default passwords
- [ ] Implement strong authentication
- [ ] Use role-based access control
- [ ] Enable audit logging
- [ ] Regular access reviews

#### Monitoring
- [ ] Monitor failed authentication attempts
- [ ] Alert on suspicious activities
- [ ] Log all administrative actions
- [ ] Regular security assessments
- [ ] Incident response procedures

## Security Maintenance

### Regular Security Tasks

#### Daily
- Review security alerts and logs
- Monitor failed authentication attempts
- Check for unusual network activity

#### Weekly
- Update container images
- Review access logs
- Scan for vulnerabilities
- Backup security configurations

#### Monthly
- Security assessment
- Update security policies
- Review and rotate secrets
- Compliance audit
- Incident response drill

### Security Updates

#### Automated Updates
```bash
# Watchtower for automatic container updates
docker run -d \
  --name watchtower \
  -v /var/run/docker.sock:/var/run/docker.sock \
  containrrr/watchtower \
  --schedule "0 0 4 * * *"  # Daily at 4 AM
```

#### Manual Update Process
```bash
# 1. Check for updates
docker-compose pull

# 2. Test in staging
docker-compose -f docker-compose.staging.yml up -d

# 3. Deploy to production
docker-compose up -d

# 4. Verify security
./scripts/security-check.sh
```