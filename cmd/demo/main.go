package main

import (
	"context"
	"fmt"
	"go-demo/internal/linkedlist"
	"go-demo/pkg/concurrency"
	"time"
)

func main() {
	fmt.Println("Go Demo Project - Data Structures & Concurrency")
	fmt.Println("================================================")
	fmt.Println()
	
	runLinkedListDemo()
	fmt.Println()
	
	runConcurrencyDemo()
	
	fmt.Println()
	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}

func runLinkedListDemo() {
	fmt.Println("📋 Linked List Demo")
	fmt.Println("-------------------")
	
	ll := linkedlist.New()
	
	for i := 1; i <= 5; i++ {
		ll.Append(i * 10)
	}
	
	fmt.Printf("List: %v\n", ll.ToSlice())
	fmt.Printf("Size: %d\n", ll.Size())
	
	ll.Prepend(5)
	fmt.Printf("After prepend(5): %v\n", ll.ToSlice())
	
	err := ll.InsertAt(3, 25)
	if err == nil {
		fmt.Printf("After insert at index 3: %v\n", ll.ToSlice())
	}
	
	node, found := ll.Find(30)
	if found {
		fmt.Printf("Found value 30: %v\n", node.Value)
	}
	
	ll.Reverse()
	fmt.Printf("After reverse: %v\n", ll.ToSlice())
}

func runConcurrencyDemo() {
	fmt.Println("⚡ Concurrency Demo")
	fmt.Println("------------------")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	fmt.Println("\n1. Worker Pool Example:")
	wp := concurrency.NewWorkerPool(3)
	
	worker := func(id int, data interface{}) error {
		job := data.(int)
		fmt.Printf("   Worker %d: Processing job %d\n", id, job)
		time.Sleep(200 * time.Millisecond)
		return nil
	}
	
	wp.Start(ctx, worker)
	
	for i := 1; i <= 6; i++ {
		wp.Submit(i)
	}
	
	wp.Close()
	fmt.Println("   ✓ All jobs completed")
	
	fmt.Println("\n2. Pipeline Example:")
	demonstratePipeline(ctx)
	
	fmt.Println("\n3. Rate Limiter Example:")
	demonstrateRateLimiter()
}

func demonstratePipeline(ctx context.Context) {
	stage1 := func(ctx context.Context, input <-chan interface{}) <-chan interface{} {
		output := make(chan interface{})
		go func() {
			defer close(output)
			for val := range input {
				select {
				case <-ctx.Done():
					return
				case output <- val.(int) * 2:
				}
			}
		}()
		return output
	}
	
	stage2 := func(ctx context.Context, input <-chan interface{}) <-chan interface{} {
		output := make(chan interface{})
		go func() {
			defer close(output)
			for val := range input {
				select {
				case <-ctx.Done():
					return
				case output <- val.(int) + 10:
				}
			}
		}()
		return output
	}
	
	pipeline := concurrency.NewPipeline(stage1, stage2)
	
	input := make(chan interface{})
	go func() {
		defer close(input)
		for i := 1; i <= 3; i++ {
			input <- i
		}
	}()
	
	output := pipeline.Execute(ctx, input)
	
	fmt.Print("   Input [1,2,3] → *2 → +10 = ")
	results := []int{}
	for result := range output {
		results = append(results, result.(int))
	}
	fmt.Printf("%v\n", results)
}

func demonstrateRateLimiter() {
	ctx := context.Background()
	rl := concurrency.NewRateLimiter(2)
	defer rl.Stop()
	
	fmt.Println("   Rate: 2 requests/second")
	start := time.Now()
	
	for i := 1; i <= 4; i++ {
		err := rl.Wait(ctx)
		if err != nil {
			fmt.Printf("   Error: %v\n", err)
			return
		}
		elapsed := time.Since(start).Round(100 * time.Millisecond)
		fmt.Printf("   Request %d executed at %v\n", i, elapsed)
	}
}

