package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// 工作线程
func worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Worker received cancellation signal.")
			return
		default:
			// Simulate some work
			time.Sleep(500 * time.Millisecond)
			fmt.Println("Working...")
		}
	}
}

func main() {
	parentCtx, cancel := context.WithCancel(context.Background())
	go worker(parentCtx)
	// Simulate main program execution
	time.Sleep(2 * time.Second)
	// Cancel the context to stop the worker
	cancel()

	// Wait for the worker to finish
	time.Sleep(1 * time.Second)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	operationWithTimeout(timeoutCtx)

	deadline := time.Now().Add(5 * time.Second)
	deadlineCtx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	operationWithDeadline(deadlineCtx)

	parentCtx = context.WithValue(context.Background(), "userID", 123)

	var wg sync.WaitGroup

	// Simulating multiple requests
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(requestID int) {
			// Creating a child context for each request
			childCtx := context.WithValue(parentCtx, "requestID", requestID)
			processRequest(childCtx, requestID)
			wg.Done()
		}(i)
	}

	wg.Wait()

}

// 超时时间
func operationWithTimeout(ctx context.Context) {
	select {
	case <-time.After(3 * time.Second): // Simulate some long operation
		fmt.Println("Operation completed.")
	case <-ctx.Done():
		fmt.Println("Operation canceled due to timeout.")
	}
}

// deadline 死亡时间
func operationWithDeadline(ctx context.Context) {
	deadline, ok := ctx.Deadline()
	if ok {
		fmt.Printf("Operation must be completed before: %s\n", deadline)
	} else {
		fmt.Println("No specific deadline for the operation.")
	}
	// Simulate some operation
	time.Sleep(10 * time.Second)
	select {
	case <-ctx.Done():
		fmt.Println("Operation canceled due to context deadline.")
	default:
		fmt.Println("Operation completed with in the deadline.")
	}

}

// 传递变量
func processRequest(ctx context.Context, requestID int) {
	// Accessing request-scoped value from the context
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		fmt.Println("Failed to get userID from context.")
		return
	}

	fmt.Printf("Processing request %d for user %d\n", requestID, userID)
}
