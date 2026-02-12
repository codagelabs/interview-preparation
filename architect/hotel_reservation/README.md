# Hotel Reservation System Architecture

## 1. Problem Statement
Design a scalable hotel reservation system (similar to Booking.com, Airbnb, or OYO) that allows users to search for hotels, view details, and book rooms. The system must effectively handle high concurrency during peak seasons, ensure zero double-bookings (strict consistency), and provide a low-latency search experience.

## 2. Requirement Gathering

### a. Functional Requirements
*   **User Users:**
    *   Search for hotels based on location, check-in/out dates, and filters (price, capacity, amenities).
    *   View hotel details (images, description, reviews).
    *   Book a room (Reserve -> Pay -> Confirm).
    *   Cancel a booking.
    *   View booking history.
*   **Hotel Managers / Admin:**
    *   Onboard new hotels/rooms.
    *   Manage room inventory and availability.
    *   Update pricing.
*   **System Internal:**
    *   Handle payments securely.
    *   Send notifications (Email/SMS).
    *   Housekeeping tasks (releasing temporary holds if payment fails).

### b. Non-Functional Requirements
*   **Consistency:** **Critical**. We cannot double-book a room. ACID properties are mandatory for the booking flow.
*   **Availability:** High. The system (especially Search) should always be up. For Booking, we might trade slight availability for consistency (CP over AP in CAP theorem for the booking shard).
*   **Latency:**
    *   Search: Extremely low (< 200ms).
    *   Booking: Acceptable to take 1-2 seconds as it involves transactions/payments.
*   **Scalability:** Must handle traffic spikes during holidays.
*   **Reliability:** Booking data must never be lost.

### c. Back-of-the-Envelope Calculations
*   **Traffic Assumptions:**
    *   100 Million Daily Active Users (DAU).
    *   500,000 Hotels globally.
    *   5 Million Rooms.
    *   Average user queries search 10 times before booking.
    *   Read:Write Ratio = 1000:1 (Search heavy).
*   **QPS (Queries Per Second):**
    *   Search: 100M users * 10 queries / 86400 ?? ~11,500 QPS. Peak load (x5) ?? ~60,000 QPS.
*   **Storage Estimates:**
    *   **Hotel Metadata:** 500k * 10KB = 5 GB (Fit in memory/cache).
    *   **Images:** stored in Blob Storage (S3), not DB.
    *   **Bookings:** 1M bookings/day * 1KB = 1GB/day ?? 3.65TB/10 year. (Manageable with RDBMS partitioning).

## 3. High-Level Design (HLD)

The system is designed using a **Microservices Architecture** to scale Read (Search) and Write (Booking) paths independently.

### Core Components
1.  **CDN (Content Delivery Network):** Caches static assets (images, CSS, JS) close to the user.
2.  **Load Balancers:** Distribute traffic across service instances.
3.  **API Gateway:** Entry point for clients. Handles Authentication, Rate Limiting, and Routing.
4.  **Services:**
    *   **Search Service:** Handles complex queries. Powered by **Elasticsearch** (or Solr) for geospatial and full-text search.
    *   **Hotel Service:** Manages static content (Hotel info, Room types). Heavily cached.
    *   **Booking Service:** Core engine. Manages reservations and inventory logic.
    *   **User Service:** User profiles and authentication.
    *   **Payment Service:** Integrates with gateways (Stripe/PayPal).
5.  **Data Storage:**
    *   **RDBMS (PostgreSQL/MySQL):** Primary source of truth for **Inventory** and **Bookings**. Essential for ACID transactions and row locking.
    *   **NoSQL (Elasticsearch):** Optimized specifically for the Search Service (read-only replica of hotel data).
    *   **Redis (Cache):** Caches Hot Hotel data and User sessions.
    *   **Message Queue (Kafka/RabbitMQ):** For async tasks (Notifications, Analytics).

### Architecture Diagram Flow
```
[User Client] 
      |
      v
 [CDN / LB] 
      |
      v
 [API Gateway]
      |
      +------------------------+---------------------+-------------------+
      |                        |                     |                   |
 [Search Service]       [Booking Service]      [Hotel Service]     [Payment Service]
      |                        |                     |                   |
 [Elasticsearch]         [RDBMS (Master)]      [Redis + RDBMS]     [Payment Gateway]
      ^                        |
      | (Sync)                 |
[Inventory Sync Worker] <------+
```

## 4. Low-Level Design (LLD)

### Database Schema (Relational)

**1. Hotel Table**
*   `id` (PK), `name`, `address_json`, `location_lat`, `location_long`, `description`

