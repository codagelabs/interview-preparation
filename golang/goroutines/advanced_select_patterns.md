# Advanced Select Patterns in Go

This guide covers advanced patterns for working with `select` and multiple channels in Go, including fan-in, fan-out, and timeout handling.

## Table of Contents
- [Basic Select Pattern](#basic-select-pattern)
- [Fan-In Pattern](#fan-in-pattern)
- [Fan-Out Pattern](#fan-out-pattern)
- [Timeout Patterns](#timeout-patterns)
- [Best Practices](#best-practices)

## Basic Select Pattern

The `select` statement lets you wait on multiple channel operations simultaneously.

```go
func basicSelect(ch1, ch2 chan string, done chan bool) {
    for {
        select {
        case msg1 := <-ch1:
            fmt.Println("Received from ch1:", msg1)
        case msg2 := <-ch2:
            fmt.Println("Received from ch2:", msg2)
        case <-done:
            return
        }
    }
}
```

## Fan-In Pattern

Fan-in combines multiple input channels into a single output channel.

```go
func fanIn(ch1, ch2 <-chan string) <-chan string {
    merged := make(chan string)
    
    // Start goroutine for first input channel
    go func() {
        for msg := range ch1 {
            merged <- msg
        }
    }()
    
    // Start goroutine for second input channel
    go func() {
        for msg := range ch2 {
            merged <- msg
        }
    }()
    
    return merged
}

// Usage example
func fanInExample() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    merged := fanIn(ch1, ch2)
    
    // Now merged receives values from both ch1 and ch2
}
```

## Fan-Out Pattern

Fan-out distributes work across multiple goroutines to process data in parallel.

```go
func fanOut(input <-chan int, workers int) []<-chan int {
    outputs := make([]<-chan int, workers)
    
    for i := 0; i < workers; i++ {
        outputs[i] = worker(input)
    }
    
    return outputs
}

func worker(input <-chan int) <-chan int {
    output := make(chan int)
    
    go func() {
        defer close(output)
        for n := range input {
            // Process data
            result := n * 2
            output <- result
        }
    }()
    
    return output
}
```

## Timeout Patterns

### Request Timeout

```go
func requestWithTimeout(req chan string, timeout time.Duration) (string, error) {
    select {
    case response := <-req:
        return response, nil
    case <-time.After(timeout):
        return "", fmt.Errorf("request timed out after %v", timeout)
    }
}
```

### Periodic Timeout

```go
func periodicWork(interval time.Duration, done chan bool) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            // Do periodic work
            fmt.Println("Performing periodic task")
        case <-done:
            return
        }
    }
}
```

### Default Case

```go
func nonBlockingSelect(ch chan string) {
    select {
    case msg := <-ch:
        fmt.Println("Received:", msg)
    default:
        fmt.Println("No message available")
    }
}
```

## Best Practices

1. **Always Include Done Channel**
```go
func worker(done chan struct{}) {
    for {
        select {
        case <-done:
            return
        default:
            // Do work
        }
    }
}
```

2. **Handle Channel Closure**
```go
func handleClosure(ch <-chan int) {
    for {
        select {
        case val, ok := <-ch:
            if !ok {
                // Channel is closed
                return
            }
            fmt.Println(val)
        }
    }
}
```

3. **Combine Multiple Done Channels**
```go
func mergeSignals(done1, done2 <-chan struct{}) <-chan struct{} {
    done := make(chan struct{})
    go func() {
        defer close(done)
        select {
        case <-done1:
        case <-done2:
        }
    }()
    return done
}
```

## Common Pitfalls to Avoid

1. **Forgetting to Close Channels**
   - Always close channels when no more data will be sent
   - Use defer when appropriate

2. **Deadlocks**
   - Ensure proper channel capacity
   - Always handle all possible select cases

3. **Goroutine Leaks**
   - Always provide a way to terminate goroutines
   - Use context or done channels for cancellation

## Performance Considerations

- Use buffered channels when appropriate
- Consider channel size for performance
- Balance number of goroutines with system resources
- Profile your application to identify bottlenecks

Remember: These patterns are powerful tools but should be used judiciously based on your specific requirements and constraints. 