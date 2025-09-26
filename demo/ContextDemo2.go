package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var limitKey = "limitkey"

func WithLimit(parent context.Context, limit int) (context.Context, context.CancelFunc) {
	// Use a buffered channel with the limit as its capacity
	sem := make(chan struct{}, limit)

	// Create a new context with a cancellation function
	ctx, cancel := context.WithCancel(parent)

	// Replace the value of the key with the semaphore
	return context.WithValue(ctx, limitKey, sem), cancel
}

func performRequest(ctx context.Context, requestID int) {
	// Access semaphore from context
	sem, _ := ctx.Value(limitKey).(chan struct{})

	// Acquire the semaphore
	sem <- struct{}{}
	defer func() { <-sem }() // Release the semaphore after request completes

	// Simulating some request processing
	// doing
	time.Sleep(1 * time.Second)
	fmt.Printf("Request %d completed.\n", requestID)
}

func main() {
	requests := []int{1, 2, 3, 4, 5}
	// Creating a context with a limit on concurrent requests
	concurrentCtx, cancel := WithLimit(context.Background(), 2)
	defer cancel()

	var wg sync.WaitGroup

	for _, requestID := range requests {
		wg.Add(1)
		go func(id int) {
			// Check if the context is canceled before performing the request
			select {
			case <-concurrentCtx.Done():
				fmt.Printf("Request %d canceled due to context cancellation.\n", id)
			default:
				performRequest(concurrentCtx, id)
			}
			wg.Done()
		}(requestID)
	}

	wg.Wait()
}
