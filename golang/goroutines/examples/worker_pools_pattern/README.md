# Worker Pools Example with Unbuffered Channels

This example demonstrates how to implement a simple worker pool pattern using unbuffered channels in Go. The program creates a pool of workers that process jobs concurrently and return results.

## Overview

The example shows:
- How to create multiple worker goroutines
- How to distribute work using an unbuffered jobs channel
- How to collect results using an unbuffered results channel
- How to coordinate work completion

## Code Explanation

```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for job := range jobs {
        results <- job * 2 // Process job and send result
    }
}

func main() {
    jobs := make(chan int)
    results := make(chan int)

    // Start 3 workers
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    // Send jobs
    for j := 1; j <= 5; j++ {
        jobs <- j
    }
    close(jobs)

    // Collect results
    for i := 1; i <= 5; i++ {
        <-results
    }
}
```

### Key Components

1. **Worker Function**
   - Takes an ID for identification
   - Receives jobs from the jobs channel
   - Processes each job (in this case, multiplies by 2)
   - Sends results to the results channel

2. **Main Function**
   - Creates unbuffered channels for jobs and results
   - Spawns worker goroutines
   - Distributes jobs
   - Collects results

## Running the Example

1. Save the code in a file named `worker_pools.go`
2. Run the program:
   ```bash
   go run worker_pools.go
   ```

## Expected Output

The program will process 5 jobs using 3 workers. Each job is multiplied by 2.

## Key Learnings

- Unbuffered channels provide natural synchronization between workers and the main goroutine
- Multiple workers can read from the same jobs channel
- Closing the jobs channel signals workers to stop when no more jobs are available
- The number of workers can be adjusted based on the workload

## Common Patterns Demonstrated

1. **Fan-out**: Multiple workers reading from a single jobs channel
2. **Channel Direction**: Using directional channel parameters (`<-chan` and `chan<-`)
3. **Graceful Shutdown**: Closing the jobs channel to signal completion
4. **Work Distribution**: Even distribution of work across multiple workers

## Tips

- Adjust the number of workers based on your CPU cores and workload
- Consider adding error handling for real-world applications
- Add logging or metrics to monitor worker performance
- Consider adding timeouts for job processing 