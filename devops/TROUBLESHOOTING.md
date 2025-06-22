# Troubleshooting Guide

## Overview

This guide provides systematic approaches to diagnose and resolve common issues in the Real-Time Flight Dashboard DevOps infrastructure. It includes diagnostic commands, common problems, and step-by-step resolution procedures.

## General Troubleshooting Approach

### 1. Identify the Problem
- **Symptoms**: What is not working?
- **Scope**: Which services are affected?
- **Timeline**: When did the issue start?
- **Impact**: How many users are affected?

### 2. Gather Information
```bash
# Quick status check
./scripts/status.sh

# Check all containers
docker ps -a

# Check system resources
docker stats
df -h
free -h
```

### 3. Isolate the Issue
- Check service dependencies
- Review recent changes
- Examine logs for errors
- Test individual components

## Common Issues and Solutions

### Docker and Container Issues

#### Issue: Container Won't Start
**Symptoms**: Container exits immediately or fails to start

**Diagnosis**:
```bash
# Check container status
docker ps -a

# View container logs
docker logs container-name

# Check resource usage
docker stats

# Inspect container configuration
docker inspect container-name
```

**Common Causes & Solutions**:

1. **Port Conflicts**
```bash
# Check port usage
netstat -tulpn | grep :3000
lsof -i :3000

# Solution: Change port or stop conflicting service
docker-compose down
# Edit docker-compose.yml to change ports
docker-compose up -d
```

2. **Resource Constraints**
```bash
# Check available resources
free -h
df -h

# Solution: Free up resources or adjust limits
docker system prune -f
docker volume prune -f
```

3. **Configuration Errors**
```bash
# Validate docker-compose file
docker-compose config

# Check environment variables
docker-compose config | grep -A 5 environment
```

#### Issue: Container Keeps Restarting
**Symptoms**: Container in restart loop

**Diagnosis**:
```bash
# Check restart count
docker ps -a

# View recent logs
docker logs --tail 50 container-name

# Check health status
docker inspect container-name | grep -A 10 Health
```

**Solutions**:
```bash
# 1. Fix configuration issues
docker-compose down
# Edit configuration files
docker-compose up -d

# 2. Increase resource limits
# Edit docker-compose.yml
deploy:
  resources:
    limits:
      memory: 2G
      cpus: '1.0'

# 3. Check dependencies
# Ensure dependent services are running
docker-compose up -d dependency-service
```

### Service-Specific Issues

#### Grafana Issues

**Issue: Cannot Access Grafana Dashboard**
```bash
# Check Grafana container
docker logs observability-grafana-1

# Check port binding
docker port observability-grafana-1

# Test connectivity
curl -f http://localhost:3000

# Common solutions:
# 1. Port conflict (change to 3001)
# 2. Container not running (restart)
# 3. Network issues (check docker networks)
```

**Issue: Grafana Login Problems**
```bash
# Reset admin password
docker exec -it observability-grafana-1 grafana-cli admin reset-admin-password newpassword

# Check configuration
docker exec observability-grafana-1 cat /etc/grafana/grafana.ini | grep -A 5 security
```

#### Prometheus Issues

**Issue: Prometheus Not Scraping Targets**
```bash
# Check Prometheus targets
curl http://localhost:9090/api/v1/targets

# Check configuration
docker exec observability-prometheus-1 cat /etc/prometheus/prometheus.yml

# Common solutions:
# 1. Service discovery issues
# 2. Network connectivity problems
# 3. Incorrect target configuration
```

**Issue: High Memory Usage**
```bash
# Check Prometheus memory usage
docker stats observability-prometheus-1

# Reduce retention period
# Edit prometheus.yml
global:
  retention: 7d  # Reduce from default 15d

# Restart Prometheus
docker-compose restart prometheus
```

#### Elasticsearch Issues

**Issue: Elasticsearch Won't Start**
```bash
# Check logs
docker logs observability-elasticsearch-1

# Common issues:
# 1. Insufficient memory
# 2. Disk space issues
# 3. Permission problems
```

**Solutions**:
```bash
# 1. Increase memory limit
# docker-compose.yml
environment:
  - "ES_JAVA_OPTS=-Xms1g -Xmx1g"

# 2. Fix permissions
sudo chown -R 1000:1000 elasticsearch-data/

# 3. Check disk space
df -h
docker system prune -f
```

**Issue: Elasticsearch Cluster Red Status**
```bash
# Check cluster health
curl http://localhost:9200/_cluster/health

# Check indices
curl http://localhost:9200/_cat/indices?v

# Solutions:
# 1. Restart Elasticsearch
docker-compose restart elasticsearch

# 2. Delete problematic indices
curl -X DELETE http://localhost:9200/problematic-index
```

