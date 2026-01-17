# sync.WaitGroup

WaitGroup is used to wait for a collection of goroutines to finish executing.

## Key Concepts
- `Add(delta)`: Adds delta to the WaitGroup counter
- `Done()`: Decrements the WaitGroup counter by 1
- `Wait()`: Blocks until the counter is zero

## Common Use Cases
1. Parallel processing of items in a collection
2. Coordinating multiple worker goroutines
3. Ensuring all goroutines complete before program exit

## Example Usage
```go
var wg sync.WaitGroup

for i := 0; i < 3; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        // Do work here
    }(i)
}

wg.Wait() // Wait for all goroutines to finish
```

## Best Practices
1. Always call `Done()` using `defer`
2. Add to WaitGroup before starting goroutines
3. Never copy a WaitGroup after first use 

# sync.WaitGroup Deep Dive

## Internal Implementation
WaitGroup maintains an internal counter using atomic operations and semaphores.

### Structure
```go
type WaitGroup struct {
    noCopy noCopy
    state1 [3]uint32 // Contains counter and waiter count
}
```

## Advanced Usage Patterns

### 1. Dynamic Worker Pool
```go
func WorkerPool(tasks []Task, maxWorkers int) {
    wg := sync.WaitGroup{}
    taskCh := make(chan Task, len(tasks))

    // Launch workers
    for i := 0; i < maxWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for task := range taskCh {
                task.Execute()
            }
        }()
    }

    // Send tasks
    for _, task := range tasks {
        taskCh <- task
    }
    close(taskCh)

    wg.Wait()
}
```

### 2. Pipeline Pattern
```go
func Pipeline(data []int) []int {
    var wg sync.WaitGroup
    stage1 := make(chan int)
    stage2 := make(chan int)
    result := make(chan int)

    // Stage 1: Square numbers
    wg.Add(1)
    go func() {
        defer wg.Done()
        defer close(stage1)
        for _, n := range data {
            stage1 <- n * n
        }
    }()

    // Stage 2: Double numbers
    wg.Add(1)
    go func() {
        defer wg.Done()
        defer close(stage2)
        for n := range stage1 {
            stage2 <- n * 2
        }
    }()

    // Collect results
    var results []int
    wg.Add(1)
    go func() {
        defer wg.Done()
        for n := range stage2 {
            results = append(results, n)
        }
    }()

    wg.Wait()
    return results
}
```

### 3. Batch Processing with Error Handling
```go
type Result struct {
    Value interface{}
    Err   error
}

func BatchProcess(items []Item) []Result {
    results := make([]Result, len(items))
    var wg sync.WaitGroup

    for i, item := range items {
        wg.Add(1)
        go func(i int, item Item) {
            defer wg.Done()
            results[i] = processItem(item)
        }(i, item)
    }

    wg.Wait()
    return results
}
```

## Advanced Concepts

### 1. WaitGroup with Context
```go
func WaitGroupWithContext(ctx context.Context, wg *sync.WaitGroup) error {
    ch := make(chan struct{})
    go func() {
        wg.Wait()
        close(ch)
    }()

    select {
    case <-ch:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

### 2. Dynamic Task Addition
```go
type TaskManager struct {
    wg      sync.WaitGroup
    tasksMu sync.Mutex
    tasks   []Task
}

func (tm *TaskManager) AddTask(t Task) {
    tm.tasksMu.Lock()
    tm.tasks = append(tm.tasks, t)
    tm.wg.Add(1)
    tm.tasksMu.Unlock()

    go func() {
        defer tm.wg.Done()
        t.Execute()
    }()
}
```

### 3. Hierarchical WaitGroups
```go
type HierarchicalTask struct {
    parentWg *sync.WaitGroup
    localWg  sync.WaitGroup
    subtasks []Task
}

func (ht *HierarchicalTask) Execute() {
    defer ht.parentWg.Done()

    for _, task := range ht.subtasks {
        ht.localWg.Add(1)
        go func(t Task) {
            defer ht.localWg.Done()
            t.Execute()
        }(task)
    }

    ht.localWg.Wait()
}
```

## Common Pitfalls and Solutions

### 1. Negative Counter
```go
// WRONG
wg.Done() // Called before Add()

// RIGHT
wg.Add(1)
go func() {
    defer wg.Done()
    // work
}()
```

### 2. WaitGroup Reuse
```go
// WRONG
wg.Add(n)
wg.Wait()
wg.Add(m) // Reusing without reset

// RIGHT
wg = sync.WaitGroup{} // Create new instance
wg.Add(m)
```

### 3. Copying WaitGroup
```go
// WRONG
func doWork(wg sync.WaitGroup) { // WaitGroup copied
    wg.Done()
}

// RIGHT
func doWork(wg *sync.WaitGroup) {
    wg.Done()
}
```

## Performance Considerations

### 1. Granularity
- Too fine-grained: Overhead from goroutine creation
- Too coarse-grained: Reduced parallelism
- Find the right balance for your use case

### 2. Memory Usage
- Each goroutine has overhead
- Consider pooling for frequent, short-lived tasks

### 3. Monitoring
```go
func monitorGoroutines() {
    ticker := time.NewTicker(time.Second)
    for range ticker.C {
        fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())
    }
}
```

I'll continue with the enhanced versions of `cond`, `once`, and `atomic` if you'd like. Each will follow a similar pattern of including advanced topics, real-world examples, and best practices. Would you like me to proceed with those as well? 