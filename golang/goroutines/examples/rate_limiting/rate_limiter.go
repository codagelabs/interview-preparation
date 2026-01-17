package main

import (
	"fmt"
	"time"
)

// RateLimiter controls the rate of operations using a buffered channel
type RateLimiter struct {
	tokens chan struct{} // Buffered channel to hold tokens
}

// NewRateLimiter creates a new rate limiter with specified capacity
func NewRateLimiter(maxRequests int, refillInterval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		tokens: make(chan struct{}, maxRequests),
	}

	// Initially fill the token bucket
	for i := 0; i < maxRequests; i++ {
		rl.tokens <- struct{}{}
	}

	// Start token refill goroutine
	go func() {
		ticker := time.NewTicker(refillInterval)
		defer ticker.Stop()

		for range ticker.C {
			select {
			case rl.tokens <- struct{}{}:
				// Added a token
			default:
				// Bucket is full, skip
			}
		}
	}()

	return rl
}

// Allow blocks until a token is available
func (rl *RateLimiter) Allow() {
	<-rl.tokens
}

func main() {
	// Create a rate limiter: 3 requests per second
	limiter := NewRateLimiter(3, time.Second)

	// Simulate multiple requests
	for i := 1; i <= 10; i++ {
		go func(requestID int) {
			fmt.Printf("Request %d waiting for rate limiter...\n", requestID)
			limiter.Allow()
			fmt.Printf("Request %d processed at %v\n", requestID, time.Now().Format("15:04:05"))
		}(i)
	}

	// Wait to see the results
	time.Sleep(5 * time.Second)
}
