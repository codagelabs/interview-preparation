package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {
	//runtime.GOMAXPROCS(2) // Limiting to 1 CPU core
	// To test behaviour of GOMAXPROCS change the setting

	wg.Add(3)
	go doSomething(10)
	go doSomething(10)
	go doSomething(10)

	wg.Wait()
}

func doSomething(n int) {
	defer wg.Done()

	for i := 0; i < n; i++ {
		fmt.Print(i)
	}
	fmt.Println()
}

//with runtime.GOMAXPROCS(1) output will be
//0123456789
//0123456789
//0123456789
