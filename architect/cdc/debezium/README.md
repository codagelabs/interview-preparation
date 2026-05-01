# Debezium CDC with PostgreSQL and Kafka

This repository contains a full local development environment for implementing Change Data Capture (CDC) using Debezium, PostgreSQL, and Kafka.

The stack automatically provisions the infrastructure, configures the database, and registers the Debezium source connector to immediately start streaming row-level database changes into Kafka.

## Architecture

The environment is orchestrated via Docker Compose and includes the following components:

*   **PostgreSQL**: The source database. It is configured with `wal_level=logical` to enable Debezium to read the Write-Ahead Log (WAL).
*   **Zookeeper**: Coordinates the Kafka broker.
*   **Kafka**: The event streaming platform that receives the CDC events.
*   **Kafka Connect (Debezium)**: The worker that runs the Debezium Postgres Source Connector.
*   **Connector Setup (Init Container)**: An ephemeral container that waits for Kafka Connect to become healthy, automatically registers the Postgres connector using the JSON configuration, and then cleanly exits.
*   **Kafka UI**: A visual web interface to monitor Kafka brokers, topics, and messages.

## Prerequisites

*   [Docker](https://docs.docker.com/get-docker/) & Docker Compose
*   `make` (Optional, but highly recommended for using the provided commands)

## Quick Start

### 1. Start the Stack

Bring up the entire stack. The init container will automatically handle registering the Debezium connector for you.

```bash
make up
```

*(Note: Zookeeper and Kafka are deliberately configured as ephemeral to prevent cluster ID mismatch errors upon restart. Only PostgreSQL data is persisted to a volume).*

### 2. Run Database Migrations & Seed Data

Once the stack is healthy, initialize the `customers` table and insert some seed data. The tables are created with `REPLICA IDENTITY FULL` to ensure Debezium captures the entire before/after state of rows on updates and deletes.

```bash
make migrate
make seed
```

### 3. Verify the CDC Pipeline

You can use the provided verification script to automatically insert a row and listen for the event:

```bash
sh scripts/verify.sh
```

Alternatively, you can test it manually:

1.  Open a terminal and start listening to the topic:
    ```bash
    make consume
    ```
2.  Open a second terminal and insert a random row:
    ```bash
    make pg-insert
    ```
3.  Watch the first terminal instantly print the CDC JSON event!

## Important URLs & Ports

*   **PostgreSQL**: `localhost:5433` (User: `postgres`, Password: `postgres`, DB: `inventory`)
*   **Kafka (Host Access)**: `localhost:29092`
*   **Kafka Connect REST API**: `http://localhost:8083`
*   **Kafka UI**: [http://localhost:8090](http://localhost:8090)

## Internal Kafka Topics

If you run `make topics`, you will notice several topics created alongside your `cdc.public.customers` topic. These are essential for the system's operation:

*   `debezium_connect_configs`: Stores the connector configuration so it survives restarts.
*   `debezium_connect_offsets`: Stores the exact PostgreSQL WAL offset (LSN) Debezium has processed, preventing data loss or duplication on restart.
*   `debezium_connect_statuses`: Tracks whether the connector is running, paused, or failed.
*   `__debezium-heartbeat.cdc`: A heartbeat ping sent every 5 seconds (configured in the JSON) to keep the Postgres WAL advancing and prevent disk space exhaustion.
*   `__consumer_offsets`: Standard Kafka topic tracking consumer group read positions.

## Troubleshooting

**Kafka "InconsistentClusterIdException" on restart:**
If you manually tear down containers or delete volumes inconsistently, Kafka and Zookeeper might disagree on the cluster ID. To fix this, run a clean reset which wipes the volumes and forces them to re-sync:

```bash
make reset
```
*(After a reset, remember to re-run `make migrate` as the Postgres volume is also wiped).*
