# sync.Map Internals in Go

## Overview

`sync.Map` is a concurrent map implementation in Go that uses a sophisticated two-level data structure to achieve high performance for specific access patterns. Understanding its internals helps in making informed decisions about when and how to use it.

## Core Data Structure

`sync.Map` uses a two-level structure:

1. **Read Map (readOnly)**: A read-only map that's optimized for fast reads
2. **Dirty Map**: A writable map that contains new entries and modified entries

```go
type Map struct {
    mu sync.Mutex
    
    // read contains the portion of the map's contents that are safe for
    // concurrent access (with or without mu held for reading).
    read atomic.Value // readOnly
    
    // dirty contains the portion of the map's contents that require mu to be
    // held. To ensure that the dirty map can be promoted to the read map if
    // needed, it also includes all of the read map's unreachable entries.
    dirty map[interface{}]*entry
    
    // misses counts the number of loads since the read map was last updated
    // that needed to lock mu to determine whether the key was present.
    misses int
}

type readOnly struct {
    m       map[interface{}]*entry
    amended bool // true if the dirty map contains some key not in m.
}

type entry struct {
    p unsafe.Pointer // *interface{}
}
```

## Key Components

### 1. Entry Structure

The `entry` struct is crucial for atomic operations:

```go
type entry struct {
    p unsafe.Pointer // *interface{}
}

// Possible states of an entry:
// - nil: deleted
// - expunged: deleted and won't be stored in dirty map
// - actual value: the real value
```

### 2. Read-Only Map

The read-only map is stored in an `atomic.Value` for lock-free access:

```go
type readOnly struct {
    m       map[interface{}]*entry
    amended bool // indicates if dirty map has keys not in read map
}
```

### 3. Dirty Map

The dirty map requires the mutex to be held and contains:
- All new entries
- All modified entries from the read map
- All unreachable entries from the read map

## Internal Operations

### 1. Load Operation

```go
func (m *Map) Load(key interface{}) (value interface{}, ok bool) {
    read, _ := m.read.Load().(readOnly)
    e, ok := read.m[key]
    if !ok && read.amended {
        m.mu.Lock()
        // Avoid reporting a spurious miss if m.dirty was promoted to m.read
        // between the call to m.read.Load() and the call to m.mu.Lock().
        read, _ = m.read.Load().(readOnly)
        e, ok = read.m[key]
        if !ok && read.amended {
            e, ok = m.dirty[key]
            // Regardless of whether the entry was present, record a miss:
            // this key will take the slow path until the dirty map is promoted
            // to the read map.
            m.missLocked()
        }
        m.mu.Unlock()
    }
    if !ok {
        return nil, false
    }
    return e.load()
}
```

**Load Process:**
1. First, try to read from the read-only map (lock-free)
2. If not found and `amended` is true, acquire the mutex
3. Double-check the read map (double-checked locking pattern)
4. If still not found, check the dirty map
5. Record a miss to potentially promote dirty map to read map

### 2. Store Operation

```go
func (m *Map) Store(key, value interface{}) {
    read, _ := m.read.Load().(readOnly)
    if e, ok := read.m[key]; ok && e.tryStore(&value) {
        return
    }
    
    m.mu.Lock()
    read, _ = m.read.Load().(readOnly)
    if e, ok := read.m[key]; ok {
        if e.unexpungeLocked() {
            // The entry was previously expunged, which means that the value
            // for e is in m.dirty and was not in m.read.
            m.dirty[key] = e
        }
        e.storeLocked(&value)
    } else if e, ok := m.dirty[key]; ok {
        e.storeLocked(&value)
    } else {
        if !read.amended {
            // We're adding the first new key to the dirty map.
            // Make sure it is allocated and mark the read-only map as incomplete.
            m.dirtyLocked()
            m.read.Store(readOnly{m: read.m, amended: true})
        }
        m.dirty[key] = newEntry(value)
    }
    m.mu.Unlock()
}
```

**Store Process:**
1. Try to update existing entry in read map (lock-free)
2. If successful, return
3. Otherwise, acquire mutex and check read map again
4. If entry exists in read map, update it and potentially move to dirty map
5. If entry exists in dirty map, update it
6. If new entry, add to dirty map and mark read map as amended

### 3. Delete Operation

```go
func (m *Map) Delete(key interface{}) {
    m.LoadAndDelete(key)
}

func (m *Map) LoadAndDelete(key interface{}) (value interface{}, loaded bool) {
    read, _ := m.read.Load().(readOnly)
    e, ok := read.m[key]
    if !ok && read.amended {
        m.mu.Lock()
        read, _ = m.read.Load().(readOnly)
        e, ok = read.m[key]
        if !ok && read.amended {
            e, ok = m.dirty[key]
            delete(m.dirty, key)
            // Regardless of whether the entry was present, record a miss:
            // this key will take the slow path until the dirty map is promoted
            // to the read map.
            m.missLocked()
        }
        m.mu.Unlock()
    }
    if ok {
        return e.delete()
    }
    return nil, false
}
```

## Memory Management

### Entry States

Entries can be in three states:

1. **nil**: Entry is deleted
2. **expunged**: Entry is deleted and won't be stored in dirty map
3. **actual value**: Entry contains a real value

```go
func (e *entry) tryStore(i *interface{}) bool {
    for {
        p := atomic.LoadPointer(&e.p)
        if p == expunged {
            return false
        }
        if atomic.CompareAndSwapPointer(&e.p, p, unsafe.Pointer(i)) {
            return true
        }
    }
}

func (e *entry) unexpungeLocked() (wasExpunged bool) {
    return atomic.CompareAndSwapPointer(&e.p, expunged, nil)
}
```

### Promotion Process

