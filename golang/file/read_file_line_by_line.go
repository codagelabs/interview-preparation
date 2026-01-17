package main

import (
	"bufio"
	"fmt"
	"log"
	"os"	
)

func main() {
	filePath := "file/example.txt"
	channel := make(chan string)

	go readFileLineByLine(channel, filePath)
	processLines(channel)
}


func readFileLineByLine(channel chan string, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		channel <- scanner.Text()
	}
	close(channel)
}


func processLines(channel chan string) {
	for line := range channel {
		fmt.Println(line)
	}
}