#### Redis Issues

**Issue: Redis Connection Refused**
```bash
# Test Redis connectivity
docker exec infrastructure-redis-1 redis-cli ping

# Check Redis logs
docker logs infrastructure-redis-1

# Test from host
redis-cli -h localhost -p 6379 ping
```

**Solutions**:
```bash
# 1. Restart Redis
docker-compose restart redis

# 2. Check network connectivity
docker network ls
docker network inspect infrastructure-network

# 3. Verify port binding
docker port infrastructure-redis-1
```

#### Kafka Issues

**Issue: Kafka Cluster Not Ready**
```bash
# Check Kafka status
kubectl get kafka -n kafka

# Check pods
kubectl get pods -n kafka

# Check events
kubectl get events -n kafka --sort-by='.lastTimestamp'
```

**Solutions**:
```bash
# 1. Wait for deployment (can take 5-10 minutes)
kubectl wait --for=condition=ready kafka/flight-tracker-kafka -n kafka --timeout=600s

# 2. Check Strimzi operator
kubectl logs deployment/strimzi-cluster-operator -n kafka

# 3. Restart cluster
kubectl delete kafka flight-tracker-kafka -n kafka
kubectl apply -f kafka-cluster.yaml
```

**Issue: Cannot Connect to Kafka Externally**
```bash
# Check external service
kubectl get svc -n kafka | grep external

# Check NodePort
kubectl get svc flight-tracker-kafka-kafka-external-bootstrap -n kafka

# Test connectivity
telnet localhost 32092
```

**Solutions**:
```bash
# 1. Port forward for testing
kubectl port-forward svc/flight-tracker-kafka-kafka-bootstrap 9092:9092 -n kafka

# 2. Check Kind cluster port mapping
kind get clusters
docker port flight-tracker-control-plane

# 3. Recreate external service
kubectl delete svc flight-tracker-kafka-kafka-external-bootstrap -n kafka
# Wait for Strimzi to recreate
```

### Network Issues

#### Issue: Services Cannot Communicate
**Diagnosis**:
```bash
# Check Docker networks
docker network ls

# Inspect network configuration
docker network inspect observability-network

# Test connectivity between containers
docker exec container1 ping container2
docker exec container1 telnet container2 port
```

**Solutions**:
```bash
# 1. Ensure services are on same network
# docker-compose.yml
services:
  service1:
    networks:
      - shared-network
  service2:
    networks:
      - shared-network

# 2. Recreate networks
docker-compose down
docker network prune -f
docker-compose up -d
```

#### Issue: DNS Resolution Problems
```bash
# Test DNS resolution
docker exec container nslookup service-name

# Check /etc/hosts
docker exec container cat /etc/hosts

# Solutions:
# 1. Use service names instead of IPs
# 2. Restart Docker daemon
sudo systemctl restart docker
```

### Performance Issues

#### Issue: High CPU Usage
**Diagnosis**:
```bash
# Check container CPU usage
docker stats

# Check host CPU
top
htop

# Check specific processes
docker exec container top
```

**Solutions**:
```bash
# 1. Set CPU limits
# docker-compose.yml
deploy:
  resources:
    limits:
      cpus: '0.5'

# 2. Scale services
docker-compose up -d --scale service-name=3

# 3. Optimize application code
```

#### Issue: High Memory Usage
**Diagnosis**:
```bash
# Check memory usage
docker stats
free -h

# Check for memory leaks
docker exec container ps aux --sort=-%mem
```

**Solutions**:
```bash
# 1. Increase memory limits
# docker-compose.yml
deploy:
  resources:
    limits:
      memory: 2G

# 2. Restart containers
docker-compose restart

# 3. Clean up unused resources
docker system prune -f
```

#### Issue: Slow Response Times
**Diagnosis**:
```bash
# Test response times
curl -w "@curl-format.txt" -o /dev/null -s http://localhost:3000

# Check service metrics
curl http://localhost:9090/api/v1/query?query=http_request_duration_seconds

# Check logs for slow queries
# Kibana: response_time:>1000
```

**Solutions**:
```bash
# 1. Scale services
docker-compose up -d --scale api-gateway=3

# 2. Add caching
# Configure Redis caching

# 3. Optimize database queries
# Review slow query logs
```

### Storage Issues

#### Issue: Disk Space Full
**Diagnosis**:
```bash
# Check disk usage
df -h

# Check Docker space usage
docker system df

# Find large files
du -sh /var/lib/docker/*
```

