# Go Best Practices Applied in This Project

This document describes all the Go best practices implemented in this project, with explanations and examples.

## Table of Contents

1. [Project Structure](#project-structure)
2. [Error Handling](#error-handling)
3. [Concurrency Patterns](#concurrency-patterns)
4. [Testing](#testing)
5. [Code Organization](#code-organization)
6. [Performance Optimization](#performance-optimization)
7. [Documentation](#documentation)
8. [Code Quality Tools](#code-quality-tools)

---

## 1. Project Structure

### Standard Go Layout

```
go-demo/
├── cmd/              # Main applications
│   └── demo/        # Entry point for demo application
├── internal/        # Private application code (cannot be imported by other projects)
│   └── linkedlist/  # Internal data structure implementation
├── pkg/             # Public libraries (can be imported by other projects)
│   └── concurrency/ # Reusable concurrency patterns
├── examples/        # Example programs
├── docs/            # Documentation
├── go.mod           # Module definition (at root)
└── go.sum          # Dependency checksums
```

**Why this structure?**

- `cmd/` - Contains entry points, keeps main packages minimal
- `internal/` - Go compiler prevents external imports, ensures encapsulation
- `pkg/` - Explicitly shows what's meant to be reusable
- One `go.mod` at root - proper module management

**Example:**
```go
// ✅ Good: minimal main function
func main() {
    if err := run(); err != nil {
        log.Fatal(err)
    }
}

// Business logic separated
func run() error {
    // actual implementation
}
```

---

## 2. Error Handling

### Custom Error Types

**File:** `internal/linkedlist/linkedlist.go`

```go
var (
    ErrEmptyList = errors.New("list is empty")
    ErrIndexOutOfRange = errors.New("index out of range")
)
```

**Why?**
- Exported errors can be checked with `errors.Is()`
- Provides clear, typed error conditions
- Consumers can handle specific errors differently

### Always Check Errors

```go
// ✅ Good: error handling
func (ll *LinkedList) InsertAt(index, value int) error {
    if index < 0 || index > ll.size {
        return ErrIndexOutOfRange
    }
    // ... implementation
    return nil
}

// Usage
if err := list.InsertAt(5, 100); err != nil {
    if errors.Is(err, linkedlist.ErrIndexOutOfRange) {
        // Handle specific error
    }
    return fmt.Errorf("failed to insert: %w", err)
}
```

### Error Wrapping

```go
// ✅ Good: wrap errors with context
if err != nil {
    return fmt.Errorf("failed to process data: %w", err)
}
```

**Benefits:**
- Preserves original error with `%w`
- Adds contextual information
- Allows unwrapping with `errors.Unwrap()`

---

## 3. Concurrency Patterns

### Context for Cancellation and Timeouts

**File:** `pkg/concurrency/patterns.go`

```go
// ✅ Good: context-aware worker
func (wp *WorkerPool) runWorker(ctx context.Context, id int, worker Worker) {
    defer wp.wg.Done()
    
    for {
        select {
        case <-ctx.Done():
            return  // Graceful shutdown
        case job, ok := <-wp.jobs:
            if !ok {
                return  // Channel closed
            }
            err := worker(id, job)
            wp.results <- err
        }
    }
}
```

**Why?**
- Allows cancellation of long-running operations
- Prevents goroutine leaks
- Standard way to propagate deadlines

### Always Close Channels

```go
// ✅ Good: proper channel lifecycle
func (wp *WorkerPool) Close() {
    close(wp.jobs)      // Signal workers to stop
    wp.wg.Wait()        // Wait for all workers
    close(wp.results)   // Close results after workers done
}
```

### Use WaitGroup for Goroutine Synchronization

```go
// ✅ Good: structured concurrency
type WorkerPool struct {
    wg sync.WaitGroup  // Tracks active goroutines
}

func (wp *WorkerPool) Start(ctx context.Context, worker Worker) {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.runWorker(ctx, i, worker)
    }
}

func (wp *WorkerPool) Close() {
    close(wp.jobs)
    wp.wg.Wait()  // Ensures all goroutines finished
}
```

### Avoid Global State

```go
// ❌ Bad: global mutable state
var wg sync.WaitGroup  // Global variable

// ✅ Good: encapsulated state
type WorkerPool struct {
    wg sync.WaitGroup  // Part of struct
}
```

### Select for Non-Blocking Operations

```go
// ✅ Good: non-blocking send with context
select {
case output <- result:
    // Sent successfully
case <-ctx.Done():
    return  // Context cancelled
}
```

---

## 4. Testing

### Table-Driven Tests

**File:** `internal/linkedlist/linkedlist_test.go`

```go
func TestLinkedList_Append(t *testing.T) {
    tests := []struct {
        name   string
        values []int
        want   []int
    }{
        {
            name:   "append to empty list",
            values: []int{1},
            want:   []int{1},
        },
        {
            name:   "append multiple values",
            values: []int{1, 2, 3, 4, 5},
            want:   []int{1, 2, 3, 4, 5},
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ll := New()
            for _, v := range tt.values {
                ll.Append(v)
            }
            
            got := ll.ToSlice()
            if !slicesEqual(got, tt.want) {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

**Why?**
- Easy to add new test cases
- Clear test structure
- Subtests for better reporting
- Reduces code duplication

### Test Helper Functions

```go
// ✅ Good: helper with t.Helper()
func createList(values []int) *LinkedList {
    ll := New()
    for _, v := range values {
        ll.Append(v)
    }
    return ll
}
```

### Context Cancellation Tests

```go
func TestWorkerPool_ContextCancellation(t *testing.T) {
    ctx, cancel := context.WithCancel(context.Background())
    
    wp := NewWorkerPool(2)
    wp.Start(ctx, worker)
    
    cancel()  // Test cancellation behavior
    
    time.Sleep(100 * time.Millisecond)
    // Verify workers stopped
}
```

### Race Condition Detection

```bash
# Run tests with race detector
go test -race ./...
```

---

## 5. Code Organization

### One Responsibility Per Function

```go
// ✅ Good: focused functions
func (ll *LinkedList) Append(value int) {
    newNode := &Node{Value: value}
    
    if ll.Head == nil {
        ll.Head = newNode
        ll.Tail = newNode
    } else {
        ll.Tail.Next = newNode
        ll.Tail = newNode
    }
    ll.size++
}
```

### Early Returns

```go
// ✅ Good: early return reduces nesting
func (ll *LinkedList) Delete(value int) error {
    if ll.Head == nil {
        return ErrEmptyList
    }
    
    if ll.Head.Value == value {
        ll.Head = ll.Head.Next
        ll.size--
        return nil
    }
    
    // Continue with main logic
}
```

### Package-Level Constructors

```go
// ✅ Good: New() constructor
func New() *LinkedList {
    return &LinkedList{
        Head: nil,
        Tail: nil,
        size: 0,
    }
}

// Usage
ll := linkedlist.New()
```

### Exported vs Unexported

```go
// ✅ Good naming
type LinkedList struct {  // Exported type
    Head *Node           // Exported field
    Tail *Node           // Exported field
    size int             // Unexported (internal state)
}

// Exported method for size access
func (ll *LinkedList) Size() int {
    return ll.size
}
```

---

## 6. Performance Optimization

### Pointer Receivers for Mutating Methods

```go
// ✅ Good: pointer receiver for mutation
func (ll *LinkedList) Append(value int) {
    // Modifies ll in-place
}

// ✅ Good: pointer receiver for large structs
func (wp *WorkerPool) Start(ctx context.Context, worker Worker) {
    // WorkerPool has channels and WaitGroup - expensive to copy
}
```

### Preallocate Slices When Size Known

```go
// ✅ Good: preallocate with capacity
func (ll *LinkedList) ToSlice() []int {
    if ll.size == 0 {
        return []int{}
    }
    
    result := make([]int, 0, ll.size)  // Preallocate capacity
    current := ll.Head
    
    for current != nil {
        result = append(result, current.Value)
        current = current.Next
    }
    
    return result
}
```

**Why?**
- Avoids multiple allocations
- Reduces memory copying
- Better performance for known sizes

### Tail Pointer Optimization

```go
type LinkedList struct {
    Head *Node
    Tail *Node  // ✅ Maintains tail pointer
    size int
}

// O(1) append instead of O(n)
func (ll *LinkedList) Append(value int) {
    newNode := &Node{Value: value}
    
    if ll.Tail == nil {
        ll.Head = newNode
        ll.Tail = newNode
    } else {
        ll.Tail.Next = newNode
        ll.Tail = newNode  // Update tail - O(1)
    }
    ll.size++
}
```

### Size Tracking

```go
type LinkedList struct {
    Head *Node
    Tail *Node
    size int  // ✅ Track size instead of counting
}

// O(1) size check instead of O(n)
func (ll *LinkedList) Size() int {
    return ll.size
}
```

---

## 7. Documentation

### Godoc Format

```go
// LinkedList represents a singly linked list data structure.
type LinkedList struct {
    Head *Node
    Tail *Node
    size int
}

// Append adds a new node with the given value to the end of the list.
// Time complexity: O(1)
func (ll *LinkedList) Append(value int) {
    // implementation
}
```

**Rules:**
- Start with the name being documented
- Be concise but informative
- Include complexity when relevant
- Document all exported types, functions, and constants

### Package Comments

```go
// Package linkedlist provides a singly linked list implementation
// with efficient operations and proper error handling.
package linkedlist
```

### Example Code in Documentation

View examples in `examples/` directory:
- `linkedlist_example.go` - Shows real usage
- `concurrency_example.go` - Demonstrates patterns

---

## 8. Code Quality Tools

### golangci-lint Configuration

**File:** `.golangci.yml`

Enabled linters:
- `errcheck` - Checks that errors are checked
- `gosimple` - Suggests code simplifications
- `govet` - Reports suspicious constructs
- `staticcheck` - Advanced static analysis
- `gofmt` - Checks formatting
- `goimports` - Checks imports
- `goconst` - Finds repeated strings (candidates for constants)
- `gocritic` - Comprehensive checks
- `misspell` - Spelling errors
- `revive` - Fast, configurable linter

**Usage:**
```bash
# Install
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run
golangci-lint run

# Or via Makefile
make lint
```

### Makefile for Automation

**File:** `Makefile`

```makefile
test:
    go test -v ./...

test-coverage:
    go test -cover ./...
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

test-race:
    go test -race ./...

lint:
    golangci-lint run

fmt:
    gofmt -w .
    goimports -w .
```

### Git Ignore Best Practices

**File:** `.gitignore`

```
# Binaries
*.exe
*.dll
*.so

# Test artifacts
*.test
*.out

# Build
bin/
dist/

# IDE
.idea/
.vscode/

# Project specific
.cursorrules  # Don't commit project-specific AI rules
```

---

## Summary: Key Takeaways

### ✅ DO

1. **Structure**
   - Use standard project layout (cmd/, internal/, pkg/)
   - One go.mod at root
   - Keep main packages minimal

2. **Error Handling**
   - Always check errors
   - Create custom error types
   - Wrap errors with context

3. **Concurrency**
   - Use context for cancellation
   - Always close channels (defer close())
   - No goroutine leaks (use WaitGroup)
   - Avoid global mutable state

4. **Testing**
   - Write table-driven tests
   - Test error cases
   - Use race detector
   - Meaningful test names

5. **Performance**
   - Pointer receivers for mutations
   - Preallocate slices
   - Track metadata (size, tail) for O(1) operations

6. **Documentation**
   - Godoc format for all exports
   - Include examples
   - Document complexity

7. **Tools**
   - Use golangci-lint
   - Automate with Makefile
   - Run tests with -race flag

### ❌ DON'T

1. Don't use panic for normal errors
2. Don't ignore context cancellation
3. Don't use global mutable state
4. Don't share memory by communicating
5. Don't skip error checks
6. Don't use magic numbers/strings

---

## References

- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)

---

## Examples in This Project

### Best Practice → Implementation

| Best Practice | File | Line/Function |
|--------------|------|---------------|
| Custom Errors | `internal/linkedlist/linkedlist.go` | Lines 8-13 |
| Table-Driven Tests | `internal/linkedlist/linkedlist_test.go` | All test functions |
| Context Usage | `pkg/concurrency/patterns.go` | `WorkerPool.runWorker()` |
| Worker Pool | `pkg/concurrency/patterns.go` | `WorkerPool` type |
| Pipeline Pattern | `pkg/concurrency/patterns.go` | `Pipeline.Execute()` |
| Pointer Receiver | `internal/linkedlist/linkedlist.go` | All mutating methods |
| Preallocated Slice | `internal/linkedlist/linkedlist.go` | `ToSlice()` method |
| Godoc Comments | All `.go` files | All exported types |
| Minimal Main | `cmd/demo/main.go` | `main()` function |

---

*This document is maintained alongside the codebase. When applying new best practices, update this documentation.*

