package main

import (
	"fmt"
	"sync"
)

// demonstrateUnbufferedChannel shows behavior of closed unbuffered channels
func demonstrateUnbufferedChannel(wg *sync.WaitGroup) {
	fmt.Println("\n=== Unbuffered Channel Example ===")
	ch := make(chan int)

	// Close channel immediately to demonstrate reading from closed channel
	close(ch)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Try reading from closed channel
		val, ok := <-ch
		if !ok {
			fmt.Println("Unbuffered channel is closed")
		}
		fmt.Printf("Received from unbuffered: Value=%d, Channel Open=%v\n", val, ok)
	}()
}

// demonstrateBufferedChannel shows behavior of closed buffered channels
func demonstrateBufferedChannel(wg *sync.WaitGroup) {
	fmt.Println("\n=== Buffered Channel Example ===")
	ch := make(chan int, 2)

	// Demonstrate writing to buffer before closing
	ch <- 1
	ch <- 2
	fmt.Println("Wrote 1 and 2 to buffer")

	close(ch)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Read all values from buffer
		for i := 0; i < 3; i++ {
			val, ok := <-ch
			if !ok {
				fmt.Println("Buffered channel is closed")
			}
			fmt.Printf("Read #%d: Value=%d, Channel Open=%v\n", i+1, val, ok)
		}
	}()
}

func main() {
	var wg sync.WaitGroup

	// Demonstrate both channel types
	demonstrateUnbufferedChannel(&wg)
	wg.Wait()

	demonstrateBufferedChannel(&wg)
	wg.Wait()

	fmt.Println("\nDemonstration completed!")
}