When the number of misses exceeds the size of the dirty map, the dirty map is promoted to the read map:

```go
func (m *Map) missLocked() {
    m.misses++
    if m.misses < len(m.dirty) {
        return
    }
    m.read.Store(readOnly{m: m.dirty})
    m.dirty = nil
    m.misses = 0
}
```

## Performance Characteristics

### Read Performance

- **Best case**: O(1) - direct access to read map (lock-free)
- **Worst case**: O(1) + mutex overhead - access to dirty map

### Write Performance

- **Best case**: O(1) - update existing entry in read map (lock-free)
- **Worst case**: O(1) + mutex overhead - update dirty map or create new entry

### Memory Overhead

- **Read map**: Contains all entries that are frequently accessed
- **Dirty map**: Contains new/modified entries + unreachable entries from read map
- **Entry overhead**: Each entry is a pointer (8 bytes on 64-bit systems)

## Optimization Strategies

### 1. Miss Tracking

The `misses` counter tracks how many times we had to fall back to the dirty map. When misses exceed the dirty map size, promotion occurs.

### 2. Amended Flag

The `amended` flag indicates whether the dirty map contains keys not in the read map, avoiding unnecessary dirty map checks.

### 3. Double-Checked Locking

Used in Load operations to avoid unnecessary mutex acquisition when the key exists in the read map.

## Race Conditions and Concurrency

### Lock-Free Reads

Read operations are lock-free when accessing the read map:

```go
// This is safe without locks
read, _ := m.read.Load().(readOnly)
e, ok := read.m[key]
```

### Mutex Protection

Write operations and dirty map access are protected by the mutex:

```go
m.mu.Lock()
// ... modify dirty map or promote
m.mu.Unlock()
```

### Atomic Operations

Entry modifications use atomic operations:

```go
func (e *entry) load() (value interface{}, ok bool) {
    p := atomic.LoadPointer(&e.p)
    if p == nil || p == expunged {
        return nil, false
    }
    return *(*interface{})(p), true
}
```

## Memory Layout Visualization

```
sync.Map
├── mu: sync.Mutex
├── read: atomic.Value
│   └── readOnly
│       ├── m: map[interface{}]*entry
│       └── amended: bool
├── dirty: map[interface{}]*entry
└── misses: int

entry
└── p: unsafe.Pointer
    ├── nil (deleted)
    ├── expunged (deleted, won't be in dirty)
    └── *interface{} (actual value)
```

## When sync.Map is Optimal

### Ideal Use Cases

1. **Write-once, read-many**: Keys are written once and read frequently
2. **Disjoint key sets**: Different goroutines work with different keys
3. **Infrequent writes**: Write operations are much less common than reads

### Performance Characteristics

- **High read concurrency**: Multiple goroutines can read simultaneously
- **Low write contention**: Writes are serialized but infrequent
- **Memory efficient**: Unused entries are eventually cleaned up

## Comparison with Other Approaches

### vs Regular Map + Mutex

```go
// Regular map with mutex
type RegularMap struct {
    mu sync.RWMutex
    m  map[interface{}]interface{}
}

// sync.Map
type SyncMap struct {
    // ... internal structure
}
```

**sync.Map advantages:**
- Lock-free reads
- Better performance for read-heavy workloads
- Automatic memory management

**Regular map advantages:**
- Simpler implementation
- Better for write-heavy workloads
- More predictable performance

### vs Regular Map + RWMutex

```go
// Regular map with RWMutex
type RWMap struct {
    mu sync.RWMutex
    m  map[interface{}]interface{}
}
```

**sync.Map advantages:**
- Truly lock-free reads
- No reader-writer contention
- Better scalability

**RWMutex advantages:**
- Simpler to understand
- Better for balanced read/write workloads
- More predictable memory usage

## Debugging and Profiling

### Race Detection

```bash
go run -race your_program.go
go test -race ./...
```

### Memory Profiling

```bash
go test -memprofile=mem.prof ./...
go tool pprof mem.prof
```

### Mutex Profiling

```bash
go test -mutexprofile=mutex.prof ./...
go tool pprof mutex.prof
```

## Common Implementation Patterns

### 1. Type-Safe Wrapper

```go
type TypedMap[K comparable, V any] struct {
    m sync.Map
}

func (tm *TypedMap[K, V]) Store(key K, value V) {
    tm.m.Store(key, value)
}

func (tm *TypedMap[K, V]) Load(key K) (V, bool) {
    if value, ok := tm.m.Load(key); ok {
        return value.(V), true
    }
    var zero V
    return zero, false
}
```

### 2. Metrics Collection

```go
type MetricsMap struct {
    m        sync.Map
    hits     int64
    misses   int64
    stores   int64
    deletes  int64
}

func (mm *MetricsMap) Load(key interface{}) (interface{}, bool) {
    value, ok := mm.m.Load(key)
    if ok {
        atomic.AddInt64(&mm.hits, 1)
    } else {
        atomic.AddInt64(&mm.misses, 1)
    }
    return value, ok
}
```

## Best Practices for sync.Map

1. **Use for the right workload**: Read-heavy, write-rarely patterns
2. **Avoid frequent iteration**: Range() is expensive
3. **Handle type assertions safely**: Always check return values
4. **Monitor memory usage**: Dirty map can grow large
5. **Consider alternatives**: For write-heavy workloads, regular maps might be better

## Conclusion

`sync.Map` is a sophisticated data structure that trades complexity for performance in specific scenarios. Understanding its internals helps in:

- Making informed decisions about when to use it
- Optimizing performance for your specific use case
- Debugging issues related to concurrent map access
- Understanding the trade-offs involved

The key insight is that `sync.Map` is not a general-purpose concurrent map but rather a specialized tool optimized for read-heavy workloads with infrequent writes.
