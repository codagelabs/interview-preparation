# Usage of sync.Map in Go

## Overview

`sync.Map` is Go's built-in concurrent map implementation that's optimized for specific use cases. Unlike regular Go maps, `sync.Map` is safe for concurrent use without additional synchronization.

## Key Characteristics

- **Thread-safe**: Can be used concurrently from multiple goroutines
- **Optimized for specific patterns**: Best when keys are only written once but read many times
- **Interface-based**: Uses `interface{}` for both keys and values
- **No range support**: Cannot iterate over all key-value pairs

## Basic Implementation

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var m sync.Map
    
    // Store values
    m.Store("key1", "value1")
    m.Store("key2", "value2")
    m.Store(123, "numeric key")
    
    // Load values
    if value, ok := m.Load("key1"); ok {
        fmt.Printf("Found: %v\n", value)
    }
    
    // Load or store (atomic operation)
    actual, loaded := m.LoadOrStore("key3", "value3")
    fmt.Printf("Actual: %v, Loaded: %v\n", actual, loaded)
    
    // Delete
    m.Delete("key2")
    
    // Check if key exists after deletion
    if _, ok := m.Load("key2"); !ok {
        fmt.Println("key2 was deleted")
    }
}
```

## Advanced Implementation Patterns

### 1. Cache Implementation

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type CacheItem struct {
    Value     interface{}
    ExpiresAt time.Time
}

type ConcurrentCache struct {
    data sync.Map
}

func NewConcurrentCache() *ConcurrentCache {
    return &ConcurrentCache{}
}

func (c *ConcurrentCache) Set(key string, value interface{}, ttl time.Duration) {
    item := CacheItem{
        Value:     value,
        ExpiresAt: time.Now().Add(ttl),
    }
    c.data.Store(key, item)
}

func (c *ConcurrentCache) Get(key string) (interface{}, bool) {
    if value, ok := c.data.Load(key); ok {
        item := value.(CacheItem)
        if time.Now().Before(item.ExpiresAt) {
            return item.Value, true
        }
        // Expired, remove it
        c.data.Delete(key)
    }
    return nil, false
}

func (c *ConcurrentCache) Delete(key string) {
    c.data.Delete(key)
}

func (c *ConcurrentCache) Cleanup() {
    c.data.Range(func(key, value interface{}) bool {
        item := value.(CacheItem)
        if time.Now().After(item.ExpiresAt) {
            c.data.Delete(key)
        }
        return true
    })
}

func main() {
    cache := NewConcurrentCache()
    
    // Set with TTL
    cache.Set("user:123", "John Doe", 5*time.Second)
    cache.Set("session:abc", "active", 10*time.Second)
    
    // Get values
    if user, ok := cache.Get("user:123"); ok {
        fmt.Printf("User: %v\n", user)
    }
    
    // Wait for expiration
    time.Sleep(6 * time.Second)
    
    if _, ok := cache.Get("user:123"); !ok {
        fmt.Println("User data expired")
    }
}
```

### 2. Counter Implementation

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

type ConcurrentCounter struct {
    data sync.Map
}

func NewConcurrentCounter() *ConcurrentCounter {
    return &ConcurrentCounter{}
}

func (c *ConcurrentCounter) Increment(key string) int64 {
    for {
        if value, ok := c.data.Load(key); ok {
            // Key exists, increment atomically
            counter := value.(*int64)
            newVal := atomic.AddInt64(counter, 1)
            return newVal
        } else {
            // Key doesn't exist, create new counter
            newCounter := int64(1)
            if actual, loaded := c.data.LoadOrStore(key, &newCounter); loaded {
                // Another goroutine created it, increment that one
                counter := actual.(*int64)
                return atomic.AddInt64(counter, 1)
            } else {
                // We created it successfully
                return 1
            }
        }
    }
}

func (c *ConcurrentCounter) Get(key string) int64 {
    if value, ok := c.data.Load(key); ok {
        return atomic.LoadInt64(value.(*int64))
    }
    return 0
}

func (c *ConcurrentCounter) Reset(key string) {
    c.data.Delete(key)
}

