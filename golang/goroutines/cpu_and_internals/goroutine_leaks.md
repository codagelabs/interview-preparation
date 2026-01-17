# Avoiding Goroutine Leaks in a Long-Running Go Application

In this example, we'll build a simple server that handles client requests. We'll demonstrate how to use contexts to manage goroutine lifecycles and ensure that they are properly terminated to avoid leaks

## What is Gorutine Leaks
Goroutine leaks occur when goroutines are not properly terminated after they have completed their work. This can happen if a goroutine is blocked indefinitely, waiting for a resource or a signal that never arrives. Over time, these leaked goroutines can accumulate, consuming memory and other system resources, which can lead to performance degradation or even application crashes.

Common causes of goroutine leaks include:

1. **Unreleased Resources**: Goroutines waiting on channels or other synchronization primitives that are never closed or signaled.
2. **Infinite Loops**: Goroutines stuck in infinite loops without proper exit conditions.
3. **Uncanceled Contexts**: Goroutines that are not properly canceled when their context is no longer needed.

To avoid goroutine leaks, it is important to ensure that all goroutines have a clear exit strategy and that resources are properly released when they are no longer needed. Using `context.Context` is a common pattern to manage the lifecycle of goroutines and ensure they can be canceled when appropriate.

## Common Causes of Goroutine Leaks

Let's look at some common causes of goroutine leaks and how to avoid them:

1. **Unreleased Resources**: Goroutines waiting on channels or other synchronization primitives that are never closed or signaled.

   ```go
   // Leak: This goroutine will be blocked forever
   func leakyGoroutine() {
       ch := make(chan int)
       go func() {
           val := <-ch // Will block forever as no one sends to this channel
           fmt.Println(val)
       }()
       // Channel is never closed or written to
   }
   ```
   Corrected version 

   ```go
   // Corrected: This goroutine will not be blocked forever
   func fixedGoroutine() {
       ch := make(chan int)
       go func() {
           val := <-ch // Will receive a value and not block
           fmt.Println(val)
       }()
       ch <- 42 // Send a value to the channel
       close(ch) // Optionally close the channel if no more values will be sent
   }
   ```

2. **Infinite Loops**: Goroutines stuck in infinite loops without proper exit conditions.
   ```go
   // Leak: This goroutine will run forever
   func leakyLoop() {
       go func() {
           for {
               // No exit condition
               time.Sleep(1 * time.Second)
           }
       }()
   }
   ```
   Corrected version

   ```go
   // Corrected: This goroutine will exit after a condition is met
   func fixedLoop() {
       done := make(chan bool)
       go func() {
           for {
               select {
               case <-done:
                   return // Exit the loop when done is signaled
               default:
                   time.Sleep(1 * time.Second)
               }
           }
       }()
       // Simulate some work
       time.Sleep(5 * time.Second)
       done <- true // Signal the goroutine to exit
   }
   ```
   Explanation 
    - Exit Condition: A done channel is used to signal when the goroutine should exit. This prevents the goroutine from running indefinitely.
    - Select Statement: The select statement listens for a signal on the done channel. When a signal is received, the goroutine exits the loop and terminates.
    - Simulated Work: The main function simulates some work by sleeping for 5 seconds before signaling the goroutine to exit.
    This approach ensures that the goroutine has a clear exit strategy, preventing it from leaking resources.


3. **Uncanceled Contexts**: Goroutines that are not properly canceled when their context is no longer needed.
   ```go
   // Leak: Context is never canceled
   func leakyContext() {
       ctx := context.Background() // No cancellation mechanism
       go func(ctx context.Context) {
           <-ctx.Done() // Will block forever
       }(ctx)
   }
   ```
   Corrected version:

    ```go
    // Corrected: Context is properly canceled
    func fixedContext() {
        ctx, cancel := context.WithCancel(context.Background()) // Create a cancellable context
        defer cancel() // Ensure the context is canceled when done

        go func(ctx context.Context) {
            select {
            case <-ctx.Done():
                fmt.Println("Context canceled")
                return
            }
        }(ctx)

        // Simulate some work
        time.Sleep(2 * time.Second)
        cancel() // Cancel the context to signal the goroutine to exit
    }
    ```
    Explanation:
    - Cancellable Context: A cancellable context is created using `context.WithCancel`, which returns a cancel function.
    - Deferred Cancellation: The cancel function is deferred to ensure the context is canceled when the function completes, preventing resource leaks.
    - Goroutine Exit: The goroutine listens for the ctx.Done() signal, which indicates that the context has been canceled. When this happens, the goroutine exits gracefully.
    - Simulated Work: The main function simulates some work by sleeping for 2 seconds before canceling the context, signaling the goroutine to exit.

4. **Blocked Channel Operations**: Goroutines blocked on send operations on unbuffered channels with no receivers.
   ```go
   // Leak: Sending to a channel with no receivers
   func leakyChannelSend() {
       ch := make(chan int) // Unbuffered channel
       go func() {
           ch <- 42 // Will block forever if no one receives
       }()
       // No receiver for the channel
   }
   ```

   // Corrected: 
   ```go 
   // Corrected: Channel send operation with proper receiver
    func fixedChannelSend() {
        ch := make(chan int)    // Unbuffered channel
        done := make(chan bool) // Channel to signal completion

        // Sender goroutine
        go func() {
            select {
            case ch <- 42: // Try to send value
                fmt.Println("Value sent successfully")
            case <-time.After(time.Second): // Timeout if no receiver
                fmt.Println("Send operation timed out")
            }
            done <- true // Signal that we're done
        }()

        // Receiver goroutine
        go func() {
            val := <-ch
            fmt.Printf("Received value: %d\n", val)
        }()

        <-done // Wait for sender to complete
    }
    ```
 
    Explanation:
    - Receiver Goroutine: Added a receiver goroutine to consume the value sent on the channel, preventing the sender from blocking indefinitely.
    - Timeout Mechanism: Used a select statement with a timeout to prevent the sender from blocking forever if the receiver is not available.
    - Completion Signal: Added a done channel to properly coordinate the completion of the sender goroutine.
    Error Handling: Included basic error handling through timeout and logging to help diagnose issues.
    
