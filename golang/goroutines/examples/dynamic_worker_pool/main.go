package main

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Task represents a unit of work
type Task struct {
	ID       int
	Load     int // Simulated load (milliseconds)
	Priority int
}

// Result represents the task processing result
type Result struct {
	TaskID    int
	Duration  time.Duration
	WorkerID  int
	Timestamp time.Time
}

// DynamicPool manages a pool of workers that scales based on workload
type DynamicPool struct {
	// Channels
	tasks   chan Task
	results chan Result

	// Pool configuration
	minWorkers     int32
	maxWorkers     int32
	currentWorkers int32

	// Metrics
	queueLength    int32
	processingTime atomic.Value // stores *time.Duration

	// Scaling configuration
	scaleUpThreshold   float64 // queue utilization threshold to scale up
	scaleDownThreshold float64 // queue utilization threshold to scale down
	scaleCooldown      time.Duration

	// Control
	mu       sync.RWMutex
	shutdown chan struct{}
	metrics  *PoolMetrics
}

// PoolMetrics tracks pool performance
type PoolMetrics struct {
	AverageProcessingTime time.Duration
	QueueUtilization      float64
	ActiveWorkers         int32
	TotalTasksProcessed   int64
	LastScalingEvent      time.Time
}

// NewDynamicPool creates a new dynamic worker pool
func NewDynamicPool(config Config) *DynamicPool {
	dp := &DynamicPool{
		tasks:              make(chan Task, config.QueueSize),
		results:            make(chan Result, config.QueueSize),
		minWorkers:         int32(config.MinWorkers),
		maxWorkers:         int32(config.MaxWorkers),
		scaleUpThreshold:   0.75, // Scale up when queue is 75% full
		scaleDownThreshold: 0.25, // Scale down when queue is 25% full
		scaleCooldown:      time.Second * 5,
		shutdown:           make(chan struct{}),
		metrics:            &PoolMetrics{},
	}

	dp.processingTime.Store(new(time.Duration))
	return dp
}

// worker processes tasks and automatically scales based on load
func (dp *DynamicPool) worker(ctx context.Context, workerID int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer atomic.AddInt32(&dp.currentWorkers, -1)

	for {
		select {
		case task, ok := <-dp.tasks:
			if !ok {
				return
			}

			start := time.Now()

			// Simulate task processing
			time.Sleep(time.Duration(task.Load) * time.Millisecond)

			duration := time.Since(start)

			// Update processing metrics
			dp.updateMetrics(duration)

			// Send result
			dp.results <- Result{
				TaskID:    task.ID,
				Duration:  duration,
				WorkerID:  workerID,
				Timestamp: time.Now(),
			}

			atomic.AddInt32(&dp.queueLength, -1)

		case <-ctx.Done():
			return
		}
	}
}

// updateMetrics updates pool performance metrics
func (dp *DynamicPool) updateMetrics(duration time.Duration) {
	// Update average processing time
	current := dp.processingTime.Load().(*time.Duration)
	if *current == 0 {
		dp.processingTime.Store(&duration)
	} else {
		newDuration := (*current + duration) / 2
		dp.processingTime.Store(&newDuration)
	}

	atomic.AddInt64(&dp.metrics.TotalTasksProcessed, 1)
}

// scaleWorkers adjusts the number of workers based on workload
func (dp *DynamicPool) scaleWorkers(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			dp.evaluateScaling(ctx, wg)
		case <-ctx.Done():
			return
		}
	}
}