func main() {
    counter := NewConcurrentCounter()
    
    // Increment counters concurrently
    var wg sync.WaitGroup
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Increment("requests")
            counter.Increment("errors")
        }()
    }
    
    wg.Wait()
    
    fmt.Printf("Requests: %d\n", counter.Get("requests"))
    fmt.Printf("Errors: %d\n", counter.Get("errors"))
}
```

### 3. Registry Pattern

```go
package main

import (
    "fmt"
    "sync"
)

type Service interface {
    Name() string
    Start() error
    Stop() error
}

type ServiceRegistry struct {
    services sync.Map
}

func NewServiceRegistry() *ServiceRegistry {
    return &ServiceRegistry{}
}

func (r *ServiceRegistry) Register(service Service) error {
    if _, loaded := r.services.LoadOrStore(service.Name(), service); loaded {
        return fmt.Errorf("service %s already registered", service.Name())
    }
    return nil
}

func (r *ServiceRegistry) Get(name string) (Service, bool) {
    if value, ok := r.services.Load(name); ok {
        return value.(Service), true
    }
    return nil, false
}

func (r *ServiceRegistry) Unregister(name string) {
    r.services.Delete(name)
}

func (r *ServiceRegistry) List() []string {
    var names []string
    r.services.Range(func(key, value interface{}) bool {
        names = append(names, key.(string))
        return true
    })
    return names
}

// Example service implementation
type DatabaseService struct {
    name string
}

func (d *DatabaseService) Name() string {
    return d.name
}

func (d *DatabaseService) Start() error {
    fmt.Printf("Starting %s\n", d.name)
    return nil
}

func (d *DatabaseService) Stop() error {
    fmt.Printf("Stopping %s\n", d.name)
    return nil
}

func main() {
    registry := NewServiceRegistry()
    
    // Register services
    db1 := &DatabaseService{name: "primary-db"}
    db2 := &DatabaseService{name: "replica-db"}
    
    registry.Register(db1)
    registry.Register(db2)
    
    // List all services
    fmt.Printf("Registered services: %v\n", registry.List())
    
    // Get a specific service
    if service, ok := registry.Get("primary-db"); ok {
        service.Start()
    }
}
```

### 4. Configuration Store

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type ConfigStore struct {
    data      sync.Map
    listeners sync.Map // map[string][]chan interface{}
}

func NewConfigStore() *ConfigStore {
    return &ConfigStore{}
}

func (c *ConfigStore) Set(key string, value interface{}) {
    c.data.Store(key, value)
    c.notifyListeners(key, value)
}

func (c *ConfigStore) Get(key string) (interface{}, bool) {
    return c.data.Load(key)
}

func (c *ConfigStore) GetString(key string) (string, bool) {
    if value, ok := c.data.Load(key); ok {
        if str, ok := value.(string); ok {
            return str, true
        }
    }
    return "", false
}

func (c *ConfigStore) GetInt(key string) (int, bool) {
    if value, ok := c.data.Load(key); ok {
        if i, ok := value.(int); ok {
            return i, true
        }
    }
    return 0, false
}

func (c *ConfigStore) Subscribe(key string) <-chan interface{} {
    ch := make(chan interface{}, 1)
    
    // Get current value if exists
    if value, ok := c.data.Load(key); ok {
        ch <- value
    }
    
    // Add to listeners
    if listeners, ok := c.listeners.Load(key); ok {
        listenerSlice := listeners.([]chan interface{})
        listenerSlice = append(listenerSlice, ch)
        c.listeners.Store(key, listenerSlice)
    } else {
        c.listeners.Store(key, []chan interface{}{ch})
    }
    
    return ch
}

func (c *ConfigStore) notifyListeners(key string, value interface{}) {
    if listeners, ok := c.listeners.Load(key); ok {
        for _, ch := range listeners.([]chan interface{}) {
            select {
            case ch <- value:
            default:
                // Channel is full, skip
            }
        }
    }
}

func main() {
    config := NewConfigStore()
    
    // Set initial values
    config.Set("database_url", "localhost:5432")
    config.Set("max_connections", 100)
    
    // Subscribe to changes
    dbCh := config.Subscribe("database_url")
    
    // Start a goroutine to listen for changes
    go func() {
        for value := range dbCh {
            fmt.Printf("Database URL changed to: %v\n", value)
        }
    }()
    
    // Simulate configuration changes
    time.Sleep(1 * time.Second)
    config.Set("database_url", "prod-db:5432")
    
    time.Sleep(1 * time.Second)
    config.Set("max_connections", 200)
    
    // Get values
    if url, ok := config.GetString("database_url"); ok {
        fmt.Printf("Current DB URL: %s\n", url)
    }
    
    if maxConn, ok := config.GetInt("max_connections"); ok {
        fmt.Printf("Max connections: %d\n", maxConn)
    }
}
```

