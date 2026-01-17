# Go Channel Directions

In Go, channels can be specified with directional constraints to make the code's intent clearer and provide better type safety. There are three types of channel declarations:

## 1. Bidirectional Channels
```go
chan T     // Can be used for both sending and receiving
```

## 2. Send-Only Channels
```go
chan<- T   // Can only send values into the channel
```
The arrow points in the direction of data flow. Think of it as "sending into the channel".

Example:
```go
func send(ch chan<- int) {
    ch <- 42    // Valid: can send
    // x := <-ch // Invalid: cannot receive from send-only channel
}
```

## 3. Receive-Only Channels
```go
<-chan T   // Can only receive values from the channel
```
The arrow points away from chan, indicating data flowing out of the channel.

Example:
```go
func receive(ch <-chan int) {
    x := <-ch   // Valid: can receive
    // ch <- 42  // Invalid: cannot send to receive-only channel
}
```

## Common Usage Pattern

```go
func main() {
    ch := make(chan int) // Create a bidirectional channel
    
    // Pass as send-only to goroutine
    go producer(ch)
    
    // Pass as receive-only to consumer
    consumer(ch)
}

func producer(out chan<- int) {
    out <- 42
}

func consumer(in <-chan int) {
    value := <-in
    fmt.Println(value)
}
```

## Key Points

1. A bidirectional channel can be passed to a function expecting a send-only or receive-only channel
2. A send-only channel cannot be converted to a receive-only channel (and vice versa)
3. Using directional channels makes code intent clearer and prevents accidental misuse
4. The compiler will enforce these directional constraints

## Best Practices

- Use the most restrictive channel type that satisfies your needs
- Functions should accept send-only or receive-only channels when they only need to send or receive
- Keep bidirectional channels mainly in the code that creates and manages the channels 