#!/bin/bash

echo "Checking connector status..."
curl -s http://localhost:8083/connectors/inventory-connector/status
echo ""
echo ""

echo "Inserting test row into postgres..."
TEST_ID=$((RANDOM % 9000 + 1000))
docker exec debezium-postgres psql -U postgres -d inventory -c "INSERT INTO public.customers (name, email) VALUES ('CDC Test $TEST_ID', 'cdc$TEST_ID@test.com');"

echo "Waiting 3 seconds for CDC event..."
sleep 3

echo "Checking Kafka topic for the event..."
docker exec debezium-kafka /kafka/bin/kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 \
  --topic cdc.public.customers \
  --from-beginning \
  --max-messages 1 \
  --timeout-ms 10000

echo ""
echo "Verification complete!"