## Performance Considerations

### When to Use sync.Map

**Use sync.Map when:**
- Keys are written once but read many times
- Multiple goroutines read, write, and overwrite entries for disjoint sets of keys
- You don't need to iterate over all key-value pairs

**Don't use sync.Map when:**
- You need to iterate over all keys
- You have a small number of keys
- You need type safety (sync.Map uses interface{})

### Performance Comparison

```go
package main

import (
    "fmt"
    "sync"
    "testing"
    "time"
)

// Benchmark regular map with mutex
func BenchmarkRegularMap(b *testing.B) {
    var mu sync.RWMutex
    m := make(map[string]int)
    
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            key := fmt.Sprintf("key_%d", time.Now().UnixNano()%1000)
            
            // 90% reads, 10% writes
            if time.Now().UnixNano()%10 == 0 {
                mu.Lock()
                m[key]++
                mu.Unlock()
            } else {
                mu.RLock()
                _ = m[key]
                mu.RUnlock()
            }
        }
    })
}

// Benchmark sync.Map
func BenchmarkSyncMap(b *testing.B) {
    var m sync.Map
    
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            key := fmt.Sprintf("key_%d", time.Now().UnixNano()%1000)
            
            // 90% reads, 10% writes
            if time.Now().UnixNano()%10 == 0 {
                if value, ok := m.Load(key); ok {
                    m.Store(key, value.(int)+1)
                } else {
                    m.Store(key, 1)
                }
            } else {
                m.Load(key)
            }
        }
    })
}
```

## Best Practices

1. **Use for the right use case**: sync.Map is optimized for specific patterns
2. **Handle type assertions carefully**: Always check if type assertions succeed
3. **Consider memory usage**: sync.Map may use more memory than regular maps
4. **Use LoadOrStore for atomic operations**: This prevents race conditions
5. **Clean up expired entries**: Use Range() to iterate and clean up

## Common Pitfalls

```go
// WRONG: Not handling type assertion
func badExample(m sync.Map, key string) {
    value := m.Load(key) // Returns interface{}
    str := value.(string) // Panic if value is not string!
}

// RIGHT: Handle type assertion safely
func goodExample(m sync.Map, key string) {
    if value, ok := m.Load(key); ok {
        if str, ok := value.(string); ok {
            // Use str safely
            fmt.Println(str)
        }
    }
}

// WRONG: Not using LoadOrStore for atomic operations
func badAtomic(m sync.Map, key string) {
    if _, ok := m.Load(key); !ok {
        m.Store(key, "value") // Race condition!
    }
}

// RIGHT: Use LoadOrStore for atomic operations
func goodAtomic(m sync.Map, key string) {
    m.LoadOrStore(key, "value")
}
```

## API Reference

### Core Methods

- `Store(key, value interface{})` - Store a key-value pair
- `Load(key interface{}) (value interface{}, ok bool)` - Load a value by key
- `LoadOrStore(key, value interface{}) (actual interface{}, loaded bool)` - Load existing or store new
- `LoadAndDelete(key interface{}) (value interface{}, loaded bool)` - Load and delete atomically
- `Delete(key interface{})` - Delete a key
- `Range(func(key, value interface{}) bool)` - Iterate over all key-value pairs

### When to Use Each Method

- **Store**: When you want to set a value unconditionally
- **Load**: When you want to read a value safely
- **LoadOrStore**: When you want atomic "get or create" semantics
- **LoadAndDelete**: When you want atomic "get and remove" semantics
- **Delete**: When you want to remove a key unconditionally
- **Range**: When you need to iterate (use sparingly as it's expensive)

This comprehensive guide covers the essential aspects of `sync.Map` implementations in Go, from basic usage to advanced patterns and performance considerations.
