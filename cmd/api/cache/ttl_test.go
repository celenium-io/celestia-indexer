// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTTLCache_SetGet(t *testing.T) {
	t.Run("set and get key from cache", func(t *testing.T) {
		c := NewTTLCache(Config{MaxEntitiesCount: 2}, time.Second)
		c.Set("test", []byte{0, 1, 2, 3})

		got, ok := c.Get("test")
		require.True(t, ok)
		require.Equal(t, []byte{0, 1, 2, 3}, got)

		_, exists := c.Get("unknown")
		require.False(t, exists)
	})

	t.Run("overflow set queue", func(t *testing.T) {
		c := NewTTLCache(Config{MaxEntitiesCount: 2}, time.Second)
		for i := 0; i < 100; i++ {
			c.Set(fmt.Sprintf("%d", i), []byte{byte(i)})
		}

		require.Len(t, c.queue, 2)
		require.Len(t, c.m, 2)

		got, ok := c.Get("99")
		require.True(t, ok)
		require.Equal(t, []byte{99}, got)

		got1, ok1 := c.Get("98")
		require.True(t, ok1)
		require.Equal(t, []byte{98}, got1)

		_, exists := c.Get("0")
		require.False(t, exists)
	})

	t.Run("overflow set queue multithread", func(t *testing.T) {
		c := NewTTLCache(Config{MaxEntitiesCount: 2}, time.Second)

		var wg sync.WaitGroup

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(c *TTLCache, wg *sync.WaitGroup) {
				defer wg.Done()
				for i := 0; i < 100; i++ {
					c.Set(fmt.Sprintf("%d", i), []byte{byte(i)})
				}
			}(c, &wg)
		}

		wg.Wait()

		require.Len(t, c.queue, 2)
		require.Len(t, c.m, 2)
	})

	t.Run("get expired value", func(t *testing.T) {
		c := NewTTLCache(Config{MaxEntitiesCount: 4}, time.Second)
		c.Set("test", []byte{0})
		c.Set("test2", []byte{0})
		c.Set("test3", []byte{0})
		c.m["test"].expiredAt = time.Now().Add(-time.Hour)

		require.Len(t, c.m, 3)
		require.Len(t, c.queue, 4)

		_, exists := c.Get("test")
		require.False(t, exists)

		require.Len(t, c.queue, 4)
		require.Len(t, c.m, 2)
	})

	t.Run("get expired value 2", func(t *testing.T) {
		c := NewTTLCache(Config{MaxEntitiesCount: 4}, time.Second)
		c.Set("test2", []byte{0})
		c.Set("test3", []byte{0})
		c.Set("test", []byte{0})
		c.m["test"].expiredAt = time.Now().Add(-time.Hour)

		require.Len(t, c.m, 3)
		require.Len(t, c.queue, 4)

		_, exists := c.Get("test")
		require.False(t, exists)

		require.Len(t, c.queue, 4)
		require.Len(t, c.m, 2)
	})

	t.Run("get expired value 3", func(t *testing.T) {
		c := NewTTLCache(Config{MaxEntitiesCount: 4}, time.Second)
		c.Set("test2", []byte{0})
		c.Set("test", []byte{0})
		c.Set("test3", []byte{0})
		c.m["test"].expiredAt = time.Now().Add(-time.Hour)

		require.Len(t, c.m, 3)
		require.Len(t, c.queue, 4)

		_, exists := c.Get("test")
		require.False(t, exists)

		require.Len(t, c.queue, 4)
		require.Len(t, c.m, 2)
	})
}

func TestTTLCache_Clear(t *testing.T) {
	t.Run("clear cache", func(t *testing.T) {
		c := NewTTLCache(Config{MaxEntitiesCount: 2}, time.Second)
		for i := 0; i < 100; i++ {
			c.Set(fmt.Sprintf("%d", i), []byte{byte(i)})
		}
		c.Clear()

		require.Len(t, c.queue, 2)
		for i := 0; i < 2; i++ {
			require.EqualValues(t, c.queue[i], "")
		}
		require.Len(t, c.m, 0)
	})
}
