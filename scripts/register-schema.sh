#!/bin/bash

# Register Avro schema with Schema Registry

SCHEMA_REGISTRY_URL="http://localhost:8081"
SUBJECT="flights-value"
SCHEMA_FILE="schemas/flight-data.avsc"

echo "Registering schema for subject: $SUBJECT"

# Read schema file and escape quotes
SCHEMA=$(cat $SCHEMA_FILE | jq -c . | sed 's/"/\\"/g')

# Register schema
curl -X POST \
  -H "Content-Type: application/vnd.schemaregistry.v1+json" \
  --data "{\"schema\":\"$SCHEMA\"}" \
  "$SCHEMA_REGISTRY_URL/subjects/$SUBJECT/versions"

echo ""
echo "Schema registration complete!"