# URL Shortener 

## What is URL Shortener ?
A URL shortener converts a long URL into a short, unique alias and redirects users to the original URL when short URL is accessed.

## Functional Requirement
1. Shorten a long URL
2. Redirect short URL -> long URL
3. Optional custom alias, expiry, analytics

## Non-Functional Requirement
1. Low Latency Redirects 
2. High Availability 

## calculation 
### Write Throughput
Assume 100 shorten urls per sec
```
100 writes/sec = 8,640,000 URLs/day
```
Peak Traffic 
10 times more per second 
```
100 × 10 = 1,000 writes/sec
```

✅ Needs horizontally scalable write path
✅ ID generation must be fast & collision-free

### Read Throughput
Average redirects per URL per day: 50
```
8,640,000 URLs/day * 50 redirects = 432,000,000 redirects/day
```
Reads per second 
```
432,000,000 redirects/day / 60*60*24 ~ 5000/RPS
```
✅ Extremely read-heavy system

### Storage Calculation

Assume that retention period is of 5 Years
```
8,640,000 × 365 × 5 ≈ 15.8 billion URLs
```
#### Per-record size (approx)
- Field	Size
- short_code	8 bytes
- long_url (avg)	200 bytes
- metadata + indexes	100 bytes

```
15.8B × 300 bytes ≈ 4.7 TB

```

### Cache (Redis) Sizing

Assumptions:

80% traffic hits top 20% URLs

Cache hot URLs only
```
Hot URLs ≈ 20% of 15.8B ≈ 3.16B
```

If we cache only top 10% active URLs:
```
~1.5B × 300 bytes ≈ 450 GB
```

👉 Use Redis Cluster
👉 Eviction + TTL is mandatory

## High Level Architecture
```
Client
  |
  v
+------------------+
|  API Gateway     |
+------------------+
  |        |
  |        v
  |   +-----------+
  |   |  Cache    |  (Redis)
  |   +-----------+
  |        |
  v        v
+------------------+
| URL Shortener    |
| Service          |
+------------------+
  |
  v
+------------------+
| Database         |
| (KV / SQL)       |
+------------------+

```
### Data Model (Simple)

```
short_code (PK) | long_url | created_at | expires_at | click_count

```

## Low Level Design (API Design)
### Create Short URL
```
POST /shorten
```
Request:
```
{
  "long_url": "https://example.com/product/123"
}
```

Response:
```
{
  "short_url": "https://sho.rt/aZ9kL2"
}
```

Redirect URL
```
GET /{short_code}
```
Flow:
```
Lookup short_code in cache

On miss → fetch from DB and update cache

Redirect (HTTP 302)

Publish analytics event asynchronously
```



## Short Code Generation Strategies – URL Shortener

This document describes **all possible ways to generate short codes** for a URL shortener, along with their **pros, cons, and suitability** for a system handling **~100 URL shortens per second** and **billions of records over time**.



### 1. Auto-Increment ID + Base62 Encoding (Recommended)

### Description
A unique numeric ID is generated using a database sequence or distributed ID generator and then encoded using Base62 (`a–zA–Z0–9`).

Example:
ID: 12583921 → Base62: aZ9kL


### Pros
- No collisions
- Very fast generation
- Short, URL-safe codes
- Simple decoding and debugging
- Scales well to trillions of URLs

### Cons
- Sequential and predictable
- Requires distributed ID generation at scale
- Slight information leakage (growth pattern)

### Suitability at Current Scale
✅ **Best choice**  
- Use **8-character Base62**
- Combine with **Snowflake ID / segment allocator**

---

## 2. Random String Generation

### Description
Generate a random alphanumeric string (Base62) of fixed length and check for collisions in the database.

Example:
       X9fK2A


### Pros
- Simple to implement
- Non-predictable codes
- No dependency on sequences

### Cons
- Collision probability increases with scale
- Requires retry and DB existence checks
- Higher write latency

### Suitability at Current Scale
⚠️ Possible but not optimal  
- Requires **8–9 character codes**
- Extra DB read per write

---

## 3. Hash-Based (MD5 / SHA)

### Description
Hash the original URL and truncate the hash to generate a short code.

Example:
SHA256(longURL) → first 8 chars

### Pros
- Same URL always maps to same short code
- No ID generator required

### Cons
- Collisions still possible
- Hard to resolve conflicts
- Not reversible
- Poor control over code length

### Suitability at Current Scale
❌ **Not recommended**
- Collision handling becomes complex at billions of URLs

---

## 4. Pre-Generated Code Pool

### Description
Pre-generate a large pool of unique short codes offline and allocate them at runtime.

### Pros
- No collisions
- Very fast allocation
- No DB existence check

### Cons
- Requires storing and managing large pools
- Operationally complex
- Risk of pool exhaustion

### Suitability at Current Scale
⚠️ Viable but operationally heavy  
- Used mainly by very large systems

---

## 5. Timestamp + Counter Encoding

### Description
Generate IDs using a combination of timestamp and per-node counter, then encode to Base62.

Example:
(timestamp << bits) + counter → Base62


### Pros
- Collision-free
- Time-sortable IDs
- High throughput
- No central DB sequence

### Cons
- Requires clock synchronization
- Slightly complex logic

### Suitability at Current Scale
✅ Good alternative  
- Works well for **100–1000 URLs/sec**
- Similar to Snowflake ID strategy

---

## Comparison Summary

| Strategy | Collisions | Code Length | Complexity | Recommendation |
|--------|-----------|-------------|------------|----------------|
| Auto-Increment + Base62 | None | 8 | Low | ⭐ Best |
| Random String | Medium | 8–9 | Low | ⚠️ |
| Hash-Based | Medium | 8 | Low | ❌ |
| Pre-Generated Pool | None | 6–8 | High | ⚠️ |
| Timestamp + Counter | None | 8 | Medium | ⭐ |

---

## Final Recommendation

For a system generating **100 short URLs per second** and storing **billions of records**, the most reliable and scalable approach is:

> **Base62 encoding over a distributed unique ID (Snowflake or timestamp+counter)**

This ensures **collision-free IDs, predictable performance, and long-term scalability**.

---

## Next Steps
- Implement Base62 encoder in Go
- Add Snowflake-style ID generator
- Design DB sharding using short_code
- Add rate limiting and abuse protection
