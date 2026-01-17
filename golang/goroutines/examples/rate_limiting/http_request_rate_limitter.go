package main

import (
	"fmt"
	"net/http"
	"time"
)

// RateLimiter controls the rate of HTTP requests
type RateLimiter struct {
	tokens chan struct{}
	client *http.Client
}

// NewRateLimiter creates a new rate limiter with specified capacity
func NewRateLimiter(maxRequests int, refillInterval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		tokens: make(chan struct{}, maxRequests),
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
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
				fmt.Println("Token added")
			default:
				// Bucket is full, skip
			}
		}
	}()

	return rl
}

// MakeRequest performs a rate-limited HTTP request
func (rl *RateLimiter) MakeRequest(url string) (*http.Response, error) {
	// Wait for available token
	<-rl.tokens
	fmt.Printf("Making request to %s at %v\n", url, time.Now().Format("15:04:05"))
	return rl.client.Get(url)
}

func main() {
	// Create a rate limiter: 2 requests per second
	limiter := NewRateLimiter(2, time.Second)

	// Example URLs to test
	urls := []string{
		"https://api.github.com",
		"https://api.github.com/users",
		"https://api.github.com/repos",
		"https://api.github.com/gists",
	}

	// Make multiple requests
	for i, url := range urls {
		go func(requestID int, requestURL string) {
			fmt.Printf("Request %d waiting for rate limiter...\n", requestID)

			resp, err := limiter.MakeRequest(requestURL)
			if err != nil {
				fmt.Printf("Request %d failed: %v\n", requestID, err)
				return
			}
			defer resp.Body.Close()

			fmt.Printf("Request %d completed with status: %s\n",
				requestID, resp.Status)
		}(i+1, url)
	}

	// Wait to see the results
	time.Sleep(5 * time.Second)
}

// Helper function to demonstrate error handling
func handleResponse(resp *http.Response, err error) {
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Response status: %s\n", resp.Status)
}
