package main

import (
	"fmt"
	"time"
)

func main() {

	// wg := sync.WaitGroup{}
	// wg.Add(1)
	ch := make(chan int)
	go func() {
		// defer wg.Done()
		<-ch // Will block forever as no one sends to this channel
		fmt.Println("Done")

	}()

	// wg.Wait()
	fmt.Println("Exiting")
	close(ch)
	time.Sleep(5 * time.Second)
	// Channel is never closed or written to
}


func createGoroutine() {
	go func() {
		fmt.Println("Hello, World!")
	}()
}