5. **Deadlocks**: Multiple goroutines waiting for each other, creating a circular dependency.
   ```go
   // Leak: Deadlock between two goroutines
   func leakyDeadlock() {
       ch1 := make(chan int)
       ch2 := make(chan int)
       go func() {
           ch1 <- 1 // Blocks until someone receives
           <-ch2    // Waits for value from ch2
       }()
       go func() {
           ch2 <- 1 // Blocks until someone receives
           <-ch1    // Waits for value from ch1
       }()
       // Both goroutines are now deadlocked
   }
   ```

   Corrected: Deadlock prevention using buffered channels and proper synchronization

   ```go
    func fixedDeadlock() {
        ch1 := make(chan int, 1) // Buffered channel
        ch2 := make(chan int, 1) // Buffered channel
        done := make(chan bool)   // Channel to signal completion

        // First goroutine
        go func() {
            defer func() { done <- true }()
            // Send on ch1 won't block due to buffer
            ch1 <- 1
            // Receive from ch2
            <-ch2
            fmt.Println("First goroutine completed")
        }()

        // Second goroutine
        go func() {
            defer func() { done <- true }()
            // Send on ch2 won't block due to buffer
            ch2 <- 1
            // Receive from ch1
            <-ch1
            fmt.Println("Second goroutine completed")
        }()

        // Wait for both goroutines to complete
        <-done
        <-done
    }
    ```

    Explanation:
    - Receiver Goroutine: Added a receiver goroutine to consume the value sent on the channel, preventing the sender from blocking indefinitely.
    - Timeout Mechanism: Used a select statement with a timeout to prevent the sender from blocking forever if the receiver is not available.
    - Completion Signal: Added a done channel to properly coordinate the completion of the sender goroutine.
    - Error Handling: Included basic error handling through timeout and logging to help diagnose issues.

6. **Missing Done Signals**: Worker goroutines in pools that don't signal completion properly.
   ```go
   // Leak: WaitGroup never completed
   func leakyWaitGroup() {
       var wg sync.WaitGroup
       wg.Add(1)
       go func() {
           // Work is done but we forgot to call wg.Done()
           // wg.Done() is missing
       }()
       wg.Wait() // Will block forever
   }
   ```

   Corrected version:
   ```go
   // Corrected: WaitGroup with proper completion signaling
   func fixedWaitGroup() {
       var wg sync.WaitGroup
       wg.Add(1)
       
       go func() {
           defer wg.Done() // Ensures wg.Done() is called even if panic occurs
           
           // Simulate some work
           time.Sleep(time.Second)
           fmt.Println("Work completed")
       }()
       
       // Wait for goroutine to complete
       wg.Wait()
       fmt.Println("All work done")
   }
   ```

   Explanation:
   - **Proper Signaling**: The `wg.Done()` is called to signal completion of the goroutine's work
   - **Use of defer**: `defer wg.Done()` ensures the signal is sent even if the goroutine panics or returns early
   - **Structured Flow**: Clear pattern of initializing WaitGroup, adding count, launching work, and waiting
   - **Safety**: The WaitGroup won't block forever because completion is properly signaled
   
   This pattern is essential for:
   - Coordinating multiple goroutines
   - Ensuring all work completes before proceeding
   - Proper cleanup and resource management
   - Preventing goroutine leaks in worker pools

7. **Timer/Ticker Leaks**: Creating timers or tickers without properly stopping them.
   ```go
   // Leak: Ticker is never stopped
   func leakyTicker() {
       ticker := time.NewTicker(time.Second)
       go func() {
           for {
               <-ticker.C
               // Process tick
           }
       }()
       // ticker.Stop() is never called
   }
   ```
   Corrected version
    ```go
    // Corrected: Proper ticker management with cleanup
    func fixedTicker() {
        ticker := time.NewTicker(time.Second)
        defer ticker.Stop() // Ensure ticker is cleaned up

        done := make(chan bool)
        
        go func() {
            for {
                select {
                case <-done:
                    return
                case <-ticker.C:
                    fmt.Println("Tick processed at:", time.Now())
                    // Process tick
                }
            }
        }()

        // Simulate some work
        time.Sleep(5 * time.Second)
        done <- true // Signal goroutine to exit
    }
    ```

    Eplaination
    - Proper Cleanup:
    defer ticker.Stop() ensures the ticker is properly stopped when the function exits
    This prevents resource leaks and unnecessary CPU usage
    - Controlled Exit:
    Added a done channel to gracefully signal the goroutine to exit
    The select statement allows the goroutine to respond to either ticker events or the exit signal
    - Resource Management:
    The ticker is properly managed throughout its lifecycle
    No goroutines are left running after the function completes
    System resources are properly released
    - Best Practices:
    Using select for handling multiple channels
    Proper signal handling for goroutine termination
    Clear cleanup pattern with defer
