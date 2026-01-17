# Change Data Capture(CDC)
## What is CDC ?
 - A process that indentify and tarck changes made in database
 - It captures INSERT, UPDATE and DELETE operations.
 - Delivers only this changes to other systems in near to real time

## How CDC Works (High Level)

1. Source Database
Data changes happen (INSERT / UPDATE / DELETE)

2. CDC Mechanism
Reads database transaction logs (binlog, WAL, redo logs)
Or uses triggers / timestamps (less preferred)

3. Streaming / Sync Layer
Publishes changes to systems like Kafka, queues, or APIs

4. Target Systems
  - Data warehouse (Snowflake, BigQuery)
  - Search engine (OpenSearch, Elasticsearch)
  - Cache (Redis)
  - Another database or service

## CDC Architecture (Line & Box Diagram)

```
+-------------------+
|  Source Database  |
| (INSERT/UPDATE/   |
|      DELETE)      |
+---------+---------+
          |
          | Transaction Logs
          | (WAL / Binlog)
          v
+-------------------+
|  CDC Connector    |
| (Debezium / DMS)  |
+---------+---------+
          |
          | Change Events
          v
+-------------------+
| Streaming Platform|
| (Kafka / Queue)   |
+---------+---------+
          |
          | Fan-out
          v
+-------------------+     +-------------------+     +-------------------+
| Data Warehouse    |     | Search Index     |     | Cache / Services  |
| (Analytics)       |     | (Elastic / OS)   |     | (Redis / APIs)    |
+-------------------+     +-------------------+     +-------------------+

```

### Notes
- CDC reads **transaction logs**, not live tables.
- Events are delivered in **correct order**.
- Multiple consumers can independently process the same change events.
- Supports **real-time sync**, **analytics**, and **event-driven systems**.

## Advantages of CDC (Change Data Capture)

### 🚀 1. Near Real-Time Data Sync
- Changes are captured immediately  
- Ideal for real-time analytics, dashboards, and monitoring  

### 💾 2. Low Database Load
- No full-table scans  
- Reads transaction logs instead of production tables  
- Safer for high-traffic systems  

### 🔁 3. Incremental Data Movement
- Only changed data is transferred  
- Saves network, storage, and compute cost  

### 🧩 4. Loose Coupling Between Systems
- Source system does not need to know about consumers  
- Enables event-driven architecture  

### 📊 5. Accurate Data Replication
- Preserves the order of operations  
- Captures exact changes (before and after state)  

### 🔄 6. Enables Multiple Use Cases
- Data migration  
- Audit and compliance  
- Cache invalidation  
- Search index updates  
- Microservice data sharing

## When Should You Use CDC (Change Data Capture)?

### ✔ Migrating Data from Legacy → New System
CDC is ideal for migrations where downtime must be minimized.  
It allows you to:
- Perform an initial bulk data load
- Continuously capture and apply ongoing changes from the legacy system
- Keep both systems in sync until the final cutover

This ensures a **zero-downtime or near zero-downtime migration**.

---

### ✔ Syncing Database → Search Index
CDC enables real-time synchronization between a database and a search engine.
- Every insert, update, or delete is immediately propagated
- Search indexes remain consistent with the source database
- Eliminates expensive re-indexing jobs

This is commonly used for **Elasticsearch/OpenSearch-backed search features**.

---

### ✔ Real-Time Analytics
CDC streams changes as they happen, making it suitable for:
- Live dashboards
- Operational monitoring
- Fraud detection and alerting

Instead of waiting for batch ETL jobs, analytics systems always work with **fresh data**.

---

### ✔ Microservices Sharing Data
In microservice architectures, direct database access between services is discouraged.
CDC helps by:
- Publishing data changes as events
- Allowing services to maintain their own read models
- Reducing tight coupling and synchronous dependencies

This enables **event-driven and scalable system design**.

---

### ✔ Audit Trails and History Tracking
CDC captures every data change along with metadata such as:
- Operation type (INSERT, UPDATE, DELETE)
- Timestamp
- Before and after values

This makes it useful for:
- Compliance and regulatory requirements
- Debugging and root-cause analysis
- Reconstructing historical states of data
