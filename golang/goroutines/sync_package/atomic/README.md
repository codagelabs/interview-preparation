# sync/atomic

The atomic package provides low-level atomic memory primitives useful for implementing synchronization algorithms.

## Key Operations
1. Add
2. Load
3. Store
4. Swap
5. CompareAndSwap (CAS)

## Supported Types
- Int32, Int64
- Uint32, Uint64
- Uintptr
- Pointer

## Example Usage
```go
var counter uint64

// Increment
atomic.AddUint64(&counter, 1)

// Load
value := atomic.LoadUint64(&counter)

// Store
atomic.StoreUint64(&counter, 10)

// Compare and Swap
swapped := atomic.CompareAndSwapUint64(&counter, 10, 20)
```

## Best Practices
1. Use for simple atomic operations
2. Prefer higher-level synchronization for complex scenarios
3. Useful for implementing lock-free data structures
4. Consider performance implications of atomic operations 