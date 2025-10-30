package concurrency

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Worker represents a worker function that processes data.
type Worker func(id int, data interface{}) error

// WorkerPool manages a pool of goroutines for concurrent task processing.
type WorkerPool struct {
	workers int
	jobs    chan interface{}
	results chan error
	wg      sync.WaitGroup
}

// NewWorkerPool creates a new worker pool with the specified number of workers.
func NewWorkerPool(workers int) *WorkerPool {
	return &WorkerPool{
		workers: workers,
		jobs:    make(chan interface{}, workers*2),
		results: make(chan error, workers*2),
	}
}

// Start begins processing jobs with the given worker function.
// The context can be used to cancel all workers.
func (wp *WorkerPool) Start(ctx context.Context, worker Worker) {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.runWorker(ctx, i, worker)
	}
}

// runWorker processes jobs from the jobs channel until context is cancelled or channel is closed.
func (wp *WorkerPool) runWorker(ctx context.Context, id int, worker Worker) {
	defer wp.wg.Done()
	
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-wp.jobs:
			if !ok {
				return
			}
			err := worker(id, job)
			wp.results <- err
		}
	}
}

// Submit adds a new job to the worker pool.
func (wp *WorkerPool) Submit(job interface{}) {
	wp.jobs <- job
}

// Close closes the jobs channel and waits for all workers to finish.
func (wp *WorkerPool) Close() {
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
}

// Results returns the results channel.
func (wp *WorkerPool) Results() <-chan error {
	return wp.results
}

// Pipeline demonstrates a pipeline pattern with multiple stages.
type Pipeline struct {
	stages []func(context.Context, <-chan interface{}) <-chan interface{}
}

// NewPipeline creates a new pipeline.
func NewPipeline(stages ...func(context.Context, <-chan interface{}) <-chan interface{}) *Pipeline {
	return &Pipeline{stages: stages}
}

// Execute runs the pipeline with the given input channel.
func (p *Pipeline) Execute(ctx context.Context, input <-chan interface{}) <-chan interface{} {
	out := input
	for _, stage := range p.stages {
		out = stage(ctx, out)
	}
	return out
}

// FanOut distributes work from a single channel to multiple workers.
// Returns a slice of output channels, one per worker.
func FanOut(ctx context.Context, input <-chan interface{}, workers int, fn func(interface{}) interface{}) []<-chan interface{} {
	outputs := make([]<-chan interface{}, workers)
	
	for i := 0; i < workers; i++ {
		outputs[i] = worker(ctx, input, fn)
	}
	
	return outputs
}

// FanIn combines multiple input channels into a single output channel.
func FanIn(ctx context.Context, inputs ...<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	output := make(chan interface{})
	
	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-c:
				if !ok {
					return
				}
				select {
				case output <- val:
				case <-ctx.Done():
					return
				}
			}
		}
	}
	
	wg.Add(len(inputs))
	for _, c := range inputs {
		go multiplex(c)
	}
	
	go func() {
		wg.Wait()
		close(output)
	}()
	
	return output
}

// worker is a helper function that processes data from input channel.
func worker(ctx context.Context, input <-chan interface{}, fn func(interface{}) interface{}) <-chan interface{} {
	output := make(chan interface{})
	
	go func() {
		defer close(output)
		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-input:
				if !ok {
					return
				}
				result := fn(val)
				select {
				case output <- result:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	
	return output
}

// RateLimiter limits the rate of operations using a token bucket algorithm.
type RateLimiter struct {
	tokens chan struct{}
	rate   time.Duration
	done   chan struct{}
}

// NewRateLimiter creates a new rate limiter with the specified rate.
func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	rl := &RateLimiter{
		tokens: make(chan struct{}, requestsPerSecond),
		rate:   time.Second / time.Duration(requestsPerSecond),
		done:   make(chan struct{}),
	}
	
	for i := 0; i < requestsPerSecond; i++ {
		rl.tokens <- struct{}{}
	}
	
	go rl.refill()
	return rl
}

// refill adds tokens to the bucket at the specified rate.
func (rl *RateLimiter) refill() {
	ticker := time.NewTicker(rl.rate)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			select {
			case rl.tokens <- struct{}{}:
			default:
			}
		case <-rl.done:
			return
		}
	}
}

// Wait blocks until a token is available or context is cancelled.
func (rl *RateLimiter) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-rl.tokens:
		return nil
	}
}

// Stop stops the rate limiter.
func (rl *RateLimiter) Stop() {
	close(rl.done)
}

// Broadcast sends a message to multiple subscribers.
type Broadcast struct {
	mu          sync.RWMutex
	subscribers map[string]chan interface{}
}

// NewBroadcast creates a new broadcast instance.
func NewBroadcast() *Broadcast {
	return &Broadcast{
		subscribers: make(map[string]chan interface{}),
	}
}

// Subscribe adds a new subscriber with the given ID.
func (b *Broadcast) Subscribe(id string, bufferSize int) <-chan interface{} {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	ch := make(chan interface{}, bufferSize)
	b.subscribers[id] = ch
	return ch
}

// Unsubscribe removes a subscriber.
func (b *Broadcast) Unsubscribe(id string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if ch, ok := b.subscribers[id]; ok {
		close(ch)
		delete(b.subscribers, id)
	}
}

// Send broadcasts a message to all subscribers.
func (b *Broadcast) Send(ctx context.Context, msg interface{}) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	for _, ch := range b.subscribers {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- msg:
		default:
			return fmt.Errorf("subscriber channel full")
		}
	}
	
	return nil
}

// Close closes all subscriber channels.
func (b *Broadcast) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	for _, ch := range b.subscribers {
		close(ch)
	}
	b.subscribers = make(map[string]chan interface{})
}

