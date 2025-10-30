package main

import (
	"context"
	"fmt"
	"go-demo/pkg/concurrency"
	"time"
)

func concurrencyExample() {
	fmt.Println("=== Concurrency Examples ===")
	fmt.Println()
	
	workerPoolExample()
	fmt.Println()
	
	pipelineExample()
	fmt.Println()
	
	fanOutFanInExample()
	fmt.Println()
	
	rateLimiterExample()
	fmt.Println()
	
	broadcastExample()
}

func workerPoolExample() {
	fmt.Println("1. Worker Pool:")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	wp := concurrency.NewWorkerPool(3)
	
	worker := func(id int, data interface{}) error {
		job := data.(int)
		fmt.Printf("   Worker %d processing job %d\n", id, job)
		time.Sleep(100 * time.Millisecond)
		return nil
	}
	
	wp.Start(ctx, worker)
	
	for i := 1; i <= 5; i++ {
		wp.Submit(i)
	}
	
	wp.Close()
	fmt.Println("   All jobs completed")
}

func pipelineExample() {
	fmt.Println("2. Pipeline Pattern:")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
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
	
	fmt.Print("   Results: ")
	for result := range output {
		fmt.Printf("%d ", result)
	}
	fmt.Println()
}

func fanOutFanInExample() {
	fmt.Println("3. Fan-Out/Fan-In Pattern:")
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	input := make(chan interface{})
	go func() {
		defer close(input)
		for i := 1; i <= 5; i++ {
			input <- i
		}
	}()
	
	square := func(val interface{}) interface{} {
		n := val.(int)
		return n * n
	}
	
	outputs := concurrency.FanOut(ctx, input, 3, square)
	result := concurrency.FanIn(ctx, outputs...)
	
	fmt.Print("   Squared values: ")
	for val := range result {
		fmt.Printf("%d ", val)
	}
	fmt.Println()
}

func rateLimiterExample() {
	fmt.Println("4. Rate Limiter:")
	
	ctx := context.Background()
	rl := concurrency.NewRateLimiter(2)
	defer rl.Stop()
	
	start := time.Now()
	
	for i := 1; i <= 5; i++ {
		err := rl.Wait(ctx)
		if err != nil {
			fmt.Printf("   Error: %v\n", err)
			return
		}
		fmt.Printf("   Request %d at %v\n", i, time.Since(start).Round(100*time.Millisecond))
	}
}

func broadcastExample() {
	fmt.Println("5. Broadcast Pattern:")
	
	ctx := context.Background()
	b := concurrency.NewBroadcast()
	defer b.Close()
	
	sub1 := b.Subscribe("subscriber1", 10)
	sub2 := b.Subscribe("subscriber2", 10)
	
	go func() {
		for i := 1; i <= 3; i++ {
			msg := fmt.Sprintf("Message %d", i)
			err := b.Send(ctx, msg)
			if err != nil {
				fmt.Printf("   Error sending: %v\n", err)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
	
	go func() {
		for msg := range sub1 {
			fmt.Printf("   Subscriber 1 received: %v\n", msg)
		}
	}()
	
	for msg := range sub2 {
		fmt.Printf("   Subscriber 2 received: %v\n", msg)
	}
	
	time.Sleep(500 * time.Millisecond)
}

