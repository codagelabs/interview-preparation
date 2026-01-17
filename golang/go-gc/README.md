# Go Garbage Collector (GC)

> **Note**: This guide is based on the [official Go GC guide](https://tip.golang.org/doc/gc-guide) and describes the garbage collector as of Go 1.19+. This document applies to the standard Go toolchain (`gc` compiler) and may not apply to other implementations like gccgo or Gollvm.

## Table of Contents
1. [Introduction](#introduction)
2. [Core Components](#core-components)
3. [How GC Works](#how-gc-works)
4. [GC Cycle and Phases](#gc-cycle-and-phases)
5. [Performance Optimization](#performance-optimization)
6. [Best Practices](#best-practices)
7. [Monitoring and Debugging](#monitoring-and-debugging)
8. [Advanced Topics](#advanced-topics)

## Introduction

Garbage collection is a form of automatic memory management that aims to reclaim memory occupied by objects that are no longer in use by the program.  Go's garbage collector is a **concurrent, tri-color, mark-sweep collector** that automatically manages memory allocation and deallocation iTS A PART OF GO RUNTIME The GC is designed to:

- **Minimize latency**: Keep pause times as low as possible archive using concurancy while sweeping memeory
- **Maximize throughput**: Efficiently reclaim memory
- **Scale with hardware**: Utilize multiple CPU cores
- **Be tunable**: Configurable via environment variables

### Where Go Values Live

Before diving into GC details, it's important to understand where Go values are stored:

#### Stack Allocation
- **Non-pointer values** in local variables are typically allocated on the goroutine stack
- The Go compiler can predetermine when this memory can be freed
- More efficient than heap allocation as cleanup is automatic

#### Heap Allocation (GC Managed)
- Values that **escape to the heap** are managed by the GC
- Occurs when the compiler cannot determine the lifetime of a value
- Examples: dynamically sized slices, values referenced by pointers that escape

#### Escape Analysis
The Go compiler performs escape analysis to determine if a value should be stack-allocated or heap-allocated:

```go
// Stack allocated (no escape)
func stackExample() int {
    x := 42
    return x
}

// Heap allocated (escapes)
func heapExample() *int {
    x := 42
    return &x  // x escapes to heap
}
```

### Detailed Examples: Stack vs Heap Allocation

#### **Stack Allocation Examples**

**1. Simple Values (Stack)**
```go
// These values are allocated on the stack
var (
    integer int = 42
    float   float64 = 3.14
    boolean bool = true
    str     string = "hello"
    array   [5]int = [5]int{1, 2, 3, 4, 5}
    point   struct {
        x, y int
        name string
    } = struct {
        x, y int
        name string
    }{10, 20, "point"}
)
```

**2. Small Slices with Known Size (Stack)**
```go
// Small slices with known size at compile time
smallSlice := make([]int, 10)
for i := range smallSlice {
    smallSlice[i] = i
}
```

**3. Function Parameters and Return Values (Stack)**
```go
func functionParameterExample(value int) int {
    // Parameter 'value' is stack allocated
    // Return value is also stack allocated
    return value * 2
}
```

#### **Heap Allocation Examples**

**1. Values That Escape (Heap)**
```go
func heapAllocationExample() *int {
    // This value escapes to heap because it's returned
    value := 42
    return &value // This causes the value to escape to heap
}
```

**2. Dynamic Slices (Heap)**
```go
func heapAllocationDynamicSlices(size int) {
    // This slice escapes to heap because size is determined at runtime
    dynamicSlice := make([]int, size)
    for i := range dynamicSlice {
        dynamicSlice[i] = i
    }
}
```

**3. Large Values (Heap)**
```go
func heapAllocationLargeValues() {
    // Large arrays escape to heap
    largeArray := [1000]int{}
    for i := range largeArray {
        largeArray[i] = i
    }
}
```

**4. Values Referenced by Pointers (Heap)**
```go
func heapAllocationPointerReferences() {
    type Person struct {
        Name string
        Age  int
    }

    // This struct escapes to heap because we take its address
    person := Person{Name: "Alice", Age: 30}
    personPtr := &person
}
```

**5. Interface Values (Heap)**
```go
func interfaceAllocationExample() {
    // Interface values are heap allocated
    var interfaceValue interface{} = "hello"
    var interfaceSlice []interface{} = []interface{}{1, "world", 3.14}
}
```

**6. Channels and Maps (Heap)**
```go
func channelMapAllocationExample() {
    // Channels and maps are always heap allocated
    ch := make(chan int, 10)
    m := make(map[string]int)
    m["key"] = 42
}
```

#### **Mixed Allocation Examples**

**1. Taking Addresses Forces Heap Allocation**
```go
func mixedAllocationExample() {
    // Stack allocated
    stackValue := 100
    stackArray := [5]int{1, 2, 3, 4, 5}

    // Heap allocated (because we take their addresses)
    heapValue := 200
    heapArray := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

    // Taking addresses forces heap allocation
    stackPtr := &stackValue
    stackArrayPtr := &stackArray
    heapPtr := &heapValue
    heapArrayPtr := &heapArray
}
```

**2. Slice Allocation Patterns**
```go
func sliceAllocationPatterns() {
    // Stack allocated (small, known size)
    smallSlice := make([]int, 5)
    
    // Heap allocated (large size)
    largeSlice := make([]int, 1000)
    
    // Heap allocated (dynamic size)
    dynamicSlice := make([]int, runtime.NumCPU())
    
    // Heap allocated (capacity > size)
    capacitySlice := make([]int, 5, 100)
}
```

#### **Escape Analysis Demonstration**

```go
func escapeAnalysisDemo() {
    // Case 1: No escape - value stays on stack
    noEscape := 42
    fmt.Printf("No escape: %d (stack allocated)\n", noEscape)

    // Case 2: Escape - value moves to heap
    escape := 42
    escapePtr := &escape
    fmt.Printf("Escape: %d (heap allocated, ptr: %p)\n", *escapePtr, escapePtr)

    // Case 3: Interface escape
    var interfaceValue interface{} = 42
    fmt.Printf("Interface escape: %v (heap allocated)\n", interfaceValue)

    // Case 4: Slice escape
    slice := make([]int, 10)
    slicePtr := &slice
    fmt.Printf("Slice escape: %v (heap allocated, ptr: %p)\n", *slicePtr, slicePtr)
}
```

#### **Memory Usage Monitoring**

```go
func memoryUsageDemo() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Initial heap stats:\n")
    fmt.Printf("  Heap Alloc: %d bytes\n", m.HeapAlloc)
    fmt.Printf("  Heap Sys: %d bytes\n", m.HeapSys)
    fmt.Printf("  Heap Objects: %d\n", m.HeapObjects)

    // Allocate some heap memory
    heapSlice := make([]int, 10000)
    for i := range heapSlice {
        heapSlice[i] = i
    }

    runtime.ReadMemStats(&m)
    fmt.Printf("\nAfter heap allocation:\n")
    fmt.Printf("  Heap Alloc: %d bytes\n", m.HeapAlloc)
    fmt.Printf("  Heap Sys: %d bytes\n", m.HeapSys)
    fmt.Printf("  Heap Objects: %d\n", m.HeapObjects)

    // Force garbage collection
    runtime.GC()

    runtime.ReadMemStats(&m)
    fmt.Printf("\nAfter garbage collection:\n")
    fmt.Printf("  Heap Alloc: %d bytes\n", m.HeapAlloc)
    fmt.Printf("  Heap Sys: %d bytes\n", m.HeapSys)
    fmt.Printf("  Heap Objects: %d\n", m.HeapObjects)
}
```

#### **Key Allocation Rules**

| **Stack Allocation** | **Heap Allocation** |
|---------------------|-------------------|
| Small, simple values | Large values |
| Known size at compile time | Dynamic size |
| No escape from function | Values that escape |
| Function parameters | Returned pointers |
| Small arrays/slices | Large arrays/slices |
| Small structs | Large structs |
| | Interface values |
| | Channels and maps |
| | Values with pointers |
| | Concatenated strings |

#### **Running the Examples**

To see these examples in action, run the provided sample program:

```bash
cd go-gc
go run stack_vs_heap_example.go
```

This will demonstrate:
- Stack vs heap allocation patterns
- Escape analysis in action
- Memory usage monitoring
- Size calculations for different types
- Pointer and interface behavior

## Core Components

### 1. Tracing Garbage Collection

Go uses **tracing garbage collection**, which identifies live objects by following pointers transitively from **roots**.

#### Key Terms:
- **Object**: A dynamically allocated piece of memory containing one or more Go values
- **Pointer**: A memory address that references a value within an object
- **Object Graph**: The network of objects connected by pointers
- **Roots**: Pointers that identify definitely in-use objects (local variables, global variables)
- **Reachable**: Objects that can be discovered by scanning from roots

### 2. Mark-Sweep Algorithm

Go's GC uses a **non-moving mark-sweep** approach:

#### Mark Phase
- Walks the object graph starting from roots
- Marks all reachable objects as "live"
- Uses **tri-color marking** for efficiency

#### Sweep Phase
- Scans all heap memory
- Reclaims memory from unmarked (dead) objects
- Makes reclaimed memory available for new allocations

### 3. Tri-Color Marking

The GC uses three colors to track object states during marking:

- **White**: Potentially unreachable objects (initial state)
- **Gray**: Objects being processed (discovered but not fully scanned)
- **Black**: Reachable objects (fully scanned and marked as live)

```go
// Example of tri-color marking process
// 1. All objects start as WHITE
// 2. Roots are marked GRAY
// 3. GRAY objects are scanned, their children marked GRAY
// 4. Fully scanned objects become BLACK
// 5. After marking, WHITE objects are garbage
```

### 4. Write Barriers

Write barriers ensure consistency during concurrent marking:

- **Tracks pointer writes** during the mark phase
- **Maintains object reachability** when pointers are modified
- **Prevents live objects from being missed** due to concurrent updates

```go
// Write barrier example
// During marking, when a pointer is written:
oldPtr := obj.field
obj.field = newPtr  // Write barrier tracks this change
// GC ensures newPtr is marked if it's reachable
```

## How GC Works

### 1. GC Triggers

The GC is triggered by several conditions:

```go
// GC triggers when:
// 1. Heap size reaches GOGC threshold (default: 100%)
// 2. Manual trigger via runtime.GC()
// 3. Memory pressure from OS
// 4. Time-based triggers (Go 1.19+)
```

#### GOGC Parameter
The `GOGC` environment variable controls when GC runs:

```bash
export GOGC=100  # Default: GC when heap grows by 100%
export GOGC=50   # More aggressive: GC when heap grows by 50%
export GOGC=200  # Less aggressive: GC when heap grows by 200%
```

**Mathematical Relationship**:
- Doubling `GOGC` doubles heap memory overhead
- Doubling `GOGC` halves GC CPU cost
- Formula: `Target heap = Live heap + (Live heap + GC roots) × GOGC / 100`

### 2. Memory Management Model

#### Heap Target Calculation
The GC maintains a target heap size based on live memory:

```
Target heap memory = Live heap + (Live heap + GC roots) × GOGC / 100
```

#### GC CPU Cost Model
The cost of a GC cycle can be modeled as:

```
GC CPU cost per cycle = (Live heap + GC roots) × (Cost per byte) + Fixed cost
```

### 3. Concurrent Execution

Go's GC is **concurrent**, meaning it runs alongside your program:

- **Mark phase**: Runs concurrently with application code
- **Sweep phase**: Runs concurrently with application code
- **Minimal STW pauses**: Only brief "stop-the-world" pauses for setup/termination

## GC Cycle and Phases

### 1. GC Cycle Overview

The GC operates in a continuous cycle with three phases:

1. **Sweep Phase**: Reclaiming memory from previous cycle
2. **Off Phase**: No GC work (application runs normally)
3. **Mark Phase**: Identifying live objects

### 2. Detailed Phase Breakdown

#### Phase 1: Mark Setup (STW)
- **Duration**: Very brief (microseconds)
- **Activities**:
  - Stop all goroutines
  - Enable write barriers
  - Initialize marking state
  - Start concurrent marking

#### Phase 2: Concurrent Mark
- **Duration**: Variable (depends on heap size)
- **Activities**:
  - Scan roots (goroutine stacks, globals)
  - Follow pointers through object graph
  - Mark reachable objects
  - Handle write barriers for pointer updates

#### Phase 3: Mark Termination (STW)
- **Duration**: Brief (microseconds to milliseconds)
- **Activities**:
  - Complete marking process
  - Disable write barriers
  - Prepare for sweep phase

#### Phase 4: Concurrent Sweep
- **Duration**: Variable (depends on dead object count)
- **Activities**:
  - Scan heap for unmarked objects
  - Reclaim memory from dead objects
  - Update free lists for allocation

### 3. GC Timing and Coordination

#### Mark Phase Timing
The mark phase is divided into sub-phases:

1. **Mark Setup**: Brief STW to initialize
2. **Concurrent Mark**: Main marking work
3. **Mark Termination**: Brief STW to complete

#### Sweep Phase Timing
Sweeping happens concurrently and is coordinated with allocation:

- **Background sweeping**: Continuous cleanup
- **Allocation-triggered sweeping**: Immediate cleanup when allocating

### 4. GC Tuning Parameters

```bash
# Core GC environment variables
export GOGC=100        # GC trigger ratio (default: 100)
export GOMEMLIMIT=1GiB # Memory limit (Go 1.19+)
export GOMAXPROCS=4    # CPU cores for GC work

# Debug and monitoring
export GODEBUG=gctrace=1    # Enable GC trace output
export GODEBUG=scavtrace=1  # Enable scavenger trace
```

## Performance Optimization

### 1. Understanding GC Costs

#### Memory Overhead
- **Live heap**: Memory your application actually needs
- **New heap**: Additional memory for allocation during GC cycle
- **GC roots**: Small amount of memory for root scanning

#### CPU Overhead
- **Marking cost**: Proportional to live heap size
- **Scanning cost**: Proportional to pointer density
- **Write barrier cost**: Proportional to pointer writes

### 2. GC Performance Metrics

#### Key Metrics to Monitor:
- **GC Frequency**: How often GC runs
- **GC Duration**: Time spent in GC phases
- **GC CPU Usage**: CPU time used by GC
- **Memory Allocation Rate**: Objects allocated per second
- **Heap Size**: Current and peak heap usage
- **Pause Times**: STW pause durations

### 3. GC Optimization Techniques

#### Reduce Allocation Rate
```go
// Bad: High allocation rate
func processItems(items []int) {
    for _, item := range items {
        result := make([]int, 0) // New allocation each time
        // process item
    }
}

// Good: Reuse allocations
func processItems(items []int) {
    result := make([]int, 0, len(items)) // Pre-allocate
    for _, item := range items {
        result = result[:0] // Reset slice without reallocation
        // process item
    }
}
```

#### Object Pooling for High-Frequency Objects
```go
import "sync"

type ObjectPool struct {
    pool sync.Pool
}

func NewObjectPool() *ObjectPool {
    return &ObjectPool{
        pool: sync.Pool{
            New: func() interface{} {
                return &MyObject{}
            },
        },
    }
}

func (p *ObjectPool) Get() *MyObject {
    return p.pool.Get().(*MyObject)
}

func (p *ObjectPool) Put(obj *MyObject) {
    // Reset object state
    obj.Reset()
    p.pool.Put(obj)
}
```

#### Optimize Data Structures
```go
// Bad: Pointer-heavy structures
type Node struct {
    Data    *Data
    Next    *Node
    Prev    *Node
    Parent  *Node
    Children []*Node
}

// Good: Reduce pointer density
type Node struct {
    Data    Data      // Embed instead of pointer
    Next    *Node     // Keep only necessary pointers
    Prev    *Node
    // Remove unnecessary pointers
}
```

### 4. Memory Profiling and Analysis

#### Using pprof for Heap Analysis
```go
import (
    "net/http"
    _ "net/http/pprof"
    "runtime/pprof"
)

func main() {
    // Start pprof server
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()
    
    // Your application code
}

// Generate heap profile
func generateHeapProfile() {
    f, err := os.Create("heap.prof")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    pprof.WriteHeapProfile(f)
}
```

#### Runtime Statistics Monitoring
```go
import "runtime"

func printMemStats() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Heap Alloc: %d MB\n", m.HeapAlloc/1024/1024)
    fmt.Printf("Heap Sys: %d MB\n", m.HeapSys/1024/1024)
    fmt.Printf("Heap Idle: %d MB\n", m.HeapIdle/1024/1024)
    fmt.Printf("Heap Inuse: %d MB\n", m.HeapInuse/1024/1024)
    fmt.Printf("Heap Released: %d MB\n", m.HeapReleased/1024/1024)
    fmt.Printf("Heap Objects: %d\n", m.HeapObjects)
    fmt.Printf("GC Cycles: %d\n", m.NumGC)
    fmt.Printf("GC CPU Fraction: %.2f%%\n", m.GCCPUFraction*100)
    fmt.Printf("GC Pause Total: %d ms\n", m.PauseTotalNs/1000000)
}
```

## Best Practices

### 1. Memory Management

#### Pre-allocate Slices and Maps
```go
// Bad: Dynamic growth causes reallocations
slice := make([]int, 0)
for i := 0; i < 1000; i++ {
    slice = append(slice, i) // Multiple reallocations
}

// Good: Pre-allocate capacity
slice := make([]int, 0, 1000)
for i := 0; i < 1000; i++ {
    slice = append(slice, i) // No reallocations
}
```

#### Use Structs Instead of Maps for Small Data
```go
// Bad: Map overhead for small data
type Config struct {
    data map[string]string
}

// Good: Struct fields for small data
type Config struct {
    Host     string
    Port     int
    Username string
    Password string
}
```

### 2. Avoid Memory Leaks

#### Close Resources Properly
```go
// Always close files, connections, etc.
file, err := os.Open("file.txt")
if err != nil {
    return err
}
defer file.Close()
```

#### Clear Large Object References
```go
// Clear large objects when done
func processLargeData() {
    data := make([]byte, 1000000)
    // Process data...
    
    // Clear reference when done
    data = nil
}
```

### 3. GC-Friendly Patterns

#### Use Value Receivers for Small Structs
```go
type Point struct {
    X, Y int
}

// Good: Value receiver for small struct
func (p Point) Distance(other Point) float64 {
    // Implementation
}

// Bad: Pointer receiver for small struct
func (p *Point) Distance(other *Point) float64 {
    // Implementation
}
```

#### Group Pointer Fields in Structs
```go
// Good: Pointers at beginning of struct
type User struct {
    Name     *string  // Pointer fields first
    Email    *string
    Age      int      // Non-pointer fields after
    Active   bool
}

// This helps GC scanning efficiency
```

#### Minimize Pointer Indirections
```go
// Bad: Multiple pointer indirections
type Config struct {
    *DatabaseConfig
    *CacheConfig
}

// Good: Embed structs directly
type Config struct {
    DatabaseConfig
    CacheConfig
}
```

### 4. Concurrent Programming

#### Use Channels for Communication
```go
// Good: Use channels for goroutine communication
func worker(jobs <-chan int, results chan<- int) {
    for job := range jobs {
        results <- process(job)
    }
}
```

#### Avoid Shared Mutable State
```go
// Bad: Shared mutable state
var globalCounter int
var mu sync.Mutex

func increment() {
    mu.Lock()
    globalCounter++
    mu.Unlock()
}

// Good: Use channels or local state
func worker(id int, results chan<- int) {
    localCounter := 0
    // Process work...
    results <- localCounter
}
```

## Monitoring and Debugging

### 1. GC Debugging

#### Enable GC Debugging
```bash
# Enable GC debugging
export GODEBUG=gctrace=1

# Run your program
go run main.go
```

#### GC Trace Output Interpretation
```
gc 1 @0.123s 0%: 0.002+0.5+0.002 ms clock, 0.008+0/0.002/0.5+0.010 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
```

**Trace Format**:
- `gc 1`: GC cycle number
- `@0.123s`: Time since program start
- `0%`: CPU percentage used by GC
- `0.002+0.5+0.002 ms`: STW + concurrent + STW times
- `4->4->0 MB`: Heap before -> after -> live heap
- `5 MB goal`: Target heap size
- `4 P`: Number of processors

### 2. Memory Profiling Tools

#### Using go tool pprof
```bash
# Get heap profile
curl -o heap.prof http://localhost:6060/debug/pprof/heap

# Analyze heap profile
go tool pprof heap.prof

# Interactive analysis
go tool pprof -http=:8080 heap.prof

# Get goroutine profile
curl -o goroutine.prof http://localhost:6060/debug/pprof/goroutine
```

#### Using go tool trace
```go
import "runtime/trace"

func main() {
    f, err := os.Create("trace.out")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    
    trace.Start(f)
    defer trace.Stop()
    
    // Your application code
}

// Analyze trace
// go tool trace trace.out
```

### 3. Performance Monitoring

#### Real-time Monitoring
```go
import (
    "time"
    "runtime"
)

func monitorGC() {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        
        fmt.Printf("GC CPU: %.2f%%, Heap: %d MB, Pause: %.2f ms\n", 
            m.GCCPUFraction*100, 
            m.HeapAlloc/1024/1024,
            float64(m.PauseNs[(m.NumGC+255)%256])/1000000)
    }
}
```

### 4. GC Tuning Guidelines

#### For Low Latency Applications
```bash
# Reduce GC frequency for lower latency
export GOGC=50

# Set memory limit to prevent OOM
export GOMEMLIMIT=512MiB

# Monitor pause times
export GODEBUG=gctrace=1
```

#### For High Throughput Applications
```bash
# Allow more memory usage for higher throughput
export GOGC=200

# Increase memory limit
export GOMEMLIMIT=2GiB

# Use more CPU cores for GC
export GOMAXPROCS=8
```

## Advanced Topics

### 1. Implementation-Specific Optimizations

#### Pointer-Free Value Segregation
The GC segregates pointer-free values from other values, reducing cache pressure:

```go
// Good: Reduce pointer density
type Data struct {
    ID       int64   // Non-pointer fields first
    Timestamp int64
    Value     float64
    Name      string  // String is pointer-heavy
    Metadata  *Meta   // Pointer fields last
}
```

#### Pointer Field Grouping
Group pointer fields at the beginning of structs to help GC scanning:

```go
// Good: Pointers at beginning
type User struct {
    Name     *string  // Pointer fields first
    Email    *string
    Age      int      // Non-pointer fields after
    Active   bool
}
```

### 2. Linux Transparent Huge Pages (THP)

For production Go applications on Linux, THP can improve performance:

#### Benefits
- **Bigger heaps (1 GiB+)**: Up to 10% throughput improvement
- **Smaller heaps**: May use 50% more memory

#### Recommended Settings
```bash
# Enable THP
echo always > /sys/kernel/mm/transparent_hugepage/enabled

# Set defrag to defer
echo defer > /sys/kernel/mm/transparent_hugepage/defrag

# Set max_ptes_none to 0 (Go 1.21+)
echo 0 > /sys/kernel/mm/transparent_hugepage/khugepaged/max_ptes_none
```

#### Disable THP if Needed
```bash
# Disable THP for heap memory (Go 1.21.6+)
export GODEBUG=disablethp=1

# Or disable at process level
prctl --thp-disable
```

### 3. GC Mathematical Model

#### Heap Target Formula
```
Target heap memory = Live heap + (Live heap + GC roots) × GOGC / 100
```

#### GC CPU Cost Formula
```
Total GC CPU cost = (Allocation rate) / (GOGC / 100) × (Cost per byte) × T
```

#### Trade-off Relationship
- **Doubling GOGC**: Doubles heap memory overhead, halves GC CPU cost
- **Reducing allocation rate**: Reduces both memory and CPU overhead
- **Optimizing data structures**: Reduces scanning cost per byte

## Summary

Go's garbage collector is designed to be:
- **Concurrent**: Runs alongside your program with minimal pauses
- **Non-generational**: Uses mark-and-sweep algorithm
- **Tunable**: Configurable via environment variables
- **Low-latency**: Sub-millisecond pause times
- **Scalable**: Utilizes multiple CPU cores

### Key Takeaways:
1. **Monitor GC performance** using built-in tools (pprof, trace, runtime stats)
2. **Reduce allocations** and use object pooling for high-frequency objects
3. **Pre-allocate slices and maps** when possible
4. **Use appropriate data structures** to minimize pointer density
5. **Profile memory usage** regularly to identify optimization opportunities
6. **Tune GC parameters** based on your application's latency/throughput requirements
7. **Consider THP settings** for production Linux deployments

### Performance Optimization Hierarchy:
1. **Measure first**: Use profiling tools to identify bottlenecks
2. **Reduce allocation rate**: Most effective optimization
3. **Optimize data structures**: Reduce pointer density
4. **Use object pooling**: For frequently allocated objects
5. **Tune GC parameters**: Fine-tune for your specific use case

## Additional Resources

- [Official Go GC Guide](https://tip.golang.org/doc/gc-guide) - Comprehensive technical guide
- [Go GC Documentation](https://golang.org/pkg/runtime/) - Runtime package documentation
- [Go Memory Management](https://golang.org/doc/effective_go.html#memory) - Effective Go memory section
- [pprof Documentation](https://golang.org/pkg/net/http/pprof/) - Profiling tools
- [Go Performance Blog](https://blog.golang.org/profiling-go-programs) - Performance optimization articles
- [Go GC Tuning](https://blog.golang.org/go15gc) - GC improvements in Go 1.5+

## How GC Works

### 1. GC Triggers
```go
// GC is triggered when:
// - Heap size reaches GOGC threshold (default: 100%)
// - Manual trigger via runtime.GC()
// - Memory pressure
```

### 2. GC Phases

#### Phase 1: Mark Setup
- Stops the world (STW) - very brief
- Enables write barriers
- Starts concurrent marking

#### Phase 2: Concurrent Mark
- Runs alongside your program
- Marks reachable objects
- Uses write barriers to track pointer changes

#### Phase 3: Mark Termination
- Brief STW pause
- Completes marking
- Disables write barriers

#### Phase 4: Concurrent Sweep
- Reclaims unreachable objects
- Runs concurrently with program

### 3. GC Tuning Parameters

```bash
# Environment variables for GC tuning
export GOGC=100        # GC trigger ratio (default: 100)
export GOMEMLIMIT=1GiB # Memory limit
export GOMAXPROCS=4    # CPU cores for GC
```

## Performance Optimization

### 1. GC Performance Metrics

#### Key Metrics to Monitor:
- **GC Frequency**: How often GC runs
- **GC Duration**: Time spent in GC
- **GC CPU Usage**: CPU time used by GC
- **Memory Allocation Rate**: Objects allocated per second
- **Heap Size**: Current and peak heap usage

### 2. GC Optimization Techniques

#### Reduce Allocations
```go
// Bad: Creates new slice on each iteration
func processItems(items []int) {
    for _, item := range items {
        result := make([]int, 0) // New allocation each time
        // process item
    }
}

// Good: Reuse slice
func processItems(items []int) {
    result := make([]int, 0, len(items)) // Pre-allocate
    for _, item := range items {
        result = result[:0] // Reset slice without reallocation
        // process item
    }
}
```

#### Object Pooling
```go
import "sync"

type ObjectPool struct {
    pool sync.Pool
}

func NewObjectPool() *ObjectPool {
    return &ObjectPool{
        pool: sync.Pool{
            New: func() interface{} {
                return &MyObject{}
            },
        },
    }
}

func (p *ObjectPool) Get() *MyObject {
    return p.pool.Get().(*MyObject)
}

func (p *ObjectPool) Put(obj *MyObject) {
    // Reset object state
    obj.Reset()
    p.pool.Put(obj)
}
```

#### Avoid Unnecessary Allocations
```go
// Bad: String concatenation creates new strings
func buildString(items []string) string {
    result := ""
    for _, item := range items {
        result += item // Creates new string each time
    }
    return result
}

// Good: Use strings.Builder
func buildString(items []string) string {
    var builder strings.Builder
    for _, item := range items {
        builder.WriteString(item)
    }
    return builder.String()
}
```

### 3. Memory Profiling

#### Using pprof
```go
import (
    "net/http"
    _ "net/http/pprof"
    "runtime/pprof"
)

func main() {
    go func() {
        http.ListenAndServe("localhost:6060", nil)
    }()
    
    // Your application code
}

// Generate heap profile
func generateHeapProfile() {
    f, err := os.Create("heap.prof")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    pprof.WriteHeapProfile(f)
}
```

#### Runtime Statistics
```go
import "runtime"

func printMemStats() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Heap Alloc: %d MB\n", m.HeapAlloc/1024/1024)
    fmt.Printf("Heap Sys: %d MB\n", m.HeapSys/1024/1024)
    fmt.Printf("Heap Idle: %d MB\n", m.HeapIdle/1024/1024)
    fmt.Printf("Heap Inuse: %d MB\n", m.HeapInuse/1024/1024)
    fmt.Printf("Heap Released: %d MB\n", m.HeapReleased/1024/1024)
    fmt.Printf("Heap Objects: %d\n", m.HeapObjects)
    fmt.Printf("GC Cycles: %d\n", m.NumGC)
    fmt.Printf("GC CPU Fraction: %.2f%%\n", m.GCCPUFraction*100)
}
```

## Best Practices

### 1. Memory Management

#### Pre-allocate Slices and Maps
```go
// Bad: Dynamic growth
slice := make([]int, 0)
for i := 0; i < 1000; i++ {
    slice = append(slice, i) // Multiple reallocations
}

// Good: Pre-allocate capacity
slice := make([]int, 0, 1000)
for i := 0; i < 1000; i++ {
    slice = append(slice, i) // No reallocations
}
```

#### Use Structs Instead of Maps for Small Data
```go
// Bad: Map for small data
type Config struct {
    data map[string]string
}

// Good: Struct fields
type Config struct {
    Host     string
    Port     int
    Username string
    Password string
}
```

### 2. Avoid Memory Leaks

#### Close Resources
```go
// Always close files, connections, etc.
file, err := os.Open("file.txt")
if err != nil {
    return err
}
defer file.Close()
```

#### Clear References
```go
// Clear large objects when done
func processLargeData() {
    data := make([]byte, 1000000)
    // Process data...
    
    // Clear reference when done
    data = nil
}
```

### 3. GC-Friendly Patterns

#### Use Value Receivers for Small Structs
```go
type Point struct {
    X, Y int
}

// Good: Value receiver for small struct
func (p Point) Distance(other Point) float64 {
    // Implementation
}

// Bad: Pointer receiver for small struct
func (p *Point) Distance(other *Point) float64 {
    // Implementation
}
```

#### Minimize Pointer Indirections
```go
// Bad: Multiple pointer indirections
type Config struct {
    *DatabaseConfig
    *CacheConfig
}

// Good: Embed structs directly
type Config struct {
    DatabaseConfig
    CacheConfig
}
```

### 4. Concurrent Programming

#### Use Channels for Communication
```go
// Good: Use channels for goroutine communication
func worker(jobs <-chan int, results chan<- int) {
    for job := range jobs {
        results <- process(job)
    }
}
```

#### Avoid Shared Mutable State
```go
// Bad: Shared mutable state
var globalCounter int
var mu sync.Mutex

func increment() {
    mu.Lock()
    globalCounter++
    mu.Unlock()
}

// Good: Use channels or local state
func worker(id int, results chan<- int) {
    localCounter := 0
    // Process work...
    results <- localCounter
}
```

## Monitoring and Debugging

### 1. GC Debugging

#### Enable GC Debugging
```bash
# Enable GC debugging
export GODEBUG=gctrace=1

# Run your program
go run main.go
```

#### GC Trace Output
```
gc 1 @0.123s 0%: 0.002+0.5+0.002 ms clock, 0.008+0/0.002/0.5+0.010 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
```

### 2. Memory Profiling Tools

#### Using go tool pprof
```bash
# Get heap profile
curl -o heap.prof http://localhost:6060/debug/pprof/heap

# Analyze heap profile
go tool pprof heap.prof

# Get goroutine profile
curl -o goroutine.prof http://localhost:6060/debug/pprof/goroutine
```

#### Using go tool trace
```go
import "runtime/trace"

func main() {
    f, err := os.Create("trace.out")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()
    
    trace.Start(f)
    defer trace.Stop()
    
    // Your application code
}
```

### 3. Performance Monitoring

#### Real-time Monitoring
```go
import (
    "time"
    "runtime"
)

func monitorGC() {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        
        fmt.Printf("GC CPU: %.2f%%, Heap: %d MB\n", 
            m.GCCPUFraction*100, m.HeapAlloc/1024/1024)
    }
}
```

### 4. GC Tuning Guidelines

#### For Low Latency Applications
```bash
# Reduce GC frequency
export GOGC=50

# Set memory limit
export GOMEMLIMIT=512MiB
```

#### For High Throughput Applications
```bash
# Allow more memory usage
export GOGC=200

# Increase memory limit
export GOMEMLIMIT=2GiB
```

## Summary

Go's garbage collector is designed to be:
- **Concurrent**: Runs alongside your program
- **Non-generational**: Uses mark-and-sweep algorithm
- **Tunable**: Configurable via environment variables
- **Low-latency**: Minimal pause times

Key takeaways:
1. Monitor GC performance using built-in tools
2. Reduce allocations and use object pooling
3. Pre-allocate slices and maps when possible
4. Use appropriate data structures
5. Profile memory usage regularly
6. Tune GC parameters based on your application needs

## Additional Resources

- [Go GC Documentation](https://golang.org/pkg/runtime/)
- [Go Memory Management](https://golang.org/doc/effective_go.html#memory)
- [pprof Documentation](https://golang.org/pkg/net/http/pprof/)
- [Go Performance Blog](https://blog.golang.org/profiling-go-programs) 