package main

import (
	"fmt"
	"time"
)


func main() {
	// buffered channel
	// buffered channel is a channel that can hold a limited number of values without a corresponding receiver for those values.
	// buffered channel is created using the make function
	// the second argument to the make function is the buffer size	
	// buffered channel is useful when you want to send multiple values to a channel without having to wait for a receiver to be available.

	
	ch := make(chan string, 2)

	// send values to the channel
	ch <- "Hello"
	ch <- "World"

	fmt.Println(<-ch)
	fmt.Println(<-ch)



}