**Solutions**:
```bash
# 1. Clean up Docker resources
docker system prune -f
docker volume prune -f
docker image prune -a -f

# 2. Clean up logs
sudo journalctl --vacuum-time=7d
docker logs container-name --tail 0

# 3. Rotate logs
# Configure log rotation in docker-compose.yml
logging:
  driver: "json-file"
  options:
    max-size: "10m"
    max-file: "3"
```

#### Issue: Volume Mount Problems
**Diagnosis**:
```bash
# Check volume mounts
docker inspect container-name | grep -A 10 Mounts

# Check permissions
ls -la /host/path
docker exec container ls -la /container/path
```

**Solutions**:
```bash
# 1. Fix permissions
sudo chown -R 1000:1000 /host/path

# 2. Use named volumes instead of bind mounts
# docker-compose.yml
volumes:
  data-volume:
    driver: local

services:
  service:
    volumes:
      - data-volume:/data
```

## Diagnostic Commands Reference

### Container Diagnostics
```bash
# Container status and resource usage
docker ps -a
docker stats
docker system df

# Container logs
docker logs container-name
docker logs --tail 50 --follow container-name

# Container inspection
docker inspect container-name
docker exec -it container-name /bin/sh
```

### Network Diagnostics
```bash
# Network information
docker network ls
docker network inspect network-name

# Connectivity testing
docker exec container ping target
docker exec container telnet host port
docker exec container nslookup hostname
```

### Service-Specific Diagnostics
```bash
# Prometheus
curl http://localhost:9090/api/v1/targets
curl http://localhost:9090/api/v1/query?query=up

# Elasticsearch
curl http://localhost:9200/_cluster/health
curl http://localhost:9200/_cat/indices?v

# Redis
docker exec redis-container redis-cli ping
docker exec redis-container redis-cli info

# Kafka
kubectl get kafka,kafkatopic -n kafka
kubectl logs -f deployment/strimzi-cluster-operator -n kafka
```

### System Diagnostics
```bash
# System resources
free -h
df -h
top
iostat

# Docker daemon
sudo systemctl status docker
sudo journalctl -u docker --tail 50

# Kubernetes (for Kafka)
kubectl cluster-info
kubectl get nodes
kubectl get pods --all-namespaces
```

## Emergency Procedures

### Complete System Recovery
```bash
# 1. Stop all services
./scripts/stop-all.sh

# 2. Clean up resources
docker system prune -a -f
docker volume prune -f
docker network prune -f

# 3. Restart Docker daemon
sudo systemctl restart docker

# 4. Start services
./scripts/start-all.sh
```

### Data Recovery
```bash
# 1. Backup current state
docker run --rm -v volume-name:/data -v $(pwd):/backup alpine tar czf /backup/backup.tar.gz /data

# 2. Restore from backup
docker run --rm -v volume-name:/data -v $(pwd):/backup alpine tar xzf /backup/backup.tar.gz -C /

# 3. Restart services
docker-compose restart
```

### Service Isolation
```bash
# Isolate problematic service
docker network disconnect network-name container-name

# Run service in debug mode
docker run -it --rm image-name /bin/sh

# Restart single service
docker-compose restart service-name
```

## Monitoring and Alerting

### Health Check Scripts
```bash
#!/bin/bash
# health-check.sh

services=("grafana" "prometheus" "elasticsearch" "redis")
for service in "${services[@]}"; do
    if curl -f http://localhost:${ports[$service]} >/dev/null 2>&1; then
        echo "✅ $service is healthy"
    else
        echo "❌ $service is unhealthy"
        # Send alert
        curl -X POST webhook-url -d "Service $service is down"
    fi
done
```

### Log Analysis
```bash
# Find error patterns
docker logs container-name 2>&1 | grep -i error | tail -20

# Monitor logs in real-time
docker logs -f container-name | grep -E "(error|warn|fail)"

# Analyze log volume
docker logs container-name | wc -l
```

### Performance Monitoring
```bash
# Monitor resource usage over time
while true; do
    echo "$(date): $(docker stats --no-stream --format 'table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}')"
    sleep 60
done > resource-usage.log
```

## Prevention and Best Practices

### Proactive Monitoring
- Set up comprehensive alerting
- Regular health checks
- Capacity planning
- Performance baselines

### Maintenance Procedures
- Regular updates and patches
- Log rotation and cleanup
- Backup verification
- Documentation updates

### Testing Procedures
- Disaster recovery drills
- Load testing
- Chaos engineering
- Security assessments