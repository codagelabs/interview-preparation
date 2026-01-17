package main

import (
	"fmt"
	"time"
)

// ConsumerProducer is an interface that defines methods for producing and consuming tasks
type ConsumerProducer interface {
	Produce(ch chan<- Task) // Method to produce tasks and send them to a channel
	Consume(ch <-chan Task) // Method to consume tasks received from a channel
}

// Task is a struct that represents a task to be processed
type Task struct {
	ID   int    // Unique identifier for the task
	Data string // Data associated with the task
}

// producer function sends tasks to the provided channel
func producer(ch chan<- Task) {
	ch <- Task{ID: 1, Data: "Task 1"} // Send Task 1 to the channel
	ch <- Task{ID: 2, Data: "Task 2"} // Send Task 2 to the channel
	ch <- Task{ID: 3, Data: "Task 3"} // Send Task 3 to the channel
}

// consumer function receives tasks from the provided channel and processes them
func consumer(ch <-chan Task) {
	for task := range ch { // Continuously receive tasks from the channel
		fmt.Println(task) // Print the task
	}
}

func main() {
	ch := make(chan Task) // Create a channel for Task communication

	go producer(ch) // Start the producer goroutine
	go consumer(ch) // Start the consumer goroutine

	time.Sleep(time.Second * 1) // Allow some time for goroutines to complete
}

