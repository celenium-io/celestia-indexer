package pool

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	value int
	data  string
}

func TestPool_Basic(t *testing.T) {
	// Test with simple types
	intPool := New(func() int { return 0 })

	// Get from pool (should create new)
	val := intPool.Get()
	assert.Equal(t, 0, val)

	// Put modified value back
	intPool.Put(42)

	// Get from pool (should reuse)
	val = intPool.Get()
	assert.Equal(t, 42, val)
}

func TestPool_Struct(t *testing.T) {
	// Test with struct type
	structPool := New(func() *testStruct {
		return &testStruct{value: 0, data: ""}
	})

	// Get from pool
	obj := structPool.Get()
	assert.NotNil(t, obj)
	assert.Equal(t, 0, obj.value)
	assert.Equal(t, "", obj.data)

	// Modify and put back
	obj.value = 100
	obj.data = "test"
	structPool.Put(obj)

	// Get again (should be same object with modified values)
	obj2 := structPool.Get()
	assert.Equal(t, 100, obj2.value)
	assert.Equal(t, "test", obj2.data)
}

func TestPool_Pointer(t *testing.T) {
	// Test with pointer types
	type counter struct {
		count int
	}

	pool := New(func() *counter {
		return &counter{count: 0}
	})

	c1 := pool.Get()
	c1.count = 10
	pool.Put(c1)

	c2 := pool.Get()
	assert.Equal(t, 10, c2.count)
}

func TestPool_Slice(t *testing.T) {
	// Test with slice type
	slicePool := New(func() []byte {
		return make([]byte, 0, 1024)
	})

	slice := slicePool.Get()
	assert.NotNil(t, slice)
	assert.Equal(t, 0, len(slice))
	assert.GreaterOrEqual(t, cap(slice), 1024)

	// Use slice
	slice = append(slice, 1, 2, 3)
	assert.Equal(t, 3, len(slice))

	// Put back
	slicePool.Put(slice)

	// Get again
	slice2 := slicePool.Get()
	assert.Equal(t, 3, len(slice2))
}

func TestPool_Reset(t *testing.T) {
	// Test pool with reset functionality
	type buffer struct {
		data []byte
	}

	pool := New(func() *buffer {
		return &buffer{data: make([]byte, 0, 1024)}
	})

	// Get buffer
	buf := pool.Get()
	buf.data = append(buf.data, 1, 2, 3, 4, 5)

	// Reset before putting back
	buf.data = buf.data[:0]
	pool.Put(buf)

	// Get again, should have empty data but preserved capacity
	buf2 := pool.Get()
	assert.Equal(t, 0, len(buf2.data))
	assert.GreaterOrEqual(t, cap(buf2.data), 1024)
}

func TestPool_Concurrent(t *testing.T) {
	// Test concurrent access
	pool := New(func() *testStruct {
		return &testStruct{}
	})

	var wg sync.WaitGroup
	concurrency := 100

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Get from pool
			obj := pool.Get()
			assert.NotNil(t, obj)

			// Modify
			obj.value = id
			obj.data = "concurrent"

			// Put back
			pool.Put(obj)
		}(i)
	}

	wg.Wait()
}
