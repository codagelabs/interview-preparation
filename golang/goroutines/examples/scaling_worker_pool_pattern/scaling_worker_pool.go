package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// WorkerPool manages a dynamic pool of workers
type WorkerPool struct {
	// Task management
	tasks    chan int
	results  chan string
	
	// Worker management
	workerCount int32
	maxWorkers  int32
	
	// Load tracking
	taskCount   int32
	activeCount int32
	
	// Control
	wg       sync.WaitGroup
	stopChan chan struct{}
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(initialWorkers, maxWorkers int) *WorkerPool {
	return &WorkerPool{
		tasks:      make(chan int, 100),
		results:    make(chan string, 100),
		maxWorkers: int32(maxWorkers),
		stopChan:   make(chan struct{}),
	}
}

// worker processes tasks
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	
	for {
		select {
		case task, ok := <-wp.tasks:
			if !ok {
				return
			}
			
			// Track active workers
			atomic.AddInt32(&wp.activeCount, 1)
			
			// Process task
			time.Sleep(time.Duration(task) * time.Millisecond)
			result := fmt.Sprintf("Worker %d completed task: %d", id, task)
			wp.results <- result
			
			atomic.AddInt32(&wp.activeCount, -1)
			atomic.AddInt32(&wp.taskCount, -1)
			
		case <-wp.stopChan:
			return
		}
	}
}

// addWorker adds a new worker to the pool
func (wp *WorkerPool) addWorker() {
	if atomic.LoadInt32(&wp.workerCount) >= wp.maxWorkers {
		return
	}
	
	wp.wg.Add(1)
	workerID := int(atomic.AddInt32(&wp.workerCount, 1))
	go wp.worker(workerID)
	log.Printf("Added worker %d. Total workers: %d", workerID, atomic.LoadInt32(&wp.workerCount))
}

// removeWorker signals a worker to stop
func (wp *WorkerPool) removeWorker() {
	if atomic.LoadInt32(&wp.workerCount) > 1 {
		atomic.AddInt32(&wp.workerCount, -1)
		wp.stopChan <- struct{}{}
		log.Printf("Removed a worker. Total workers: %d", atomic.LoadInt32(&wp.workerCount))
	}
}

// monitorLoad adjusts the worker count based on workload
func (wp *WorkerPool) monitorLoad() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		taskCount := atomic.LoadInt32(&wp.taskCount)
		workerCount := atomic.LoadInt32(&wp.workerCount)
		activeWorkers := atomic.LoadInt32(&wp.activeCount)

		// Scale up if there are more tasks than workers
		if taskCount > workerCount && workerCount < wp.maxWorkers {
			wp.addWorker()
		}

		// Scale down if there are too many idle workers
		if activeWorkers < workerCount/2 && workerCount > 1 {
			wp.removeWorker()
		}

		log.Printf("Status - Workers: %d, Active: %d, Tasks: %d",
			workerCount, activeWorkers, taskCount)
	}
}

// Start begins the worker pool
func (wp *WorkerPool) Start() {
	// Start with one worker
	wp.addWorker()
	
	// Start load monitoring
	go wp.monitorLoad()
	
	// Start result processing
	go func() {
		for result := range wp.results {
			log.Println(result)
		}
	}()
}

// Submit adds a task to the pool
func (wp *WorkerPool) Submit(taskDuration int) error {
	select {
	case wp.tasks <- taskDuration:
		atomic.AddInt32(&wp.taskCount, 1)
		return nil
	default:
		return fmt.Errorf("task queue is full")
	}
}

func main() {
	// Create pool with max 10 workers
	pool := NewWorkerPool(1, 10)
	pool.Start()

	// Simulate varying workload
	go func() {
		for i := 0; i < 50; i++ {
			// Random task duration between 100ms and 1s
			taskDuration := rand.Intn(900) + 100
			
			if err := pool.Submit(taskDuration); err != nil {
				log.Printf("Failed to submit task: %v", err)
			}
			
			// Random delay between tasks
			time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
		}
	}()

	// Run for 30 seconds
	time.Sleep(30 * time.Second)
}