// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package cache

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTTLCache_SetGet(t *testing.T) {
	c, err := NewTTLCache(time.Second)
	require.NoError(t, err)

	t.Run("set and get key from cache", func(t *testing.T) {
		c.Set("value1", []byte{0, 1, 2, 3})

		got, ok := c.Get("value1")
		require.True(t, ok)
		require.Equal(t, []byte{0, 1, 2, 3}, got)

		_, exists := c.Get("unknown")
		require.False(t, exists)
	})

	t.Run("many set and get", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			c.Set(fmt.Sprintf("%d", i), []byte{byte(i)})
		}

		got, ok := c.Get("99")
		require.True(t, ok)
		require.Equal(t, []byte{99}, got)

		got1, ok1 := c.Get("98")
		require.True(t, ok1)
		require.Equal(t, []byte{98}, got1)

		got2, exists := c.Get("0")
		require.True(t, exists)
		require.Equal(t, []byte{0}, got2)
	})

	t.Run("set multithread", func(t *testing.T) {
		var wg sync.WaitGroup

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(c *TTLCache, wg *sync.WaitGroup) {
				defer wg.Done()
				for j := 200; j < 300; j++ {
					c.Set(fmt.Sprintf("%d", j), []byte{byte(j)})
				}
			}(c, &wg)
		}

		wg.Wait()
	})

	t.Run("get expired value", func(t *testing.T) {
		c.Set("test", []byte{0})

		time.Sleep(time.Second)

		_, exists := c.Get("test")
		require.False(t, exists)
	})

	t.Run("multithread", func(t *testing.T) {
		var wg sync.WaitGroup
		set := func(wg *sync.WaitGroup) {
			wg.Done()

			for i := 0; i < 100; i++ {
				val, err := rand.Int(rand.Reader, big.NewInt(255))
				require.NoError(t, err)
				c.Set(val.String(), []byte{byte(i)})
			}
		}
		get := func(wg *sync.WaitGroup) {
			wg.Done()

			for i := 0; i < 100; i++ {
				c.Get(fmt.Sprintf("%d", i))
			}
		}

		for i := 0; i < 100; i++ {
			wg.Add(2)
			set(&wg)
			get(&wg)
		}

		wg.Wait()
	})
}

func TestTTLCache_Clear(t *testing.T) {
	c, err := NewTTLCache(time.Second)
	require.NoError(t, err)

	t.Run("clear cache", func(t *testing.T) {

		for i := 0; i < 100; i++ {
			c.Set(fmt.Sprintf("%d", i), []byte{byte(i)})
		}
		c.Clear()
	})
}
