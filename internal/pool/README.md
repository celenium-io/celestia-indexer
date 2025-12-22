# Pool

Type-safe generic wrapper around `sync.Pool` for Go.

## Overview

The `pool` package provides a type-safe wrapper around Go's `sync.Pool` using generics. This eliminates the need for type assertions and provides compile-time type safety when working with object pools.

## Usage in the Project

Currently used in:
- **[pkg/node/rpc/block.go](../../pkg/node/rpc/block.go)** - Pools for request/response slices in `BlockBulkData` function to reduce allocations during bulk block data fetching

## Features

- **Type Safety**: No more type assertions or runtime panics
- **Zero Overhead**: Performance identical to `sync.Pool`
- **Simple API**: Easy to use with a familiar interface
- **Generic**: Works with any Go type

## Usage

### Basic Example

```go
import "github.com/celenium-io/celestia-indexer/internal/pool"

// Create a pool for byte slices
bufferPool := pool.New(func() []byte {
    return make([]byte, 0, 1024)
})

// Get a buffer from the pool
buffer := bufferPool.Get()
buffer = append(buffer, []byte("data")...)

// Reset and return to pool
buffer = buffer[:0]
bufferPool.Put(buffer)
```

### Struct Example

```go
type Request struct {
    ID   int
    Data string
}

// Create a pool for Request objects
requestPool := pool.New(func() *Request {
    return &Request{}
})

// Get and use
req := requestPool.Get()
req.ID = 1
req.Data = "example"

// Reset and return
req.ID = 0
req.Data = ""
requestPool.Put(req)
```

### Best Practices

1. **Always Reset Before Returning**: Clear object state before putting it back in the pool

```go
buffer = buffer[:0]  // Reset slice length
bufferPool.Put(buffer)
```

2. **Use Pointers for Structs**: Reduces allocations and copying

```go
pool.New(func() *MyStruct {  // Good
    return &MyStruct{}
})
```

3. **Pre-allocate Capacity**: Initialize with appropriate capacity

```go
pool.New(func() []byte {
    return make([]byte, 0, 1024)  // Pre-allocate capacity
})
```

## Performance

Benchmarks show that the generic wrapper has virtually identical performance to raw `sync.Pool`:

```
BenchmarkPool_GetPut-12          	0.9037 ns/op	0 B/op	0 allocs/op
BenchmarkSyncPool_GetPut-12      	0.8999 ns/op	0 B/op	0 allocs/op
```

## Comparison with sync.Pool

### Before (sync.Pool)

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 0, 1024)
    },
}

buffer := bufferPool.Get().([]byte)  // Type assertion required
bufferPool.Put(buffer)
```

### After (pool.Pool)

```go
var bufferPool = pool.New(func() []byte {
    return make([]byte, 0, 1024)
})

buffer := bufferPool.Get()  // Type-safe, no assertion needed
bufferPool.Put(buffer)
```

## API

### New

```go
func New[T any](factory func() T) *Pool[T]
```

Creates a new typed Pool with the given factory function.

### Get

```go
func (p *Pool[T]) Get() T
```

Retrieves an object from the pool. If the pool is empty, uses the factory function to create a new object.

### Put

```go
func (p *Pool[T]) Put(x T)
```

Adds an object back to the pool for reuse.
