package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Cargo manages a batch of tasks of any type and processes them in bulk.
type Cargo[T any] struct {
	batchSize   int           // Maximum number of tasks in a batch
	timeout     time.Duration // Time to wait before processing a batch
	tasks       []T           // Slice holding tasks
	mu          sync.Mutex    // Mutex to protect concurrent access to tasks
	processFunc func([]T)     // Function to process a batch of tasks
}

// NewCargo creates a new Cargo instance.
func NewCargo[T any](batchSize int, timeout time.Duration, processFunc func([]T)) *Cargo[T] {
	return &Cargo[T]{
		batchSize:   batchSize,
		timeout:     timeout,
		processFunc: processFunc,
	}
}

// AddTask adds a task to the cargo. If batch size is reached, the tasks are processed.
func (c *Cargo[T]) AddTask(task T) {
	c.mu.Lock()
	c.tasks = append(c.tasks, task)
	if len(c.tasks) >= c.batchSize {
		c.processBatch()
	}
	c.mu.Unlock()
}

// processBatch processes the current batch of tasks and clears the queue.
func (c *Cargo[T]) processBatch() {
	if len(c.tasks) > 0 {
		tasksToProcess := make([]T, len(c.tasks))
		copy(tasksToProcess, c.tasks)
		c.tasks = []T{}
		go c.processFunc(tasksToProcess)
	}
}

// start runs a ticker that triggers processing of tasks when the timeout is reached.
// It listens for a cancellation signal from the context for graceful shutdown.
func (c *Cargo[T]) start(ctx context.Context) {
	ticker := time.NewTicker(c.timeout)
	defer ticker.Stop() // Ensure ticker is stopped when exiting

	for {
		select {
		case <-ctx.Done():
			// Graceful shutdown: process any remaining tasks
			c.mu.Lock()
			c.processBatch()
			c.mu.Unlock()
			fmt.Println("Gracefully shutting down...")
			return
		case <-ticker.C:
			// Timeout reached, process the batch
			c.mu.Lock()
			c.processBatch()
			c.mu.Unlock()
		}
	}
}

// Example process function for a generic task
func processTasks[T any](tasks []T) {
	fmt.Printf("Processing batch of %d tasks: %v\n", len(tasks), tasks)
}

func main() {
	// Create a new context with cancel for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())

	// Create a new cargo with a batch size of 3 and a timeout of 5 seconds for integer tasks
	intCargo := NewCargo(3, 5*time.Second, processTasks[int])

	// Start the cargo processing in the background
	go intCargo.start(ctx)

	// Simulate adding tasks to the cargo
	for i := 1; i <= 10; i++ {
		intCargo.AddTask(i)
		time.Sleep(1 * time.Second) // Simulate delay between task arrivals
	}

	// Simulate a graceful shutdown after 12 seconds
	time.Sleep(12 * time.Second)
	cancel() // Signal the cargo to shutdown gracefully

	// Give some time to finish processing remaining tasks
	time.Sleep(2 * time.Second)
}

// https://chatgpt.com/c/66ff5d0c-9234-8000-8009-f248dc252954
