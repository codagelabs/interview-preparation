# Advanced Goroutines and Concurrency in Go

## 1. Goroutines Deep Dive
- **[How Goroutines work internally](cpu_and_internals/go_scheduler.md)**: Scheduler, GOMAXPROCS, work-stealing.
- **[Goroutine lifecycle](cpu_and_internals/goroutine_lifecycle.md)**: Creation, execution, termination.
- **[Stack growth in Goroutines](cpu_and_internals/stack_growth.md)**: How Go dynamically grows the stack.
- **[Best practices for Goroutine management](cpu_and_internals/goroutine_management.md)**: Guidelines and patterns for effective Goroutine handling.
- **[Detecting and avoiding Goroutine leaks](cpu_and_internals/goroutine_leaks.md)**: Avoiding Goroutine Leaks in a Long-Running Go Application.

## 2. Channels – Advanced Concepts
- **[Buffered channels](channels/buffered_channels/buffered_channels.md)**
- **[Unbuffered channels](channels/unbuffered_channels/unbuffered_channels.md)**
- **[Buffered vs Unbuffered channels](channels/channel_comparison/README.md)**: Performance differences.
- **[Closing channels](channels/closing_channals.md)**: How to safely detect a closed channel.
- **[Using `select` for multiplexing channels](channels/select/README.md)**.
- **[Channel directions](channels/channel_directions.md)**: `chan<-` for send, `<-chan` for receive.
- **[Deadlocks and race conditions](channels/deadlocks_and_races.md)** when using channels.
- **Fan-in and Fan-out patterns**.
- **[Using channels for signaling and synchronization](channels/signaling_and_sync.md)**.
- **[Best practices for using channels efficiently](channels/channel_best_practices.md)**.

## 3. The `sync` Package (Low-Level Concurrency)
- **[`sync.Mutex` and `sync.RWMutex`](sync_package/mutex/README.md)**: When to use them.
- **`sync.WaitGroup`** for coordinating multiple Goroutines.
- **`sync.Cond`** for advanced signaling between Goroutines.
- **`sync.Once`** for one-time initialization (singleton pattern).
- **`sync/atomic`** for lock-free concurrency (atomic operations).

## 4. Context Package (Managing Goroutines)
- **Using `context.Context` for Goroutine cancellation**.
- **Context timeouts and deadlines**: `context.WithTimeout`, `context.WithDeadline`.
- **Passing context through function calls**.
- **Preventing Goroutine leaks with context**.

## 5. Worker Pools and Job Queues
- **Creating a worker pool pattern using Goroutines and channels**.
- **Scaling worker pools dynamically based on workload**.
- **Implementing a job queue with workers**.
- **Managing Goroutine lifecycle within worker pools**.

## 6. Parallelism vs Concurrency
- **Understanding the difference between concurrency and parallelism**.
- **Tuning GOMAXPROCS for CPU-bound tasks**.
- **How Go scheduler manages concurrency**.
- **Profiling concurrent applications with `pprof`**.

## 7. Avoiding Common Concurrency Issues
- **Data races and how to detect them**: `go run -race`.
- **Deadlocks and livelocks**: How to prevent them.
- **Starvation and priority inversion**.
- **Best practices for avoiding Goroutine leaks**.
- **Benchmarking and optimizing concurrent applications**.

## 8. Real-World Concurrency Patterns
- **Implementing rate limiting**: Token bucket, leaky bucket.
- **Using Goroutines for event-driven programming**.
- **Designing concurrent pipelines for data processing**.
- **Handling backpressure in concurrent systems**.
- **Designing fault-tolerant concurrent systems**.
