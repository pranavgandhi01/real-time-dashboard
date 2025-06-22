# Scripts

## Available Scripts

### run-microservices-tests.sh
Runs all Go microservice tests with coverage reports.

```bash
./scripts/run-microservices-tests.sh
```

### Token Generation
```bash
# Generate secure WebSocket token
python3 -c "import secrets; print(secrets.token_hex(32))"
```

## Usage
All scripts should be run from the project root directory.