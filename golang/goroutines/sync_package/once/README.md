# sync.Once

sync.Once ensures that a function is executed only once, even across multiple goroutines.

## Key Concepts
- Provides a single method: `Do(func())`
- Guaranteed to execute exactly once
- Thread-safe initialization
- Common for singleton patterns

## Use Cases
1. Lazy initialization of shared resources
2. One-time setup operations
3. Implementing the singleton pattern

## Example Usage
```go
var (
    instance *Singleton
    once     sync.Once
)

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{}
    })
    return instance
}
```

## Best Practices
1. Don't call `Do()` with different functions
2. Use for expensive one-time initializations
3. Consider package init() for simple initializations 

# sync.Once Advanced Guide

## Internal Implementation
```go
type Once struct {
    done uint32
    m    sync.Mutex
}
```

## Advanced Patterns

### 1. Resettable Once
```go
type ResettableOnce struct {
    sync.Mutex
    done uint32
}

func (o *ResettableOnce) Do(f func()) {
    if atomic.LoadUint32(&o.done) == 0 {
        o.doSlow(f)
    }
}

func (o *ResettableOnce) doSlow(f func()) {
    o.Lock()
    defer o.Unlock()
    if o.done == 0 {
        defer atomic.StoreUint32(&o.done, 1)
        f()
    }
}

func (o *ResettableOnce) Reset() {
    o.Lock()
    defer o.Unlock()
    atomic.StoreUint32(&o.done, 0)
}
```

### 2. Once with Result
```go
type OnceValue[T any] struct {
    sync.Once
    value T
    err   error
}

func (v *OnceValue[T]) Get(f func() (T, error)) (T, error) {
    v.Once.Do(func() {
        v.value, v.err = f()
    })
    return v.value, v.err
}
```

### 3. Panic Recovery Pattern
```go
type SafeOnce struct {
    sync.Once
    recovered interface{}
}

func (o *SafeOnce) DoSafe(f func()) (recovered interface{}) {
    o.Once.Do(func() {
        defer func() {
            if r := recover(); r != nil {
                o.recovered = r
            }
        }()
        f()
    })
    return o.recovered
}
```

## Advanced Use Cases

### 1. Lazy Connection Pool
```go
type ConnectionPool struct {
    once sync.Once
    pool *Pool
}

func (cp *ConnectionPool) getPool() *Pool {
    cp.once.Do(func() {
        cp.pool = &Pool{
            connections: make([]*Connection, 0),
            maxSize:    100,
        }
    })
    return cp.pool
}
```

### 2. Complex Initialization Pattern
```go
type Service struct {
    initOnce   sync.Once
    startOnce  sync.Once
    stopOnce   sync.Once
    config     *Config
    client     *Client
    cache      *Cache
}

func (s *Service) Init(cfg *Config) error {
    var initErr error
    s.initOnce.Do(func() {
        if err := s.validateConfig(cfg); err != nil {
            initErr = err
            return
        }
        s.config = cfg
        s.client = NewClient(cfg)
        s.cache = NewCache(cfg)
    })
    return initErr
}
```

### 3. Singleton with Cleanup
```go
type CleanupSingleton struct {
    once     sync.Once
    initOnce sync.Once
    cleanup  func()
}

func (s *CleanupSingleton) Init() {
    s.initOnce.Do(func() {
        // Initialize resources
        s.cleanup = func() {
            // Cleanup resources
        }
    })
}

func (s *CleanupSingleton) Close() {
    s.once.Do(func() {
        if s.cleanup != nil {
            s.cleanup()
        }
    })
}
```

## Performance Considerations

### 1. Fast Path Optimization
```go
type OptimizedOnce struct {
    done uint32
    m    sync.Mutex
}

func (o *OptimizedOnce) Do(f func()) {
    if atomic.LoadUint32(&o.done) == 1 {
        return
    }
    o.doSlow(f)
}

func (o *OptimizedOnce) doSlow(f func()) {
    o.m.Lock()
    defer o.m.Unlock()
    if o.done == 0 {
        defer atomic.StoreUint32(&o.done, 1)
        f()
    }
}
```

### 2. Memory Ordering
```go
type MemoryEfficientOnce struct {
    done uint32
    m    sync.Mutex
    val  atomic.Value
}

func (o *MemoryEfficientOnce) Do(f func() interface{}) interface{} {
    if atomic.LoadUint32(&o.done) == 0 {
        o.m.Lock()
        defer o.m.Unlock()
        if o.done == 0 {
            v := f()
            o.val.Store(v)
            atomic.StoreUint32(&o.done, 1)
            return v
        }
    }
    return o.val.Load()
}
```

## Common Pitfalls

### 1. Recursive Initialization
```go
// WRONG - Deadlock
type BadSingleton struct {
    once sync.Once
    instance *BadSingleton
}

func (s *BadSingleton) Instance() *BadSingleton {
    s.once.Do(func() {
        s.instance = s.Instance() // Deadlock!
    })
    return s.instance
}

// RIGHT
func (s *BadSingleton) Instance() *BadSingleton {
    s.once.Do(func() {
        s.instance = &BadSingleton{}
    })
    return s.instance
}
```

### 2. Copying Once
```go
// WRONG
func (s *Service) Copy() Service {
    return Service{once: s.once} // Once should not be copied

// RIGHT
func (s *Service) Copy() *Service {
    return &Service{} // New Once instance
}
```

### 3. Error Handling
```go
type OnceError struct {
    sync.Once
    err error
}

func (o *OnceError) Do(f func() error) error {
    o.Once.Do(func() {
        o.err = f()
    })
    return o.err
}
```

Would you like me to continue with the enhanced version of the `atomic` package README as well? 