package main

import (
	"fmt"
)

func main() {
	ch := make(chan int, 2)

	// Non-blocking send
	for i := 0; i < 3; i++ {

		select {
		case ch <- i:
			fmt.Println("Sent value successfully")
		default:
			fmt.Println("Channel full - couldn't send")
		}
	}

	// Non-blocking receive
	for i := 0; i < 3; i++ {
	select {
	case value := <-ch:
		fmt.Println("Got value:", value)
	default:
			fmt.Println("No value available")
		}
	}

}
