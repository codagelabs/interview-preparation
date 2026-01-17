package main

import (
	"fmt"
	"os"
)

func main() {

	file, err := os.OpenFile("example.json", os.O_APPEND, 0644)

	if err != nil {
		fmt.Println("Could not open example.txt")
		return
	}

	defer file.Close()

	_, err2 := file.Write([]byte("Appending some text to example.txt"))

	if err2 != nil {
		fmt.Println("Could not write text to example.txt")

	} else {
		fmt.Println("Operation successful! Text has been appended to example.txt")
	}
}
