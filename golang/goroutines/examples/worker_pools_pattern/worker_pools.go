package main

import (
	"fmt"
	"time"
)

// worker processes jobs from the jobs channel and sends results
// to the results channel
func worker(id int, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, job)
		time.Sleep(100 * time.Millisecond) // Simulate work
		results <- job * 2
		fmt.Printf("Worker %d completed job %d\n", id, job)
	}
}

func main() {
	numJobs := 5
	numWorkers := 3

	jobs := make(chan int)
	results := make(chan int)

	// Start workers
	fmt.Printf("Starting %d workers\n", numWorkers)
	for w := 1; w <= numWorkers; w++ {
		go worker(w, jobs, results)
	}

	// Send jobs
	fmt.Println("Sending jobs...")
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	// Collect results
	fmt.Println("Collecting results...")
	for i := 1; i <= numJobs; i++ {
		result := <-results
		fmt.Printf("Result: %d\n", result)
	}

	fmt.Println("All jobs completed")
}
