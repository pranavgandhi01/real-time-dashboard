# Contributing to DevOps Infrastructure

Thank you for your interest in contributing to the Real-Time Flight Dashboard DevOps infrastructure! This guide will help you understand our development process, coding standards, and how to submit contributions.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Documentation Standards](#documentation-standards)
- [Submission Process](#submission-process)
- [Review Process](#review-process)

## Getting Started

### Prerequisites

Before contributing, ensure you have:

- **Docker**: Version 20.10 or higher
- **Docker Compose**: Version 2.0 or higher
- **Git**: For version control
- **Text Editor**: VS Code, Vim, or your preferred editor
- **Optional**: Kind and kubectl for Kafka development

### Development Environment Setup

1. **Clone the Repository**
```bash
git clone <repository-url>
cd real-time-dashboard/devops
```

2. **Start Development Environment**
```bash
./scripts/start-all.sh
```

3. **Verify Setup**
```bash
./scripts/status.sh
```

4. **Access Services**
- Grafana: http://localhost:3000 (admin/admin)
- Prometheus: http://localhost:9090
- Jaeger: http://localhost:16686

## Development Workflow

### Branch Strategy

We use a feature branch workflow:

```
main
├── feature/monitoring-improvements
├── feature/kafka-security
├── hotfix/grafana-config
└── release/v1.1.0
```

### Branch Naming Convention

- **Features**: `feature/description-of-feature`
- **Bug Fixes**: `bugfix/description-of-fix`
- **Hotfixes**: `hotfix/critical-issue-description`
- **Releases**: `release/v1.1.0`

### Development Process

1. **Create Feature Branch**
```bash
git checkout -b feature/your-feature-name
```

2. **Make Changes**
- Follow coding standards
- Update documentation
- Add tests where applicable

3. **Test Changes**
```bash
# Test your changes
./scripts/start-all.sh
./scripts/status.sh

# Run any specific tests
docker-compose config  # Validate configurations
```

4. **Commit Changes**
```bash
git add .
git commit -m "feat: add monitoring dashboard for Kafka metrics"
```

5. **Push and Create PR**
```bash
git push origin feature/your-feature-name
# Create pull request via GitHub/GitLab
```

## Coding Standards

### Docker Compose Files

#### Structure and Formatting
```yaml
# Use consistent indentation (2 spaces)
services:
  service-name:
    image: image:tag  # Always use specific tags
    ports:
      - "host:container"  # Always quote port mappings
    environment:
      - ENV_VAR=value
    networks:
      - network-name
    volumes:
      - volume-name:/path
    healthcheck:
      test: ["CMD", "command"]
      interval: 30s
      timeout: 10s
      retries: 3

networks:
  network-name:
    driver: bridge

volumes:
  volume-name:
    driver: local
```

#### Best Practices
- **Specific Tags**: Never use `latest` tag in production configurations
- **Resource Limits**: Always define resource constraints
- **Health Checks**: Include health checks for all services
- **Networks**: Use dedicated networks for service isolation
- **Volumes**: Use named volumes instead of bind mounts when possible

### Shell Scripts

#### Script Header
```bash
#!/bin/bash
set -e  # Exit on error

# Script description
# Usage: ./script.sh [options]
# Author: Your Name
# Date: YYYY-MM-DD
```

#### Function Structure
```bash
# Function description
function_name() {
    local param1=$1
    local param2=$2
    
    # Function logic
    echo "Processing $param1"
    
    return 0
}
```

#### Error Handling
```bash
# Check prerequisites
if ! command -v docker &> /dev/null; then
    echo "❌ Docker not found"
    exit 1
fi

# Validate parameters
if [ $# -eq 0 ]; then
    echo "Usage: $0 <parameter>"
    exit 1
fi
```

### Configuration Files

#### YAML Files
```yaml
# Use consistent formatting
# Comments should explain why, not what
global:
  scrape_interval: 15s  # Balance between accuracy and performance
  
scrape_configs:
  - job_name: 'service-name'
    static_configs:
      - targets: ['service:port']
    metrics_path: '/metrics'
    scrape_interval: 30s  # Override global for specific needs
```

#### Environment Variables
```bash
# Use descriptive names
GRAFANA_ADMIN_PASSWORD=secure-password
PROMETHEUS_RETENTION_DAYS=7
ELASTICSEARCH_HEAP_SIZE=1g

# Group related variables
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=dashboard

# Monitoring Configuration
METRICS_ENABLED=true
METRICS_PORT=9090
```

## Testing Guidelines

### Configuration Testing

#### Docker Compose Validation
```bash
# Validate syntax
docker-compose config

# Test service startup
docker-compose up -d service-name
docker-compose ps
docker-compose logs service-name
```

#### Service Health Testing
```bash
# Create test script
#!/bin/bash
test_service_health() {
    local service_url=$1
    local service_name=$2
    
    if curl -f "$service_url" >/dev/null 2>&1; then
        echo "✅ $service_name is healthy"
        return 0
    else
        echo "❌ $service_name is unhealthy"
        return 1
    fi
}

# Test all services
test_service_health "http://localhost:9090" "Prometheus"
test_service_health "http://localhost:3000" "Grafana"
```

### Integration Testing

#### End-to-End Testing
```bash
#!/bin/bash
# e2e-test.sh

echo "🧪 Running end-to-end tests"

# 1. Start all services
./scripts/start-all.sh

# 2. Wait for services to be ready
sleep 30

# 3. Test service connectivity
./scripts/status.sh

# 4. Test data flow
# Send test metrics to Prometheus
# Verify data appears in Grafana
# Check logs in Kibana

# 5. Cleanup
./scripts/stop-all.sh
```

### Performance Testing

#### Load Testing
```bash
# Test service under load
docker run --rm -i loadimpact/k6 run - <<EOF
import http from 'k6/http';
import { check } from 'k6';

export let options = {
  stages: [
    { duration: '2m', target: 100 },
    { duration: '5m', target: 100 },
    { duration: '2m', target: 0 },
  ],
};

export default function () {
  let response = http.get('http://localhost:9090/api/v1/targets');
  check(response, {
    'status is 200': (r) => r.status === 200,
    'response time < 500ms': (r) => r.timings.duration < 500,
  });
}
EOF
```

## Documentation Standards

### README Files

#### Structure
```markdown
# Service Name

Brief description of the service and its purpose.

## Quick Start
- Minimal steps to get started

## Configuration
- Key configuration options

## Troubleshooting
- Common issues and solutions

## API Reference
- If applicable
```

#### Writing Style
- Use clear, concise language
- Include code examples
- Provide context for decisions
- Update documentation with code changes

### Code Comments

#### Docker Compose Comments
```yaml
services:
  prometheus:
    image: prom/prometheus:v2.40.0
    ports:
      - "9090:9090"  # Prometheus web UI
    volumes:
      # Configuration file with scrape targets
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    command:
      # Enable web admin API for configuration reloads
      - '--web.enable-admin-api'
      # Set data retention to 7 days to manage disk usage
      - '--storage.tsdb.retention.time=7d'
```

#### Script Comments
```bash
#!/bin/bash
# Comprehensive health check for all DevOps services
# Checks service availability, response times, and basic functionality
# Usage: ./health-check.sh [--verbose]

set -e

# Configuration
TIMEOUT=5  # Connection timeout in seconds
VERBOSE=false

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --verbose)
            VERBOSE=true
            shift
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done
```

### Architecture Documentation

#### Diagrams
Use ASCII art or mermaid diagrams for architecture:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Application   │───▶│   API Gateway   │───▶│    Services     │
│   (Frontend)    │    │   (Port 8080)   │    │  (8081-8083)    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Observability Stack                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────────┐ │
│  │   Jaeger    │  │ Prometheus  │  │      Elasticsearch      │ │
│  │  (Tracing)  │  │ (Metrics)   │  │        (Logs)           │ │
│  └─────────────┘  └─────────────┘  └─────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

## Submission Process

### Pull Request Guidelines

#### PR Title Format
```
type(scope): brief description

Examples:
feat(monitoring): add Kafka metrics dashboard
fix(scripts): resolve startup race condition
docs(security): update authentication procedures
```

#### PR Description Template
```markdown
## Description
Brief description of changes and motivation.

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
- [ ] Tested locally with `./scripts/start-all.sh`
- [ ] All services pass health checks
- [ ] Configuration files validated
- [ ] Documentation updated

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No breaking changes (or clearly documented)
```

### Commit Message Format

Follow conventional commits:

```
type(scope): description

body (optional)

footer (optional)
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(kafka): add external access configuration

Add NodePort service for external Kafka access to enable
application connectivity from outside the Kubernetes cluster.

Closes #123

fix(grafana): resolve dashboard provisioning issue

The dashboard provisioning was failing due to incorrect
file permissions. Updated the volume mount to use proper
ownership settings.

docs(deployment): update prerequisites section

Added memory requirements and clarified Docker version
compatibility for better user experience.
```

## Review Process

### Review Criteria

Reviewers will check for:

1. **Functionality**
   - Changes work as intended
   - No breaking changes to existing functionality
   - Proper error handling

2. **Code Quality**
   - Follows coding standards
   - Appropriate comments and documentation
   - No security vulnerabilities

3. **Testing**
   - Adequate testing coverage
   - All tests pass
   - Manual testing completed

4. **Documentation**
   - Documentation updated
   - Clear commit messages
   - Proper PR description

### Review Timeline

- **Initial Review**: Within 2 business days
- **Follow-up Reviews**: Within 1 business day
- **Approval**: Requires at least 1 approval from maintainer
- **Merge**: After approval and CI checks pass

### Addressing Review Comments

1. **Make Requested Changes**
```bash
# Make changes based on feedback
git add .
git commit -m "address review comments: fix configuration syntax"
git push origin feature/your-feature-name
```

2. **Respond to Comments**
- Mark resolved comments as resolved
- Explain your approach if different from suggestion
- Ask for clarification if needed

3. **Request Re-review**
- Use GitHub/GitLab re-review request feature
- Comment when ready for re-review

## Getting Help

### Communication Channels

- **Issues**: Use GitHub/GitLab issues for bugs and feature requests
- **Discussions**: Use discussion forums for questions and ideas
- **Documentation**: Check existing documentation first

### Issue Templates

#### Bug Report
```markdown
**Describe the bug**
A clear description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Run command '...'
2. Check service '....'
3. See error

**Expected behavior**
What you expected to happen.

**Environment:**
- OS: [e.g. Ubuntu 20.04]
- Docker version: [e.g. 20.10.0]
- Docker Compose version: [e.g. 2.0.0]

**Additional context**
Add any other context about the problem here.
```

#### Feature Request
```markdown
**Is your feature request related to a problem?**
A clear description of what the problem is.

**Describe the solution you'd like**
A clear description of what you want to happen.

**Describe alternatives you've considered**
Alternative solutions or features you've considered.

**Additional context**
Add any other context or screenshots about the feature request here.
```

Thank you for contributing to the Real-Time Flight Dashboard DevOps infrastructure!