// evaluateScaling decides whether to scale up or down
func (dp *DynamicPool) evaluateScaling(ctx context.Context, wg *sync.WaitGroup) {
	dp.mu.Lock()
	defer dp.mu.Unlock()

	if time.Since(dp.metrics.LastScalingEvent) < dp.scaleCooldown {
		return
	}

	queueUtilization := float64(atomic.LoadInt32(&dp.queueLength)) / float64(cap(dp.tasks))
	currentWorkers := atomic.LoadInt32(&dp.currentWorkers)

	switch {
	case queueUtilization >= dp.scaleUpThreshold && currentWorkers < dp.maxWorkers:
		// Scale up
		workersToAdd := min(dp.maxWorkers-currentWorkers, 2) // Add up to 2 workers at a time
		for i := int32(0); i < workersToAdd; i++ {
			wg.Add(1)
			atomic.AddInt32(&dp.currentWorkers, 1)
			workerID := int(atomic.LoadInt32(&dp.currentWorkers))
			go dp.worker(ctx, workerID, wg)
		}
		dp.metrics.LastScalingEvent = time.Now()
		log.Printf("Scaled up to %d workers (Queue utilization: %.2f%%)\n",
			atomic.LoadInt32(&dp.currentWorkers), queueUtilization*100)

	case queueUtilization <= dp.scaleDownThreshold && currentWorkers > dp.minWorkers:
		// Scale down
		workersToRemove := min(currentWorkers-dp.minWorkers, 1) // Remove 1 worker at a time
		atomic.AddInt32(&dp.currentWorkers, -workersToRemove)
		dp.metrics.LastScalingEvent = time.Now()
		log.Printf("Scaled down to %d workers (Queue utilization: %.2f%%)\n",
			atomic.LoadInt32(&dp.currentWorkers), queueUtilization*100)
	}
}

// Start begins processing tasks and managing workers
func (dp *DynamicPool) Start(ctx context.Context) error {
	var wg sync.WaitGroup

	// Start initial workers
	for i := 0; i < int(dp.minWorkers); i++ {
		wg.Add(1)
		atomic.AddInt32(&dp.currentWorkers, 1)
		go dp.worker(ctx, i+1, &wg)
	}

	// Start scaling manager
	go dp.scaleWorkers(ctx, &wg)

	// Start metrics reporter
	go dp.reportMetrics(ctx)

	return nil
}

// Submit adds a task to the pool
func (dp *DynamicPool) Submit(task Task) error {
	select {
	case dp.tasks <- task:
		atomic.AddInt32(&dp.queueLength, 1)
		return nil
	default:
		return fmt.Errorf("task queue is full")
	}
}

// reportMetrics periodically logs pool metrics
func (dp *DynamicPool) reportMetrics(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			avgTime := dp.processingTime.Load().(*time.Duration)
			log.Printf("Pool Metrics - Workers: %d, Queue Length: %d, Avg Processing Time: %v\n",
				atomic.LoadInt32(&dp.currentWorkers),
				atomic.LoadInt32(&dp.queueLength),
				*avgTime)
		case <-ctx.Done():
			return
		}
	}
}

// Example usage
func main() {
	config := Config{
		MinWorkers: runtime.NumCPU(),
		MaxWorkers: runtime.NumCPU() * 4,
		QueueSize:  1000,
	}

	pool := NewDynamicPool(config)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the pool
	if err := pool.Start(ctx); err != nil {
		log.Fatalf("Failed to start pool: %v", err)
	}

	// Simulate varying workload
	go func() {
		for i := 0; i < 1000; i++ {
			// Simulate varying load
			load := 100 // Base load 100ms
			if i%100 == 0 {
				load = 500 // Occasional high load
			}

			task := Task{
				ID:       i,
				Load:     load,
				Priority: i % 3,
			}

			if err := pool.Submit(task); err != nil {
				log.Printf("Failed to submit task: %v", err)
			}

			// Vary submission rate
			time.Sleep(time.Duration(50+i%100) * time.Millisecond)
		}
	}()

	// Process results
	go func() {
		for result := range pool.results {
			log.Printf("Task %d completed by Worker %d in %v\n",
				result.TaskID, result.WorkerID, result.Duration)
		}
	}()

	// Let it run for a while
	time.Sleep(time.Minute)
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

// Config holds pool configuration
type Config struct {
	MinWorkers int
	MaxWorkers int
	QueueSize  int
}
