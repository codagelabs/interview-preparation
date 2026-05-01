#!/bin/bash

# Start the stack
echo "Starting Debezium stack..."
docker compose up -d

# Wait for Kafka Connect to be ready
echo "Waiting for Kafka Connect..."
until curl -s http://localhost:8083/connectors > /dev/null; do
  sleep 5
done
echo "Kafka Connect is ready!"

# Register the connector
echo "Registering Postgres connector..."
curl -X POST -H "Content-Type: application/json" -d @connectors/register-postgres.json http://localhost:8083/connectors

echo ""
echo "Setup complete!"
