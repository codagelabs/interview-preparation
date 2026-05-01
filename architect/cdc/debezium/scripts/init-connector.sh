#!/bin/bash
# Runs inside the 'connector-setup' docker container

echo "Waiting for Kafka Connect to become available..."
until curl -s http://connect:8083/connectors > /dev/null 2>&1; do 
  sleep 2
done

echo "Registering connector..."
curl -s -X POST -H "Content-Type: application/json" -d @/register-postgres.json http://connect:8083/connectors || true

echo ""
echo "Verifying connector status..."
sleep 3
curl -s http://connect:8083/connectors/inventory-connector/status

echo ""
echo "Connector setup complete! Exiting."
