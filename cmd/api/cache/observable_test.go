// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package cache

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestObservableCache_SetGet(t *testing.T) {
	t.Run("set and get key from cache", func(t *testing.T) {
		c := NewObservableCache(Config{MaxEntitiesCount: 2}, nil)
		c.Set("test", []byte{0, 1, 2, 3})

		got, ok := c.Get("test")
		require.True(t, ok)
		require.Equal(t, []byte{0, 1, 2, 3}, got)

		_, exists := c.Get("unknown")
		require.False(t, exists)
	})

	t.Run("overflow set queue", func(t *testing.T) {
		c := NewObservableCache(Config{MaxEntitiesCount: 2}, nil)
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
		c := NewObservableCache(Config{MaxEntitiesCount: 2}, nil)

		var wg sync.WaitGroup

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(c *ObservableCache, wg *sync.WaitGroup) {
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
}

func TestObservableCache_Clear(t *testing.T) {
	t.Run("clear cache", func(t *testing.T) {
		c := NewObservableCache(Config{MaxEntitiesCount: 100}, nil)
		for i := 0; i < 100; i++ {
			c.Set(fmt.Sprintf("%d", i), []byte{byte(i)})
		}
		c.Clear()

		require.Len(t, c.queue, 100)
		for i := 0; i < 100; i++ {
			require.EqualValues(t, c.queue[i], "")
		}
		require.Len(t, c.m, 0)
	})
}
