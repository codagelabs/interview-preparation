package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// DataChunk represents a batch of data to be processed
type DataChunk struct {
	ID    int
	Items []int
}

// ProcessingResult represents the result of processing a data chunk
type ProcessingResult struct {
	ChunkID int
	Sum     int
	Average float64
}

// BatchProcessor handles parallel processing of data chunks
type BatchProcessor struct {
	inputChan  chan DataChunk
	resultChan chan ProcessingResult
	workerWg   sync.WaitGroup
	resultWg   sync.WaitGroup
}

// NewBatchProcessor creates a new batch processor_unused with specified number of workers
func NewBatchProcessor(numWorkers int) *BatchProcessor {
	bp := &BatchProcessor{
		inputChan:  make(chan DataChunk, numWorkers),
		resultChan: make(chan ProcessingResult, numWorkers),
	}

	// Start workers
	for i := 0; i < numWorkers; i++ {
		bp.workerWg.Add(1)
		go bp.worker(i)
	}

	return bp
}

// worker processes data chunks and produces results
func (bp *BatchProcessor) worker(id int) {
	defer bp.workerWg.Done()

	for chunk := range bp.inputChan {
		// Simulate complex processing
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

		// Process the chunk
		sum := 0
		for _, item := range chunk.Items {
			sum += item
		}
		average := float64(sum) / float64(len(chunk.Items))

		// Send result
		result := ProcessingResult{
			ChunkID: chunk.ID,
			Sum:     sum,
			Average: average,
		}
		bp.resultChan <- result

		fmt.Printf("Worker %d processed chunk %d: Sum = %d, Average = %.2f\n",
			id, chunk.ID, sum, average)
	}
}

// ProcessBatches processes multiple data chunks and collects results
func (bp *BatchProcessor) ProcessBatches(chunks []DataChunk) []ProcessingResult {
	results := make([]ProcessingResult, 0, len(chunks))
	resultsMutex := sync.Mutex{}

	// Start result collector
	bp.resultWg.Add(1)
	go func() {
		defer bp.resultWg.Done()
		for result := range bp.resultChan {
			resultsMutex.Lock()
			results = append(results, result)
			resultsMutex.Unlock()
		}
	}()

	// Send chunks for processing
	for _, chunk := range chunks {
		bp.inputChan <- chunk
	}

	// Close input channel after all chunks are sent
	close(bp.inputChan)

	// Wait for all workers to finish
	bp.workerWg.Wait()

	// Close result channel and wait for result collection
	close(bp.resultChan)
	bp.resultWg.Wait()

	return results
}

// generateTestData creates sample data chunks for processing
func generateTestData(numChunks, chunkSize int) []DataChunk {
	chunks := make([]DataChunk, numChunks)
	for i := 0; i < numChunks; i++ {
		items := make([]int, chunkSize)
		for j := 0; j < chunkSize; j++ {
			items[j] = rand.Intn(100)
		}
		chunks[i] = DataChunk{
			ID:    i,
			Items: items,
		}
	}
	return chunks
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Create test data
	numChunks := 10
	chunkSize := 1000
	chunks := generateTestData(numChunks, chunkSize)

	// Create batch processor_unused with 4 workers
	processor := NewBatchProcessor(4)

	// Process batches and collect results
	fmt.Println("Starting batch processing...")
	startTime := time.Now()

	results := processor.ProcessBatches(chunks)

	// Calculate total statistics
	totalSum := 0
	totalAverage := 0.0
	for _, result := range results {
		totalSum += result.Sum
		totalAverage += result.Average
	}
	totalAverage /= float64(len(results))

	// Print final results
	fmt.Printf("\nProcessing completed in %v\n", time.Since(startTime))
	fmt.Printf("Total chunks processed: %d\n", len(results))
	fmt.Printf("Total sum: %d\n", totalSum)
	fmt.Printf("Overall average: %.2f\n", totalAverage)
}
