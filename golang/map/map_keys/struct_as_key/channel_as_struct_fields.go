package main

import "fmt"

// Define a struct with a channel field
type Worker struct {
    ID      int
    Name    string
    JobChan chan string
}

func main() {
    // Create some workers with channels
    worker1 := Worker{
        ID:      1,
        Name:    "Worker 1",
        JobChan: make(chan string),
    }

    worker2 := Worker{
        ID:      2, 
        Name:    "Worker 2",
        JobChan: make(chan string),
    }

    // Create a map using the Worker struct as key
    workerStatus := make(map[Worker]string)

    // Add workers to the map
    workerStatus[worker1] = "Available"
    workerStatus[worker2] = "Busy"

    // Print the map
    fmt.Println("Worker Status:")
    for worker, status := range workerStatus {
        fmt.Printf("Worker %d (%s): %s\n", worker.ID, worker.Name, status)
    }

    // Channels are comparable, so this works fine
    if _, exists := workerStatus[worker1]; exists {
        fmt.Println("\nWorker 1 exists in the map")
    }

    // Create a new worker with same ID and Name but different channel
    worker3 := Worker{
        ID:      1,
        Name:    "Worker 1", 
        JobChan: make(chan string), // New channel
    }

    // This will be treated as a different key because the channel is different
    workerStatus[worker3] = "On Break"

    fmt.Println("\nAfter adding worker with same ID/Name but different channel:")
    for worker, status := range workerStatus {
        fmt.Printf("Worker %d (%s): %s\n", worker.ID, worker.Name, status)
    }

    // Demonstrate that channels in the struct are still usable
    go func() {
        worker1.JobChan <- "Task 1"
    }()

    // Receive from the channel
    select {
    case task := <-worker1.JobChan:
        fmt.Printf("\nWorker 1 received task: %s\n", task)
    default:
        fmt.Println("\nNo task received")
    }
}


