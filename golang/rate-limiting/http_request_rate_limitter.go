package main

import (
	"fmt"
	"net/http"
	"time"
)

// SimpleRateLimiter uses a buffered channel to limit concurrent requests
type SimpleRateLimiter struct {
	limiter chan struct{}
	client  *http.Client
}

// NewSimpleRateLimiter creates a rate limiter with a specified concurrency limit
func NewSimpleRateLimiter(limit int) *SimpleRateLimiter {
	return &SimpleRateLimiter{
		limiter: make(chan struct{}, limit), // Buffered channel with limit
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// MakeRequest performs a rate-limited HTTP request
func (rl *SimpleRateLimiter) MakeRequest(url string) (*http.Response, error) {
	rl.limiter <- struct{}{} // Acquire token
	defer func() {
		<-rl.limiter // Release token after request completes
	}()

	fmt.Printf("Making request to %s at %v\n", url, time.Now().Format("15:04:05"))
	return rl.client.Get(url)
}

func main() {
	// Create a rate limiter with concurrent limit of 2
	limiter := NewSimpleRateLimiter(2)

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
