# Unbuffered Channels in Go

Unbuffered channels are a fundamental concurrency primitive in Go that provide synchronous communication between goroutines. Unlike buffered channels, unbuffered channels have no capacity to store values and require both a sender and receiver to be ready at the same time for communication to occur.

## Basic Usage

### Creating an Unbuffered Channel
An unbuffered channel is created by omitting the buffer size parameter:

```go
// Create an unbuffered channel of type int
ch := make(chan int)
```

### Sending and Receiving
With unbuffered channels, the sending and receiving operations block until both sides are ready:

```go
func main() {
    ch := make(chan string)

    // Start a goroutine to send a message
    go func() {
        ch <- "Hello!" // This blocks until someone receives
    }()

    // Receive the message in the main goroutine
    msg := <-ch // This blocks until someone sends
    fmt.Println(msg)
}
```

## Key Characteristics

1. **Synchronization**: Unbuffered channels provide guaranteed synchronization between sender and receiver.
2. **Blocking Nature**: Send operations block until there's a receiver, and receive operations block until there's a sender.
3. **Zero Capacity**: No messages can be stored in the channel.

## Common Use Cases

### 1. [Worker Pools](../../examples/worker_pools/worker_pools.go)
Unbuffered channels are excellent for coordinating worker pools where you want precise control over task distribution:

```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for job := range jobs {
        results <- job * 2 // Process job and send result
    }
}

func main() {
    jobs := make(chan int)
    results := make(chan int)

    // Start 3 workers
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    // Send jobs
    for j := 1; j <= 5; j++ {
        jobs <- j
    }
    close(jobs)

    // Collect results
    for i := 1; i <= 5; i++ {
        <-results
    }
}
```

### 2. Signaling
Unbuffered channels are perfect for signaling between goroutines:

```go
func process(done chan bool) {
    // Do some work...
    time.Sleep(time.Second)
    done <- true // Signal completion
}

func main() {
    done := make(chan bool)
    go process(done)
    <-done // Wait for process to complete
}
```

## Best Practices

1. **Always Close Channels from the Sender Side**
```go
ch := make(chan int)
go func() {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch) // Sender closes the channel
}()

// Receiver can use range to receive until channel is closed
for num := range ch {
    fmt.Println(num)
}
```

2. **Use Directional Channel Parameters**
Specify channel direction in function parameters for better type safety:
```go
// Receive-only channel
func receive(ch <-chan int) {
    val := <-ch
}

// Send-only channel
func send(ch chan<- int) {
    ch <- 42
}
```

## Common Pitfalls

### 1. Deadlocks
Unbuffered channels can cause deadlocks if not used carefully:

```go
func main() {
    ch := make(chan int)
    ch <- 1  // Deadlock! No receiver available
    // This will panic
}
```

### 2. Goroutine Leaks
Always ensure there's a way to exit goroutines to prevent leaks:

```go
func main() {
    done := make(chan bool)
    ch := make(chan int)

    go func() {
        select {
        case <-done:
            return // Exit when done signal received
        case ch <- 42:
            // Send value
        }
    }()

    // Cleanup
    close(done)
}
```

## Select Statement with Unbuffered Channels
The `select` statement is often used with unbuffered channels for handling multiple channel operations:

```go
func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    go func() {
        ch1 <- "Hello"
    }()

    go func() {
        ch2 <- "World"
    }()

    select {
    case msg1 := <-ch1:
        fmt.Println("Received from ch1:", msg1)
    case msg2 := <-ch2:
        fmt.Println("Received from ch2:", msg2)
    case <-time.After(time.Second):
        fmt.Println("Timeout")
    }
}
```

## Conclusion

Unbuffered channels are powerful tools for synchronization and communication between goroutines in Go. They provide guaranteed delivery and synchronization points in your concurrent programs. Understanding their blocking nature and proper usage patterns is crucial for writing efficient and correct concurrent Go programs.

Remember:
- Use unbuffered channels when you need synchronization between goroutines
- Always handle channel closure appropriately
- Be aware of potential deadlocks
- Consider using select statements for more complex channel operations
- Use directional channels when possible to make code intentions clear