**2. Room_Type Table**
*   `id` (PK), `hotel_id` (FK), `name` (e.g., 'Deluxe'), `base_price`, `capacity`

**3. Room Table** (Physical Rooms)
*   `id` (PK), `hotel_id`, `room_type_id`, `room_number`

**4. Inventory Table** (Crucial for Availability)
*   `id` (PK), `room_type_id`, `date`, `total_inventory`, `reserved_count`
*   *Note: Instead of tracking every single room ID for availability, we often track counts per type per date for scalability.*

**5. Booking Table**
*   `id` (PK), `user_id`, `hotel_id`, `room_type_id`, `start_date`, `end_date`, `status` (PENDING, CONFIRMED, CANCELLED, REFUNDED), `created_at`

### Concurrency Control (Preventing Double Bookings)
When User A and User B try to book the last room simultaneously:

**Approach: Pessimistic Locking (Recommended for High Consistency)**
To book a room for a specific date range:
1.  Start Transaction.
2.  `SELECT * FROM inventory WHERE room_type_id = X AND date BETWEEN start AND end FOR UPDATE;`
3.  Check if `total_inventory - reserved_count > 0` for all days.
4.  If yes, `UPDATE inventory SET reserved_count = reserved_count + 1`.
5.  Insert into `Booking` table with status `PENDING`.
6.  Commit Transaction.
7.  *If any step fails, Rollback.*

**Why not Optimistic Locking?**
Optimistic locking works well for low contention. For high contention (hot hotels), it causes many user failures/retries. Pessimistic locking (row-level) guarantees the slot once the query returns.

## 5. Explain Flow

**Flow: The Booking Journey**

1.  **Search Phase:**
    *   User sends request: "Hotels in NY, Dec 25-30".
    *   **Search Service** queries **Elasticsearch**.
    *   Returns list of hotels with *cached* prices.
2.  **Details Phase:**
    *   User clicks Hotel A.
    *   **Hotel Service** fetches details from **Redis** (Cache Hit).
    *   Client calls **Booking Service** to get *real-time* availability (Inventory DB check).
3.  **Booking Phase (The Transaction):**
    *   User clicks "Book Now".
    *   **Booking Service** initiates the DB Transaction (as described in LLD).
    *   Locks Inventory rows -> Decrements availability -> Creates PENDING booking.
    *   A temporal timer (e.g., 10 mins) starts (via Redis or internal scheduler).
4.  **Payment Phase:**
    *   Client is redirected to Payment Gateway.
    *   On success, Payment Gateway calls webhook -> **Payment Service**.
    *   **Booking Service** updates Booking status: PENDING -> CONFIRMED.
    *   If payment fails or times out, the "Housekeeper" processes release the Inventory lock (decrement `reserved_count`).
5.  **Notification:**
    *   On CONFIRMED status, an event is pushed to **Kafka**.
    *   **Notification Service** consumes event and emails the user the receipt.

## 6. Feature Development (Extensibility)

*   **1. Dynamic Pricing:**
    *   Introduce a **Pricing Service**. Before showing results, Search/Hotel services query Pricing Service.
    *   Logic: If `reserved_count` > 80% or search volume is high, apply multiplier to `base_price`.
*   **2. Reviews and Ratings:**
    *   Use a NoSQL store (Cassandra/DynamoDB) because reviews are append-only huge datasets and eventual consistency is fine.
    *   Rating Aggregation: Compute averages asynchronously (e.g., Spark job or serialized stream processing) and store in Hotel DB.
*   **3. Hotel Recommendations:**
    *   Capture user click streams via **Analytics Service**.
    *   Feed data to ML models to rank "Recommended for You" in search results.

---

## Interview Wrap-Up: Justifications & Reasons

**Why separate Search and Booking DBs?**
*   **Reason:** Search requires complex queries (fuzzy match, geo-spatial) which SQL is slow at. Booking requires strict transactions (ACID) which NoSQL/Search engines are weak at.
*   **CQRS Pattern**: We effectively divide the system into Command (Booking - SQL) and Query (Search - Elastic) responsibilities to pick the best tool for each job.

**Why Redis?**
*   **Reason:** Hotel data (Names, Amenities) is "Read-Mostly". Fetching from SQL every time is wasteful. Redis provides <5ms access, handling the massive read load.

**How to handle "Hold" expiration?**
*   **Reason:** When a user starts booking, we hold the room. We can use Redis keys with TTL (Time To Live) matching the hold duration (e.g., 10 mins). When keys expire, a listener can trigger a DB rollback if the booking wasn't confirmed.
