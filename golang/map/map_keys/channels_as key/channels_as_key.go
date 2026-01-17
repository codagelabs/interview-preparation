package main

import "fmt"

func main() {
	// Create channels
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	// Map with channels as keys
	m := make(map[chan int]string)
	m[ch1] = "Channel 1"
	m[ch2] = "Channel 2"
	m[ch3] = "Channel 3"

	// Accessing map values by channel keys
	fmt.Println("ch1:", m[ch1]) // Output: Channel 1
	fmt.Println("ch2:", m[ch2]) // Output: Channel 2
	fmt.Println("ch3:", m[ch3]) // Output: Channel 3

	// Checking if a channel exists in the map
	if _, ok := m[ch1]; ok {
		fmt.Println("ch1 exists in the map")
	}
}
