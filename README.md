# Go Demo Project

A learning and demonstration project focusing on Go fundamentals, data structures, and concurrency patterns.

## 🎯 Project Goals

- Learn and practice Go best practices
- Implement common data structures
- Explore concurrency patterns
- Write comprehensive tests

## 📁 Project Structure

```
go-demo/
├── cmd/
│   └── demo/           # Main application entry point
├── internal/
│   └── linkedlist/     # Linked list implementation
├── pkg/
│   └── concurrency/    # Concurrency patterns
├── examples/           # Standalone examples
├── docs/               # Documentation
├── .golangci.yml       # Linter configuration
├── go.mod              # Module definition
└── README.md
```

## 🚀 Quick Start

### Prerequisites

- Go 1.24.1 or higher
- golangci-lint (optional, for linting)

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd go-demo

# Install dependencies
go mod download

# Run the main demo
go run cmd/demo/main.go
```

### Running Examples

```bash
# Linked list examples
go run examples/linkedlist_example.go

# Concurrency examples
go run examples/concurrency_example.go
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...

# Verbose output
go test -v ./...
```

## 📚 Features

### Data Structures

#### Linked List (`internal/linkedlist`)
- ✅ Append, Prepend operations (O(1))
- ✅ Insert at index (O(n))
- ✅ Delete by value/index (O(n))
- ✅ Find, Get operations
- ✅ Reverse in-place
- ✅ Comprehensive error handling
- ✅ Full test coverage

### Concurrency Patterns (`pkg/concurrency`)

#### Worker Pool
Manages a pool of goroutines for concurrent task processing with context support.

```go
wp := concurrency.NewWorkerPool(3)
wp.Start(ctx, workerFunc)
wp.Submit(job)
wp.Close()
```

#### Pipeline
Multi-stage data processing with context cancellation.

```go
pipeline := concurrency.NewPipeline(stage1, stage2, stage3)
output := pipeline.Execute(ctx, input)
```

#### Fan-Out/Fan-In
Distribute work across multiple workers and combine results.

```go
outputs := concurrency.FanOut(ctx, input, 5, processFn)
result := concurrency.FanIn(ctx, outputs...)
```

#### Rate Limiter
Token bucket algorithm for rate limiting operations.

```go
rl := concurrency.NewRateLimiter(10) // 10 requests/second
rl.Wait(ctx)
```

#### Broadcast
Send messages to multiple subscribers.

```go
b := concurrency.NewBroadcast()
sub := b.Subscribe("id", 10)
b.Send(ctx, message)
```

## 🧪 Testing

All packages include comprehensive table-driven tests following Go best practices:

- Unit tests for all public APIs
- Error case coverage
- Context cancellation tests
- Race condition detection
- Edge case handling

Run specific package tests:

```bash
go test ./internal/linkedlist -v
go test ./pkg/concurrency -v
```

## 🔧 Development

### Linting

Install golangci-lint:

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

Run linter:

```bash
golangci-lint run
```

### Code Formatting

```bash
# Format all files
gofmt -w .

# With imports
goimports -w .
```

### Building

```bash
# Build main application
go build -o bin/demo cmd/demo/main.go

# Build with optimizations
go build -ldflags="-s -w" -o bin/demo cmd/demo/main.go
```

## 📖 Documentation

- [Best Practices](docs/best_practices.md) - **Comprehensive guide to all applied Go best practices**
- [Concurrency Patterns](docs/concurrency_patterns.md) - **Detailed explanation of all concurrency patterns**
- [Go Pointers](docs/go_pointers.md) - Understanding pointers in Go
- [Go Types](docs/go_types.md) - Type system overview

## 🎓 Learning Resources

This project demonstrates:

- ✅ Proper Go project structure
- ✅ Package organization (internal/ vs pkg/)
- ✅ Error handling patterns
- ✅ Table-driven tests
- ✅ Context usage
- ✅ Channel patterns
- ✅ Goroutine management
- ✅ Memory efficiency (struct alignment)
- ✅ Code documentation (godoc format)

## 🚦 Best Practices Applied

1. **Error Handling**: All errors are properly checked and wrapped
2. **Concurrency**: Context for cancellation, proper channel closing
3. **Testing**: Comprehensive test coverage with table-driven tests
4. **Documentation**: Godoc-style comments for all exports
5. **Code Quality**: Linted with golangci-lint
6. **Memory Safety**: No goroutine leaks, proper resource cleanup

## 📝 License

This is a learning project - feel free to use and modify as needed.

## 🤝 Contributing

This is a personal learning project, but suggestions and improvements are welcome!

## 📧 Contact

For questions or suggestions, please open an issue.

