package concurrency

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	wp := NewWorkerPool(3)
	
	var processed int32
	worker := func(id int, data interface{}) error {
		atomic.AddInt32(&processed, 1)
		time.Sleep(10 * time.Millisecond)
		return nil
	}
	
	wp.Start(ctx, worker)
	
	for i := 0; i < 10; i++ {
		wp.Submit(i)
	}
	
	wp.Close()
	
	if atomic.LoadInt32(&processed) != 10 {
		t.Errorf("expected 10 jobs processed, got %d", processed)
	}
}

func TestWorkerPool_WithError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	wp := NewWorkerPool(2)
	
	expectedErr := errors.New("worker error")
	worker := func(id int, data interface{}) error {
		if data.(int) == 5 {
			return expectedErr
		}
		return nil
	}
	
	wp.Start(ctx, worker)
	
	for i := 0; i < 10; i++ {
		wp.Submit(i)
	}
	
	wp.Close()
	
	errorCount := 0
	for err := range wp.Results() {
		if err != nil {
			errorCount++
		}
	}
	
	if errorCount != 1 {
		t.Errorf("expected 1 error, got %d", errorCount)
	}
}

func TestWorkerPool_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	
	wp := NewWorkerPool(2)
	
	var processed int32
	worker := func(id int, data interface{}) error {
		atomic.AddInt32(&processed, 1)
		time.Sleep(100 * time.Millisecond)
		return nil
	}
	
	wp.Start(ctx, worker)
	
	for i := 0; i < 10; i++ {
		wp.Submit(i)
	}
	
	time.Sleep(50 * time.Millisecond)
	cancel()
	
	time.Sleep(200 * time.Millisecond)
	
	finalProcessed := atomic.LoadInt32(&processed)
	if finalProcessed >= 10 {
		t.Error("expected workers to stop processing after context cancellation")
	}
}

func TestPipeline(t *testing.T) {
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
	
	pipeline := NewPipeline(stage1, stage2)
	
	input := make(chan interface{})
	go func() {
		defer close(input)
		for i := 1; i <= 5; i++ {
			input <- i
		}
	}()
	
	output := pipeline.Execute(ctx, input)
	
	expected := []int{12, 14, 16, 18, 20}
	i := 0
	for result := range output {
		if result.(int) != expected[i] {
			t.Errorf("expected %d, got %d", expected[i], result)
		}
		i++
	}
}

func TestFanOutFanIn(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	input := make(chan interface{})
	go func() {
		defer close(input)
		for i := 1; i <= 10; i++ {
			input <- i
		}
	}()
	
	square := func(val interface{}) interface{} {
		n := val.(int)
		return n * n
	}
	
	outputs := FanOut(ctx, input, 3, square)
	result := FanIn(ctx, outputs...)
	
	sum := 0
	for val := range result {
		sum += val.(int)
	}
	
	expected := 385
	if sum != expected {
		t.Errorf("expected sum %d, got %d", expected, sum)
	}
}

func TestRateLimiter(t *testing.T) {
	ctx := context.Background()
	rl := NewRateLimiter(5)
	defer rl.Stop()
	
	start := time.Now()
	
	for i := 0; i < 10; i++ {
		err := rl.Wait(ctx)
		if err != nil {
			t.Fatalf("Wait() error = %v", err)
		}
	}
	
	elapsed := time.Since(start)
	
	minDuration := 1 * time.Second
	if elapsed < minDuration {
		t.Errorf("rate limiter too fast: expected at least %v, got %v", minDuration, elapsed)
	}
}

func TestRateLimiter_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	rl := NewRateLimiter(1)
	defer rl.Stop()
	
	cancel()
	
	err := rl.Wait(ctx)
	if err == nil {
		t.Error("expected error when context is cancelled")
	}
}

func TestBroadcast(t *testing.T) {
	ctx := context.Background()
	b := NewBroadcast()
	defer b.Close()
	
	sub1 := b.Subscribe("sub1", 10)
	sub2 := b.Subscribe("sub2", 10)
	sub3 := b.Subscribe("sub3", 10)
	
	messages := []string{"msg1", "msg2", "msg3"}
	
	go func() {
		for _, msg := range messages {
			err := b.Send(ctx, msg)
			if err != nil {
				t.Errorf("Send() error = %v", err)
			}
		}
	}()
	
	received1 := collectMessages(sub1, len(messages))
	received2 := collectMessages(sub2, len(messages))
	received3 := collectMessages(sub3, len(messages))
	
	if len(received1) != len(messages) {
		t.Errorf("sub1: expected %d messages, got %d", len(messages), len(received1))
	}
	if len(received2) != len(messages) {
		t.Errorf("sub2: expected %d messages, got %d", len(messages), len(received2))
	}
	if len(received3) != len(messages) {
		t.Errorf("sub3: expected %d messages, got %d", len(messages), len(received3))
	}
}

func TestBroadcast_Unsubscribe(t *testing.T) {
	ctx := context.Background()
	b := NewBroadcast()
	defer b.Close()
	
	sub1 := b.Subscribe("sub1", 10)
	b.Subscribe("sub2", 10)
	
	b.Unsubscribe("sub1")
	
	err := b.Send(ctx, "test")
	if err != nil {
		t.Errorf("Send() error = %v", err)
	}
	
	select {
	case _, ok := <-sub1:
		if ok {
			t.Error("expected channel to be closed after unsubscribe")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("channel should be closed immediately")
	}
}

// Helper functions

func collectMessages(ch <-chan interface{}, count int) []interface{} {
	messages := make([]interface{}, 0, count)
	timeout := time.After(1 * time.Second)
	
	for i := 0; i < count; i++ {
		select {
		case msg := <-ch:
			messages = append(messages, msg)
		case <-timeout:
			return messages
		}
	}
	
	return messages
}

