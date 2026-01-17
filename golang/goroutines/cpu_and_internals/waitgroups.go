package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// Task represents a URL to be processed
type Task struct {
	URL string
	ID  int
}

// Result represents the HTTP request result
type Result struct {
	URL      string
	Status   string
	Error    error
	Duration time.Duration
}

// Config holds application configuration
type Config struct {
	URLs           []string
	Workers        int
	RequestTimeout time.Duration
	ProcessTimeout time.Duration
}

// DefaultConfig provides default configuration values
var DefaultConfig = Config{
	URLs: []string{
		"https://facebook1.com",
		"https://golang.org",
		"https://amazon.com",
		"https://apple.com",
	},
	Workers:        2,
	RequestTimeout: 10 * time.Second,
	ProcessTimeout: 30 * time.Second,
}

// HTTPProcessor handles concurrent HTTP requests
type HTTPProcessor struct {
	config  Config
	tasks   chan Task
	results chan Result
}

// NewHTTPProcessor creates a new HTTPProcessor instance
func NewHTTPProcessor(config Config) *HTTPProcessor {
	return &HTTPProcessor{
		config:  config,
		tasks:   make(chan Task, len(config.URLs)),
		results: make(chan Result, len(config.URLs)),
	}
}

// processRequest handles individual HTTP requests
func (hp *HTTPProcessor) processRequest(ctx context.Context, task Task) Result {
	client := &http.Client{
		Timeout: hp.config.RequestTimeout,
	}

	start := time.Now()

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", task.URL, nil)
	if err != nil {
		return Result{
			URL:   task.URL,
			Error: fmt.Errorf("failed to create request: %w", err),
		}
	}

	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		return Result{
			URL:      task.URL,
			Error:    fmt.Errorf("request failed: %w", err),
			Duration: duration,
		}
	}
	defer resp.Body.Close()

	return Result{
		URL:      task.URL,
		Status:   resp.Status,
		Duration: duration,
	}
}

// worker processes tasks from the task channel
func (hp *HTTPProcessor) worker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case task, ok := <-hp.tasks:
			if !ok {
				return
			}
			result := hp.processRequest(ctx, task)
			hp.results <- result

		case <-ctx.Done():
			return
		}
	}
}

// resultCollector collects and processes results
func (hp *HTTPProcessor) resultCollector(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case result, ok := <-hp.results:
			if !ok {
				return
			}
			if result.Error != nil {
				log.Printf("Error processing %s: %v\n", result.URL, result.Error)
			} else {
				log.Printf("Success: %s - Status: %s (took: %v)\n",
					result.URL, result.Status, result.Duration)
			}

		case <-ctx.Done():
			return
		}
	}
}

// Process starts the HTTP processing
func (hp *HTTPProcessor) Process(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, hp.config.ProcessTimeout)
	defer cancel()

	var workersWg sync.WaitGroup
	var collectorWg sync.WaitGroup

	// Start workers
	for i := 0; i < hp.config.Workers; i++ {
		workersWg.Add(1)
		go hp.worker(ctx, &workersWg)
	}

	// Start result collector
	collectorWg.Add(1)
	go hp.resultCollector(ctx, &collectorWg)

	// Send tasks
	for i, url := range hp.config.URLs {
		select {
		case hp.tasks <- Task{URL: url, ID: i}:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	// Close tasks channel and wait for workers
	close(hp.tasks)
	workersWg.Wait()

	// Close results channel and wait for collector
	close(hp.results)
	collectorWg.Wait()

	return nil
}

func main() {
	// Set up logging
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	// Create processor_unused with default config
	processor := NewHTTPProcessor(DefaultConfig)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start processing
	log.Println("Starting HTTP processing...")
	if err := processor.Process(ctx); err != nil {
		log.Fatalf("Processing failed: %v", err)
	}
	log.Println("Processing completed successfully")
}
