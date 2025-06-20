# Utility Scripts

## generate-token.sh
Generates secure WebSocket authentication tokens and updates all environment files.

```bash
./scripts/generate-token.sh
```

## setup-env.sh
Sets up production environment files and generates secure tokens.

```bash
./scripts/setup-env.sh
```

## Manual Token Generation
```bash
python3 -c "import secrets; print(secrets.token_hex(32))"
```