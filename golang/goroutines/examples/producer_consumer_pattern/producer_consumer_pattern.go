package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Message represents a unit of work
type Message struct {
	ID        int
	Data      string
	Timestamp time.Time
}

// Buffer represents a thread-safe message queue with separate producer and consumer management
type Buffer struct {
	messages         chan Message
	producerWg       sync.WaitGroup
	consumerWg       sync.WaitGroup
	producersRunning int
	mutex            sync.Mutex
}

// NewBuffer creates a new buffer with specified capacity
func NewBuffer(capacity int) *Buffer {
	return &Buffer{
		messages: make(chan Message, capacity),
	}
}

// StartProducer adds and starts a new producer
func (b *Buffer) StartProducer(id int, producerMessageCount int) {
	b.producerWg.Add(1)
	b.mutex.Lock()
	b.producersRunning++
	b.mutex.Unlock()

	go func() {
		defer b.producerWg.Done()
		defer func() {
			b.mutex.Lock()
			b.producersRunning--
			if b.producersRunning == 0 {
				close(b.messages)
			}
			b.mutex.Unlock()
		}()

		for i := 0; i < producerMessageCount; i++ {
			msg := Message{
				ID:        i,
				Data:      fmt.Sprintf("Message from Producer-%d", id),
				Timestamp: time.Now(),
			}
			b.messages <- msg
			fmt.Printf("Producer-%d produced message %d at %v\n",
				id, i, msg.Timestamp.Format(time.StampMilli))

			// Simulate fast production (100-300ms)
			time.Sleep(time.Duration(rand.Intn(200)+100) * time.Millisecond)
		}
	}()
}

// StartConsumer adds and starts a new consumer
func (b *Buffer) StartConsumer(id int) {
	b.consumerWg.Add(1)
	go func() {
		defer b.consumerWg.Done()
		for msg := range b.messages {
			// Simulate slower processing (300-700ms)
			processingTime := time.Duration(rand.Intn(400)+300) * time.Millisecond
			time.Sleep(processingTime)

			latency := time.Since(msg.Timestamp)
			fmt.Printf("Consumer-%d processed message %d after %v latency\n",
				id, msg.ID, latency.Round(time.Millisecond))
		}
	}()
}

// WaitForProducers waits for all producers to finish
func (b *Buffer) WaitForProducers() {
	b.producerWg.Wait()
}

// WaitForConsumers waits for all consumers to finish
func (b *Buffer) WaitForConsumers() {
	b.consumerWg.Wait()
}

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Create a buffer with capacity of 50 messages
	buffer := NewBuffer(50)

	// Start multiple producers
	numProducers := 3
	messagesPerProducer := 5
	for i := 0; i < numProducers; i++ {
		buffer.StartProducer(i, messagesPerProducer)
	}

	// Start multiple consumers
	numConsumers := 2
	for i := 0; i < numConsumers; i++ {
		buffer.StartConsumer(i)
	}

	// Wait for all producers to finish
	buffer.WaitForProducers()

	// Wait for all consumers to finish processing messages
	// This ensures all messages have been properly handled before program exit
	// Without this, the program might terminate while consumers are still processing
	buffer.WaitForConsumers()

	fmt.Println("All processing completed!")
}
