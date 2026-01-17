# sync.Cond

Condition variables provide a way for goroutines to efficiently wait for arbitrary conditions.

## Key Concepts
- Built on top of `sync.Mutex` or `sync.RWMutex`
- Methods: `Wait()`, `Signal()`, `Broadcast()`
- Used for complex synchronization scenarios

## Use Cases
1. Producer-Consumer patterns
2. Implementing wait queues
3. Complex state-dependent synchronization

## Example Usage
```go
var (
    condition = sync.NewCond(&sync.Mutex{})
    ready     = false
)

// Consumer
go func() {
    condition.L.Lock()
    for !ready {
        condition.Wait()
    }
    condition.L.Unlock()
}()

// Producer
condition.L.Lock()
ready = true
condition.Signal() // or Broadcast()
condition.L.Unlock()
```

## Best Practices
1. Always protect condition state with the associated lock
2. Use `for` loops with Wait() to handle spurious wakeups
3. Consider using channels for simpler synchronization needs 

# sync.Cond Advanced Guide

## Internal Implementation
A condition variable is built on top of a mutex/rwmutex and maintains a linked list of waiting goroutines.

## Core Concepts

### 1. Basic Structure
```go
type Cond struct {
    L Locker           // Underlying mutex/rwmutex
    notify  notifyList // List of waiting goroutines
    checker copyChecker
}
```

### 2. Primary Operations
- `Wait()`: Atomically unlocks and waits for signal
- `Signal()`: Wakes one waiting goroutine
- `Broadcast()`: Wakes all waiting goroutines

## Advanced Patterns

### 1. Bounded Queue Implementation
```go
type BoundedQueue struct {
    cond *sync.Cond
    data []interface{}
    size int
}

func NewBoundedQueue(capacity int) *BoundedQueue {
    return &BoundedQueue{
        cond: sync.NewCond(&sync.Mutex{}),
        data: make([]interface{}, 0, capacity),
        size: capacity,
    }
}

func (q *BoundedQueue) Put(item interface{}) {
    q.cond.L.Lock()
    defer q.cond.L.Unlock()

    for len(q.data) == q.size {
        q.cond.Wait() // Wait for space
    }
    q.data = append(q.data, item)
    q.cond.Signal() // Signal consumers
}

func (q *BoundedQueue) Get() interface{} {
    q.cond.L.Lock()
    defer q.cond.L.Unlock()

    for len(q.data) == 0 {
        q.cond.Wait() // Wait for items
    }
    item := q.data[0]
    q.data = q.data[1:]
    q.cond.Signal() // Signal producers
    return item
}
```

### 2. Multiple Condition Variables
```go
type ResourcePool struct {
    empty    *sync.Cond
    full     *sync.Cond
    mu       sync.Mutex
    resources []interface{}
    capacity int
}

func NewResourcePool(capacity int) *ResourcePool {
    rp := &ResourcePool{
        capacity:  capacity,
        resources: make([]interface{}, 0, capacity),
    }
    rp.empty = sync.NewCond(&rp.mu)
    rp.full = sync.NewCond(&rp.mu)
    return rp
}
```

### 3. Timeout Pattern with Context
```go
func (c *Cond) WaitWithTimeout(timeout time.Duration) bool {
    ch := make(chan struct{})
    go func() {
        c.Wait()
        close(ch)
    }()

    select {
    case <-ch:
        return true
    case <-time.After(timeout):
        return false
    }
}
```

### 4. Priority Condition Pattern
```go
type PriorityCondition struct {
    cond     *sync.Cond
    priority int
    data     interface{}
}

func (pc *PriorityCondition) WaitUntilPriority(minPriority int) {
    pc.cond.L.Lock()
    defer pc.cond.L.Unlock()

    for pc.priority < minPriority {
        pc.cond.Wait()
    }
}
```

## Advanced Use Cases

### 1. Barrier Implementation
```go
type Barrier struct {
    cond    *sync.Cond
    count   int
    waiting int
}

func (b *Barrier) Wait() {
    b.cond.L.Lock()
    defer b.cond.L.Unlock()
    
    b.waiting++
    if b.waiting < b.count {
        b.cond.Wait()
    } else {
        b.waiting = 0
        b.cond.Broadcast()
    }
}
```

### 2. Reader-Writer Preference
```go
type RWPreference struct {
    cond       *sync.Cond
    readers    int
    writers    int
    writeWait  int
    preference string // "reader" or "writer"
}

func (rw *RWPreference) StartRead() {
    rw.cond.L.Lock()
    defer rw.cond.L.Unlock()

    if rw.preference == "writer" {
        for rw.writers > 0 || rw.writeWait > 0 {
            rw.cond.Wait()
        }
    }
    rw.readers++
}
```

## Performance Considerations

### 1. Broadcast vs Signal
```go
// Efficient use of Signal
func (q *Queue) Get() interface{} {
    q.cond.L.Lock()
    defer q.cond.L.Unlock()

    for len(q.items) == 0 {
        q.cond.Wait()
    }
    item := q.items[0]
    q.items = q.items[1:]
    
    if len(q.items) > 0 {
        q.cond.Signal() // Signal only if more items available
    }
    return item
}
```

### 2. Lock Contention Monitoring
```go
type MonitoredCond struct {
    *sync.Cond
    waitCount int64
}

func (mc *MonitoredCond) Wait() {
    atomic.AddInt64(&mc.waitCount, 1)
    mc.Cond.Wait()
    atomic.AddInt64(&mc.waitCount, -1)
}
```

## Common Pitfalls

### 1. Lost Wakeups
```go
// WRONG
if condition {
    cond.Signal()
}
mutex.Lock()

// RIGHT
mutex.Lock()
if condition {
    cond.Signal()
}
```

### 2. Spurious Wakeups
```go
// WRONG
cond.L.Lock()
cond.Wait()
// Process condition

// RIGHT
cond.L.Lock()
for !condition() {
    cond.Wait()
}
// Process condition
```

### 3. Deadlock Prevention
```go
func SafeCondOperation(cond *sync.Cond, timeout time.Duration) error {
    done := make(chan struct{})
    go func() {
        cond.L.Lock()
        defer cond.L.Unlock()
        defer close(done)
        
        cond.Wait()
    }()

    select {
    case <-done:
        return nil
    case <-time.After(timeout):
        return errors.New("operation timed out")
    }
}
``` 