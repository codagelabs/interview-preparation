# sync.Mutex and sync.RWMutex

## Introduction
Mutexes are fundamental synchronization primitives in Go that prevent multiple goroutines from concurrently accessing a shared resource.

## sync.Mutex Deep Dive

### Basic Operations
```go
var mu sync.Mutex

mu.Lock()   // Acquire lock
// Critical section
mu.Unlock() // Release lock
```

### Advanced Concepts

#### 1. Internal Implementation
- Uses runtime's semaphore implementation
- Maintains a wait queue for blocked goroutines
- Implements spinning before parking goroutines

#### 2. Performance Considerations
- Lock contention impacts performance
- High contention scenarios may benefit from:
  - Reducing critical section size
  - Using buffered channels
  - Implementing lock striping

#### 3. Common Patterns

##### Embedded Mutex
```go
type SafeCounter struct {
    sync.Mutex  // embedded mutex
    count map[string]int
}

func (c *SafeCounter) Increment(key string) {
    c.Lock()
    defer c.Unlock()
    c.count[key]++
}
```

##### Mutex with Composition
```go
type SafeCounter struct {
    mu    sync.Mutex
    count map[string]int
}
```

### Advanced Use Cases
1. **Lock Striping**
```go
type StripedMap struct {
    locks    []*sync.Mutex
    segments []map[string]interface{}
}

func (sm *StripedMap) getSegment(key string) int {
    hash := hash(key)
    return hash % len(sm.locks)
}
```

2. **Try Lock Pattern**
```go
func TryLock(mu *sync.Mutex) bool {
    return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(mu)), 0, 1)
}
```

## sync.RWMutex Advanced Topics

### 1. Internal Mechanics
- Maintains separate read and write semaphores
- Uses atomic operations for reader count
- Implements write preference

### 2. Performance Characteristics
- Read locks are cheaper than write locks
- High read contention scenarios benefit from RWMutex
- Write locks block all readers

### 3. Advanced Patterns

#### Downgrading Write Lock to Read Lock
```go
type RWMap struct {
    sync.RWMutex
    data map[string]interface{}
}

func (m *RWMap) LoadOrCompute(key string, compute func() interface{}) interface{} {
    m.RLock()
    if val, ok := m.data[key]; ok {
        m.RUnlock()
        return val
    }
    m.RUnlock()
    
    m.Lock()
    defer m.Unlock()
    if val, ok := m.data[key]; ok {
        return val
    }
    val := compute()
    m.data[key] = val
    return val
}
```

### Best Practices and Gotchas

1. **Lock Ordering**
- Always acquire locks in the same order to prevent deadlocks
- Use lock leveling when possible

2. **Copying Mutexes**
```go
// DON'T do this
func (c Counter) Copy() Counter {
    return Counter{count: c.count}  // Mutex is copied!
}

// DO this
func (c *Counter) Copy() *Counter {
    return &Counter{count: c.count}
}
```

3. **Context Cancellation with Mutex**
```go
func DoWithTimeout(ctx context.Context, mu *sync.Mutex, fn func()) error {
    ch := make(chan struct{})
    go func() {
        mu.Lock()
        defer mu.Unlock()
        fn()
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

### Performance Monitoring

1. **Using Runtime Statistics**
```go
var mutex sync.Mutex
var mutexProfile = runtime.MutexProfile

func analyzeMutexContention() {
    // Get mutex profile
    p := pprof.Lookup("mutex")
    p.WriteTo(os.Stdout, 1)
}
```

2. **Mutex Profiling**
```bash
go test -bench=. -mutexprofile=mutex.prof
go tool pprof mutex.prof
```
```

```markdown:sync_package/waitgroup/README.md
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