# Rate Limiting with Buffered Channels

This example demonstrates a simple rate limiting implementation using Go's buffered channels. Rate limiting is useful for controlling the pace of operations, such as API requests or resource-intensive tasks.

## How It Works

The rate limiter uses a buffered channel as a token bucket:
- The channel has a fixed capacity (maximum number of concurrent operations)
- Each token represents permission to perform an operation
- When the bucket is empty, operations must wait for new tokens
- Tokens are automatically refilled at a fixed interval

## Key Components

1. **RateLimiter struct**:
   - Contains a buffered channel for tokens
   - Controls access to resources

2. **Token System**:
   - Empty struct (`struct{}`) used as tokens (zero memory overhead)
   - Channel buffer size determines maximum concurrent operations

3. **Automatic Refill**:
   - Background goroutine refills tokens periodically
   - Uses `time.Ticker` for consistent intervals

## Usage Example

```go
// Create a rate limiter allowing 3 operations per second
limiter := NewRateLimiter(3, time.Second)

// Use the rate limiter
limiter.Allow() // Blocks until a token is available
// Perform your rate-limited operation here
```

## Sample Output 