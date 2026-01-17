# Go Channels: Deadlocks and Race Conditions

Understanding and avoiding deadlocks and race conditions is crucial when working with Go channels and concurrent programming.

## Deadlocks

A deadlock occurs when all goroutines are blocked waiting for something that can never happen. The Go runtime will panic with a "deadlock" message when it detects this situation. Deadlocks are particularly dangerous because they completely halt program execution.

### Common Deadlock Scenarios

1. **Unbuffered Channel with No Receiver**
   - Description: This is the most common deadlock scenario. When you try to send on an unbuffered channel, the sender blocks until there's a receiver. If there's no receiver, the program deadlocks.
   - Impact: Program immediately panics with "fatal error: all goroutines are asleep - deadlock!"

    ```go
    func main() {
        ch := make(chan int)
        ch <- 1  // Deadlock: no goroutine available to receive
    }
    ```

    Fix - Create a receiver in a separate goroutine:
    
    ```go
    func main() {
        ch := make(chan int)
        go func() {
            ch <- 1  // Send in separate goroutine
        }()
        fmt.Println(<-ch)  // Receive in main goroutine
    }
    ```

2. **Waiting for Send on Full Buffer**
   - Description: Buffered channels block when full. If you try to send to a full channel and there's no receiver to free up space, you get a deadlock.
   - Impact: Program blocks indefinitely when buffer is full and no receivers are available.

    ```go
    func main() {
        ch := make(chan int, 1)
        ch <- 1  // Buffer full
        ch <- 2  // Deadlock: buffer full and no receiver
    }
    ```

3. **Circular Wait**
   - Description: Occurs when two or more goroutines are waiting for each other in a circular fashion.
   - Impact: All involved goroutines become permanently blocked, causing a system-wide deadlock.

    ```go
    func main() {
        ch1 := make(chan int)
        ch2 := make(chan int)
        
        go func() {
            v := <-ch1  // Wait for ch1
            ch2 <- v    // Then send to ch2
        }()
        
        v := <-ch2    // Wait for ch2
        ch1 <- v      // Then send to ch1
        // Deadlock: both goroutines wait for each other
    }
    ```

### Preventing Deadlocks

1. **Always ensure there's a receiver for every sender**
   - Use select statements with default cases for non-blocking sends
   - Implement timeouts to prevent indefinite blocking

2. **Use buffered channels when appropriate**
   - Buffer size should match expected channel usage patterns
   - Don't use buffers to avoid designing proper synchronization

3. **Implement timeouts using `select` and `time.After()`**
   - Prevents infinite waiting
   - Allows graceful error handling

4. **Use context for cancellation**
    ```go
    func worker(ctx context.Context, ch chan int) {
        select {
        case value := <-ch:
            // Process value
        case <-ctx.Done():
            return // Cancellation requested
        }
    }
    ```

## Race Conditions

A race condition occurs when multiple goroutines access shared data concurrently, and at least one of them is writing. Race conditions can cause unpredictable behavior and subtle bugs that are hard to reproduce.

### Common Race Scenarios

1. **Shared Variable Access**
   - Description: Multiple goroutines trying to modify the same variable without synchronization
   - Impact: Inconsistent results, lost updates, and corrupted data

    ```go
    func main() {
        counter := 0
        
        // RACE CONDITION: multiple goroutines updating counter
        for i := 0; i < 1000; i++ {
            go func() {
                counter++  // Unsafe concurrent access
            }()
        }
    }
    ```

    Fix using channels for synchronized access:
    ```go
    func main() {
        counter := 0
        ch := make(chan int, 1000)
        
        for i := 0; i < 1000; i++ {
            go func() {
                ch <- 1  // Synchronized send
            }()
        }
        
        // Count sequentially
        for i := 0; i < 1000; i++ {
            counter += <-ch  // Synchronized receive
        }
    }
    ```

2. **Concurrent Map Access**
   - Description: Go maps are not safe for concurrent access. Multiple goroutines writing to a map can cause memory corruption.
   - Impact: Program may panic with "concurrent map writes" or exhibit undefined behavior

    ```go
    // RACE CONDITION: concurrent map writes
    func main() {
        m := make(map[int]int)
        
        for i := 0; i < 100; i++ {
            go func(i int) {
                m[i] = i  // Unsafe concurrent map access
            }(i)
        }
    }
    ```

    Fix using mutex for synchronized access:
    ```go
    func main() {
        var mu sync.Mutex
        m := make(map[int]int)
        
        for i := 0; i < 100; i++ {
            go func(i int) {
                mu.Lock()
                m[i] = i  // Protected by mutex
                mu.Unlock()
            }(i)
        }
    }
    ```

### Detecting Race Conditions

1. Use Go's race detector:
    ```bash
    go run -race main.go    # Detect races during execution
    go test -race ./...     # Detect races during testing
    ```

### Best Practices to Avoid Race Conditions

1. **Use channels to communicate between goroutines**
   - Channels provide synchronized communication
   - Follow "Don't communicate by sharing memory; share memory by communicating"

2. **When needed, use sync package primitives**
   - sync.Mutex for exclusive access
   - sync.RWMutex for read/write scenarios
   - atomic package for simple operations

3. **Make data immutable when possible**
   - Immutable data can't cause race conditions
   - Pass copies instead of pointers when practical

4. **Minimize shared state**
   - Keep shared data to a minimum
   - Use local variables when possible

## General Tips

1. **Always use `go vet` and the race detector during development**
   - Catches common concurrent programming mistakes
   - Should be part of CI/CD pipeline

2. **Implement timeouts for channel operations**
   - Prevents indefinite blocking
   - Makes system more resilient

3. **Use `select` statements for safe channel handling**
   - Manages multiple channels safely
   - Provides timeout and cancellation options

4. **Close channels properly**
   - Signal completion to receivers
   - Prevent goroutine leaks

5. **Use worker pools for concurrent operations**
   - Controls resource usage
   - Prevents goroutine explosion

6. **Consider using errgroup for concurrent error handling**
   - Manages multiple goroutines
   - Propagates errors properly 