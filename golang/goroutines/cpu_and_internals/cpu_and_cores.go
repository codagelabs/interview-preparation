package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println(runtime.NumCPU())
	fmt.Println(runtime.GOMAXPROCS(1))
	fmt.Println(runtime.GOMAXPROCS(0))
}
