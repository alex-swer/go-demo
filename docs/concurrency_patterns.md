# Concurrency Patterns in Go

Detailed explanation of all concurrency patterns implemented in `pkg/concurrency/`.

## Table of Contents

1. [Worker Pool](#1-worker-pool)
2. [Pipeline](#2-pipeline)
3. [Fan-Out/Fan-In](#3-fan-outfan-in)
4. [Rate Limiter](#4-rate-limiter)
5. [Broadcast](#5-broadcast)
6. [Common Patterns](#common-patterns)

---

## 1. Worker Pool

**Pattern:** Distribute work among a fixed number of goroutines.

### When to Use

- Fixed-size pool of workers
- Controlled resource usage (e.g., max database connections)
- Process many tasks concurrently
- Graceful shutdown required

### Implementation

**File:** `pkg/concurrency/patterns.go`

```go
type WorkerPool struct {
    workers int
    jobs    chan interface{}
    results chan error
    wg      sync.WaitGroup
}

func NewWorkerPool(workers int) *WorkerPool {
    return &WorkerPool{
        workers: workers,
        jobs:    make(chan interface{}, workers*2),
        results: make(chan error, workers*2),
    }
}
```

### Key Features

✅ **Context Cancellation**
```go
select {
case <-ctx.Done():
    return  // Stop when context cancelled
case job, ok := <-wp.jobs:
    // Process job
}
```

✅ **Graceful Shutdown**
```go
func (wp *WorkerPool) Close() {
    close(wp.jobs)      // Signal no more jobs
    wp.wg.Wait()        // Wait for workers to finish
    close(wp.results)   // Safe to close results
}
```

✅ **Buffered Channels**
- Jobs channel buffered to prevent blocking
- Results channel buffered for async collection

### Usage Example

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

wp := concurrency.NewWorkerPool(3)

worker := func(id int, data interface{}) error {
    // Process data
    fmt.Printf("Worker %d processing %v\n", id, data)
    return nil
}

wp.Start(ctx, worker)

for i := 0; i < 100; i++ {
    wp.Submit(i)
}

wp.Close()

// Collect results
for err := range wp.Results() {
    if err != nil {
        log.Printf("Error: %v", err)
    }
}
```

### Best Practices

1. ✅ Always close jobs channel when done
2. ✅ Wait for all workers before closing results
3. ✅ Use context for cancellation
4. ✅ Buffer channels appropriately
5. ✅ Handle errors from workers

---

## 2. Pipeline

**Pattern:** Multi-stage data processing where output of one stage feeds the next.

### When to Use

- Sequential transformations
- Each stage can run concurrently
- Clear separation of concerns
- Streaming data processing

### Implementation

```go
type Pipeline struct {
    stages []func(context.Context, <-chan interface{}) <-chan interface{}
}

func NewPipeline(stages ...func(context.Context, <-chan interface{}) <-chan interface{}) *Pipeline {
    return &Pipeline{stages: stages}
}

func (p *Pipeline) Execute(ctx context.Context, input <-chan interface{}) <-chan interface{} {
    out := input
    for _, stage := range p.stages {
        out = stage(ctx, out)
    }
    return out
}
```

### Pattern: Generator → Processor → Consumer

```go
// Stage 1: Generator
func generator(ctx context.Context, nums ...int) <-chan interface{} {
    out := make(chan interface{})
    go func() {
        defer close(out)
        for _, n := range nums {
            select {
            case <-ctx.Done():
                return
            case out <- n:
            }
        }
    }()
    return out
}

// Stage 2: Processor (multiply by 2)
func multiply(ctx context.Context, input <-chan interface{}) <-chan interface{} {
    out := make(chan interface{})
    go func() {
        defer close(out)
        for val := range input {
            select {
            case <-ctx.Done():
                return
            case out <- val.(int) * 2:
            }
        }
    }()
    return out
}

// Stage 3: Processor (add 10)
func add(ctx context.Context, input <-chan interface{}) <-chan interface{} {
    out := make(chan interface{})
    go func() {
        defer close(out)
        for val := range input {
            select {
            case <-ctx.Done():
                return
            case out <- val.(int) + 10:
            }
        }
    }()
    return out
}
```

### Usage Example

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

pipeline := concurrency.NewPipeline(multiply, add)

input := generator(ctx, 1, 2, 3, 4, 5)
output := pipeline.Execute(ctx, input)

// Input: [1, 2, 3, 4, 5]
// After multiply: [2, 4, 6, 8, 10]
// After add: [12, 14, 16, 18, 20]

for result := range output {
    fmt.Println(result)
}
```

### Key Features

✅ **Composable Stages** - Each stage is independent  
✅ **Concurrent Processing** - All stages run in parallel  
✅ **Context Propagation** - Cancellation flows through all stages  
✅ **Backpressure** - Channels naturally provide flow control

### Best Practices

1. ✅ Always defer close(out) in each stage
2. ✅ Check context in both send and receive
3. ✅ Each stage is a pure function
4. ✅ Stages should not share state

---

## 3. Fan-Out/Fan-In

**Pattern:** Distribute work to multiple workers (fan-out), then combine results (fan-in).

### When to Use

- Independent tasks that can run in parallel
- CPU-bound operations
- Need to scale horizontally
- Aggregate results from multiple sources

### Implementation

#### Fan-Out: One Input → Many Workers

```go
func FanOut(ctx context.Context, input <-chan interface{}, workers int, 
            fn func(interface{}) interface{}) []<-chan interface{} {
    
    outputs := make([]<-chan interface{}, workers)
    
    for i := 0; i < workers; i++ {
        outputs[i] = worker(ctx, input, fn)
    }
    
    return outputs
}

func worker(ctx context.Context, input <-chan interface{}, 
            fn func(interface{}) interface{}) <-chan interface{} {
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
```

#### Fan-In: Many Workers → One Output

```go
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
```

### Usage Example

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

input := make(chan interface{})
go func() {
    defer close(input)
    for i := 1; i <= 100; i++ {
        input <- i
    }
}()

// Fan-out to 5 workers, each squares the number
square := func(val interface{}) interface{} {
    n := val.(int)
    return n * n
}

outputs := concurrency.FanOut(ctx, input, 5, square)

// Fan-in results from all workers
result := concurrency.FanIn(ctx, outputs...)

// Collect all squared numbers
for val := range result {
    fmt.Println(val)
}
```

### Visual Representation

```
         ┌─────────┐
Input ──→│ Worker 1│──┐
         └─────────┘  │
         ┌─────────┐  │
    ──→  │ Worker 2│──┤
         └─────────┘  │  ┌────────┐
         ┌─────────┐  ├─→│ Output │
    ──→  │ Worker 3│──┤  └────────┘
         └─────────┘  │
         ┌─────────┐  │
    ──→  │ Worker 4│──┘
         └─────────┘
```

### Best Practices

1. ✅ Workers should not share state
2. ✅ Use context for cancellation
3. ✅ Close output channel after all workers finish
4. ✅ Handle slow workers with buffered channels

---

## 4. Rate Limiter

**Pattern:** Limit the rate of operations using token bucket algorithm.

### When to Use

- API rate limiting
- Prevent overwhelming external services
- Throttle user actions
- Control resource consumption

### Implementation

```go
type RateLimiter struct {
    tokens chan struct{}
    rate   time.Duration
    done   chan struct{}
}

func NewRateLimiter(requestsPerSecond int) *RateLimiter {
    rl := &RateLimiter{
        tokens: make(chan struct{}, requestsPerSecond),
        rate:   time.Second / time.Duration(requestsPerSecond),
        done:   make(chan struct{}),
    }
    
    // Initial tokens
    for i := 0; i < requestsPerSecond; i++ {
        rl.tokens <- struct{}{}
    }
    
    go rl.refill()
    return rl
}

func (rl *RateLimiter) refill() {
    ticker := time.NewTicker(rl.rate)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            select {
            case rl.tokens <- struct{}{}:
            default:  // Bucket full, drop token
            }
        case <-rl.done:
            return
        }
    }
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-rl.tokens:
        return nil
    }
}
```

### Token Bucket Algorithm

```
Bucket Capacity: N tokens
Refill Rate: R tokens per second

[●●●●○] ← Bucket (4 tokens available)
    ↓
Request takes 1 token
    ↓
[●●●○○] ← Bucket (3 tokens left)
    ↑
  Refill adds tokens at rate R
```

### Usage Example

```go
ctx := context.Background()
rl := concurrency.NewRateLimiter(10)  // 10 requests/second
defer rl.Stop()

for i := 0; i < 100; i++ {
    if err := rl.Wait(ctx); err != nil {
        log.Fatal(err)
    }
    
    // Make API call
    makeAPIRequest(i)
}
```

### Real-World Example: API Client

```go
type APIClient struct {
    limiter *concurrency.RateLimiter
    client  *http.Client
}

func NewAPIClient() *APIClient {
    return &APIClient{
        limiter: concurrency.NewRateLimiter(100), // 100 req/sec
        client:  &http.Client{},
    }
}

func (c *APIClient) Get(ctx context.Context, url string) (*http.Response, error) {
    if err := c.limiter.Wait(ctx); err != nil {
        return nil, err
    }
    
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    return c.client.Do(req)
}
```

### Best Practices

1. ✅ Use context for timeout/cancellation
2. ✅ Stop the refill goroutine when done
3. ✅ Choose appropriate buffer size (burst capacity)
4. ✅ Handle context errors properly

---

## 5. Broadcast

**Pattern:** Send messages to multiple subscribers (pub/sub).

### When to Use

- Event notification system
- Live updates to multiple clients
- Message broadcasting
- Observer pattern

### Implementation

```go
type Broadcast struct {
    mu          sync.RWMutex
    subscribers map[string]chan interface{}
}

func NewBroadcast() *Broadcast {
    return &Broadcast{
        subscribers: make(map[string]chan interface{}),
    }
}

func (b *Broadcast) Subscribe(id string, bufferSize int) <-chan interface{} {
    b.mu.Lock()
    defer b.mu.Unlock()
    
    ch := make(chan interface{}, bufferSize)
    b.subscribers[id] = ch
    return ch
}

func (b *Broadcast) Unsubscribe(id string) {
    b.mu.Lock()
    defer b.mu.Unlock()
    
    if ch, ok := b.subscribers[id]; ok {
        close(ch)
        delete(b.subscribers, id)
    }
}

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
```

### Usage Example

```go
ctx := context.Background()
b := concurrency.NewBroadcast()
defer b.Close()

// Subscriber 1
sub1 := b.Subscribe("user1", 10)
go func() {
    for msg := range sub1 {
        fmt.Printf("User1 received: %v\n", msg)
    }
}()

// Subscriber 2
sub2 := b.Subscribe("user2", 10)
go func() {
    for msg := range sub2 {
        fmt.Printf("User2 received: %v\n", msg)
    }
}()

// Broadcast messages
for i := 1; i <= 5; i++ {
    msg := fmt.Sprintf("Message %d", i)
    if err := b.Send(ctx, msg); err != nil {
        log.Printf("Error: %v", err)
    }
}

time.Sleep(time.Second)
```

### Thread Safety

```go
// ✅ RWMutex for concurrent access
type Broadcast struct {
    mu          sync.RWMutex  // Protects subscribers map
    subscribers map[string]chan interface{}
}

// Read lock for sending (multiple senders OK)
func (b *Broadcast) Send(ctx context.Context, msg interface{}) error {
    b.mu.RLock()
    defer b.mu.RUnlock()
    // ...
}

// Write lock for modifying subscribers
func (b *Broadcast) Subscribe(id string, bufferSize int) <-chan interface{} {
    b.mu.Lock()
    defer b.mu.Unlock()
    // ...
}
```

### Best Practices

1. ✅ Use buffered channels to prevent blocking
2. ✅ Handle slow subscribers (full channels)
3. ✅ Use RWMutex for read-heavy workloads
4. ✅ Clean up subscribers properly

---

## Common Patterns

### 1. Defer Close Pattern

```go
func process(ctx context.Context, input <-chan int) <-chan int {
    output := make(chan int)
    go func() {
        defer close(output)  // ✅ Always close on exit
        for val := range input {
            // process
            output <- val * 2
        }
    }()
    return output
}
```

### 2. Select with Context

```go
select {
case <-ctx.Done():
    return ctx.Err()
case result := <-ch:
    return result
case <-time.After(timeout):
    return errors.New("timeout")
}
```

### 3. Done Channel Pattern

```go
type Service struct {
    done chan struct{}
}

func (s *Service) Stop() {
    close(s.done)
}

func (s *Service) Run() {
    for {
        select {
        case <-s.done:
            return
        default:
            // do work
        }
    }
}
```

### 4. Timeout Pattern

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result := make(chan Result)
go doWork(result)

select {
case r := <-result:
    return r, nil
case <-ctx.Done():
    return nil, ctx.Err()
}
```

---

## Testing Concurrency

### Test Context Cancellation

```go
func TestWorkerPool_ContextCancellation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    wp := NewWorkerPool(2)
    
    wp.Start(ctx, worker)
    
    // Submit jobs
    for i := 0; i < 10; i++ {
        wp.Submit(i)
    }
    
    cancel()  // Cancel context
    
    time.Sleep(100 * time.Millisecond)
    // Verify workers stopped
}
```

### Test Race Conditions

```bash
go test -race ./pkg/concurrency
```

### Test Goroutine Leaks

```go
func TestNoGoroutineLeaks(t *testing.T) {
    before := runtime.NumGoroutine()
    
    // Run concurrent code
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()
    
    wp := NewWorkerPool(5)
    wp.Start(ctx, worker)
    wp.Close()
    
    time.Sleep(100 * time.Millisecond)
    
    after := runtime.NumGoroutine()
    if after > before+1 {  // +1 for test goroutine
        t.Errorf("Goroutine leak: before=%d, after=%d", before, after)
    }
}
```

---

## Performance Tips

### 1. Buffer Channels Appropriately

```go
// ❌ Bad: unbuffered, blocks on every send
jobs := make(chan Job)

// ✅ Good: buffered, reduces contention
jobs := make(chan Job, 100)
```

### 2. Use sync.Pool for Frequent Allocations

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

buf := bufferPool.Get().(*bytes.Buffer)
defer bufferPool.Put(buf)
```

### 3. Right-Size Worker Pools

```go
// CPU-bound: workers = number of CPUs
workers := runtime.NumCPU()

// I/O-bound: can have more workers
workers := runtime.NumCPU() * 2
```

---

## References

- [Go Concurrency Patterns (Google I/O 2012)](https://www.youtube.com/watch?v=f6kdp27TYZs)
- [Go Blog: Pipelines](https://go.dev/blog/pipelines)
- [Go Blog: Context](https://go.dev/blog/context)
- [Effective Go: Concurrency](https://golang.org/doc/effective_go#concurrency)

---

*All patterns implemented in `pkg/concurrency/` with comprehensive tests.*

